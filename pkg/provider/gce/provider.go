package gce

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	compute "google.golang.org/api/compute/v1"
)

// Provider Holds DO account info.
type Provider struct {
	Core   *core.Core
	Client func(*model.Kube) (*compute.Service, error)
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

// CreateVolume createss a Volume on DO for Kubernetes
func (p *Provider) CreateVolume(m *model.Volume, action *core.Action) error {

	// fetch client.
	client, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	disk := &compute.Disk{
		Name:   m.Name,
		SizeGb: int64(m.Size),
	}
	_, err = client.Disks.Insert(m.Kube.CloudAccount.Credentials["project_id"], m.Kube.GCEConfig.Zone, disk).Do()
	if err != nil {
		return err
	}

	m.ProviderID = m.Name
	return p.Core.DB.Save(m)
}

func (p *Provider) KubernetesVolumeDefinition(m *model.Volume) *kubernetes.Volume {
	return &kubernetes.Volume{
		Name: m.Name,
		GcePersistentDisk: &kubernetes.GcePersistentDisk{
			PDName: m.ProviderID,
			FSType: "ext4",
		},
	}
}

// ResizeVolume re-sizes volume on DO kubernetes cluster.
func (p *Provider) ResizeVolume(m *model.Volume, action *core.Action) error {

	// fetch client.
	client, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	disksresizerequest := &compute.DisksResizeRequest{
		SizeGb: int64(m.Size),
	}
	_, err = client.Disks.Resize(
		m.Kube.CloudAccount.Credentials["project_id"],
		m.Kube.GCEConfig.Zone,
		m.ProviderID,
		disksresizerequest,
	).Do()
	if err != nil {
		return err
	}

	return p.Core.DB.Save(m)
}

// WaitForVolumeAvailable waits for DO volume to become available.
func (p *Provider) WaitForVolumeAvailable(m *model.Volume, action *core.Action) error {
	// fetch client.
	client, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	desc := fmt.Sprintf("GCE disk %s to be available or deleted", m.Name)
	return util.WaitFor(desc, 5*time.Minute, 10*time.Second, func() (bool, error) {
		resp, err := client.Disks.Get(
			m.Kube.CloudAccount.Credentials["project_id"],
			m.Kube.GCEConfig.Zone,
			m.ProviderID,
		).Do()
		if err != nil {
			return false, err
		}

		state := resp.Status
		return state == "READY", nil
	})
}

// DeleteVolume deletes a DO volume.
func (p *Provider) DeleteVolume(m *model.Volume, action *core.Action) error {
	// fetch client.
	client, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	_, err = client.Disks.Delete(
		m.Kube.CloudAccount.Credentials["project_id"],
		m.Kube.GCEConfig.Zone,
		m.ProviderID,
	).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") || strings.Contains(err.Error(), "notFound") {
			// it does not exist,notFound
			return nil
		}
		return err
	}

	return p.Core.DB.Save(m)
}

// CreateEntrypoint creates a new Load Balancer for Kubernetes in DO
func (p *Provider) CreateEntrypoint(m *model.Entrypoint, action *core.Action) error {
	return nil
}

// DeleteEntrypoint deletes load balancer from DO.
func (p *Provider) DeleteEntrypoint(m *model.Entrypoint, action *core.Action) error {
	return nil
}

func (p *Provider) CreateEntrypointListener(m *model.EntrypointListener, action *core.Action) error {
	return nil
}

func (p *Provider) DeleteEntrypointListener(m *model.EntrypointListener, action *core.Action) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

// Client creates the client for the provider.
func Client(kube *model.Kube) (*compute.Service, error) {

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

	client := conf.Client(oauth2.NoContext)

	computeService, err := compute.New(client)
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
