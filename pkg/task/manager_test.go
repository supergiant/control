package task

import (
	"testing"

	"time"

	"context"
	"io/ioutil"
	"os"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/supergiant/pkg/runner/command"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
)

func TestNewManager(t *testing.T) {
	cnf := &config.Config{
		Broker:        "eager",
		ResultBackend: "eager",
	}
	mgr, err := NewManager(cnf)
	require.NoError(t, err)

	called := false
	mgr.RegisterTask("test", func(arg string) (string, error) {
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

	r, err := mgr.srv.SendTask(sig)
	require.NoError(t, err)
	require.True(t, called)

	val, err := r.Get(1 * time.Second)
	require.NoError(t, err)

	require.Equal(t, 1, len(val))

	res, ok := val[0].Interface().(string)
	require.True(t, ok)

	require.Equal(t, "hello, world", res)

	mgr.RegisterTask("1", func() error {
		key, err := ioutil.ReadFile("/home/yegor/.ssh/id_rsa")
		if err != nil {
			return err
		}
		runner, err := ssh.NewRunner(os.Stdout, os.Stderr, &ssh.Config{
			User:    "root",
			Host:    "209.97.135.160",
			Timeout: int(1 * time.Second),
			Port:    22,
			Key:     key,
		})
		if err != nil {
			return err
		}
		cmd := command.NewCommand(context.Background(), "ls", []string{" -al"}, os.Stdout, os.Stderr)
		return runner.Run(*cmd)
	})

	mgr.SendTask(&tasks.Signature{
		Name: "1",
	})
}
