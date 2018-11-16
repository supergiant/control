package gce

import (
	"context"
	"io"

	"golang.org/x/oauth2/jwt"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const DeleteClusterStepName = "gce_delete_cluster"

type DeleteClusterStep struct {
	// Client creates the client for the provider.
	getClient func(context.Context, string, string, string) (*compute.Service, error)
}

func NewDeleteClusterStep() (steps.Step, error) {
	return &DeleteClusterStep{
		getClient: func(ctx context.Context, email, privateKey, tokenUri string) (*compute.Service, error) {
			clientScopes := []string{
				compute.ComputeScope,
				compute.CloudPlatformScope,
				dns.NdevClouddnsReadwriteScope,
				compute.DevstorageFullControlScope,
			}

			conf := jwt.Config{
				Email:      email,
				PrivateKey: []byte(privateKey),
				Scopes:     clientScopes,
				TokenURL:   tokenUri,
			}

			client := conf.Client(ctx)
			computeService, err := compute.New(client)
			if err != nil {
				return nil, err
			}
			return computeService, nil
		},
	}, nil
}

func (s *DeleteClusterStep) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	// fetch client.)
	client, err := s.getClient(ctx, config.GCEConfig.ClientEmail,
		config.GCEConfig.PrivateKey, config.GCEConfig.TokenURI)
	if err != nil {
		return err
	}

	for _, master := range config.GetMasters() {
		logrus.Debugf("Delete node %s", master.Name)

		_, serr := client.Instances.Delete(config.GCEConfig.ProjectID,
			config.GCEConfig.Zone,
			master.Name).Do()

		if serr != nil {
			return errors.Wrap(serr, "GCE delete instance")
		}
	}

	for _, node := range config.GetNodes() {
		logrus.Debugf("Delete node %s", node.Name)
		_, serr := client.Instances.Delete(config.GCEConfig.ProjectID,
			config.GCEConfig.Zone,
			node.Name).Do()

		if serr != nil {
			return errors.Wrap(serr, "GCE delete instance")
		}
	}

	return nil
}

func (s *DeleteClusterStep) Name() string {
	return DeleteClusterStepName
}

func (s *DeleteClusterStep) Depends() []string {
	return nil
}

func (s *DeleteClusterStep) Description() string {
	return "Google compute engine step for creating instance"
}

func (s *DeleteClusterStep) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}
