package pool

import (
	"testing"
	"time"
)

func TestTaskFinished(t *testing.T) {
	cnt := 4
	pool := NewPool(cnt, cnt)
	go pool.Run()

	tasks := make([]*Task, 0, cnt)

	for i := 0; i < len(tasks); i++ {
		result := make(chan interface{})

		f := func() error {
			time.Sleep(time.Second * 1)
			close(result)
			return nil
		}

		t := NewTask(f, result)
		pool.Submit(t)
		tasks = append(tasks, t)
	}

	pool.Stop()

	for _, task := range tasks {
		if task.err != nil {
			t.Errorf("Unexpected error while running task %v", task.err)
		}

		select {
		case <-task.result:
		case <-time.After(time.Second * 2):
			t.Error("Time limit exceeded for completing task")
			return
		}
	}
}
