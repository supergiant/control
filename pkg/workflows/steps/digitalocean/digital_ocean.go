package digitalocean

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/digitalocean/godo"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"log"
	"github.com/sirupsen/logrus"
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

	// Create
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
				config.Node = node.Node{
					Id: string(droplet.ID),
					CreatedAt: time.Now().Unix(),
					Provider: clouds.DigitalOcean,
					Region: droplet.Region.String(),
					PublicIp: getPublicIpPort(droplet.Networks.V4),
					PrivateIp: getPrivateIpPort(droplet.Networks.V4),
				}
				logrus.Println(config.Node)
				cfg := ssh.Config{
					Host: getPublicIpPort(droplet.Networks.V4),
					Port: "22",
					User: "root",
					Timeout: 10,
					Key: []byte(`id_rsa
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`),
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