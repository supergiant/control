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

const StepName = "digitalOcean"

var (
	// TODO(stgleb): We need global error for timeout exceeding
	ErrTimeoutExceeded = errors.New("timeout exceeded")
)

type DropletService interface {
	Get(int) (*godo.Droplet, *godo.Response, error)
	Create(*godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error)
}

type TagService interface {
	TagResources(string, *godo.TagResourcesRequest) (*godo.Response, error)
}

type Step struct {
	DropletTimeout time.Duration
	CheckPeriod    time.Duration
}

func Init() {
	steps.RegisterStep(StepName, New(time.Minute*5, time.Second*5))
}

func New(dropletTimeout, checkPeriod time.Duration) *Step {
	return &Step{
		DropletTimeout: dropletTimeout,
		CheckPeriod:    checkPeriod,
	}
}

func (t *Step) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	c := digitaloceanSDK.New(config.DigitalOceanConfig.AccessToken).GetClient()

	config.DigitalOceanConfig.Name = util.MakeNodeName(config.ClusterName, config.IsMaster)

	var fingers []godo.DropletCreateSSHKey
	fingers = append(fingers, godo.DropletCreateSSHKey{
		Fingerprint: config.DigitalOceanConfig.Fingerprint,
	})

	dropletRequest := &godo.DropletCreateRequest{
		Name:              config.DigitalOceanConfig.Name,
		Region:            config.DigitalOceanConfig.Region,
		Size:              config.DigitalOceanConfig.Size,
		PrivateNetworking: true,
		SSHKeys:           fingers,
		Image: godo.DropletCreateImage{
			Slug: config.DigitalOceanConfig.Image,
		},
	}

	tags := []string{fmt.Sprintf("Kubernetes-Cluster-%s", config.ClusterName),
		config.DigitalOceanConfig.Name}
	droplet, _, err := c.Droplets.Create(ctx, dropletRequest)

	if err != nil {
		return errors.Wrap(err, "dropletService has returned an error in Run job")
	}

	// NOTE(stgleb): ignore droplet tagging error, it always fails
	t.tagDroplet(ctx, c.Tags, droplet.ID, tags)

	after := time.After(t.DropletTimeout)
	ticker := time.NewTicker(t.CheckPeriod)

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
					Id:          fmt.Sprintf("%d", droplet.ID),
					CreatedAt:   time.Now().Unix(),
					Role:        role,
					Provider:    clouds.DigitalOcean,
					Region:      droplet.Region.Name,
					PublicIp:    getPublicIpPort(droplet.Networks.V4),
					PrivateIp:   getPrivateIpPort(droplet.Networks.V4),
					ClusterName: config.ClusterName,
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

func (t *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (t *Step) tagDroplet(ctx context.Context, tagService godo.TagsService, dropletId int, tags []string) error {
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

func (t *Step) Name() string {
	return StepName
}

func (s *Step) Depends() []string {
	return nil
}

func (t *Step) Description() string {
	return ""
}

// Returns private ip
func getPrivateIpPort(networks []godo.NetworkV4) string {
	for _, network := range networks {
		if network.Type == "private" {
			return network.IPAddress
		}
	}

	return ""
}

// Returns public ip
func getPublicIpPort(networks []godo.NetworkV4) string {
	for _, network := range networks {
		if network.Type == "public" {
			return network.IPAddress
		}
	}

	return ""
}
