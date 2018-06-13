package task

type Status string

const (
	StatusCreated    Status = "created"
	StatusSuccess    Status = "success"
	StatusFailed     Status = "failed"
	StatusInProgress Status = "in progress"
	StatusCanceled   Status = "canceled"
)

type Task struct {
	ID       string  `json:"id"`
	Status   Status  `json:"status"`
	Payload  *string `json:"payload"`
}

type Processor interface {
	Process(*Task) error
}

type ProcessorFunc func(*Task) error

func (f ProcessorFunc) Process(t *Task) error {
	return f(t)
}
