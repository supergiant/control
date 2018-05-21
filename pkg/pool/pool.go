package pool

import (
	"github.com/satori/go.uuid"
	"sync"
)

type Pool struct {
	submitChan chan *Task
	doneChan   chan struct{}

	wg sync.WaitGroup

	workerCount       int
	activeWorkerCount int
}

func NewPool(workerCount, bufferSize int) *Pool {
	if bufferSize <= 0 {
		bufferSize = 64
	}

	return &Pool{
		submitChan:  make(chan *Task, bufferSize),
		workerCount: workerCount,
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.workerCount; i++ {

	}

	for {
		select {
		case task := <-p.submitChan:
			task.Result()
		case <-p.doneChan:
		}
	}
}

func (p *Pool) Submit(t *Task) {
	t.id = uuid.NewV4().String()
	p.submitChan <- t
}

func (p *Pool) Stop() {
	close(p.doneChan)
}
