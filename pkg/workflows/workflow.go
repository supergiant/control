package workflows

import (
	"github.com/supergiant/control/pkg/workflows/steps/install_app"
	"github.com/supergiant/control/pkg/workflows/steps/openstack"
	"sync"

	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/addons"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"github.com/supergiant/control/pkg/workflows/steps/apply"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedkeys"
	"github.com/supergiant/control/pkg/workflows/steps/azure"
	"github.com/supergiant/control/pkg/workflows/steps/bootstraptoken"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/control/pkg/workflows/steps/configmap"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/drain"
	"github.com/supergiant/control/pkg/workflows/steps/evacuate"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/network"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
	"github.com/supergiant/control/pkg/workflows/steps/prometheus"
	"github.com/supergiant/control/pkg/workflows/steps/provider"
	"github.com/supergiant/control/pkg/workflows/steps/ssh"
	"github.com/supergiant/control/pkg/workflows/steps/storageclass"
	"github.com/supergiant/control/pkg/workflows/steps/tiller"
	"github.com/supergiant/control/pkg/workflows/steps/uncordon"
	"github.com/supergiant/control/pkg/workflows/steps/upgrade"
)

// StepStatus aggregates data that is needed to track progress
// of step to persistent storage.
type StepStatus struct {
	Status   statuses.Status `json:"status"`
	StepName string          `json:"stepName"`
	ErrMsg   string          `json:"errorMessage"`
}

// Workflow is a template for doing some actions
type Workflow []steps.Step

const (
	Prefix = "tasks"

	PostProvision     = "PostProvision"
	Infra             = "Infra"
	AwsInfra          = "awsInfra"
	DigitalOceanInfra = "digitaloceanInfra"
	GCEInfra          = "gceInfra"
	AzureInfra        = "azureInfra"
	OpenstackInfra    = "openstackInfra"
	InstallApp        = "installApp"

	ProvisionMaster = "ProvisionMaster"
	ProvisionNode   = "ProvisionNode"
	DeleteNode      = "DeleteNode"
	DeleteCluster   = "DeleteCluster"
	ImportCluster   = "ImportCluster"
	Upgrade         = "Upgrade"
	ApplyYaml       = "ApplyYaml"
)

type WorkflowSet struct {
	PreProvision    string
	ProvisionMaster string
	ProvisionNode   string
	DeleteNode      string
	DeleteCluster   string
}

var (
	m           sync.RWMutex
	workflowMap map[string]Workflow
)

func Init() {
	workflowMap = make(map[string]Workflow)

	awsInfra := []steps.Step{
		steps.GetStep(amazon.StepFindAMI),
		steps.GetStep(amazon.StepCreateVPC),
		steps.GetStep(amazon.StepCreateSecurityGroups),
		steps.GetStep(amazon.StepNameCreateInstanceProfiles),
		steps.GetStep(amazon.ImportKeyPairStepName),
		steps.GetStep(amazon.StepCreateInternetGateway),
		steps.GetStep(amazon.StepCreateSubnets),
		steps.GetStep(amazon.StepCreateRouteTable),
		steps.GetStep(amazon.StepAssociateRouteTable),
		steps.GetStep(amazon.StepCreateLoadBalancer),
	}

	digitalOceanInfra := []steps.Step{
		steps.GetStep(digitalocean.CreateLoadBalancerStepName),
	}

	gceInfra := []steps.Step{
		steps.GetStep(gce.CreateNetworksStepName),
		steps.GetStep(gce.CreateIPAddressStepName),
		steps.GetStep(gce.CreateTargetPullStepName),
		steps.GetStep(gce.CreateInstanceGroupsStepName),
		steps.GetStep(gce.CreateHealthCheckStepName),
		steps.GetStep(gce.CreateBackendServiceStepName),
		steps.GetStep(gce.CreateForwardingRulesStepName),
	}

	azureInfra := []steps.Step{
		steps.GetStep(azure.GetAuthorizerStepName),
		steps.GetStep(azure.CreateGroupStepName),
		steps.GetStep(azure.CreateVNetAndSubnetsStepName),
		steps.GetStep(azure.CreateSecurityGroupStepName),
		steps.GetStep(azure.CreateLBStepName),
	}

	openstackInfra := []steps.Step{
		steps.GetStep(openstack.CreateSecurityGroupStepName),
		steps.GetStep(openstack.CreateNetworkStepName),
		steps.GetStep(openstack.CreateSubnetStepName),
		steps.GetStep(openstack.CreateRouterStepName),
		steps.GetStep(openstack.CreateLoadBalancerStepName),
		steps.GetStep(openstack.CreatePoolStepName),
		steps.GetStep(openstack.CreateHealthCheckStepName),
		steps.GetStep(openstack.CreateKeyPairStepName),
		steps.GetStep(openstack.FindImageStepName),
	}

	masterWorkflow := []steps.Step{
		// TODO(stgleb): Provider steps should also register itsels it step map
		provider.StepCreateMachine{},
		&provider.RegisterInstanceToLoadBalancer{},
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedkeys.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(kubeadm.StepName),
		steps.GetStep(bootstraptoken.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(poststart.StepName),
		steps.GetStep(network.StepName),
		steps.GetStep(clustercheck.StepName),
	}

	nodeWorkflow := []steps.Step{
		// TODO(stgleb): Provider steps should also register their self it step map
		provider.StepCreateMachine{},
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedkeys.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(kubeadm.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(poststart.StepName),
	}

	postProvision := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(storageclass.StepName),
		steps.GetStep(tiller.StepName),
		steps.GetStep(prometheus.StepName),
		steps.GetStep(configmap.StepName),
		addons.Step{},
		provider.StepPostStartCluster{},
	}

	importClusterWorkflow := []steps.Step{
		provider.ImportClusterStep{},
	}

	deleteMachineWorkflow := []steps.Step{
		steps.GetStep(drain.StepName),
		provider.StepDeleteMachine{},
	}

	deleteClusterWorkflow := []steps.Step{
		provider.DeleteCluster{},
	}

	upgradeNode := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(evacuate.StepName),
		steps.GetStep(upgrade.StepName),
		steps.GetStep(uncordon.StepName),
	}

	apply := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(apply.StepName),
	}

	installApp := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(install_app.StepName),
	}

	m.Lock()
	defer m.Unlock()

	workflowMap[AwsInfra] = awsInfra
	workflowMap[DigitalOceanInfra] = digitalOceanInfra
	workflowMap[GCEInfra] = gceInfra
	workflowMap[AzureInfra] = azureInfra
	workflowMap[OpenstackInfra] = openstackInfra

	workflowMap[ProvisionMaster] = masterWorkflow
	workflowMap[ProvisionNode] = nodeWorkflow
	workflowMap[DeleteNode] = deleteMachineWorkflow
	workflowMap[DeleteCluster] = deleteClusterWorkflow
	workflowMap[PostProvision] = postProvision
	workflowMap[ImportCluster] = importClusterWorkflow
	workflowMap[Upgrade] = upgradeNode
	workflowMap[ApplyYaml] = apply
	workflowMap[InstallApp] = installApp
}

func RegisterWorkFlow(workflowName string, workflow Workflow) {
	m.Lock()
	defer m.Unlock()
	workflowMap[workflowName] = workflow
}

func GetWorkflow(workflowName string) Workflow {
	m.RLock()
	defer m.RUnlock()
	return workflowMap[workflowName]
}
