package runner

import "github.com/supergiant/supergiant/pkg/runner/command"

type Runner interface {
	Run(command command.Command)
}