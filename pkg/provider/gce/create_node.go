package gce

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"google.golang.org/api/compute/v1"

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

	// find the core os image.
	image, err := client.Images.GetFromFamily("coreos-cloud", "coreos-stable").Do()
	if err != nil {
		return err
	}

	// get master machine type.
	instType, err := client.MachineTypes.Get(m.Kube.CloudAccount.Credentials["project_id"], m.Kube.GCEConfig.Zone, m.Size).Do()
	if err != nil {
		return err
	}

	prefix := "https://www.googleapis.com/compute/v1/projects/" + m.Kube.CloudAccount.Credentials["project_id"]

	procedure.AddStep("Creating Kubernetes Minion Node...", func() error {
		err := err

		m.Name = m.Kube.Name + "-minion" + "-" + strings.ToLower(util.RandomString(5))
		// Build template
		minionUserDataTemplate, err := bindata.Asset("config/providers/gce/minion.yaml")
		if err != nil {
			return err
		}
		minionTemplate, err := template.New("minion_template").Parse(string(minionUserDataTemplate))
		if err != nil {
			return err
		}
		var minionUserData bytes.Buffer
		if err = minionTemplate.Execute(&minionUserData, m); err != nil {
			return err
		}
		userData := string(minionUserData.Bytes())

		// launch master.
		role := "minion"

		instance := &compute.Instance{
			Name:         m.Name,
			Description:  "Kubernetes minion node for cluster:" + m.Name,
			MachineType:  instType.SelfLink,
			CanIpForward: true,
			Tags: &compute.Tags{
				Items: []string{"https-server", "kubernetes", "kubelet", "kubernetes-minion"},
			},
			Metadata: &compute.Metadata{
				Items: []*compute.MetadataItems{
					{
						Key:   "KubernetesCluster",
						Value: &m.Name,
					},
					{
						Key:   "Role",
						Value: &role,
					},
					{
						Key:   "user-data",
						Value: &userData,
					},
				},
			},
			Disks: []*compute.AttachedDisk{
				{
					AutoDelete: true,
					Boot:       true,
					Type:       "PERSISTENT",
					InitializeParams: &compute.AttachedDiskInitializeParams{
						DiskName:    m.Name + "-root-pd",
						SourceImage: image.SelfLink,
					},
				},
			},
			NetworkInterfaces: []*compute.NetworkInterface{
				{
					AccessConfigs: []*compute.AccessConfig{
						{
							Type: "ONE_TO_ONE_NAT",
							Name: "External NAT",
						},
					},
					Network: prefix + "/global/networks/default",
				},
			},
			ServiceAccounts: []*compute.ServiceAccount{
				{
					Email: m.Kube.CloudAccount.Credentials["client_email"],
					Scopes: []string{
						compute.DevstorageFullControlScope,
						compute.ComputeScope,
						"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
					},
				},
			},
		}

		// create the instance.
		_, serr := client.Instances.Insert(m.Kube.CloudAccount.Credentials["project_id"], m.Kube.GCEConfig.Zone, instance).Do()
		if serr != nil {
			return serr
		}

		return action.CancellableWaitFor("Kubernetes Minion Launch", 5*time.Minute, 3*time.Second, func() (bool, error) {
			resp, serr := client.Instances.Get(m.Kube.CloudAccount.Credentials["project_id"], m.Kube.GCEConfig.Zone, instance.Name).Do()
			if serr != nil {
				return false, serr
			}

			// Save Master info when ready
			if resp.Status == "RUNNING" {
				m.ProviderID = resp.SelfLink
				m.Name = resp.Name
				m.ProviderCreationTimestamp = time.Now()
				if serr := p.Core.DB.Save(m); serr != nil {
					return false, serr
				}
			}
			return resp.Status == "RUNNING", nil
		})
	})

	procedure.AddStep("Adding Kubernetes Minion to Minion Instance Group...", func() error {
		fmt.Println("Adding self link:", m.ProviderID)
		_, err = client.InstanceGroups.AddInstances(
			m.Kube.CloudAccount.Credentials["project_id"],
			m.Kube.GCEConfig.Zone,
			m.Kube.Name+"-kubernetes-minions",
			&compute.InstanceGroupsAddInstancesRequest{
				Instances: []*compute.InstanceReference{
					{
						Instance: m.ProviderID,
					},
				},
			},
		).Do()

		if err != nil {
			return err
		}
		return nil
	})

	return procedure.Run()
}
