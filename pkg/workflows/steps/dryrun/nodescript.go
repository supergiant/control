package dryrun

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/authorizedKeys"
	"github.com/supergiant/control/pkg/workflows/steps/certificates"
	"github.com/supergiant/control/pkg/workflows/steps/cni"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
	"github.com/supergiant/control/pkg/workflows/steps/downloadk8sbinary"
	"github.com/supergiant/control/pkg/workflows/steps/flannel"
	"github.com/supergiant/control/pkg/workflows/steps/kubelet"
	"github.com/supergiant/control/pkg/workflows/steps/manifest"
)

func NodeScript(config *steps.Config) (string, error) {
	nodeSetup := []steps.Step{
		steps.GetStep(authorizedKeys.StepName),
		steps.GetStep(downloadk8sbinary.StepName),
		steps.GetStep(certificates.StepName),
		steps.GetStep(manifest.StepName),
		steps.GetStep(flannel.StepName),
		steps.GetStep(docker.StepName),
		steps.GetStep(kubelet.StepName),
		steps.GetStep(cni.StepName),
	}
	return ScriptFor(nodeSetup, config)
}

// TODO: this is a hack, implement dryRun for workflow
func ScriptFor(setupWorkflow []steps.Step, config *steps.Config) (string, error) {
	if setupWorkflow == nil {
		return "", errors.New("workflow is empty")
	}
	if config == nil {
		return "", errors.New("config is empty")
	}

	dryRunOld := config.DryRun
	defer func() {
		config.DryRun = dryRunOld
	}()
	config.DryRun = true

	script := &bytes.Buffer{}
	script.WriteString("# Node provisioning script\n\n")

	w := &bytes.Buffer{}
	for _, step := range setupWorkflow {
		w.Reset()
		if err := step.Run(context.Background(), w, config); err != nil {
			return "", errors.New(fmt.Sprintf("run %s step: %s", step.Name(), err))
		}

		script.WriteString(fmt.Sprintf("## Step: %s\n%s\n\n", step.Name(), w.String()))
	}

	return script.String(), nil
}
