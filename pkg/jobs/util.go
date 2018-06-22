package jobs

import (
	"bytes"
	"context"
	"io"
	"text/template"

	"github.com/supergiant/supergiant/pkg/runner"
)

func runTemplate(ctx context.Context, tpl *template.Template, r runner.Runner, output io.Writer, cfg interface{}) error {
	buffer := new(bytes.Buffer)
	err := tpl.Execute(buffer, cfg)

	if err != nil {
		return err
	}

	cmd := runner.NewCommand(ctx, buffer.String(), output)
	err = r.Run(cmd)

	if err != nil {
		return err
	}

	return nil
}
