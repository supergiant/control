package pool

type Task struct {
	id string

	fn func() error

	result chan interface{}
	err    chan error
}

func NewTask(action func() error, resultChan chan interface{}, errChan chan error) *Task {
	return &Task{
		fn:     action,
		result: resultChan,
		err:    errChan,
	}
}

func (t *Task) GetID() string {
	return t.id
}

func (t *Task) Result() <-chan interface{} {
	return t.result
}

func (t *Task) Err() <-chan error {
	return t.err
}
