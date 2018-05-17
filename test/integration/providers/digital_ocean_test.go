package providers

import (
	"testing"
	"os"
	"fmt"
	"strings"
	"github.com/supergiant/supergiant/pkg/util"
	"time"
	"context"
)

func TestDigitalOcean(t *testing.T) {
	region := os.Getenv("DO_REGION")
	fingerPrint := os.Getenv("DO_KEY_FINGER_PRINT")
	token := os.Getenv("DO_TOKEN")

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

	for _, k8sVersion := range k8sVersions {
		t.Run(fmt.Sprintf("Test Digital Ocean k8s version %s", k8sVersion), func(t *testing.T){
			kube, err := createKubeDO(client,
				fmt.Sprintf("test-%s", strings.ToLower(util.RandomString(5))),
				region,
				fingerPrint,
				token,
				k8sVersion)

			if err != nil {
				t.Error(err)
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
			client.Kubes.Delete(kube.ID, kube)
		})
	}
}