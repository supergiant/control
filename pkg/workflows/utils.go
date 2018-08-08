package workflows

import (
	"context"
	"encoding/json"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type cloudAccountGetter interface {
	Get(context.Context, string) (*model.CloudAccount, error)
}

// bind params uses json serializing and reflect package that is underneath
// to avoid direct access to map for getting appropriate field values.
func bindParams(params map[string]string, object interface{}) error {
	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, object)

	if err != nil {
		return err
	}

	return nil
}

// Gets cloud account from storage and fills config object with those credentials
func fillCloudAccountCredentials(ctx context.Context, getter cloudAccountGetter, config *steps.Config) error {
	cloudAccount, err := getter.Get(ctx, config.CloudAccountName)

	if err != nil {
		return nil
	}

	switch cloudAccount.Provider {
	case clouds.AWS:
	case clouds.GCE:
	case clouds.DigitalOcean:
		return bindParams(cloudAccount.Credentials, &config.DigitalOceanConfig)
	case clouds.Packet:
	}

	return nil
}
