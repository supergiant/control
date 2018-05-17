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
	"github.com/supergiant/supergiant/pkg/util"
)

// CreateNode creates a new minion on DO kubernetes cluster.
func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {

	name := m.Kube.Name + "-node" + "-" + strings.ToLower(util.RandomString(5))
	// Build template
	mversion := strings.Split(m.Kube.KubernetesVersion, ".")
	minionUserdataTemplate, err := bindata.Asset("config/providers/common/" + mversion[0] + "." + mversion[1] + "/minion.yaml")
	if err != nil {
		return err
	}
	minionTemplate, err := template.New("minion_template").Parse(string(minionUserdataTemplate))
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

	var fingers []godo.DropletCreateSSHKey
	for _, ssh := range m.Kube.DigitalOceanConfig.SSHKeyFingerprint {
		fingers = append(fingers, godo.DropletCreateSSHKey{
			Fingerprint: ssh,
		})
	}

	dropletRequest := &godo.DropletCreateRequest{
		Name:              name,
		Region:            m.Kube.DigitalOceanConfig.Region,
		Size:              m.Size,
		PrivateNetworking: true,
		UserData:          string(minionUserdata.Bytes()),
		SSHKeys:           fingers,
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
