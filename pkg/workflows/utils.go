package workflows

import (
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/storage"
)

func DeserializeTask(data []byte, repository storage.Interface) (*Task, error) {
	task := &Task{}
	err := json.Unmarshal(data, task)

	if err != nil {
		return nil, err
	}

	// Assign repository from task handler to task and restore workflow
	task.repository = repository
	task.workflow = GetWorkflow(task.Type)

	cfg := ssh.Config{
		Host:    task.Config.Node.PublicIp,
		Port:    task.Config.SshConfig.Port,
		User:    task.Config.SshConfig.User,
		Timeout: task.Config.SshConfig.Timeout,
		Key:     []byte(task.Config.SshConfig.BootstrapPrivateKey),
	}

	task.Config.Runner, err = ssh.NewRunner(cfg)

	if err != nil {
		return nil, err
	}

	return task, nil
}
