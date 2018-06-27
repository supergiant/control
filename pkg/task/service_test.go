package task

import (
	"testing"
	"time"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	cnf := &config.Config{
		Broker:        "eager",
		ResultBackend: "eager",
	}
	svc, err := NewService(cnf)
	require.NoError(t, err)

	called := false
	svc.RegisterTaskFunction("test", func(arg string) (string, error) {
		logrus.Infof("argument is %s", arg)
		called = true
		return "hello, world", nil
	})

	sig := &tasks.Signature{
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: "supergiant-test",
				Name:  "arg",
			},
		},
		Name:       "test",
		RetryCount: 2,
	}

	r, err := svc.srv.SendTask(sig)
	require.NoError(t, err)
	require.True(t, called)

	val, err := r.Get(1 * time.Second)
	require.NoError(t, err)

	require.Equal(t, 1, len(val))

	res, ok := val[0].Interface().(string)
	require.True(t, ok)

	require.Equal(t, "hello, world", res)
}
