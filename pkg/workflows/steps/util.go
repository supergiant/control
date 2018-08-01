package steps

import (
	"bytes"
	"context"
	"io"
	"text/template"

	"github.com/supergiant/supergiant/pkg/runner"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

func RunTemplate(ctx context.Context, tpl *template.Template, r runner.Runner, output io.Writer, cfg interface{}) error {
	buffer := new(bytes.Buffer)
	err := tpl.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cmd := runner.NewCommand(ctx, buffer.String(), output, output)
	err = r.Run(cmd)

	if err != nil {
		return err
	}

	return nil
}

func NewSshRunner(user, host, key string, timeout int) (runner.Runner, error) {
	cfg := ssh.Config{
		User:    user,
		Host:    host,
		Timeout: timeout,
		Port:    "22",
		Key:     []byte(key),
	}

	r, err := ssh.NewRunner(cfg)
	if err != nil {
		return nil, err
	}

	return r, nil
}
