package digitalocean

import (
	"context"
	"time"

	"github.com/digitalocean/godo"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateMachineStepName      = "createMachineDigitalOcean"
	CreateLoadBalancerStepName = "createLoadBalancerDigitalOcean"
	RegisterInstanceToLB       = "registerInstanceToLoadBalancerDigitalOcean"

	DeleteMachineStepName      = "deleteMachineDigitalOcean"
	DeleteClusterMachines      = "deleteClusterMachineDigitalOcean"
	DeleteDeleteKeysStepName   = "deleteKeysDigitalOcean"
	DeleteLoadBalancerStepName = "deleteLoadBalancerDigitalOcean"
)

type DropletService interface {
	Get(context.Context, int) (*godo.Droplet, *godo.Response, error)
	Create(context.Context, *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error)
}

type TagService interface {
	TagResources(string, *godo.TagResourcesRequest) (*godo.Response, error)
}

type KeyService interface {
	Create(context.Context, *godo.KeyCreateRequest) (*godo.Key, *godo.Response, error)
}

type DeleteService interface {
	DeleteByTag(context.Context, string) (*godo.Response, error)
}

type LoadBalancerService interface {
	Create(context.Context, *godo.LoadBalancerRequest) (*godo.LoadBalancer, *godo.Response, error)
	Delete(context.Context, string) (*godo.Response, error)
	Get(context.Context, string) (*godo.LoadBalancer, *godo.Response, error)
}

func Init() {
	steps.RegisterStep(CreateMachineStepName, NewCreateInstanceStep(time.Minute*5, time.Second*5))
	steps.RegisterStep(DeleteMachineStepName, NewDeleteMachineStep(time.Minute*1))
	steps.RegisterStep(DeleteClusterMachines, NewDeletemachinesStep(time.Minute*1))
	steps.RegisterStep(DeleteDeleteKeysStepName, NewDeleteKeysStep())

	steps.RegisterStep(CreateLoadBalancerStepName, NewCreateLoadBalancerStep())
	steps.RegisterStep(DeleteLoadBalancerStepName, NewDeleteLoadBalancerStep())
}
