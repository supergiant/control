package task

import (
	"sync"
	"github.com/pkg/errors"
)

type SchedulerConfig struct {
}

type Scheduler struct {
	Submit chan Task

	functions sync.Map
	queues    sync.Map
}

type Queue struct {
	name string
}

func (queue *Queue) GetName() string {
	return queue.name
}

func (sch *Scheduler) NewQueue(queueName string) (*Queue, error) {
	return &Queue{}, nil
}
func (sch *Scheduler) GetQueue(queueName string) (*Queue, error) {
	return nil, nil
}

func (sch *Scheduler) NewTask(functionName string, payload interface{}) *Task {
	return nil
}

func (sch *Scheduler) RegisterProcessorFunction(fnName string, processorFunc ProcessorFunc) error {
	if fnName == "" {
		return errors.New("function name can't be empty")
	}
	sch.functions.Store(fnName, processorFunc)
	return nil
}

func NewScheduler(config SchedulerConfig) (*Scheduler, error) {
	return &Scheduler{}, nil
}
