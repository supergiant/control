package pool

type Worker struct {
	taskChan chan *Task
	doneChan chan struct{}
}

func (w *Worker) Run() {
	for {
		select {
		case task := <-w.taskChan:
			err := task.fn()

			if err != nil {
				task.result <- err
			}
			close(task.result)
		case <-w.doneChan:
			return
		}
	}
}
