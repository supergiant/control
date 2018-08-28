package gce

import (
	"context"
	"io"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/compute/v1"

	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "google_compute_engine"

type Step struct{
	// Client creates the client for the provider.
	getClient func(context.Context, string, string, string) (*compute.Service, error)
}

func New() (steps.Step, error) {
	return &Step{
		getClient:  func(ctx context.Context, email, privateKey, tokenUri string) (*compute.Service, error) {
			clientScopes := []string{
				"https://www.googleapis.com/auth/compute",
				"https://www.googleapis.com/auth/cloud-platform",
				"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
				"https://www.googleapis.com/auth/devstorage.full_control",
			}

				conf := jwt.Config{
				Email: email,
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


func (s *Step) Run(ctx context.Context, output io.Writer, config *steps.Config) error {
	return nil
}


func (s *Step) Name() string {
	return StepName
}

func (s *Step) Depends() []string {
	return nil
}

func (s *Step) Description() string {
	return ""
}