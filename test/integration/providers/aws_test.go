// +build integration

package providers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/supergiant/supergiant/pkg/util"
)

// Run test on AWS against all support versions of k8s
func TestAmazon(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration tests for short mode")
	}

	// Run test in parallel for different cloud providers
	t.Parallel()

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsAZ := os.Getenv("AWS_AVAILABILITY_ZONE")
	pubKey := os.Getenv("AWS_SSH_PUB_KEY")

	k8sVersions := []string{"1.8.7", "1.7.7", "1.6.7", "1.5.7"}
	srv, err := newServer()

	if err != nil {
		t.Errorf("Unexpected error while creating client and server %v", err)
		return
	}

	requestor := createAdmin(srv.Core)
	client := srv.Core.APIClient("token", requestor.APIToken)

	go srv.Start()
	defer srv.Stop()

	cloudAccount, err := createCloudAccount(client,
		map[string]string{
			"access_key": awsAccessKey,
			"secret_key": awsSecretKey},
		"aws")

	if err != nil {
		t.Errorf("Unexpected error while creating cloud account %v", err)
		return
	}

	for _, k8sVersion := range k8sVersions {
		t.Run(fmt.Sprintf("Test-AWS-%s", k8sVersion), func(t *testing.T) {
			kube, err := createKubeAWS(client, cloudAccount,
				fmt.Sprintf("test-%s", strings.ToLower(util.RandomString(5))),
				awsRegion,
				awsAZ,
				pubKey,
				k8sVersion)

			if err != nil {
				t.Errorf("Unexpected error while creating a kube %v", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
			defer cancel()

			err = util.WaitFor(ctx, "Wait for cluster to start", time.Second*30, func() (bool, error) {
				err := client.Kubes.Get(kube.ID, kube)

				if err != nil {
					return false, err
				}

				return kube.Ready, nil
			})

			// Clean up afterwards
			err = client.Kubes.Delete(kube.ID, kube)

			if err != nil {
				t.Errorf("Error while deleting the kube %d", kube.ID)
			}

			time.Sleep(time.Minute)
		})
	}
}
