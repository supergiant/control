package digitalocean

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/digitalocean/godo"

	"github.com/supergiant/supergiant/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type DeleteMachineStep struct {
	getDeleteService func(string) DeleteService
	timeout          time.Duration
}

func NewDeleteMachineStep(timeout time.Duration) *DeleteMachineStep {
	return &DeleteMachineStep{
		timeout: timeout,
		getDeleteService: func(accessToken string) DeleteService {
			return digitaloceansdk.New(accessToken).GetClient().Droplets
		},
	}
}

func (s *DeleteMachineStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	deleteService := s.getDeleteService(config.DigitalOceanConfig.AccessToken)
	timeout := s.timeout
	var (
		err  error
		resp *godo.Response
	)

	for i := 0; i < 3; i++ {
		resp, err = deleteService.DeleteByTag(ctx, config.Node.Name)

		if resp != nil && resp.StatusCode == http.StatusNoContent {
			return err
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	return err
}

func (s *DeleteMachineStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteMachineStep) Name() string {
	return DeleteMachineStepName
}

func (s *DeleteMachineStep) Depends() []string {
	return nil
}

func (s *DeleteMachineStep) Description() string {
	return ""
}
