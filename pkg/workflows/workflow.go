package workflows

import (
	"sync"

	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/control/pkg/workflows/steps/cni"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/drain"
	"github.com/supergiant/control/pkg/workflows/steps/etcd"
	"github.com/supergiant/control/pkg/workflows/steps/flannel"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/manifest"
	"github.com/supergiant/control/pkg/workflows/steps/network"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
	"github.com/supergiant/control/pkg/workflows/steps/prometheus"
	"github.com/supergiant/control/pkg/workflows/steps/provider"
	"github.com/supergiant/control/pkg/workflows/steps/ssh"
	"github.com/supergiant/control/pkg/workflows/steps/storageclass"
	"github.com/supergiant/control/pkg/workflows/steps/tiller"
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

	Cluster = "Cluster"

	DigitalOceanMaster        = "DigitalOceanMaster"
	DigitalOceanNode          = "DigitalOceanNode"
	DigitalOceanDeleteNode    = "DigitalOceanDeleteNode"
	DigitalOceanDeleteCluster = "DigitalOceanDeleteCluster"
	AWSMaster                 = "AWSMaster"
	AWSNode                   = "AWSNode"
	AWSPreProvision           = "AWSPreProvisionCluster"
	AWSDeleteCluster          = "AWSDeleteCluster"
	AWSDeleteNode             = "AWSDeleteNode"
	GCEMaster                 = "GCEMaster"
	GCENode                   = "GCENode"
	GCEDeleteCluster          = "GCEDeleteCluster"
	GCEDeleteNode             = "GCEDeleteNode"
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

	preProvision := []steps.Step{
		provider.StepPreProvision{},
	}

	masterWorkflow := []steps.Step{
		provider.StepCreateMachine{},
		steps.GetStep(ssh.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(etcd.StepName),
		steps.GetStep(network.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(poststart.StepName),
	}

	nodeWorkflow := []steps.Step{
		provider.StepCreateMachine{},
		steps.GetStep(ssh.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(poststart.StepName),
	}

	commonWorkflow := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(clustercheck.StepName),
		steps.GetStep(storageclass.StepName),
		steps.GetStep(tiller.StepName),
		steps.GetStep(prometheus.StepName),
	}

	deleteMachineWorkflow := []steps.Step{
		steps.GetStep(drain.StepName),
		provider.StepDeleteMachine{},
	}

	deleteClusterWorkflow := []steps.Step{
		provider.StepCleanUp{},
	}

	m.Lock()
	defer m.Unlock()

	workflowMap[Cluster] = commonWorkflow

	workflowMap[DigitalOceanMaster] = masterWorkflow
	workflowMap[DigitalOceanNode] = nodeWorkflow
	workflowMap[DigitalOceanDeleteNode] = deleteMachineWorkflow
	workflowMap[DigitalOceanDeleteCluster] = deleteClusterWorkflow

	workflowMap[AWSMaster] = masterWorkflow
	workflowMap[AWSNode] = nodeWorkflow
	workflowMap[AWSPreProvision] = preProvision
	workflowMap[AWSDeleteNode] = deleteMachineWorkflow
	workflowMap[AWSDeleteCluster] = deleteClusterWorkflow

	workflowMap[GCENode] = nodeWorkflow
	workflowMap[GCEMaster] = masterWorkflow
	workflowMap[GCEDeleteNode] = deleteMachineWorkflow
	workflowMap[GCEDeleteCluster] = deleteClusterWorkflow
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
