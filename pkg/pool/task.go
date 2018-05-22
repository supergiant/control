package pool

// Task structure that presents delayed task
type Task struct {
	fn func() error

	result chan interface{}
	err    error
}

// NewTask create new task for function and requires result chan to be passed
func NewTask(action func() error, resultChan chan interface{}) *Task {
	return &Task{
		fn:     action,
		result: resultChan,
		err:    nil,
	}
}

// Result return result channel
func (t *Task) Result() <-chan interface{} {
	return t.result
}

// Err returns err that may occur while task executes
func (t *Task) Err() error {
	return t.err
}
