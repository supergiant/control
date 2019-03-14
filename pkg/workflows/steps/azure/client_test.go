package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
)

type fakeGroupsClient struct {
	GroupsInterface
	createRes resources.Group
	createErr error
	deleteRes resources.GroupsDeleteFuture
	deleteErr error
}

func (f fakeGroupsClient) CreateOrUpdate(ctx context.Context, name string, grp resources.Group) (resources.Group, error) {
	return f.createRes, f.createErr
}
func (f fakeGroupsClient) Delete(ctx context.Context, name string) (resources.GroupsDeleteFuture, error) {
	return f.deleteRes, f.deleteErr
}
