package gce

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2/jwt"
	computev1 "google.golang.org/api/compute/v1"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// Provider Holds DO account info.
type Provider struct {
	Core   *core.Core
	Client func(*model.Kube) (*computev1.Service, error)
}

// ValidateAccount Valitades Open Stack account info.
func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	client, err := p.Client(&model.Kube{CloudAccount: m})
	if err != nil {
		return err
	}

	// find the core os image.
	_, err = client.Images.GetFromFamily("coreos-cloud", "coreos-stable").Do()
	if err != nil {
		return err
	}
	return nil
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

// Client creates the client for the provider.
func Client(kube *model.Kube) (*computev1.Service, error) {

	clientScopes := []string{
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
		"https://www.googleapis.com/auth/devstorage.full_control",
	}

	conf := jwt.Config{
		Email:      kube.CloudAccount.Credentials["client_email"],
		PrivateKey: []byte(kube.CloudAccount.Credentials["private_key"]),
		Scopes:     clientScopes,
		TokenURL:   kube.CloudAccount.Credentials["token_uri"],
	}

	client := conf.Client(context.Background())

	computeService, err := computev1.New(client)
	if err != nil {
		return nil, err
	}
	return computeService, nil
}

func convInstanceURLtoString(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}

func etcdToken(num string) (string, error) {
	resp, err := http.Get("https://discovery.etcd.io/new?size=" + num + "")
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
