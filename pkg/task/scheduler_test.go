package task

import (
	"testing"
	"github.com/stretchr/testify/require"
	"time"
	"fmt"
)

type Command struct {
	Field1 string `json:"field_1"`
	Field2 string `json:"field_2"`
}

func Test(t *testing.T) {
	schedulerConfig := SchedulerConfig{}
	taskScheduler, err := NewScheduler(schedulerConfig)
	require.NoError(t, err)

	queue, err := taskScheduler.NewQueue("test")
	require.NoError(t, err)
	require.NotNil(t, queue)

	q, err := taskScheduler.GetQueue("test")
	require.NoError(t, err)
	require.NotNil(t, q)

	require.Equal(t, queue.GetName(), q.GetName())

	workerConfig := WorkerConfig{
		QueueName: "test",
	}

	worker := NewWorker(workerConfig)
	require.NotNil(t, worker)

	for i := 0; i < 10; i++ {
		err := worker.start()
		require.NoError(t, err)
	}

	taskScheduler.RegisterProcessorFunction("printArgs", ProcessorFunc(func(t *Task) error {
		fmt.Println(t.Payload)
		return nil
	}))

	task := taskScheduler.NewTask("printArgs", &Command{
		Field1: "foo",
		Field2: "bar",
	})

	for i := 0; i < 100; i++ {
		go func() { taskScheduler.Submit <- *task }()
	}

	time.Sleep(3 * time.Second)
}
