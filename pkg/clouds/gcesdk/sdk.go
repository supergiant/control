package gcesdk

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"

	"google.golang.org/api/option"

	"github.com/supergiant/control/pkg/workflows/steps"
)

func GetClient(ctx context.Context, config steps.GCEConfig) (*compute.Service, error) {
	data, err := json.Marshal(&config.ServiceAccount)

	if err != nil {
		return nil, errors.Wrapf(err, "Error marshalling service account")
	}

	opts := option.WithCredentialsJSON(data)

	computeService, err := compute.NewService(ctx, opts)

	if err != nil {
		return nil, err
	}
	return computeService, nil
}
