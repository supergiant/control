package digitalocean2

import (
	"github.com/digitalocean/godo"
	"github.com/supergiant/supergiant/pkg/model"
	"context"
	"golang.org/x/oauth2"
	"github.com/pkg/errors"
	"net/http"
	"github.com/supergiant/supergiant/pkg/profile"
	"strings"
	"github.com/supergiant/supergiant/pkg/util"
)

type DropletID int

type DigitalOceanClient interface {
	NewDroplet(*model.Kube, string, context.Context) (*godo.Droplet, error)
	DeleteDroplet(DropletID, context.Context) (error)
}

type DOClient struct {
	Client   *godo.Client
	Profiles profile.Interface
}

func NewClient(digitalOceanToken string) *DOClient {
	token := &TokenSource{
		AccessToken: digitalOceanToken,
	}
	oauthClient := oauth2.NewClient(context.Background(), token)
	return &DOClient{
		Client:   godo.NewClient(oauthClient),
		Profiles: &profile.Service{},
	}
}

func (d *DOClient) NewDroplet(kube *model.Kube, profileName string, ctx context.Context) (*godo.Droplet, error) {
	var fingers []godo.DropletCreateSSHKey
	for _, ssh := range kube.DigitalOceanConfig.SSHKeyFingerprint {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: ssh,
		})
	}
	userData, err := d.Profiles.GetUserData(profileName)
	if err != nil {
		return nil, errors.Wrapf(err, "Profile script not found or profile %s is unknown", profileName)
	}

	name := kube.Name + "-master" + "-" + strings.ToLower(util.RandomString(5))
	//TODO Move to profiles
	dropletRequest := &godo.DropletCreateRequest{
		Name:              name,
		Region:            kube.DigitalOceanConfig.Region,
		Size:              kube.MasterNodeSize,
		PrivateNetworking: true,
		//TODO Introduce profiles
		UserData: string(userData),
		SSHKeys:  fingers,
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-16-04-x64",
		},
	}

	droplet, _, err := d.Client.Droplets.Create(dropletRequest)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create droplet for cluster %s", kube.Name)
	}

	return droplet, err
}

func (d *DOClient) DeleteDroplet(ID DropletID, ctx context.Context) (error) {
	resp, err := d.Client.Droplets.Delete(int(ID))
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return errors.Wrapf(err, "error while deleting droplet %s", string(ID))
	}
	return nil
}
