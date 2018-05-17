package openstack

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// Provider Holds DO account info.
type Provider struct {
	Core   *core.Core
	Client func(*model.Kube) (*gophercloud.ProviderClient, error)
}

var (
	publicDisabled = "disabled"
)

// ValidateAccount Valitades Open Stack account info.
func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	_, err := p.Client(&model.Kube{CloudAccount: m})
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
func Client(kube *model.Kube) (*gophercloud.ProviderClient, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: kube.CloudAccount.Credentials["identity_endpoint"],
		Username:         kube.CloudAccount.Credentials["username"],
		Password:         kube.CloudAccount.Credentials["password"],
		TenantID:         kube.CloudAccount.Credentials["tenant_id"],
		DomainID:         kube.CloudAccount.Credentials["domain_id"],
		DomainName:       kube.CloudAccount.Credentials["domain_name"],
	}

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ignoreErrors(err error) bool {
	okErrors := []string{
		"No suitable endpoint",
		"i/o timeout",
		"404",
		"host is down",
		"is not associated with instance",
		"Resource not found",
		"At least one of SubnetID and PortID must be provided",
	}

	for _, msg := range okErrors {
		if strings.Contains(err.Error(), msg) {
			return true
		}
	}
	return false
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
