package workflows

import (
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/control/pkg/workflows/steps/network"
	"sync"

	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedKeys"
	"github.com/supergiant/control/pkg/workflows/steps/bootstraptoken"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/drain"
	"github.com/supergiant/control/pkg/workflows/steps/kubeadm"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
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

	PostProvision   = "PostProvision"
	PreProvision    = "PreProvision"
	ProvisionMaster = "ProvisionMaster"
	ProvisionNode   = "ProvisionNode"
	DeleteNode      = "DeleteNode"
	DeleteCluster   = "DeleteCluster"
	ImportCluster   = "ImportCluster"
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
		// TODO(stgleb): Provider steps should also register itsels it step map
		provider.StepCreateMachine{},
		&provider.RegisterInstanceToLoadBalancer{},
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedKeys.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(kubeadm.StepName),
		steps.GetStep(network.StepName),
		steps.GetStep(bootstraptoken.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(poststart.StepName),
		steps.GetStep(clustercheck.StepName),
	}

	nodeWorkflow := []steps.Step{
		// TODO(stgleb): Provider steps should also register theirself it step map
		provider.StepCreateMachine{},
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedKeys.StepName),
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

	m.Lock()
	defer m.Unlock()

	workflowMap[PreProvision] = preProvision
	workflowMap[ProvisionMaster] = masterWorkflow
	workflowMap[ProvisionNode] = nodeWorkflow
	workflowMap[DeleteNode] = deleteMachineWorkflow
	workflowMap[DeleteCluster] = deleteClusterWorkflow
	workflowMap[PostProvision] = postProvision
	workflowMap[ImportCluster] = importClusterWorkflow
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
