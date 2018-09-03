package workflows

import (
	"errors"
	"sync"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/certificates"
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
	"github.com/supergiant/supergiant/pkg/workflows/steps/clustercheck"
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
	prefix = "tasks"

	Cluster = "Cluster"
	DigitalOceanMaster = "DigitalOceanMaster"
	DigitalOceanNode   = "DigitalOceanNode"
)

var (
	m           sync.RWMutex
	workflowMap map[string]Workflow

	ErrUnknownProviderWorkflowType = errors.New("unknown provider_workflow type")
)

func Init() {
	workflowMap = make(map[string]Workflow)

	digitalocean.Init()
	certificates.Init()
	cni.Init()
	docker.Init()
	downloadk8sbinary.Init()
	flannel.Init()
	kubelet.Init()
	manifest.Init()
	poststart.Init()
	tiller.Init()
	etcd.Init()
	ssh.Init()
	network.Init()
	clustercheck.Init()

	digitalOceanMasterWorkflow := []steps.Step{
		steps.GetStep(digitalocean.StepName),
		steps.GetStep(ssh.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(etcd.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(network.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(poststart.StepName),
		steps.GetStep(tiller.StepName),
	}
	digitalOceanNodeWorkflow := []steps.Step{
		steps.GetStep(digitalocean.StepName),
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

	clusterWorkflow := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(clustercheck.StepName),
	}

	m.Lock()
	defer m.Unlock()

	workflowMap[Cluster] = clusterWorkflow
	workflowMap[DigitalOceanMaster] = digitalOceanMasterWorkflow
	workflowMap[DigitalOceanNode] = digitalOceanNodeWorkflow
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
