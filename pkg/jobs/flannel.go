package jobs

import (
	"io"
	"text/template"

	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

type FlannelJobConfig struct {
	Version     string
	Arch        string
	Network     string
	NetworkType string
}

type FlannelJob struct {
	scriptTemplate *template.Template
	runner         runner.Runner

	out io.Writer
	err io.Writer
}

func NewFlannelJob(tpl *template.Template, outStream, errStream io.Writer, cfg *ssh.Config) (*FlannelJob, error) {
	sshRunner, err := ssh.NewRunner(cfg)

	if err != nil {
		return nil, errors.Wrap(err, "error creating ssh runner")
	}

	return &FlannelJob{
		scriptTemplate: tpl,
		runner:         sshRunner,
		out:            outStream,
		err:            errStream,
	}, nil
}

func (i *FlannelJob) InstallFlannel(config FlannelJobConfig) {

}
