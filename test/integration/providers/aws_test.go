package providers

import (
	"fmt"
	"github.com/supergiant/supergiant/pkg/model"
	"testing"
)

var (
	configFile = "config.json"
)

func TestAmazon(t *testing.T) {
	k8sVersions := []string{"1.8.7", "1.7.7", "1.6.7", "1.5.7"}

	for _, k8sVersion := range k8sVersions {

		t.Run(fmt.Sprintf("Test-aws-%s", k8sVersion), func(t *testing.T) {
			t.Parallel()
			srv, err := newTestServer()

			if err != nil {
				t.Error(err)
				return
			}

			go srv.Start()
			defer srv.Stop()

			client := newClient(configFile)
			kube, err := createKube(client, "1.8.7")

			if err != nil {
				t.Error(err)
			}

			err = client.Kubes.Provision(kube.ID, kube)

			if err != nil {
				t.Error(err)
			}

			list := &model.NodeList{}
			client.Nodes.List(list)

			if len(list.Items) != 1 {
				t.Errorf("Wrong nodes count expected %d actual %d", 1, len(list.Items))
			}
		})
	}
}
