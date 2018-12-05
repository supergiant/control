package gce

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/api/compute/v1"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/node"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
	"strings"
)

const CreateInstanceStepName = "gce_create_instance"

type CreateInstanceStep struct {
	// Client creates the client for the provider.
	getClient func(context.Context, string, string, string) (*compute.Service, error)
}

func NewCreateInstanceStep() (steps.Step, error) {
	return &CreateInstanceStep{
		getClient: GetClient,
	}, nil
}

func (s *CreateInstanceStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// fetch client.
	// TODO(stgleb):  Add UI and API for selecting image family
	config.GCEConfig.ImageFamily = "ubuntu-1604-lts"

	client, err := s.getClient(ctx, config.GCEConfig.ClientEmail,
		config.GCEConfig.PrivateKey, config.GCEConfig.TokenURI)
	if err != nil {
		return err
	}

	// TODO(stgleb): probably we want to switch between projects in future
	image, err := client.Images.GetFromFamily("ubuntu-os-cloud",
		config.GCEConfig.ImageFamily).Do()
	if err != nil {
		return err
	}

	// get master machine type.
	instType, err := client.MachineTypes.Get(config.GCEConfig.ProjectID,
		config.GCEConfig.AvailabilityZone, config.GCEConfig.Size).Do()
	if err != nil {
		return err
	}

	prefix := "https://www.googleapis.com/compute/v1/projects/" + config.GCEConfig.ProjectID

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
		config.SshConfig.User, config.SshConfig.BootstrapPublicKey)
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
				Network: prefix + "/global/networks/default",
			},
		},
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: config.GCEConfig.ClientEmail,
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
				},
			},
		},
	}

	// create the instance.
	_, serr := client.Instances.Insert(config.GCEConfig.ProjectID,
		config.GCEConfig.AvailabilityZone,
		instance).Do()

	if serr != nil {
		return serr
	}

	resp, err := client.Instances.Get(config.GCEConfig.ProjectID,
		config.GCEConfig.AvailabilityZone, name).Do()
	if serr != nil {
		return err
	}

	metadata.Fingerprint = resp.Metadata.Fingerprint
	_, err = client.Instances.SetMetadata(config.GCEConfig.ProjectID,
		config.GCEConfig.AvailabilityZone, name, metadata).Do()

	if err != nil {
		return err
	}

	nodeRole := node.RoleMaster

	if !config.IsMaster {
		nodeRole = node.RoleNode
	}

	config.Node = node.Node{
		ID:        string(resp.Id),
		Name:      name,
		CreatedAt: time.Now().Unix(),
		State:     node.StateBuilding,
		Role:      nodeRole,
		Provider:  clouds.GCE,
		Size:      config.GCEConfig.Size,
		// Note(stgleb):  This is a hack, we put az to region, because region is
		// cluster wide and we need az to delete instance.
		// TODO(stgleb): consider adding AZ to node struct
		Region:    config.GCEConfig.AvailabilityZone,
	}

	// Update node state in cluster
	config.NodeChan() <- config.Node

	ticker := time.NewTicker(time.Second * 10)
	after := time.After(time.Minute * 5)

	for {
		select {
		case <-ticker.C:
			resp, serr := client.Instances.Get(config.GCEConfig.ProjectID,
				config.GCEConfig.AvailabilityZone, instance.Name).Do()
			if serr != nil {
				continue
			}

			// Save Master info when ready
			if resp.Status == "RUNNING" {
				config.Node.PublicIp = resp.NetworkInterfaces[0].AccessConfigs[0].NatIP
				config.Node.PrivateIp = resp.NetworkInterfaces[0].NetworkIP
				config.Node.State = node.StateActive

				// Update node state in cluster
				config.NodeChan() <- config.Node

				if config.IsMaster {
					config.AddMaster(&config.Node)
				} else {
					config.AddNode(&config.Node)
				}

				return nil
			}
		case <-after:
			return sgerrors.ErrTimeoutExceeded
		}
	}

	return nil
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
