package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

var _ GroupsInterface = resources.GroupsClient{}

type Autorizerer interface {
	Authorizer() (autorest.Authorizer, error)
}

type GroupsInterface interface {
	CreateOrUpdate(ctx context.Context, name string, grp resources.Group) (resources.Group, error)
	Delete(ctx context.Context, name string) (resources.GroupsDeleteFuture, error)
}

// DEPRECATED
func BaseClientFor(authorizer Autorizerer, subscriptionID string) (resources.BaseClient, error) {
	token, err := authorizer.Authorizer()
	if err != nil {
		return resources.BaseClient{}, err
	}
	baseClient := resources.New(subscriptionID)
	baseClient.Authorizer = token
	return baseClient, nil
}

func GroupsClientFor(authorizer Autorizerer, subscriptionID string) (GroupsInterface, error) {
	token, err := authorizer.Authorizer()
	if err != nil {
		return resources.GroupsClient{}, err
	}
	baseClient := resources.New(subscriptionID)
	baseClient.Authorizer = token
	return resources.GroupsClient{BaseClient: baseClient}, nil
}

func toResourceGroupName(clusterID, clusterName string) string {
	return fmt.Sprintf("sg-%s-%s", clusterName, clusterID)
}

func toStrPtr(s string) *string {
	return &s
}
