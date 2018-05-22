package pool

type Task struct {
	id string

	fn func() error

	result chan interface{}
	err    error
}

func NewTask(action func() error, resultChan chan interface{}) *Task {
	return &Task{
		fn:     action,
		result: resultChan,	
		err:    nil,
	}
}

func (t *Task) GetID() string {
	return t.id
}

func (t *Task) Result() <-chan interface{} {
	return t.result
}

func (t *Task) Err() error {
	return t.err
}
