package packet

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/packethost/packngo"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

const (
	consumerToken = ""
)

// Provider Holds DO account info.
type Provider struct {
	Core   *core.Core
	Client func(*model.Kube) (*packngo.Client, error)
}

// ValidateAccount Valitades Open Stack account info.
func (p *Provider) ValidateAccount(m *model.CloudAccount) error {

	// fetch client.
	client, err := p.Client(&model.Kube{CloudAccount: m})
	if err != nil {
		return err
	}
	_, _, err = client.Projects.List()
	if err != nil {
		return err
	}
	return nil
}

// CreateLoadBalancer creates a new LoadBalancer
func (p *Provider) CreateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.CreateLoadBalancer(m, action)
}

// UpdateLoadBalancer updates a LoadBalancer configuration
func (p *Provider) UpdateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.UpdateLoadBalancer(m, action)
}

// DeleteLoadBalancer deletes a LoadBalancer
func (p *Provider) DeleteLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.DeleteLoadBalancer(m, action)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

// Client creates the client for the provider.
func Client(kube *model.Kube) (*packngo.Client, error) {
	return packngo.NewClient(
		consumerToken,
		kube.CloudAccount.Credentials["api_token"],
		cleanhttp.DefaultClient(),
	), nil
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

func getPlan(m *model.Kube, client *packngo.Client, name string) (string, error) {
	var planID string
	plans, _, err := client.Plans.List()
	if err != nil {
		return "", err
	}

	for _, plan := range plans {
		if plan.Name == m.MasterNodeSize && plan.Line == "baremetal" {
			planID = plan.ID
		}
	}
	return planID, nil
}

func getProject(m *model.Kube, client *packngo.Client, name string) (string, error) {
	var projectID string
	projects, _, err := client.Projects.List()
	if err != nil {
		return "", err
	}

	for _, project := range projects {
		if project.Name == m.PACKConfig.Project {
			projectID = project.ID
		}
	}

	if projectID == "" {
		projectID = m.PACKConfig.Project
	}
	return projectID, nil
}

func getDevice(m *model.Kube, client *packngo.Client, name string) (string, error) {
	var deviceID string

	project, err := getProject(m, client, m.PACKConfig.Project)
	if err != nil {
		return "", err
	}

	devices, _, err := client.Devices.List(project)
	if err != nil {
		return "", err
	}

	for _, device := range devices {
		if device.Hostname == name {
			deviceID = device.ID
		}
	}
	return deviceID, nil
}
