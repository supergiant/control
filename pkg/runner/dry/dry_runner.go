package dry

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/runner"
)

type DryRunner struct {
	output *bytes.Buffer
}

func NewDryRunner() *DryRunner {
	return &DryRunner{
		output: &bytes.Buffer{},
	}
}

func (r *DryRunner) Run(cmd *runner.Command) (err error) {
	logrus.Info(cmd.Script)
	if n, err := r.output.WriteString(cmd.Script); n < len(cmd.Script) || err != nil {
		return err
	}
	logrus.Info(r.output.String())
	return nil
}

func (r *DryRunner) GetOutput() string {
	return r.output.String()
}
