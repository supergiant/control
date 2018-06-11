package runner

import "github.com/supergiant/supergiant/pkg/runner/command"

// Runner is interface for running command in different environment
type Runner interface {
	Run(command command.Command)
}
