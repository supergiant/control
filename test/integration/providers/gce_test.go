package providers

import (
	"context"
	"fmt"
	"github.com/supergiant/supergiant/pkg/util"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGCE(t *testing.T) {
	projectId := os.Getenv("GCE_PROJECT_ID")
	region := os.Getenv("GCE_REGION")
	zone := os.Getenv("GCE_ZONE")
	pubKey := os.Getenv("GCE_PUB_KEY")
	privateKey := os.Getenv("GCE_PRIVATE_KEY")
	email := os.Getenv("GCE_EMAIL")
	tokenUri := os.Getenv("GCE_TOKEN_URI")

	k8sVersions := []string{"1.8.7", "1.7.7", "1.6.7", "1.5.7"}
	srv, err := newServer()

	if err != nil {
		t.Errorf("Unexpected error while creating client and server %v", err)
		return
	}

	requestor := createAdmin(srv.Core)
	client := srv.Core.APIClient("token", requestor.APIToken)

	go srv.Start()

	// TODO(stgleb): Need to add proper server shutdown to finish del requests
	defer srv.Stop()

	cloudAccount, err := createCloudAccount(client,
		map[string]string{
			"project_id":   projectId,
			"region":       region,
			"zone":         zone,
			"private_key":  privateKey,
			"client_email": email,
			"token_uri":    tokenUri},
		"gce")

	if err != nil {
		t.Errorf("Unexpected error while creating cloud account %v", err)
		return
	}

	for _, k8sVersion := range k8sVersions {
		t.Run(fmt.Sprintf("Run test GCE-%s", k8sVersion), func(t *testing.T) {

			kube, err := createKubeGCE(client, cloudAccount,
				fmt.Sprintf("test-%s", strings.ToLower(util.RandomString(5))),
				zone,
				pubKey,
				k8sVersion)

			if err != nil {
				t.Errorf("error creating the kube %v", err)
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
		})
	}
}
