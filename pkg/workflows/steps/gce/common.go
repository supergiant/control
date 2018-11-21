package gce

import (
	"context"

	"golang.org/x/oauth2/jwt"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

func Init() {
	createInstance, _ := NewCreateInstanceStep()
	deleteCluster, _ := NewDeleteClusterStep()
	deleteNode, _ := NewDeleteNodeStep()

	steps.RegisterStep(CreateInstanceStepName, createInstance)
	steps.RegisterStep(DeleteClusterStepName, deleteCluster)
	steps.RegisterStep(DeleteNodeStepName, deleteNode)
}

func GetClient(ctx context.Context, email, privateKey, tokenUri string) (*compute.Service, error) {
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
}
