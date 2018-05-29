package digitalocean

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// Provider Holds DO account info.
type Provider struct {
	Core   *core.Core
	Client func(*model.Kube) *godo.Client
}

// ValidateAccount Valitades DO account info.
func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	client := p.Client(&model.Kube{CloudAccount: m})

	_, _, err := client.Droplets.List(new(godo.ListOptions))
	return err
}

func (p *Provider) CreateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.CreateLoadBalancer(m, action)
}

func (p *Provider) UpdateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.UpdateLoadBalancer(m, action)
}

func (p *Provider) DeleteLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.DeleteLoadBalancer(m, action)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

func Client(kube *model.Kube) *godo.Client {
	token := &TokenSource{
		AccessToken: kube.CloudAccount.Credentials["token"],
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, token)
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

// Create droplet
func (p *Provider) createDroplet(client *godo.Client, action *core.Action, req *godo.DropletCreateRequest, tags []string) (droplet *godo.Droplet, publicIP string, err error) {
	// Create
	droplet, _, err = client.Droplets.Create(req)
	if err != nil {
		return nil, "", err
	}

	// Tag (TODO error handling needs work for atomicity / idempotence)
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
			// TODO
			p.Core.Log.Warnf("Failed to tag Droplet %d with value %s", droplet.ID, tag)
			// return nil, err
		}
	}

	// NOTE we have to reload to get the IP -- even with a looping wait, the
	// droplet returned from create resp never loads the IP.
	waitErr := action.CancellableWaitFor("master public IP assignment", 5*time.Minute, 5*time.Second, func() (bool, error) {
		if droplet, _, err = client.Droplets.Get(droplet.ID); err != nil {
			return false, err
		}
		if publicIP, err = droplet.PublicIPv4(); err != nil {
			return false, err
		}
		return publicIP != "", nil
	})
	if waitErr != nil {
		return nil, "", err
	}

	return droplet, publicIP, nil
}

func etcdToken(num string) (string, error) {
	resp, err := http.Get("https://.io/new?size=" + num + "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
