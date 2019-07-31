package digitalocean

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/digitalocean/godo"

	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type DeleteMachinesStep struct {
	getDeleteService func(string) DeleteService
	timeout          time.Duration
}

func NewDeletemachinesStep(timeout time.Duration) *DeleteMachinesStep {
	return &DeleteMachinesStep{
		timeout: timeout,
		getDeleteService: func(accessToken string) DeleteService {
			return digitaloceansdk.New(accessToken).GetClient().Droplets
		},
	}
}

func (s *DeleteMachinesStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	deleteService := s.getDeleteService(config.DigitalOceanConfig.AccessToken)

	var (
		err     error
		resp    *godo.Response
		timeout = s.timeout
	)

	for i := 0; i < 3; i++ {
		resp, err = deleteService.DeleteByTag(ctx, config.Kube.ID)

		if resp != nil && resp.StatusCode == http.StatusNoContent {
			return err
		}

		time.Sleep(timeout)
		timeout = timeout * 2
	}

	return err
}

func (s *DeleteMachinesStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteMachinesStep) Name() string {
	return DeleteClusterMachines
}

func (s *DeleteMachinesStep) Depends() []string {
	return nil
}

func (s *DeleteMachinesStep) Description() string {
	return "delete digital ocean cluster"
}
