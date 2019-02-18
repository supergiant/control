package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
)

type Autorizerer interface {
	Authorizer() (autorest.Authorizer, error)
}

func BaseClientFor(authorizer Autorizerer, subscriptionID string) (resources.BaseClient, error) {
	token, err := authorizer.Authorizer()
	if err != nil {
		return resources.BaseClient{}, err
	}
	baseClient := resources.New(subscriptionID)
	baseClient.Authorizer = token
	return baseClient, nil
}

func toResourceGroupName(clusterID, clusterName string) string {
	return fmt.Sprintf("sg-%s-%s", clusterName, clusterID)
}

func toStrPtr(s string) *string {
	return &s
}
