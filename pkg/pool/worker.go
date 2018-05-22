package pool

import "sync"

type Worker struct {
	doneChan chan struct{}
}

func (w *Worker) Run(taskChanChan chan chan *Task, wg *sync.WaitGroup) {
	taskChan := make(chan *Task)

	taskChanChan <- taskChan
	t := <-taskChan
	t.fn()

	for {
		select {
		case task := <-taskChan:
			if err := task.fn(); err != nil {
				task.err = err
			}

			close(task.result)
			wg.Done()
			taskChanChan <- taskChan
		case <-w.doneChan:
			return
		}
	}
}
