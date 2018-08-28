package gce

import (
	"context"
	"io"
	"time"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "google_compute_engine"

type Step struct {
	// Client creates the client for the provider.
	getClient func(context.Context, string, string, string) (*compute.Service, error)
}

func New() (steps.Step, error) {
	return &Step{
		getClient: func(ctx context.Context, email, privateKey, tokenUri string) (*compute.Service, error) {
			clientScopes := []string{
				"https://www.googleapis.com/auth/compute",
				"https://www.googleapis.com/auth/cloud-platform",
				"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
				"https://www.googleapis.com/auth/devstorage.full_control",
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

func (s *Step) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// fetch client.
	client, err := s.getClient(ctx, config.GCEConfig.Email, config.GCEConfig.PrivateKey, config.GCEConfig.TokenUri)
	if err != nil {
		return err
	}

	// TODO(stgleb): probably we want to switch between projects in future
	image, err := client.Images.GetFromFamily("ubuntu-os-cloud", config.GCEConfig.ImageFamily).Do()
	if err != nil {
		return err
	}

	// get master machine type.
	instType, err := client.MachineTypes.Get(config.GCEConfig.ProjectID, config.GCEConfig.Zone, config.GCEConfig.Size).Do()
	if err != nil {
		return err
	}

	prefix := "https://www.googleapis.com/compute/v1/projects/" + config.GCEConfig.ProjectID

	var instanceGroup *compute.InstanceGroup

	if config.IsMaster {

		// Create Master Instance group.
		instanceGroup = &compute.InstanceGroup{
			Name:        config.ClusterName + "-kubernetes-masters",
			Description: "Kubernetes master group for cluster:" + config.ClusterName,
		}
	} else {
		// Create Minion Instance group
		instanceGroup = &compute.InstanceGroup{
			Name:        config.ClusterName + "-kubernetes-minions",
			Description: "Kubernetes minion group for cluster:" + config.ClusterName,
		}
	}

	group, serr := client.InstanceGroups.Insert(config.GCEConfig.ProjectID, config.GCEConfig.Zone, instanceGroup).Do()

	if serr != nil {
		return serr
	}

	config.GCEConfig.InstanceGroup = group.SelfLink

	role := "master"
	name := util.MakeNodeName(config.ClusterName, config.IsMaster)

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
				Email: config.GCEConfig.Email,
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
				},
			},
		},
	}

	// create the instance.
	_, serr = client.Instances.Insert(config.GCEConfig.ProjectID,
		config.GCEConfig.Zone,
		instance).Do()

	if serr != nil {
		return serr
	}

	ticker := time.NewTicker(time.Second * 10)
	after := time.After(time.Minute * 5)

	for {
		select {
		case <-ticker.C:
			resp, serr := client.Instances.Get(config.GCEConfig.ProjectID, config.GCEConfig.Zone, instance.Name).Do()
			if serr != nil {
				continue
			}

			// Save Master info when ready
			if resp.Status == "RUNNING" {
				n := node.Node{
					Id:        string(resp.Id),
					CreatedAt: time.Now().Unix(),
					Provider:  clouds.GCE,
					Region:    resp.Zone,
					PublicIp:  resp.NetworkInterfaces[0].AccessConfigs[0].NatIP,
					PrivateIp: resp.NetworkInterfaces[0].NetworkIP,
				}

				if config.IsMaster {
					config.AddMaster(&n)

					// Add master to instance group
					_, err = client.InstanceGroups.AddInstances(
						config.GCEConfig.ProjectID,
						config.GCEConfig.Zone,
						config.ClusterName+"-kubernetes-masters",
						&compute.InstanceGroupsAddInstancesRequest{
							Instances: []*compute.InstanceReference{
								{
									Instance: instance.SelfLink,
								},
							},
						},
					).Do()

					if err != nil {
						return err
					}
				}

				config.Node = n
			}

			return nil
		case <-after:
			return sgerrors.ErrTimeoutExceeded
		}
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Depends() []string {
	return nil
}

func (s *Step) Description() string {
	return "Google compute engine step for creating instance"
}
