package workflows

import (
	"sync"

	"github.com/supergiant/control/pkg/workflows/statuses"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/amazon"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedKeys"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/clustercheck"
	"github.com/supergiant/control/pkg/workflows/steps/cni"
	"github.com/supergiant/control/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/etcd"
	"github.com/supergiant/control/pkg/workflows/steps/flannel"
	"github.com/supergiant/control/pkg/workflows/steps/gce"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/manifest"
	"github.com/supergiant/control/pkg/workflows/steps/network"
	"github.com/supergiant/control/pkg/workflows/steps/poststart"
	"github.com/supergiant/control/pkg/workflows/steps/prometheus"
	"github.com/supergiant/control/pkg/workflows/steps/ssh"
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
	digitalOceanNodeWorkflow := []steps.Step{
		steps.GetStep(digitalocean.CreateMachineStepName),
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

	awsPreProvision := []steps.Step{
		steps.GetStep(amazon.StepCreateVPC),
		steps.GetStep(amazon.StepCreateSecurityGroups),
		steps.GetStep(amazon.StepImportKeyPair),
	}

	awsMasterWorkflow := []steps.Step{
		steps.GetStep(amazon.StepCreateSubnet),
		steps.GetStep(amazon.StepNameCreateEC2Instance),
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedKeys.StepName),
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

	awsNodeWorkflow := []steps.Step{
		steps.GetStep(amazon.StepCreateSubnet),
		steps.GetStep(amazon.StepNameCreateEC2Instance),
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedKeys.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(poststart.StepName),
	}

	commonWorkflow := []steps.Step{
		steps.GetStep(ssh.StepName),
		steps.GetStep(clustercheck.StepName),
		steps.GetStep(tiller.StepName),
		steps.GetStep(prometheus.StepName),
	}

	digitalOceanDeleteWorkflow := []steps.Step{
		steps.GetStep(digitalocean.DeleteMachineStepName),
	}

	digitalOceanDeleteClusterWorkflow := []steps.Step{
		steps.GetStep(digitalocean.DeleteClusterMachines),
	}

	awsDeleteClusterWorkflow := []steps.Step{
		steps.GetStep(amazon.DeleteClusterMachinesStepName),
		steps.GetStep(amazon.DeleteSecurityGroupsStepName),
		steps.GetStep(amazon.DeleteSubnetsStepName),
		steps.GetStep(amazon.DeleteVPCStepName),
	}

	awsDeleteNodeWorkflow := []steps.Step{
		steps.GetStep(amazon.DeleteNodeStepName),
	}

	gceNodeWorkflow := []steps.Step{
		steps.GetStep(gce.CreateInstanceStepName),
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedKeys.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(poststart.StepName),
	}

	gceMasterWorkflow := []steps.Step{
		steps.GetStep(gce.CreateInstanceStepName),
		steps.GetStep(ssh.StepName),
		steps.GetStep(authorizedKeys.StepName),
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

	gceDeleteCluster := []steps.Step{
		steps.GetStep(gce.DeleteClusterStepName),
	}

	gceDeleteNode := []steps.Step{
		steps.GetStep(gce.DeleteNodeStepName),
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
	workflowMap[AWSDeleteCluster] = awsDeleteClusterWorkflow
	workflowMap[AWSDeleteNode] = awsDeleteNodeWorkflow
	workflowMap[GCENode] = gceNodeWorkflow
	workflowMap[GCEMaster] = gceMasterWorkflow
	workflowMap[GCEDeleteCluster] = gceDeleteCluster
	workflowMap[GCEDeleteNode] = gceDeleteNode
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
