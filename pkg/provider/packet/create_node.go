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
	sgtemplate "github.com/supergiant/supergiant/pkg/provider/template"
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

		var (
			userDatatemplate string
			oS               string
		)

		var minionTemplate *template.Template

		// Special config template for ARM architecture
		if m.Size == "Type 2A" {
			userDatatemplate = "config/providers/common/" + mversion[0] + "." + mversion[1] + "/arm/ubuntu/minion.yaml"
			oS = "ubuntu_17_04"

			// Build template
			minionUserdataTemplate, err := bindata.Asset(userDatatemplate)
			if err != nil {
				return err
			}
			minionTemplate, err = template.New("minion_template").Parse(string(minionUserdataTemplate))
			if err != nil {
				return err
			}
		} else {
			userDatatemplate = "config/providers/common/" + mversion[0] + "." + mversion[1] + "/minion.yaml"
			oS = "coreos_stable"
			minionTemplate = sgtemplate.Templates[userDatatemplate]
		}

		data := struct {
			*model.Node
			Token string
		}{
			m,
			m.Kube.CloudAccount.Credentials["api_token"],
		}

		var minionUserdata bytes.Buffer
		if err = minionTemplate.Execute(&minionUserdata, data); err != nil {
			return err
		}
		userData := string(minionUserdata.Bytes())

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
