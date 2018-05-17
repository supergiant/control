package digitalocean2

import (
	"context"
	"net/http"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/util"
	"golang.org/x/oauth2"
)

type DropletID int

type DigitalOceanClient interface {
	NewDroplet(context.Context, *model.Kube, *profile.NodeProfile) (*godo.Droplet, error)
	DeleteDroplet(context.Context, DropletID) error
}

type DOClient struct {
	Client       *godo.Client
	NodeProfiles profile.NodeProfileService
}

func NewClient(digitalOceanToken string) *DOClient {
	token := &TokenSource{
		AccessToken: digitalOceanToken,
	}
	oauthClient := oauth2.NewClient(context.Background(), token)
	return &DOClient{
		Client:       godo.NewClient(oauthClient),
		NodeProfiles: &profile.NodeProfiles{},
	}
}

func (d *DOClient) NewDroplet(ctx context.Context, kube *model.Kube, profile *profile.NodeProfile) (*godo.Droplet, error) {
	var fingers []godo.DropletCreateSSHKey
	for _, ssh := range kube.DigitalOceanConfig.SSHKeyFingerprint {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: ssh,
		})
	}
	name := kube.Name + "-master" + "-" + strings.ToLower(util.RandomString(5))

	dropletRequest := &godo.DropletCreateRequest{
		Name:              name,
		Region:            kube.DigitalOceanConfig.Region,
		Size:              kube.MasterNodeSize,
		PrivateNetworking: true,
		SSHKeys:           fingers,
		Image: godo.DropletCreateImage{
			Slug: profile.Name,
		},
	}

	droplet, _, err := d.Client.Droplets.Create(dropletRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create droplet for cluster %s", kube.Name)
	}

	return droplet, err
}

func (d *DOClient) DeleteDroplet(ctx context.Context, ID DropletID) error {
	resp, err := d.Client.Droplets.Delete(int(ID))
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "error while deleting droplet %s", string(ID))
	}
	return nil
}
