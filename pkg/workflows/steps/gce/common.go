package gce

import (
	"context"
	"time"

	"golang.org/x/oauth2/jwt"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"

	"github.com/supergiant/control/pkg/workflows/steps"
)

type computeService struct {
	getFromFamily       func(context.Context, steps.GCEConfig) (*compute.Image, error)
	getMachineTypes     func(context.Context, steps.GCEConfig) (*compute.MachineType, error)
	insertInstance      func(context.Context, steps.GCEConfig, *compute.Instance) (*compute.Operation, error)
	getInstance         func(context.Context, steps.GCEConfig, string) (*compute.Instance, error)
	setInstanceMetadata func(context.Context, steps.GCEConfig, string, *compute.Metadata) (*compute.Operation, error)
	deleteInstance      func(string, string, string) (*compute.Operation, error)
}

func Init() {
	createInstance, _ := NewCreateInstanceStep(time.Second*10, time.Minute*1)
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
