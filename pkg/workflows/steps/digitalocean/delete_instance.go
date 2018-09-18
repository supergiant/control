package digitalocean

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/supergiant/supergiant/pkg/clouds/digitaloceanSDK"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type DeleteInstaceStep struct {
	getDeleteService func(string) DeleteService
	timeout          time.Duration
}

func NewDeleteInstanceStep(timeout time.Duration) *DeleteInstaceStep {
	return &DeleteInstaceStep{
		timeout: timeout,
		getDeleteService: func(accessToken string) DeleteService {
			return digitaloceanSDK.New(accessToken).GetClient().Droplets
		},
	}
}

func (s *DeleteInstaceStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	deleteService := s.getDeleteService(config.DigitalOceanConfig.AccessToken)
	timeout := s.timeout

	for i := 0; i < 3; i++ {
		resp, err := deleteService.DeleteByTag(ctx, config.Node.Name)

		if resp.StatusCode == http.StatusNoContent {
			return err
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	return ErrTimeoutExceeded
}

func (s *DeleteInstaceStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteInstaceStep) Name() string {
	return DeleteInstanceStepName
}

func (s *DeleteInstaceStep) Depends() []string {
	return nil
}

func (s *DeleteInstaceStep) Description() string {
	return ""
}
