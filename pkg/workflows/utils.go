package workflows

import (
	"encoding/json"
	"github.com/supergiant/control/pkg/runner/ssh"
	"github.com/supergiant/control/pkg/storage"
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
	// NOTE(stgleb): If step has failed on machine creation state
	// public ip will be blank and lead to error when restart
	// TODO(stgleb): Move ssh runner creation to task Restart method
	if task.Config != nil && task.Config.Node.PublicIp != "" {
		cfg := ssh.Config{
			Host:    task.Config.Node.PublicIp,
			Port:    task.Config.Kube.SSHConfig.Port,
			User:    task.Config.Kube.SSHConfig.User,
			Timeout: task.Config.Kube.SSHConfig.Timeout,
			Key:     []byte(task.Config.Kube.SSHConfig.BootstrapPrivateKey),
		}

		task.Config.Runner, err = ssh.NewRunner(cfg)

		if err != nil {
			return nil, err
		}
	}

	return task, nil
}
