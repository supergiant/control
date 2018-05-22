package pool

import (
	"errors"
	"sync"
)

var (
	errPoolIsNotRunning = errors.New("pool is not running")
)

// Pool that handles execution of delayed tasks
type Pool struct {
	workers []worker

	submitChan     chan *Task
	idleWorkerChan chan chan *Task
	doneChan       chan struct{}

	wg *sync.WaitGroup

	workerCount       int
	activeWorkerCount int
	isRunning         bool
}

// NewPool creates pool with specified buffer size for new tasks and count of workers
func NewPool(workerCount, bufferSize int) *Pool {
	if bufferSize <= 0 {
		bufferSize = 64
	}

	if workerCount <= 0 {
		workerCount = 8
	}

	return &Pool{
		doneChan:       make(chan struct{}),
		submitChan:     make(chan *Task, bufferSize),
		idleWorkerChan: make(chan chan *Task, workerCount),
		workerCount:    workerCount,
		workers:        make([]worker, workerCount),
		wg:             &sync.WaitGroup{},
	}
}

// Run pool event loop
func (p *Pool) Run() {
	p.isRunning = true
	// Spawn all workers
	for i := 0; i < p.workerCount; i++ {
		// Stop all workers with the same chan as pool stops
		p.workers[i].doneChan = p.doneChan
		go p.workers[i].Run(p.idleWorkerChan, p.wg)
	}

	for {
		select {
		case task := <-p.submitChan:
			p.wg.Add(1)
			// Get task chan from idle worker and submit task there
			taskChan := <-p.idleWorkerChan
			taskChan <- task
		case <-p.doneChan:
			return
		}
	}
}

// Submit new task to pool for execution
func (p *Pool) Submit(t *Task) error {
	if !p.isRunning {
		return errPoolIsNotRunning
	}

	p.submitChan <- t
	return nil
}

// Stop pool work and waits until all tasks are done
func (p *Pool) Stop() {
	// Wait for all submitted tasks to finish
	p.wg.Wait()
	// Stop all workers
	close(p.doneChan)
}
