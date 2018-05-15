package digitalocean2

import (
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"context"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/provision"
)

type Provider struct {
	DOClient  DigitalOceanClient
	Core      *core.Core
	Provision provision.Interface
}

func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	client := p.Client(m.Credentials["token"])

	_, _, err := client.Droplets.List(new(godo.ListOptions))
	if err != nil {
		return errors.Wrap(err, "error while getting list of droplets")
	}
	return nil
}

func (p *Provider) CreateKube(m *model.Kube, ac *core.Action) error {
	ctx := context.Background()

	//TODO FIXME
	droplet, err := p.DOClient.NewDroplet(m, "digitalocean/digital_ocean_ubuntu_x64_master.sh", ctx)
	if err != nil {
		return errors.Wrapf(err, "error while provisioning droplet for cluster ID %d", *m.ID)
	}

	ip, err := droplet.PublicIPv4()
	if err != nil {
		return errors.Wrapf(err, "error while creating the IP for droplet ID %d", droplet.ID)
	}
	if ip == "" {
		return errors.Wrapf(err, "error no ip assigned for droplet ID %d", droplet.ID)
	}

	err = p.Provision.CreateMaster(m, []string{ip}, ctx)
	if err != nil {
		return errors.Wrapf(err, "error while provisioning k8s master on droplet ID %d with IP %s",
			droplet.ID, ip)
	}

	return nil
}

func (p *Provider) DeleteKube(kube *model.Kube, action *core.Action) error {
	panic("implement me")
}

func (*Provider) CreateNode(n *model.Node, action *core.Action) error {
	panic("implement me")
}

func (*Provider) DeleteNode(*model.Node, *core.Action) error {
	panic("implement me")
}

func (*Provider) CreateLoadBalancer(*model.LoadBalancer, *core.Action) error {
	panic("implement me")
}

func (*Provider) UpdateLoadBalancer(*model.LoadBalancer, *core.Action) error {
	panic("implement me")
}

func (*Provider) DeleteLoadBalancer(*model.LoadBalancer, *core.Action) error {
	panic("implement me")
}

func (p *Provider) Client(token string) *godo.Client {
	oauthClient := oauth2.NewClient(context.Background(), &TokenSource{
		AccessToken: token,
	})
	oauthClient.Timeout = util.DefaultTimeout
	return godo.NewClient(oauthClient)
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
