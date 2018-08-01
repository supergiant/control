package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/pborman/uuid"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/certificates"
	"github.com/supergiant/supergiant/pkg/workflows/steps/cni"
	"github.com/supergiant/supergiant/pkg/workflows/steps/digitalocean"
	"github.com/supergiant/supergiant/pkg/workflows/steps/docker"
	"github.com/supergiant/supergiant/pkg/workflows/steps/flannel"
	"github.com/supergiant/supergiant/pkg/workflows/steps/kubelet"
	"github.com/supergiant/supergiant/pkg/workflows/steps/kubeletconf"
	"github.com/supergiant/supergiant/pkg/workflows/steps/kubeproxy"
	"github.com/supergiant/supergiant/pkg/workflows/steps/poststart"
	"github.com/supergiant/supergiant/pkg/workflows/steps/systemd"
	"github.com/supergiant/supergiant/pkg/workflows/steps/tiller"
)

type StepStatus struct {
	Status   steps.Status `json:"status"`
	StepName string       `json:"stepName"`
	ErrMsg   string       `json:"errorMessage"`
}

type WorkFlow struct {
	Id           string       `json:"id"`
	Type         string       `json:"type"`
	Config       steps.Config `json:"config"`
	StepStatuses []StepStatus `json:"steps"`

	workflowSteps []steps.Step
	repository    storage.Interface
}

const (
	MasterWorkFlow = "master"
	NodeWorkflow   = "node"
)

var (
	digitalOceanMaster = []steps.Step{
		steps.GetStep(digitalocean.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(kubeletconf.StepName),
		steps.GetStep(kubeproxy.StepName),
		steps.GetStep(systemd.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(tiller.StepName),
		steps.GetStep(poststart.StepName),
	}
	digitalOceanNode = []steps.Step{
		steps.GetStep(digitalocean.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(kubeletconf.StepName),
		steps.GetStep(kubeproxy.StepName),
		steps.GetStep(systemd.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(cni.StepName),
		steps.GetStep(tiller.StepName),
		steps.GetStep(poststart.StepName),
	}

	ErrUnknownWorkflowType = errors.New("unknown workflow type")
	ErrUnknownProviderType = errors.New("unknown provider type")
)

// TODO(stgleb): Add ability to pass arbitrary list of steps to create workflow by your own
func New(workflowType, providerString clouds.Name, config steps.Config, repository storage.Interface) (*WorkFlow, error) {
	id := uuid.New()

	switch providerString {
	case clouds.DigitalOcean:
		if workflowType == MasterWorkFlow {
			return &WorkFlow{
				Id:     id,
				Config: config,

				workflowSteps: digitalOceanMaster,
				repository:    repository,
			}, nil
		} else if workflowType == NodeWorkflow {
			return &WorkFlow{
				Id:     id,
				Config: config,

				workflowSteps: digitalOceanNode,
				repository:    repository,
			}, nil
		}

		return nil, ErrUnknownWorkflowType
	case clouds.AWS:
	case clouds.GCE:
	case clouds.Packet:
	}

	return nil, ErrUnknownProviderType
}

// Run executes all steps of workflow and tracks the progress in persistent storage
func (w *WorkFlow) Run(ctx context.Context, out io.Writer) chan error {
	errChan := make(chan error)

	go func() {
		// Create list of statuses to track
		for _, step := range w.workflowSteps {
			w.StepStatuses = append(w.StepStatuses, StepStatus{
				Status:   steps.StatusTodo,
				StepName: step.Name(),
				ErrMsg:   "",
			})
		}
		// Start from the first step
		w.startFrom(ctx, w.Id, out, 0, errChan)
		close(errChan)
	}()

	return errChan
}

// Restart executes workflow from the last failed step
func (w *WorkFlow) Restart(ctx context.Context, id string, out io.Writer) chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		data, err := w.repository.Get(ctx, "workflows", id)

		if err != nil {
			errChan <- err
			return
		}

		err = json.Unmarshal(data, w)

		if err != nil {
			errChan <- err
			return
		}

		i := 0
		// Skip successfully finished steps
		for index, stepStatus := range w.StepStatuses {
			if stepStatus.Status == steps.StatusError {
				i = index
				break
			}
		}
		// Start from the last failed one
		w.startFrom(ctx, id, out, i, errChan)
	}()
	return errChan
}

// start workflow execution from particular step
func (w *WorkFlow) startFrom(ctx context.Context, id string, out io.Writer, i int, errChan chan error) {
	// Start workflow from the last failed step
	for index := i; index < len(w.StepStatuses); index++ {
		step := w.workflowSteps[index]
		if err := step.Run(ctx, out, w.Config); err != nil {
			// Mark step status as error
			w.StepStatuses[index].Status = steps.StatusError
			w.StepStatuses[index].ErrMsg = err.Error()
			w.sync(ctx, id)

			errChan <- err
		} else {
			// Mark step as success
			w.StepStatuses[index].Status = steps.StatusSuccess
			w.sync(ctx, id)
		}
	}
}

// synchronize state of workflow to storage
func (w *WorkFlow) sync(ctx context.Context, id string) error {
	data, err := json.Marshal(w)

	if err != nil {
		return err
	}

	return w.repository.Put(ctx, "workflows", id, data)
}
