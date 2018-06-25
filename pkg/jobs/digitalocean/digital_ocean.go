package digitalocean

import (
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	"github.com/digitalocean/godo"
	"github.com/supergiant/supergiant/pkg/util"
)

type Job struct {
}

func (j *Job) CreateDroplet(credentials map[string]string, name, k8sVersion, region, size string, fingerprints []string) error {
	name = name + "-master" + "-" + strings.ToLower(util.RandomString(5))

	var fingers []godo.DropletCreateSSHKey
	for _, ssh := range fingerprints {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: ssh,
		})
	}

	dropletRequest := &godo.DropletCreateRequest{
		Name:              name,
		Region:            region,
		Size:              size,
		PrivateNetworking: true,
		SSHKeys:           fingers,
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-stable",
		},
	}

	client := getClient(credentials)
	tags := []string{"Kubernetes-Cluster", name}
	droplet, err := createDroplet(client, dropletRequest, tags)

	if err != nil {
		return err
	}

	for {
		droplet, _, err = client.Droplets.Get(droplet.ID)

		if err != nil {
			return err
		}

		if droplet.Status == "active" {
			//TODO(stgleb): Save droplet related data to Etcd
			privateIP, _ := droplet.PrivateIPv4()
			publicIP, _ := droplet.PublicIPv4()
			return nil
		}
	}

	return nil
}

func createDroplet(client *godo.Client, req *godo.DropletCreateRequest, tags []string) (droplet *godo.Droplet, err error) {
	// Create
	droplet, _, err = client.Droplets.Create(req)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		input := &godo.TagResourcesRequest{
			Resources: []godo.Resource{
				{
					ID:   strconv.Itoa(droplet.ID),
					Type: godo.DropletResourceType,
				},
			},
		}
		if _, err = client.Tags.TagResources(tag, input); err != nil {
			return nil, err
		}
	}

	return droplet, nil
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
