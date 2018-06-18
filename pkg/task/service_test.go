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
	mgr, err := NewService(cnf)
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
	mgr.RegisterTask("1", RunRemoteScript)

	//key, err := ioutil.ReadFile("/home/yegor/.ssh/id_rsa")
	//require.NoError(t, err)
	//
	//mgr.SendTask(context.Background(), &tasks.Signature{
	//	Name: "1",
	//	Args: []tasks.Arg{
	//		{
	//			Type: "string",
	//			Name: "script",
	//			Value: `#!/bin/bash
	//					uname -a`,
	//		},
	//		{
	//			Type:  "string",
	//			Name:  "user",
	//			Value: "root",
	//		},
	//		{
	//			Type:  "string",
	//			Name:  "host",
	//			Value: "209.97.135.160",
	//		},
	//		{
	//			Type:  "string",
	//			Name:  "cert",
	//			Value: string(key),
	//		},
	//		{
	//			Type:  "int",
	//			Name:  "timeoutSec",
	//			Value: 10,
	//		},
	//	},
	//})
	//time.Sleep(1 * time.Second)
}
