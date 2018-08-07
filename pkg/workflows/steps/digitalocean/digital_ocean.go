package digitalocean

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/digitalocean/godo"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"log"
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
	c := getClient(config.DigitalOceanConfig.AccessToken)

	config.DigitalOceanConfig.Name = util.MakeNodeName(config.DigitalOceanConfig.Name, config.DigitalOceanConfig.Role)

	var fingers []godo.DropletCreateSSHKey
	for _, ssh := range config.DigitalOceanConfig.Fingerprints {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: ssh,
		})
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

	tags := []string{"Kubernetes-Cluster", config.DigitalOceanConfig.Name}
	droplet, _, err := c.Droplets.Create(dropletRequest)

	if err != nil {
		return errors.Wrap(err, "dropletService has returned an error in Run job")
	}

	// NOTE(stgleb): ignore droplet tagging error, it always fails
	t.tagDroplet(c.Tags, droplet.ID, tags)

	after := time.After(t.DropletTimeout)
	ticker := time.NewTicker(t.CheckPeriod)

	for {
		select {
		case <-ticker.C:
			droplet, _, err = c.Droplets.Get(droplet.ID)

			if err != nil {
				return err
			}
			// Wait for droplet becomes active
			if droplet.Status == "active" {
				// Get private ip ports from droplet networks
				config.KubeProxyConfig.MasterPrivateIP = getPrivateIpPort(droplet.Networks.V4)
				config.KubeletConfig.MasterPrivateIP = getPrivateIpPort(droplet.Networks.V4)
				config.ManifestConfig.MasterHost = getPrivateIpPort(droplet.Networks.V4)
				config.EtcdConfig.MasterPrivateIP = "0.0.0.0"

				config.Node = node.Node{
					Id:        fmt.Sprintf("%d", droplet.ID),
					CreatedAt: time.Now().Unix(),
					Provider:  clouds.DigitalOcean,
					Region:    droplet.Region.Name,
					PublicIp:  getPublicIpPort(droplet.Networks.V4),
					PrivateIp: getPrivateIpPort(droplet.Networks.V4),
				}
				logrus.Println(config.Node)
				cfg := ssh.Config{
					Host:    getPublicIpPort(droplet.Networks.V4),
					Port:    "22",
					User:    "root",
					Timeout: 10,
					Key:     []byte(``),
				}
				config.Runner, err = ssh.NewRunner(cfg)

				if err != nil {
					log.Fatal(err)
				}
				return nil
			}

		case <-after:
			return ErrTimeoutExceeded
		}
	}

	return nil
}

func (t *Step) tagDroplet(tagService godo.TagsService, dropletId int, tags []string) error {
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
		if _, err := tagService.TagResources(tag, input); err != nil {
			return err
		}
	}

	return nil
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func getClient(accessToken string) *godo.Client {
	token := &TokenSource{
		AccessToken: accessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, token)
	return godo.NewClient(oauthClient)
}

func (t *Step) Name() string {
	return StepName
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
