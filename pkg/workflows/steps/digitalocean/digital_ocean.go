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

const (
	StepName = "digitalOcean"
	keySize  = 4096
)

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

type KeyService interface {
	Create(context.Context, *godo.KeyCreateRequest) (*godo.Key, *godo.Response, error)
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

func (s *Step) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	c := digitaloceanSDK.New(config.DigitalOceanConfig.AccessToken).GetClient()
	config.DigitalOceanConfig.Name = util.MakeNodeName(config.ClusterName, config.IsMaster)

	fingers, err := s.createKeys(ctx, c.Keys, config)

	if err != nil {
		return err
	}

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
	s.tagDroplet(ctx, c.Tags, droplet.ID, tags)

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

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *Step) tagDroplet(ctx context.Context, tagService godo.TagsService, dropletId int, tags []string) error {
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

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Depends() []string {
	return nil
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) createKey(ctx context.Context, keyService KeyService, name, publicKey string) (*godo.Key, error) {
	req := &godo.KeyCreateRequest{
		Name:      name,
		PublicKey: publicKey,
	}

	key, _, err := keyService.Create(ctx, req)

	if err != nil {
		return nil, err
	}

	return key, err
}

func (s *Step) createKeys(ctx context.Context, keyService KeyService, config *steps.Config) ([]godo.DropletCreateSSHKey, error) {
	var fingers []godo.DropletCreateSSHKey
	privateKey, publicKey, err := generateKeyPair(keySize)

	if err != nil {
		return nil, err
	}

	// Create key for provisioning
	key, err := s.createKey(ctx, keyService, fmt.Sprintf("%s-provision", config.DigitalOceanConfig.Name), publicKey)

	if err != nil {
		return nil, errors.Wrap(err, "create provision key")
	}

	config.SshConfig.PrivateKey = privateKey

	fingers = append(fingers, godo.DropletCreateSSHKey{
		Fingerprint: key.Fingerprint,
	})

	// Create user provided key
	key, _ = s.createKey(ctx, keyService, fmt.Sprintf("%s-user", config.DigitalOceanConfig.Name), config.SshConfig.PublicKey)

	// NOTE(stgleb): In case if this key is already used by user of this account
	// just compute fingerprint and pass it
	if key == nil {
		fg, _  := fingerprint(config.SshConfig.PublicKey)
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
