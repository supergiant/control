package task

type WorkerConfig struct {
	QueueName string
}

type Worker struct {
}

func (worker *Worker) start() error {
	return nil
}

func NewWorker(config WorkerConfig) *Worker {
	return nil
}
