package digitalocean

import (
	"context"
	"errors"
	"time"

	"github.com/digitalocean/godo"

	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	CreateMachineStepName = "createMachineDigitalOcean"
	DeleteMachineStepName = "deleteMachineDigitalOcean"
	DeleteClusterMachines = "deleteClusterMachineDigitalOcean"
	DeleteDeleteKeysStep  = "deleteKeys"
)

var (
	// TODO(stgleb): We need global error for timeout exceeding
	ErrTimeoutExceeded = errors.New("timeout exceeded")
)

type DropletService interface {
	Get(int) (*godo.Droplet, *godo.Response, error)
	Create(*godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error)
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

func Init() {
	steps.RegisterStep(CreateMachineStepName, NewCreateInstanceStep(time.Minute*5, time.Second*5))
	steps.RegisterStep(DeleteMachineStepName, NewDeleteMachineStep(time.Minute*1))
	steps.RegisterStep(DeleteClusterMachines, NewDeletemachinesStep(time.Minute*1))
}
