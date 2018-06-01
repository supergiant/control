package packet

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/packethost/packngo"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/template"
	"github.com/supergiant/supergiant/pkg/util"
)

// CreateKube creates a new GCE kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {
	// setup provider steps.
	procedure := &core.Procedure{
		Core:   p.Core,
		Name:   "Create Kube",
		Model:  m,
		Action: action,
	}

	// fetch client.
	client, err := p.Client(m)
	if err != nil {
		return err
	}

	// Default master count to 1
	if m.KubeMasterCount == 0 {
		m.KubeMasterCount = 1
	}

	// provision an etcd token
	url, err := etcdToken(strconv.Itoa(m.KubeMasterCount))
	if err != nil {
		return err
	}

	// save the token
	m.ETCDDiscoveryURL = url

	err = p.Core.DB.Save(m)
	if err != nil {
		return err
	}

	for i := 1; i <= m.KubeMasterCount; i++ {
		// Create master(s)
		count := strconv.Itoa(i)

		procedure.AddStep("Creating Kubernetes Master Node "+count+"...", func() error {

			project, err := getProject(m, client, m.PACKConfig.Project)
			if err != nil {
				return err
			}
			m.PACKConfig.ProjectID = project
			plan, err := getPlan(m, client, m.MasterNodeSize)
			if err != nil {
				return err
			}

			// Master name
			name := m.Name + "-master" + "-" + strings.ToLower(util.RandomString(5))

			m.MasterName = name
			// Build template
			mversion := strings.Split(m.KubernetesVersion, ".")
			masterFileName := fmt.Sprintf("config/providers/common/%s.%s/master.yaml)", mversion[0], mversion[1])
			masterTemplate, err := template.Templates.Get(masterFileName)

			if err != nil {
				return err
			}

			var masterUserdata bytes.Buffer
			if err = masterTemplate.Execute(&masterUserdata, m); err != nil {
				return err
			}
			userData := string(masterUserdata.Bytes())

			createRequest := &packngo.DeviceCreateRequest{
				HostName:     name,
				Plan:         plan,
				Facility:     m.PACKConfig.Facility,
				OS:           "coreos_stable",
				BillingCycle: "hourly",
				ProjectID:    project,
				UserData:     userData,
				Tags:         []string{"supergiant", "kubernetes", m.Name, "master"},
			}

			server, _, err := client.Devices.Create(createRequest)
			if err != nil {
				return err
			}

			return action.CancellableWaitFor("Kubernetes master launch", 10*time.Minute, 3*time.Second, func() (bool, error) {
				resp, _, serr := client.Devices.Get(server.ID)
				if serr != nil {
					return false, serr
				}

				// Save Master info when ready
				if resp.State == "active" {
					m.MasterNodes = append(m.MasterNodes, resp.Hostname)
					m.MasterPublicIP = resp.Network[0].Address
					m.MasterPrivateIP = resp.Network[2].Address
					if serr := p.Core.DB.Save(m); serr != nil {
						return false, serr
					}
				}
				return resp.State == "active", nil
			})
		})
	}

	// Create first minion//
	procedure.AddStep("creating Kubernetes minion", func() error {
		// TODO repeated in DO provider
		if err := p.Core.DB.Find(&m.Nodes, "kube_name = ?", m.Name); err != nil {
			return err
		}
		if len(m.Nodes) > 0 {
			return nil
		} //
		node := &model.Node{
			KubeName: m.Name,
			Size:     m.NodeSizes[0],
		}
		return p.Core.Nodes.Create(node)
	})

	procedure.AddStep("waiting for Kubernetes", func() error {

		k8s := p.Core.K8S(m)

		return action.CancellableWaitFor("Kubernetes API and first minion", 20*time.Minute, time.Second, func() (bool, error) {
			nodes, err := k8s.ListNodes("")
			if err != nil {
				return false, nil
			}
			return len(nodes) > 0, nil
		})
	})

	return procedure.Run()
}
