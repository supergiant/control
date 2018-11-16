package gce

import (
	"context"
	"fmt"
	"io"
	"time"

	"golang.org/x/oauth2/jwt"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const CreateInstanceStepName = "gce_create_instance"

type CreateInstanceStep struct {
	// Client creates the client for the provider.
	getClient func(context.Context, string, string, string) (*compute.Service, error)
}

func NewCreateInstanceStep() (steps.Step, error) {
	return &CreateInstanceStep{
		getClient: func(ctx context.Context, email, privateKey, tokenUri string) (*compute.Service, error) {
			clientScopes := []string{
				compute.ComputeScope,
				compute.CloudPlatformScope,
				dns.NdevClouddnsReadwriteScope,
				compute.DevstorageFullControlScope,
			}

			conf := jwt.Config{
				Email:      email,
				PrivateKey: []byte(privateKey),
				Scopes:     clientScopes,
				TokenURL:   tokenUri,
			}

			client := conf.Client(ctx)

			computeService, err := compute.New(client)
			if err != nil {
				return nil, err
			}
			return computeService, nil
		},
	}, nil
}

func Init() {
	createInstance, _ := NewCreateInstanceStep()
	deleteCluster, _ := NewDeleteClusterStep()

	steps.RegisterStep(CreateInstanceStepName, createInstance)
	steps.RegisterStep(DeleteClusterStepName, deleteCluster)
}

func (s *CreateInstanceStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// fetch client.
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
		config.GCEConfig.Zone, config.GCEConfig.Size).Do()
	if err != nil {
		return err
	}

	prefix := "https://www.googleapis.com/compute/v1/projects/" + config.GCEConfig.ProjectID

	role := "master"

	if !config.IsMaster {
		role = "node"
	}
	name := util.MakeNodeName(config.ClusterName, config.TaskID, config.IsMaster)

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
		config.GCEConfig.Zone,
		instance).Do()

	if serr != nil {
		return serr
	}

	resp, err := client.Instances.Get(config.GCEConfig.ProjectID,
		config.GCEConfig.Zone, name).Do()
	if serr != nil {
		return err
	}

	metadata.Fingerprint = resp.Metadata.Fingerprint
	_, err = client.Instances.SetMetadata(config.GCEConfig.ProjectID,
		config.GCEConfig.Zone, name, metadata).Do()

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
		Region:    config.GCEConfig.Zone,
	}

	// Update node state in cluster
	config.NodeChan() <- config.Node

	ticker := time.NewTicker(time.Second * 10)
	after := time.After(time.Minute * 5)

	for {
		select {
		case <-ticker.C:
			resp, serr := client.Instances.Get(config.GCEConfig.ProjectID,
				config.GCEConfig.Zone, instance.Name).Do()
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
