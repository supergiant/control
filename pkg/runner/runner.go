package runner

// Runner is interface for running command in different environment
type Runner interface {
	Run(command *Command) error
}
