package pool

import "sync"

type Worker struct {
	doneChan chan struct{}
}

func (w *Worker) Run(taskChanChan chan chan *Task, wg *sync.WaitGroup) {
	taskChan := make(chan *Task)

	// Return task chan to pool for getting new tasks or finish in case of done
	select {
	case taskChanChan <- taskChan:
	case <-w.doneChan:
		return
	}
	// Wait for new tasks and execute them or stop
	for {
		select {
		case task := <-taskChan:
			if err := task.fn(); err != nil {
				task.err = err
			}

			close(task.result)
			wg.Done()
		case <-w.doneChan:
			return
		}
	}
}
