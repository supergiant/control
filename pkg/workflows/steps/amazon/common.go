package amazon

import (
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/clouds/awssdk"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

func GetSDK(cfg steps.AWSConfig) (*awssdk.SDK, error) {
	sdk, err := awssdk.New(cfg.Region, cfg.KeyID, cfg.Secret, "")
	if err != nil {
		return nil, errors.Wrap(err, "aws: failed to authorize")
	}
	return sdk, nil
}
