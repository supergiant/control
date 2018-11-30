package digitalocean

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type DeleteKeysStep struct {
	keyService KeyService
	timeout    time.Duration
}

func NewDeleteKeysStep() *DeleteKeysStep {
	return &DeleteKeysStep{}
}

func (s *DeleteKeysStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	c := digitaloceansdk.New(config.DigitalOceanConfig.AccessToken).GetClient()

	bootstrapFg, err := fingerprint(config.SshConfig.BootstrapPublicKey)

	if err != nil {
		logrus.Debugf("error computing fingerprint of bootstrap key")
	}

	resp, err := c.Keys.DeleteByFingerprint(ctx, bootstrapFg)

	if err != nil {
		logrus.Debugf("Delete bootstrap key status %s error %s",
			resp.Status, err)
	}

	publicFg, err := fingerprint(config.SshConfig.PublicKey)

	if err != nil {
		logrus.Debugf("error computing fingerprint of public key")
	}

	resp, err = c.Keys.DeleteByFingerprint(ctx, publicFg)

	if err != nil {
		logrus.Debugf("Delete bootstrap key status %s error %s",
			resp.Status, err)
	}

	return nil
}

func (s *DeleteKeysStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *DeleteKeysStep) Name() string {
	return DeleteDeleteKeysStep
}

func (s *DeleteKeysStep) Depends() []string {
	return nil
}

func (s *DeleteKeysStep) Description() string {
	return ""
}
