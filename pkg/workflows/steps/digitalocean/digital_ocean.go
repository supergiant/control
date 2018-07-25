package digitalocean

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/digitalocean/godo"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"io"
)

const StepName = "digital_ocean"

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
	storage storage.Interface

	DropletTimeout time.Duration
	CheckPeriod    time.Duration
}

func New(s storage.Interface, dropletTimeout, checkPeriod time.Duration) *Step {
	return &Step{
		storage: s,

		DropletTimeout: dropletTimeout,
		CheckPeriod:    checkPeriod,
	}
}

func (t *Step) Run(ctx context.Context, output io.Writer, config steps.Config) error {
	c := getClient(config.DOConfig.AccessToken)

	config.Name = util.MakeNodeName(config.Name, config.Role)

	var fingers []godo.DropletCreateSSHKey
	for _, ssh := range config.Fingerprints {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: ssh,
		})
	}

	dropletRequest := &godo.DropletCreateRequest{
		Name:              config.Name,
		Region:            config.Region,
		Size:              config.Size,
		PrivateNetworking: true,
		SSHKeys:           fingers,
		Image: godo.DropletCreateImage{
			Slug: config.Image,
		},
	}

	tags := []string{"Kubernetes-Cluster", config.Name}

	// Create
	droplet, _, err := c.Droplets.Create(dropletRequest)

	if err != nil {
		return errors.Wrap(err, "dropletService has returned an error in Run job")
	}

	err = t.tagDroplet(c.Tags, droplet.ID, tags)

	if err != nil {
		return errors.Wrap(err,
			fmt.Sprintf("Tagging droplet %d has failed for Run job ", droplet.ID))
	}

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
				if data, err := json.Marshal(droplet); err != nil {
					return err
				} else {
					return t.storage.Put(context.Background(), "droplet", droplet.Name, data)
				}
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
