package workflows

import (
	"sync"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/amazon"
	"github.com/supergiant/supergiant/pkg/workflows/steps/certificates"
	"github.com/supergiant/supergiant/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/supergiant/pkg/workflows/steps/cni"
	"github.com/supergiant/supergiant/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/supergiant/pkg/workflows/steps/docker"
	"github.com/supergiant/supergiant/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/supergiant/pkg/workflows/steps/etcd"
	"github.com/supergiant/supergiant/pkg/workflows/steps/flannel"
	"github.com/supergiant/supergiant/pkg/workflows/steps/kubelet"
	"github.com/supergiant/supergiant/pkg/workflows/steps/manifest"
	"github.com/supergiant/supergiant/pkg/workflows/steps/network"
	"github.com/supergiant/supergiant/pkg/workflows/steps/poststart"
	"github.com/supergiant/supergiant/pkg/workflows/steps/ssh"
	"github.com/supergiant/supergiant/pkg/workflows/steps/tiller"
)

// StepStatus aggregates data that is needed to track progress
// of step to persistent storage.
type StepStatus struct {
	Status   steps.Status `json:"status"`
	StepName string       `json:"stepName"`
	ErrMsg   string       `json:"errorMessage"`
}

// Workflow is a template for doing some actions
type Workflow []steps.Step

const (
	Prefix = "tasks"

	Cluster = "Cluster"

	DigitalOceanMaster        = "DigitalOceanMaster"
	DigitalOceanNode          = "DigitalOceanNode"
	DigitalOceanDeleteNode    = "DigitalOceanDeleteNode"
	DigitalOceanDeleteCluster = "DigitalOceanDeleteCluster"
	AWSMaster                 = "AWSMaster"
	AWSNode                   = "AWSNode"
	AWSPreProvision           = "AWSPreProvisionCluster"
)

type WorkflowSet struct {
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

	digitalOceanMasterWorkflow := []steps.Step{
		steps.GetStep(digitalocean.CreateMachineStepName),
		steps.GetStep(ssh.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(etcd.StepName),
		steps.GetStep(network.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(poststart.StepName),
	}
	digitalOceanNodeWorkflow := []steps.Step{
		steps.GetStep(digitalocean.CreateMachineStepName),
		steps.GetStep(ssh.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(poststart.StepName),
	}

	awsPreProvision := []steps.Step{
		steps.GetStep(amazon.StepCreateVPC),
		steps.GetStep(amazon.StepCreateSecurityGroups),
		steps.GetStep(amazon.StepCreateSubnet),
		steps.GetStep(amazon.StepImportKeyPair),
	}

	awsMasterWorkflow := []steps.Step{
		steps.GetStep(amazon.StepCreateMachine),
		steps.GetStep(ssh.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(etcd.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(network.StepName),
		steps.GetStep(poststart.StepName),
	}

	awsNodeWorkflow := []steps.Step{
		steps.GetStep(amazon.StepCreateMachine),
		steps.GetStep(ssh.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(poststart.StepName),
	}

	commonWorkflow := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(clustercheck.StepName),
		steps.GetStep(tiller.StepName),
	}

	digitalOceanDeleteWorkflow := []steps.Step{
		steps.GetStep(digitalocean.DeleteMachineStepName),
	}

	digitalOceanDeleteClusterWorkflow := []steps.Step{
		steps.GetStep(digitalocean.DeleteClusterStepName),
	}

	m.Lock()
	defer m.Unlock()

	workflowMap[DigitalOceanDeleteNode] = digitalOceanDeleteWorkflow
	workflowMap[Cluster] = commonWorkflow
	workflowMap[DigitalOceanMaster] = digitalOceanMasterWorkflow
	workflowMap[DigitalOceanNode] = digitalOceanNodeWorkflow
	workflowMap[DigitalOceanDeleteCluster] = digitalOceanDeleteClusterWorkflow
	workflowMap[AWSMaster] = awsMasterWorkflow
	workflowMap[AWSNode] = awsNodeWorkflow
	workflowMap[AWSPreProvision] = awsPreProvision
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
