package dry

import (
	"bytes"
	"github.com/supergiant/control/pkg/runner"
)

type DryRunner struct{
	output *bytes.Buffer
}

func NewDryRunner() *DryRunner {
	return &DryRunner{
		output: &bytes.Buffer{},
	}
}

func (r *DryRunner) Run(cmd *runner.Command) (err error) {
	if n, err := r.output.WriteString(cmd.Script); n < len(cmd.Script) || err != nil {
		return err
	}
	return nil
}