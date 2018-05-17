package packet

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/packethost/packngo"

	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

// CreateNode creates a new minion on DO kubernetes cluster.
func (p *Provider) CreateNode(m *model.Node, action *core.Action) error {

	// setup provider steps.
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
		Model:  m,
		Action: action,
	}

	// fetch client.
	client, err := p.Client(m.Kube)
	if err != nil {
		return err
	}

	project, err := getProject(m.Kube, client, m.Kube.PACKConfig.Project)
	if err != nil {
		return err
	}
	plan, err := getPlan(m.Kube, client, m.Size)
	if err != nil {
		return err
	}
	procedure.AddStep("Creating Kubernetes Minion Node...", func() error {

		m.Name = m.Kube.Name + "-minion" + "-" + strings.ToLower(util.RandomString(5))

		mversion := strings.Split(m.Kube.KubernetesVersion, ".")

		userDatatemplate := "config/providers/common/" + mversion[0] + "." + mversion[1] + "/minion.yaml"
		oS := "coreos_stable"

		if m.Size == "Type 2A" {
			userDatatemplate = "config/providers/common/" + mversion[0] + "." + mversion[1] + "/arm/ubuntu/minion.yaml"
			oS = "ubuntu_17_04"
		}
		// Build template
		masterUserdataTemplate, err := bindata.Asset(userDatatemplate)
		if err != nil {
			return err
		}
		masterTemplate, err := template.New("master_template").Parse(string(masterUserdataTemplate))
		if err != nil {
			return err
		}

		data := struct {
			*model.Node
			Token string
		}{
			m,
			m.Kube.CloudAccount.Credentials["api_token"],
		}

		var masterUserdata bytes.Buffer
		if err = masterTemplate.Execute(&masterUserdata, data); err != nil {
			return err
		}
		userData := string(masterUserdata.Bytes())

		createRequest := &packngo.DeviceCreateRequest{
			HostName:     m.Name,
			Plan:         plan,
			Facility:     m.Kube.PACKConfig.Facility,
			OS:           oS,
			BillingCycle: "hourly",
			ProjectID:    project,
			UserData:     userData,
			Tags:         []string{"supergiant", "kubernetes", m.Name, "minion"},
		}

		server, resp, err := client.Devices.Create(createRequest)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(resp.String())
			return err
		}

		return action.CancellableWaitFor("Kubernetes Minion Launch", 30*time.Minute, 3*time.Second, func() (bool, error) {
			resp, _, serr := client.Devices.Get(server.ID)
			if serr != nil {
				return false, serr
			}

			// Save Master info when ready
			if resp.State == "active" {
				m.ProviderID = resp.ID
				m.Name = resp.Hostname
				m.ProviderCreationTimestamp = time.Now()
				if serr := p.Core.DB.Save(m); serr != nil {
					return false, serr
				}
			}
			return resp.State == "active", nil
		})
	})

	return procedure.Run()
}
