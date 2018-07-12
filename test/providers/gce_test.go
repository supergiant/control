// +build integration

package providers

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/supergiant/supergiant/pkg/util"
)

func TestGCE(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration tests for short mode")
	}

	// Run test in parallel for different cloud providers
	t.Parallel()

	// The following environment variables need to be set on the machine used for testing:
	// NOTE: in order for Golang to access these in a terminal session, they need to be exported (e.g. export GCE_ENV="data")

	// The GCE Service Account's Project ID (e.g. "saber-hub-123456"):
	projectId := os.Getenv("GCE_PROJECT_ID")

	// A GCE Region to create VMs in (e.g. "us-east1"):
	region := os.Getenv("GCE_REGION")

	// A GCE Zone to create VMs in (e.g. "us-east1-b"):
	zone := os.Getenv("GCE_ZONE")

	// A Public Key to contact the VMs with (e.g. "sh-rsa AAAAB3NzaC1..."):
	pubKey := os.Getenv("GCE_PUB_KEY")

	// The GCE Service Account's Private Key. (e.g. "-----BEGIN PRIVATE KEY-----\n... ...\n-----END PRIVATE KEY-----\n"):
	rawPrivateKey := os.Getenv("GCE_PRIVATE_KEY")
	reg := regexp.MustCompile(`\\n`)
	privateKey := reg.ReplaceAllString(rawPrivateKey, "\n")

	// The GCE Service Account's email (e.g. "carlos@saber-hub-123456.iam.gserviceaccount.com"):
	email := os.Getenv("GCE_EMAIL")

	// The GCE Service Account's Token URI (e.g. "https://accounts.google.com/o/oauth2/token"):
	tokenUri := os.Getenv("GCE_TOKEN_URI")

	// The GCE Service Account's Auth URI (e.g. "https://accounts.google.com/o/oauth2/auth"):
	authUri := os.Getenv("GCE_AUTH_URI")

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
			"project_id":   projectId,
			"region":       region,
			"zone":         zone,
			"private_key":  privateKey,
			"client_email": email,
			"token_uri":    tokenUri,
			"auth_uri":     authUri},
		"gce")

	if err != nil {
		t.Errorf("Unexpected error while creating cloud account %v", err)
		return
	}

	for _, k8sVersion := range k8sVersions {
		t.Run(fmt.Sprintf("Run test GCE-%s", k8sVersion), func(t *testing.T) {

			kube, err := createKubeGCE(client,
				cloudAccount,
				fmt.Sprintf("test-%s", strings.ToLower(util.RandomString(5))),
				zone,
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
