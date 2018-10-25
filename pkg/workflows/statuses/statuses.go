package statuses

type Status string

const (
	StatusTodo      Status = "todo"
	StatusExecuting Status = "executing"
	StatusSuccess   Status = "success"
	StatusError     Status = "error"
)
