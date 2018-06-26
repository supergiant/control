package digitalocean

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"

	"github.com/coreos/etcd/clientv3"

	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
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

type Job struct {
	storage        storage.Interface
	dropletService DropletService
	tagService     TagService
}

type Config struct {
	Name         string
	K8sVersion   string
	Region       string
	Size         string
	Role         string // master/node
	Fingerprints []string

	DropletTimeout time.Duration
	CheckPeriod    time.Duration
}

func NewJob(credentials map[string]string, cfg clientv3.Config) *Job {
	c := getClient(credentials)
	s := storage.NewETCDRepository(cfg)

	return &Job{
		storage:        s,
		dropletService: c.Droplets,
		tagService:     c.Tags,
	}
}

func (j *Job) CreateDroplet(config Config) error {
	config.Name = config.Name + "-" + config.Role + "-" + strings.ToLower(util.RandomString(5))

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
			Slug: "ubuntu-stable",
		},
	}

	tags := []string{"Kubernetes-Cluster", config.Name}

	// Create
	droplet, _, err := j.dropletService.Create(dropletRequest)

	if err != nil {
		return err
	}

	err = j.tagDroplet(droplet.ID, tags)

	if err != nil {
		return err
	}

	after := time.After(config.DropletTimeout)
	ticker := time.NewTicker(config.CheckPeriod)

	for {
		select {
		case <-ticker.C:
			droplet, _, err = j.dropletService.Get(droplet.ID)

			if err != nil {
				return err
			}
			// Wait for droplet becomes active
			if droplet.Status == "active" {
				if data, err := json.Marshal(droplet); err != nil {
					return err
				} else {
					return j.storage.Put(context.Background(), "droplet", droplet.Name, data)
				}
			}

		case <-after:
			return ErrTimeoutExceeded
		}
	}

	return nil
}

func (j *Job) tagDroplet(dropletId int, tags []string) error {
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
		if _, err := j.tagService.TagResources(tag, input); err != nil {
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

func getClient(credentials map[string]string) *godo.Client {
	token := &TokenSource{
		AccessToken: credentials["token"],
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, token)
	return godo.NewClient(oauthClient)
}
