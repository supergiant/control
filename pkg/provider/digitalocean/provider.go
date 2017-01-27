package digitalocean

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/digitalocean/godo"
	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"golang.org/x/oauth2"
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

// CreateKube creates a new DO kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
		Model:  m,
		Action: action,
	}

	client := p.Client(m)

	procedure.AddStep("creating global tags for Kube", func() error {
		// These are created once, and then attached by name to created resource
		globalTags := []string{
			"Kubernetes-Cluster",
			m.Name,
			m.Name + "-master",
			m.Name + "-minion",
		}
		for _, tag := range globalTags {
			createInput := &godo.TagCreateRequest{
				Name: tag,
			}
			if _, _, err := client.Tags.Create(createInput); err != nil {
				// TODO
				p.Core.Log.Warnf("Failed to create Digital Ocean tag '%s': %s", tag, err)
			}
		}
		return nil
	})

	procedure.AddStep("creating master", func() error {
		if m.MasterPublicIP != "" {
			return nil
		}

		// Build template
		masterUserdataTemplate, err := bindata.Asset("config/providers/digitalocean/master.yaml")
		if err != nil {
			return err
		}
		masterTemplate, err := template.New("master_template").Parse(string(masterUserdataTemplate))
		if err != nil {
			return err
		}
		var masterUserdata bytes.Buffer
		if err = masterTemplate.Execute(&masterUserdata, m); err != nil {
			return err
		}

		dropletRequest := &godo.DropletCreateRequest{
			Name:              m.Name + "-master",
			Region:            m.DigitalOceanConfig.Region,
			Size:              m.MasterNodeSize,
			PrivateNetworking: false,
			UserData:          string(masterUserdata.Bytes()),
			SSHKeys: []godo.DropletCreateSSHKey{
				{
					Fingerprint: m.DigitalOceanConfig.SSHKeyFingerprint,
				},
			},
			Image: godo.DropletCreateImage{
				Slug: "coreos-stable",
			},
		}
		tags := []string{"Kubernetes-Cluster", m.Name, dropletRequest.Name}

		masterDroplet, publicIP, err := p.createDroplet(client, action, dropletRequest, tags)
		if err != nil {
			return err
		}

		m.DigitalOceanConfig.MasterID = masterDroplet.ID
		m.MasterPublicIP = publicIP
		return nil
	})

	procedure.AddStep("building Kubernetes minion", func() error {
		// Load Nodes to see if we've already created a minion
		// TODO -- I think we can get rid of a lot of this do-unless behavior if we
		// modify Procedure to save progess on Action (which is easy to implement).
		if err := p.Core.DB.Find(&m.Nodes, "kube_name = ?", m.Name); err != nil {
			return err
		}
		if len(m.Nodes) > 0 {
			return nil
		}

		node := &model.Node{
			KubeName: m.Name,
			Kube:     m,
			Size:     m.NodeSizes[0],
		}
		return p.Core.Nodes.Create(node)
	})

	// TODO repeated in provider_aws.go
	procedure.AddStep("waiting for Kubernetes", func() error {
		return action.CancellableWaitFor("Kubernetes API and first minion", 20*time.Minute, 3*time.Second, func() (bool, error) {
			k8s := p.Core.K8S(m)
			k8sNodes, err := k8s.ListNodes("")
			if err != nil {
				return false, nil
			}
			return len(k8sNodes) > 0, nil
		})
	})

	return procedure.Run()
}

// DeleteKube deletes a DO kubernetes cluster.
func (p *Provider) DeleteKube(m *model.Kube, action *core.Action) error {
	// New Client
	client := p.Client(m)
	// Step procedure
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Delete Kube",
		Model:  m,
		Action: action,
	}

	procedure.AddStep("deleting master", func() error {
		if m.DigitalOceanConfig.MasterID == 0 {
			return nil
		}
		if _, err := client.Droplets.Delete(m.DigitalOceanConfig.MasterID); err != nil && !strings.Contains(err.Error(), "404") {
			return err
		}
		m.DigitalOceanConfig.MasterID = 0
		return nil
	})

	return procedure.Run()
}

// CreateNode creates a new minion on DO kubernetes cluster.
func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {
	// Build template
	minionUserdataTemplate, err := bindata.Asset("config/providers/digitalocean/minion.yaml")
	if err != nil {
		return err
	}
	minionTemplate, err := template.New("master_template").Parse(string(minionUserdataTemplate))
	if err != nil {
		return err
	}

	data := struct {
		*model.Node
		Token string
	}{
		m,
		m.Kube.CloudAccount.Credentials["token"],
	}

	var minionUserdata bytes.Buffer
	if err = minionTemplate.Execute(&minionUserdata, data); err != nil {
		return err
	}

	dropletRequest := &godo.DropletCreateRequest{
		Name:              m.Kube.Name + "-minion",
		Region:            m.Kube.DigitalOceanConfig.Region,
		Size:              m.Size,
		PrivateNetworking: true,
		UserData:          string(minionUserdata.Bytes()),
		SSHKeys: []godo.DropletCreateSSHKey{
			{
				Fingerprint: m.Kube.DigitalOceanConfig.SSHKeyFingerprint,
			},
		},
		Image: godo.DropletCreateImage{
			Slug: "coreos-stable",
		},
	}
	tags := []string{"Kubernetes-Cluster", m.Kube.Name, dropletRequest.Name}

	minionDroplet, publicIP, err := p.createDroplet(p.Client(m.Kube), action, dropletRequest, tags)
	if err != nil {
		return err
	}

	// Parse creation timestamp
	createdAt, err := time.Parse("2006-01-02T15:04:05Z", minionDroplet.Created)
	if err != nil {
		// TODO need to return on error here
		p.Core.Log.Warnf("Could not parse Droplet creation timestamp string '%s': %s", minionDroplet.Created, err)
	}

	// Save info before waiting on IP
	m.ProviderID = strconv.Itoa(minionDroplet.ID)
	m.ProviderCreationTimestamp = createdAt
	m.ExternalIP = publicIP
	m.Name = publicIP

	return p.Core.DB.Save(m)
}

// DeleteNode deletes a minsion on a DO kubernetes cluster.
func (p *Provider) DeleteNode(m *model.Node, action *core.Action) error {
	client := p.Client(m.Kube)

	intID, err := strconv.Atoi(m.ProviderID)
	if err != nil {
		return err
	}
	_, err = client.Droplets.Delete(intID)
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
