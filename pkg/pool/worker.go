package pool

import "sync"

type worker struct {
	doneChan chan struct{}
}

// Run main event loop of worker
func (w *worker) Run(taskChanChan chan chan *Task, wg *sync.WaitGroup) {
	taskChan := make(chan *Task)

	for {
		select {
		case task := <-taskChan:
			if err := task.fn(); err != nil {
				task.err = err
			}

			close(task.result)
			wg.Done()
			// Notify pool that worker is available now
			taskChanChan <- taskChan
		case <-w.doneChan:
			return
		}
	}
}
