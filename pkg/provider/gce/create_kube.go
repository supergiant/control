package gce

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"
	"time"

	compute "google.golang.org/api/compute/v1"

	"github.com/supergiant/supergiant/bindata"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
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

	// find the core os image.
	image, err := client.Images.GetFromFamily("coreos-cloud", "coreos-stable").Do()
	if err != nil {
		return err
	}

	// get master machine type.
	instType, err := client.MachineTypes.Get(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, m.MasterNodeSize).Do()
	if err != nil {
		return err
	}

	prefix := "https://www.googleapis.com/compute/v1/projects/" + m.CloudAccount.Credentials["project_id"]

	// Create Master Instance group.
	procedure.AddStep("Creating Kubernetes Master Instance Group...", func() error {
		instanceGroup := &compute.InstanceGroup{
			Name:        m.Name + "-kubernetes-masters",
			Description: "Kubernetes master group for cluster:" + m.Name,
		}
		group, serr := client.InstanceGroups.Insert(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, instanceGroup).Do()
		if serr != nil {
			return serr
		}

		m.GCEConfig.MasterInstanceGroup = group.SelfLink
		serr = p.Core.DB.Save(m)
		if serr != nil {
			return serr
		}
		return nil
	})

	// Create Minion Instance group
	procedure.AddStep("Creating Kubernetes Minion Instance Group...", func() error {
		instanceGroup := &compute.InstanceGroup{
			Name:        m.Name + "-kubernetes-minions",
			Description: "Kubernetes minion group for cluster:" + m.Name,
		}
		group, serr := client.InstanceGroups.Insert(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, instanceGroup).Do()
		if serr != nil {
			return serr
		}

		m.GCEConfig.MinionInstanceGroup = group.SelfLink
		serr = p.Core.DB.Save(m)
		if serr != nil {
			return serr
		}
		return nil
	})

	// Stage items for master build.

	// Default master count to 1
	if m.GCEConfig.KubeMasterCount == 0 {
		m.GCEConfig.KubeMasterCount = 1
	}

	// provision an etcd token
	url, err := etcdToken(strconv.Itoa(m.GCEConfig.KubeMasterCount))
	if err != nil {
		return err
	}

	// save the token
	m.GCEConfig.ETCDDiscoveryURL = url

	err = p.Core.DB.Save(m)
	if err != nil {
		return err
	}

	for i := 1; i <= m.GCEConfig.KubeMasterCount; i++ {
		// Create master(s)
		count := strconv.Itoa(i)
		procedure.AddStep("Creating Kubernetes Master Node "+count+"...", func() error {
			err := err

			// Master name
			name := m.Name + "-master" + "-" + strings.ToLower(util.RandomString(5))
			m.GCEConfig.MasterName = name
			// Build template
			masterUserdataTemplate, err := bindata.Asset("config/providers/gce/master.yaml")
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
			userData := string(masterUserdata.Bytes())

			// launch master.
			role := "master"
			instance := &compute.Instance{
				Name:         name,
				Description:  "Kubernetes master node for cluster:" + m.Name,
				MachineType:  instType.SelfLink,
				CanIpForward: true,
				Tags: &compute.Tags{
					Items: []string{"https-server", "kubernetes"},
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
							DiskName:    name + "-root-pd",
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
						Email: m.CloudAccount.Credentials["client_email"],
						Scopes: []string{
							compute.DevstorageFullControlScope,
							compute.ComputeScope,
						},
					},
				},
			}

			// create the instance.
			_, serr := client.Instances.Insert(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, instance).Do()
			if serr != nil {
				return serr
			}

			return action.CancellableWaitFor("Kubernetes master launch", 5*time.Minute, 3*time.Second, func() (bool, error) {
				resp, serr := client.Instances.Get(m.CloudAccount.Credentials["project_id"], m.GCEConfig.Zone, instance.Name).Do()
				if serr != nil {
					return false, serr
				}

				// Save Master info when ready
				if resp.Status == "RUNNING" {
					m.GCEConfig.MasterNodes = append(m.GCEConfig.MasterNodes, resp.SelfLink)
					m.MasterPublicIP = resp.NetworkInterfaces[0].AccessConfigs[0].NatIP
					m.GCEConfig.MasterPrivateIP = resp.NetworkInterfaces[0].NetworkIP
					if serr := p.Core.DB.Save(m); serr != nil {
						return false, serr
					}
				}
				return resp.Status == "RUNNING", nil
			})
		})
	}
	procedure.AddStep("Adding Kubernetes master(s) to instance group...", func() error {
		for _, masterSelfLink := range m.GCEConfig.MasterNodes {
			_, err = client.InstanceGroups.AddInstances(
				m.CloudAccount.Credentials["project_id"],
				m.GCEConfig.Zone,
				m.Name+"-kubernetes-masters",
				&compute.InstanceGroupsAddInstancesRequest{
					Instances: []*compute.InstanceReference{
						{
							Instance: masterSelfLink,
						},
					},
				},
			).Do()

			if err != nil {
				return err
			}

		}
		return nil
	})

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
