package digitalocean

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/supergiant/supergiant/pkg/clouds/digitaloceanSDK"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type DeleteClusterStep struct {
	getDeleteService func(string) DeleteService
	timeout          time.Duration
}

func NewDeleteClusterStep(timeout time.Duration) *DeleteClusterStep {
	return &DeleteClusterStep{
		timeout: timeout,
		getDeleteService: func(accessToken string) DeleteService {
			return digitaloceanSDK.New(accessToken).GetClient().Droplets
		},
	}
}

func (s *DeleteClusterStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	deleteService := s.getDeleteService(config.DigitalOceanConfig.AccessToken)

	var (
		err     error
		resp    *godo.Response
		timeout = s.timeout
	)

	for i := 0; i < 3; i++ {
		resp, err = deleteService.DeleteByTag(ctx, config.ClusterName)

		if resp != nil && resp.StatusCode == http.StatusNoContent {
			return err
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	return err
}

func (s *DeleteClusterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteClusterStep) Name() string {
	return DeleteClusterStepName
}

func (s *DeleteClusterStep) Depends() []string {
	return nil
}

func (s *DeleteClusterStep) Description() string {
	return ""
}
