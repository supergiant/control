package statuses

type Status string

const (
	Todo      Status = "todo"
	Executing Status = "executing"
	Success   Status = "success"
	Error     Status = "error"
)
