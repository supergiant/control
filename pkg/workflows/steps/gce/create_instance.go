package gce

import (
	"context"
	"fmt"
	"github.com/supergiant/control/pkg/clouds/gcesdk"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const CreateInstanceStepName = "gce_create_instance"

type CreateInstanceStep struct {
	// Client creates the client for the provider.
	instanceTimeout time.Duration
	checkPeriod     time.Duration

	getComputeSvc func(context.Context, steps.GCEConfig) (*computeService, error)
}

func NewCreateInstanceStep(period, timeout time.Duration) *CreateInstanceStep {
	return &CreateInstanceStep{
		checkPeriod:     period,
		instanceTimeout: timeout,
		getComputeSvc: func(ctx context.Context, config steps.GCEConfig) (*computeService, error) {
			client, err := gcesdk.GetClient(ctx, config)

			if err != nil {
				return nil, err
			}

			return &computeService{
				getFromFamily: func(ctx context.Context, config steps.GCEConfig) (*compute.Image, error) {
					return client.Images.GetFromFamily("ubuntu-os-cloud", config.ImageFamily).Do()
				},
				getMachineTypes: func(ctx context.Context,
					config steps.GCEConfig) (*compute.MachineType, error) {
					return client.MachineTypes.Get(config.ServiceAccount.ProjectID,
						config.AvailabilityZone, config.Size).Do()
				},
				insertInstance: func(ctx context.Context,
					config steps.GCEConfig, instance *compute.Instance) (*compute.Operation, error) {
					return client.Instances.Insert(config.ServiceAccount.ProjectID,
						config.AvailabilityZone, instance).Do()
				},
				getInstance: func(ctx context.Context,
					config steps.GCEConfig, name string) (*compute.Instance, error) {
					return client.Instances.Get(config.ServiceAccount.ProjectID,
						config.AvailabilityZone, name).Do()
				},
				setInstanceMetadata: func(ctx context.Context, config steps.GCEConfig,
					name string, metadata *compute.Metadata) (*compute.Operation, error) {
					return client.Instances.SetMetadata(config.ServiceAccount.ProjectID,
						config.AvailabilityZone, name, metadata).Do()
				},
				addInstanceToTargetGroup: func(ctx context.Context, config steps.GCEConfig, targetPoolName string, request *compute.TargetPoolsAddInstanceRequest) (*compute.Operation, error) {
					return client.TargetPools.AddInstance(config.ServiceAccount.ProjectID, config.Region, config.TargetPoolName, request).Do()
				},
				addInstanceToInstanceGroup: func(ctx context.Context, config steps.GCEConfig, instanceGroup string, request *compute.InstanceGroupsAddInstancesRequest) (*compute.Operation, error) {
					return client.InstanceGroups.AddInstances(config.ServiceAccount.ProjectID, config.AvailabilityZone,
						instanceGroup, request).Do()
				},
			}, nil
		},
	}
}

func (s *CreateInstanceStep) Run(ctx context.Context, output io.Writer,
	config *steps.Config) error {
	logrus.Debugf("Step %s", CreateInstanceStepName)

	svc, err := s.getComputeSvc(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting service %v", err)
		return errors.Wrapf(err, "%s getting service caused", CreateInstanceStepName)
	}

	image, err := svc.getFromFamily(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting image from family %s %v",
			config.GCEConfig.ImageFamily, err)
		return errors.Wrapf(err, "Error getting image from family %s",
			config.GCEConfig.ImageFamily)
	}

	// get master machine type.
	instType, err := svc.getMachineTypes(ctx, config.GCEConfig)

	if err != nil {
		logrus.Errorf("Error getting machine type %v", err)
		return errors.Wrapf(err, "error gettting machine types")
	}

	role := "master"

	if !config.IsMaster {
		role = "node"
	}

	// NOTE(stgleb): Upper-case symbols are forbidden
	// Instance name must follow regexp: (?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)
	name := util.MakeNodeName(strings.ToLower(config.ClusterName),
		config.TaskID, config.IsMaster)

	// TODO(stgleb): also copy user provided ssh key
	publicKey := fmt.Sprintf("%s:%s",
		config.Kube.SSHConfig.User, config.Kube.SSHConfig.BootstrapPublicKey)
	// Put bootstrap key to instance metadata that allows ssh connection to the node
	metadata := &compute.Metadata{
		Items: []*compute.MetadataItems{
			{
				Key:   "ssh-keys",
				Value: &publicKey,
			},
		},
	}

	instance := &compute.Instance{
		Name:         name,
		Description:  "Kubernetes master node for cluster:" + config.ClusterName,
		MachineType:  instType.SelfLink,
		CanIpForward: true,
		Tags: &compute.Tags{
			Items: []string{"https-server", "kubernetes"},
		},
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "KubernetesCluster",
					Value: &name,
				},
				{
					Key:   "Role",
					Value: &role,
				},
			},
		},
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskName:    name + "-root-pd",
					SourceImage: image.SelfLink,
					DiskSizeGb:  30,
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
				Network: config.GCEConfig.NetworkLink,
				Subnetwork: config.GCEConfig.SubnetLink,
			},
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: config.GCEConfig.ServiceAccount.ClientEmail,
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
				},
			},
		},
	}

	// create the instance.
	_, err = svc.insertInstance(ctx, config.GCEConfig, instance)

	if err != nil {
		logrus.Errorf("inserting instance caused %v", err)
		return errors.Wrapf(err, "%s inserting instance",
			CreateInstanceStepName)
	}

	resp, err := svc.getInstance(ctx, config.GCEConfig, name)
	if err != nil {
		logrus.Errorf("getting instance caused %v", err)
		return errors.Wrapf(err, "%s getting instance",
			CreateInstanceStepName)
	}

	metadata.Fingerprint = resp.Metadata.Fingerprint
	_, err = svc.setInstanceMetadata(ctx, config.GCEConfig, name, metadata)

	if err != nil {
		logrus.Errorf("setting instance metadata caused %v", err)
		return errors.Wrapf(err, "%s setting instance metadata",
			CreateInstanceStepName)
	}

	nodeRole := model.RoleMaster

	if !config.IsMaster {
		nodeRole = model.RoleNode
	}

	config.Node = model.Machine{
		ID:        string(resp.Id),
		Name:      name,
		CreatedAt: time.Now().Unix(),
		State:     model.MachineStateBuilding,
		Role:      nodeRole,
		Provider:  clouds.GCE,
		Size:      config.GCEConfig.Size,
		// Note(stgleb):  This is a hack, we put az to region, because region is
		// cluster wide and we need az to delete instance.
		// TODO(stgleb): consider adding AZ to node struct
		Region: config.GCEConfig.AvailabilityZone,
	}

	// Update node state in cluster
	config.NodeChan() <- config.Node

	ticker := time.NewTicker(s.checkPeriod)
	after := time.After(s.instanceTimeout)

	for {
		select {
		case <-ticker.C:
			resp, _ := svc.getInstance(ctx, config.GCEConfig, instance.Name)

			// Save Master info when ready
			if resp != nil && resp.Status == "RUNNING" {
				config.Node.PublicIp = resp.NetworkInterfaces[0].AccessConfigs[0].NatIP
				config.Node.PrivateIp = resp.NetworkInterfaces[0].NetworkIP
				config.Node.State = model.MachineStateActive

				// Update node state in cluster
				config.NodeChan() <- config.Node

				if config.IsMaster {
					config.AddMaster(&config.Node)

					addInstanceToTargetPoolReq := &compute.TargetPoolsAddInstanceRequest{
						Instances: []*compute.InstanceReference{
							{
								Instance: resp.SelfLink,
							},
						},
					}

					logrus.Debugf("Add instance %s to target pool %s", config.Node.Name,
						config.GCEConfig.TargetPoolLink)

					_, err := svc.addInstanceToTargetGroup(ctx, config.GCEConfig,
						config.GCEConfig.TargetPoolName, addInstanceToTargetPoolReq)

					if err != nil {
						logrus.Errorf("error adding instance %s URL %s to target pool %s",
							resp.Name, resp.SelfLink, config.GCEConfig.TargetPoolName)
					}

					go func(){
						// NOTE(stgleb): This is stupid, but it works.
						// Add intstance to instance group only after provisioning has finished
						// Because of that https://cloud.google.com/load-balancing/docs/internal/setting-up-internal#test-from-backend-vms
						if !config.IsBootstrap {
							time.Sleep(time.Minute * 10)
						}
						req := &compute.InstanceGroupsAddInstancesRequest{
							Instances: []*compute.InstanceReference{
								{
									Instance: resp.SelfLink,
								},
							},
						}

						logrus.Debugf("Add instance %s to instance group %s", config.Node.Name,
							config.GCEConfig.InstanceGroupNames[config.GCEConfig.AvailabilityZone])
						_, err = svc.addInstanceToInstanceGroup(ctx, config.GCEConfig,
							config.GCEConfig.InstanceGroupNames[config.GCEConfig.AvailabilityZone], req)

						if err != nil {
							logrus.Errorf("error adding instance %s URL %s to instance group %s %v",
								resp.Name, resp.SelfLink, config.GCEConfig.InstanceGroupLinks[config.GCEConfig.AvailabilityZone], err)
						}
					}()

					if len(resp.NetworkInterfaces) > 0 {
						logrus.Debugf("Add instance name %s link %s with network interface %s subnetwork %s", resp.Name, resp.SelfLink,
							resp.NetworkInterfaces[0].Network, resp.NetworkInterfaces[0].Subnetwork)
					}
				} else {
					config.AddNode(&config.Node)
				}

				return nil
			}
		case <-after:
			return sgerrors.ErrTimeoutExceeded
		}
	}
}

func (s *CreateInstanceStep) Name() string {
	return CreateInstanceStepName
}

func (s *CreateInstanceStep) Depends() []string {
	return nil
}

func (s *CreateInstanceStep) Description() string {
	return "Google compute engine step for creating instance"
}

func (s *CreateInstanceStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
