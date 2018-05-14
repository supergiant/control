package providers

import (
	"fmt"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
	"strings"
	"testing"
	"time"
)

// Run test on AWS against all support versions of k8s
func TestAmazon(t *testing.T) {
	k8sVersions := []string{"1.8.7", "1.7.7", "1.6.7", "1.5.7"}
	srv, client, err := newClientServer()

	if err != nil {
		t.Error(err)
		return
	}

	session := &model.Session{
		User: &model.User{
			Username: "support",
			Password: "1234",
		},
	}
	err = client.Sessions.Create(session)

	for _, k8sVersion := range k8sVersions {
		t.Run(fmt.Sprintf("Test-AWS-%s", k8sVersion), func(t *testing.T) {
			t.Parallel()

			go srv.Start()
			defer srv.Stop()

			kube, err := createKube(client, "1.8.7")

			if err != nil {
				t.Error(err)
				return
			}

			list := &model.NodeList{}
			client.Nodes.List(list)

			err = util.WaitFor("Wait for cluster to start", time.Second*600, time.Second*1, func() (bool, error) {
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
