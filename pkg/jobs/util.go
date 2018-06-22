package jobs

import (
	"bytes"
	"context"
	"io"
	"text/template"

	"github.com/supergiant/supergiant/pkg/runner"
)

func runTemplate(ctx context.Context, tpl *template.Template, r runner.Runner, outStream, errStream io.Writer, cfg JobConfig) error {
	buffer := new(bytes.Buffer)
	err := tpl.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cmd := runner.NewCommand(ctx, buffer.String(), outStream, errStream)
	err = r.Run(cmd)

	if err != nil {
		return err
	}

	return nil
}
