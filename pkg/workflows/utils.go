package workflows

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/runner/ssh"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"fmt"
)

type cloudAccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

// Gets cloud account from storage and fills config object with those credentials
func FillCloudAccountCredentials(ctx context.Context, getter cloudAccountGetter, config *steps.Config) error {
	cloudAccount, err := getter.Get(ctx, config.CloudAccountName)

	if err != nil {
		return nil
	}

	config.ManifestConfig.ProviderString = string(cloudAccount.Provider)
	config.Provider = cloudAccount.Provider

	// Bind private key to config
	util.BindParams(cloudAccount.Credentials, &config.SshConfig)

	switch cloudAccount.Provider {
	case clouds.AWS:
		return util.BindParams(cloudAccount.Credentials, &config.AWSConfig)
	case clouds.GCE:
		return util.BindParams(cloudAccount.Credentials, &config.GCEConfig)
	case clouds.DigitalOcean:
		return util.BindParams(cloudAccount.Credentials, &config.DigitalOceanConfig)
	case clouds.Packet:
		return util.BindParams(cloudAccount.Credentials, &config.PacketConfig)
	case clouds.OpenStack:
		return util.BindParams(cloudAccount.Credentials, &config.OSConfig)
	default:
		return sgerrors.ErrUnknownProvider
	}

	return nil
}

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

func MakeKeyName(name string, isUser bool) string {
	if isUser {
		return fmt.Sprintf("%s-user", name)
	}

	return fmt.Sprintf("%s-provision", name)
}
