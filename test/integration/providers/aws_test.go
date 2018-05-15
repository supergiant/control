package providers

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
	"context"
	"github.com/supergiant/supergiant/pkg/core"
)

func createAdmin(c *core.Core) *model.User {
	admin := &model.User{
		Username: "bossman",
		Password: "password",
		Role:     "admin",
	}
	c.Users.Create(admin)
	return admin
}

// Run test on AWS against all support versions of k8s
func TestAmazon(t *testing.T) {
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
		t.Run(fmt.Sprintf("Test-AWS-%s", k8sVersion), func(t *testing.T) {
			kube, err := createKube(client, k8sVersion)
			if err != nil {
				t.Error(err)
				return
			}

			list := &model.NodeList{}
			client.Nodes.List(list)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
			defer cancel()

			err = util.WaitFor(ctx, "Wait for cluster to start", time.Second*1, func() (bool, error) {
				err := client.Kubes.Get(kube.ID, kube)

				if err != nil {
					return false, err
				}

				// TODO(stgleb): Create string constants list for cluster/node/service statuses
				if strings.Contains(kube.Status.Description, "Run") {
					return true, nil
				}

				return false, nil
			})

			if err != nil {
				t.Error(err)
			}

			if len(list.Items) != 1 {
				t.Errorf("Wrong nodes count expected %d actual %d", 1, len(list.Items))
			}
		})
	}
}
