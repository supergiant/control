package openstack

import (
	"bytes"
	"time"

	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/template"
	"github.com/supergiant/supergiant/pkg/util"
)

// CreateNode creates a new node for a kubernetes cluster.
func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {
	// fetch an authenticated provider.
	authenticatedProvider, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	// Fetch compute client.
	computeClient, err := openstack.NewComputeV2(authenticatedProvider, gophercloud.EndpointOpts{
		Region: m.Kube.OpenStackConfig.Region,
	})
	if err != nil {
		return err
	}
	m.Name = m.Kube.Name + "-node"
	// Get and fill template
	minionTemplate := template.Templates[m.Kube.KubernetesVersion]

	var minionUserdata bytes.Buffer
	if err = minionTemplate.Execute(&minionUserdata, m); err != nil {
		return err
	}
	serverCreateOpts := servers.CreateOpts{
		ServiceClient: computeClient,
		Name:          m.Name,
		FlavorName:    m.Size,
		ImageName:     m.Kube.OpenStackConfig.ImageName,
		UserData:      minionUserdata.Bytes(),
		Networks: []servers.Network{
			{UUID: m.Kube.OpenStackConfig.NetworkID},
		},
		Metadata: map[string]string{"kubernetes-cluster": m.Kube.Name, "Role": "minion"},
	}

	// Create server
	server, err := servers.Create(computeClient, serverCreateOpts).Extract()
	if err != nil {
		return err
	}
	// Save data
	m.ProviderID = server.ID
	m.ProviderCreationTimestamp = time.Now()

	p.Core.DB.Save(m)

	// Wait for IP to be assigned.
	pNetwork := m.Kube.Name + "-network"
	duration := 2 * time.Minute
	interval := 10 * time.Second

	// TODO(stgleb): Context should be inherited from higher level context
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	waitErr := util.WaitFor(ctx, "Kubernetes Minion IP assign...", interval, func() (bool, error) {
		serverObj, _ := servers.Get(computeClient, server.ID).Extract()
		if serverObj.Addresses[pNetwork] == nil {
			return false, nil
		}
		items := serverObj.Addresses[pNetwork].([]interface{})
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			m.Name = itemMap["addr"].(string)
		}
		return true, nil
	})
	if waitErr != nil {
		return waitErr
	}

	return p.Core.DB.Save(m)
}
