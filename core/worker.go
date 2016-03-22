package core

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	interval = time.Second
)

type Performable interface {
	Perform(data []byte) error
}

type Worker struct {
	c *Client
}

func NewWorker(c *Client) *Worker {
	return &Worker{c}
}

func (w *Worker) Work() {
	for {
		// TODO printing just to show we're alive
		fmt.Print(".")

		time.Sleep(interval)

		jobs, err := w.c.Jobs().List()
		if err != nil {
			// panic(err)
			// TODO -- key does not exist, just keep going
			continue
		}

		// Find first queued job, or return.
		// Claim job and return if claim fails.
		var job *Job
		for _, j := range jobs.Items {
			if j.IsQueued() {
				job = j
				break
			}
		}
		if job == nil {
			continue
		}
		if err := job.Claim(); err != nil {
			// TODO the error here is presumed to be a CompareAndSwap error; if so,
			// we should just return. If it's another error, then this is not good.
			fmt.Println(err)
			continue
		}

		var performer Performable

		switch job.Type {
		case JobTypeCreateComponent:
			performer = CreateComponent{w.c}
		case JobTypeDestroyComponent:
			performer = DestroyComponent{w.c}
		case JobTypeDestroyApp:
			performer = DestroyApp{w.c}
		default:
			panic("Could not find job type")
		}

		// TODO
		jobstr, _ := json.Marshal(job)
		fmt.Println(fmt.Sprintf("Starting job: %s", jobstr))

		// // TODO ---------------- ideally we should not be recovering, and would be more explicit when returning errors
		// defer func() {
		// 	err = recover().(error)
		// 	recordError(job, err)
		// }()

		err = performer.Perform(job.Data)

		if err == nil {
			w.c.Jobs().Delete(job.ID) // Job is successful, delete from Queue
			continue
		}

		recordError(job, err)
	}
}

// Record error, and panic if that goes wrong
func recordError(job *Job, err error) {
	if errRecordingErr := job.RecordError(err); errRecordingErr != nil {
		panic(errRecordingErr)
	}
}
