package storageclass

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

//TODO cleanup all fakerunners and move them to separate package
type fakeRunner struct {
	errMsg string
}

func (f *fakeRunner) Run(command *runner.Command) error {
	if len(f.errMsg) > 0 {
		return errors.New(f.errMsg)
	}

	_, err := io.Copy(command.Out, strings.NewReader(command.Script))
	return err
}

func TestStep_Run(t *testing.T) {
	err := templatemanager.Init("../../../../templates")
	require.NoError(t, err)

	tpl, _ := templatemanager.GetTemplate(StepName)
	output := new(bytes.Buffer)

	tt := []struct {
		provider clouds.Name
		output   string
	}{
		{
			clouds.AWS,
			"kubernetes.io/aws-ebs",
		}, {
			clouds.DigitalOcean,
			"local-storage",
		}, {
			clouds.GCE,
			"kubernetes.io/gce-pd",
		},
	}

	for _, tc := range tt {
		cfg := steps.NewConfig("", "",
			"", profile.Profile{})
		cfg.Provider = tc.provider
		cfg.Runner = &fakeRunner{}

		task := &Step{
			tpl,
		}

		err = task.Run(context.Background(), output, cfg)
		require.NoError(t, err)

		require.True(t, strings.Contains(output.String(), tc.output))
	}
}
