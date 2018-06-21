package aws

import (
	"context"
)

type Interface interface {
	CreateInstance(context.Context, InstanceConfig) error
	DeleteInstance(ctx context.Context, region, instanceID string) error
}
