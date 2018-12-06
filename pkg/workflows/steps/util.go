package steps

import (
	"bytes"
	"context"
	"io"
	"text/template"

	"github.com/supergiant/control/pkg/runner"
)

func RunTemplate(ctx context.Context, tpl *template.Template, r runner.Runner, output io.Writer, cfg interface{}, dryRun bool) error {
	cmd := new(bytes.Buffer)
	if err := tpl.Execute(cmd, cfg); err != nil {
		return err
	}
	// TODO: this is a hack. Add DryRun to the step interface?
	if dryRun {
		_, err := output.Write(cmd.Bytes())
		return err
	}
	return RunCommand(ctx, r, cmd.String(), output)
}

func RunCommand(ctx context.Context, r runner.Runner, command string, output io.Writer) error {
	resultChan := make(chan error)

	go func() {
		cmd, err := runner.NewCommand(ctx, command, output, output)

		if err != nil {
			resultChan <- err
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
