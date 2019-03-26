package gce

import (
	"context"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	"encoding/json"
	"github.com/pkg/errors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type computeService struct {
	getFromFamily       func(context.Context, steps.GCEConfig) (*compute.Image, error)
	getMachineTypes     func(context.Context, steps.GCEConfig) (*compute.MachineType, error)
	insertInstance      func(context.Context, steps.GCEConfig, *compute.Instance) (*compute.Operation, error)
	getInstance         func(context.Context, steps.GCEConfig, string) (*compute.Instance, error)
	setInstanceMetadata func(context.Context, steps.GCEConfig, string, *compute.Metadata) (*compute.Operation, error)
	deleteInstance      func(string, string, string) (*compute.Operation, error)

	insertTargetPool           func(context.Context, steps.GCEConfig, *compute.TargetPool) (*compute.Operation, error)
	insertAddress              func(context.Context, steps.GCEConfig, *compute.Address) (*compute.Operation, error)
	getAddress                 func(context.Context, steps.GCEConfig, string) (*compute.Address, error)
	insertForwardingRule       func(context.Context, steps.GCEConfig, *compute.ForwardingRule) (*compute.Operation, error)
	addInstanceToTargetGroup   func(context.Context, steps.GCEConfig, string, *compute.TargetPoolsAddInstanceRequest) (*compute.Operation, error)
	addInstanceToInstanceGroup func(context.Context, steps.GCEConfig, string, *compute.InstanceGroupsAddInstancesRequest) (*compute.Operation, error)
	insertHealthCheck          func(context.Context, steps.GCEConfig, *compute.HealthCheck) (*compute.Operation, error)
	addHealthCheckToTargetPool func(context.Context, steps.GCEConfig, string, *compute.TargetPoolsAddHealthCheckRequest) (*compute.Operation, error)
	insertInstanceGroup        func(context.Context, steps.GCEConfig, *compute.InstanceGroup) (*compute.Operation, error)
	insertBackendService       func(context.Context, steps.GCEConfig, *compute.BackendService) (*compute.Operation, error)
	getHealthCheck             func(context.Context, steps.GCEConfig, string) (*compute.HealthCheck, error)
	getInstanceGroup           func(context.Context, steps.GCEConfig, string) (*compute.InstanceGroup, error)
}

func Init() {
	createInstance, _ := NewCreateInstanceStep(time.Second*10, time.Minute*1)
	deleteCluster, _ := NewDeleteClusterStep()
	deleteNode, _ := NewDeleteNodeStep()
	createTargetPool, _ := NewCreateTargetPoolStep()
	createIPAddress, _ := NewCreateAddressStep()
	createHealthCheck, _ := NewCreateHealthCheckStep()
	createInstanceGroup, _ := NewCreateInstanceGroupStep()
	createBackendService, _ := NewCreateBackendServiceStep()

	steps.RegisterStep(CreateInstanceStepName, createInstance)
	steps.RegisterStep(DeleteClusterStepName, deleteCluster)
	steps.RegisterStep(DeleteNodeStepName, deleteNode)
	steps.RegisterStep(CreateTargetPullStepName, createTargetPool)
	steps.RegisterStep(CreateIPAddressStepName, createIPAddress)
	steps.RegisterStep(CreateHealthCheckStepName, createHealthCheck)
	steps.RegisterStep(CreateInstanceGroupStepName, createInstanceGroup)
	steps.RegisterStep(CreateBackendServiceStepName, createBackendService)
}

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
