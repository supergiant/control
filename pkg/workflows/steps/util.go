package steps

import (
	"bytes"
	"context"
	"io"
	"text/template"

	"github.com/supergiant/control/pkg/runner"
)

func RunTemplate(ctx context.Context, tpl *template.Template, r runner.Runner, output io.Writer, cfg interface{}) error {
	resultChan := make(chan error)

	go func() {
		buffer := new(bytes.Buffer)
		err := tpl.Execute(buffer, cfg)

		if err != nil {
			resultChan <- err
		}
		cmd, err := runner.NewCommand(ctx, buffer.String(), output, output)

		if err != nil {
			resultChan <- err
			return
		}

		err = r.Run(cmd)

		if err != nil {
			resultChan <- err
			return
		}

		close(resultChan)
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() != nil {
			return ctx.Err()
		}
	case err := <-resultChan:
		return err
	}

	return nil
}
