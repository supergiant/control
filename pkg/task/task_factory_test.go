package task

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	task, err := CreateTask("TEST",
		WithHostPort("localhost", 22),
		WithScript("ls -al"))

	require.NoError(t, err)

	require.Equal(t, "localhost", task.Args[0].Value)
	require.Equal(t, 22, task.Args[1].Value)
	require.Equal(t, "ls -al", task.Args[2].Value)
}
