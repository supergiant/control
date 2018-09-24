package digitalocean

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/clouds/digitaloceanSDK"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type CreateInstanceStep struct {
	DropletTimeout time.Duration
	CheckPeriod    time.Duration
}

func NewCreateInstanceStep(dropletTimeout, checkPeriod time.Duration) *CreateInstanceStep {
	return &CreateInstanceStep{
		DropletTimeout: dropletTimeout,
		CheckPeriod:    checkPeriod,
	}
}

func (s *CreateInstanceStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// TODO(stgleb): Extract getting digital ocean sdk to function that will allow it to be mocked.
	c := digitaloceanSDK.New(config.DigitalOceanConfig.AccessToken).GetClient()
	config.DigitalOceanConfig.Name = util.MakeNodeName(config.ClusterName, config.TaskId, config.IsMaster)

	// TODO(stgleb): Move keys creation for provisioning to provisioner to be able to get
	// this key on cluster check phase.
	fingers, err := s.createKeys(ctx, c.Keys, config)

	if err != nil {
		return err
	}

	// Delete provision key from cloud account
	defer c.Keys.DeleteByFingerprint(context.Background(), fingers[0].Fingerprint)

	tags := []string{fmt.Sprintf("Kubernetes-Cluster-%s", config.ClusterName),
		config.DigitalOceanConfig.Name}

	dropletRequest := &godo.DropletCreateRequest{
		Name:              config.DigitalOceanConfig.Name,
		Region:            config.DigitalOceanConfig.Region,
		Size:              config.DigitalOceanConfig.Size,
		PrivateNetworking: true,
		SSHKeys:           fingers,
		Image: godo.DropletCreateImage{
			Slug: config.DigitalOceanConfig.Image,
		},
		Tags: tags,
	}

	droplet, _, err := c.Droplets.Create(ctx, dropletRequest)

	if err != nil {
		return errors.Wrap(err, "dropletService has returned an error in Run job")
	}

	after := time.After(s.DropletTimeout)
	ticker := time.NewTicker(s.CheckPeriod)

	for {
		select {
		case <-ticker.C:
			droplet, _, err = c.Droplets.Get(ctx, droplet.ID)

			if err != nil {
				return err
			}
			// Wait for droplet becomes active
			if droplet.Status == "active" {
				// Get private ip ports from droplet networks

				role := "master"
				if !config.IsMaster {
					role = "node"
				}

				config.Node = node.Node{
					Id:        fmt.Sprintf("%d", droplet.ID),
					CreatedAt: time.Now().Unix(),
					Role:      role,
					Provider:  clouds.DigitalOcean,
					Region:    droplet.Region.Name,
					PublicIp:  getPublicIpPort(droplet.Networks.V4),
					PrivateIp: getPrivateIpPort(droplet.Networks.V4),
					Name:      droplet.Name,
				}

				if config.IsMaster {
					config.AddMaster(&config.Node)
				} else {
					config.AddNode(&config.Node)
				}

				logrus.Infof("Node has been created %v", config.Node)

				return nil
			}
		case <-after:
			return ErrTimeoutExceeded
		}
	}

	return nil
}

func (s *CreateInstanceStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *CreateInstanceStep) tagDroplet(ctx context.Context, tagService godo.TagsService, dropletId int, tags []string) error {
	// Tag droplet
	for _, tag := range tags {
		input := &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{
					ID:   strconv.Itoa(dropletId),
					Type: godo.DropletResourceType,
				},
			},
		}
		if _, err := tagService.TagResources(ctx, tag, input); err != nil {
			return err
		}
	}

	return nil
}

func (s *CreateInstanceStep) Name() string {
	return CreateMachineStepName
}

func (s *CreateInstanceStep) Depends() []string {
	return nil
}

func (s *CreateInstanceStep) Description() string {
	return ""
}

func (s *CreateInstanceStep) createKeys(ctx context.Context, keyService KeyService, config *steps.Config) ([]godo.DropletCreateSSHKey, error) {
	var fingers []godo.DropletCreateSSHKey

	// Create key for provisioning
	key, err := createKey(ctx, keyService,
		util.MakeKeyName(config.DigitalOceanConfig.Name, false),
		config.SshConfig.BootstrapPublicKey)

	if err != nil {
		return nil, errors.Wrap(err, "create provision key")
	}

	fingers = append(fingers, godo.DropletCreateSSHKey{
		Fingerprint: key.Fingerprint,
	})

	// Create user provided key
	key, _ = createKey(ctx, keyService,
		util.MakeKeyName(config.DigitalOceanConfig.Name, true),
		config.SshConfig.PublicKey)

	// NOTE(stgleb): In case if this key is already used by user of this account
	// just compute fingerprint and pass it
	if key == nil {
		fg, _ := fingerprint(config.SshConfig.PublicKey)
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: fg,
		})
	} else {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: key.Fingerprint,
		})
	}

	return fingers, nil
}
