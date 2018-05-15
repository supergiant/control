package providers

import (
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/server"
	"github.com/phayes/freeport"
	"github.com/supergiant/supergiant/pkg/core"
	"fmt"
)

func newServer() (*server.Server, error) {
	port, err := freeport.GetFreePort()

	if err != nil {
		return nil, err
	}

	settings := core.Settings{
		LogLevel:        "debug",
		PublishHost:     "localhost",
		HTTPPort:        fmt.Sprintf("%d", port),
		SQLiteFile:      "file::memory:?cache=shared",
		SupportPassword: "1234",
	}

	c := &core.Core{}
	c.Settings = settings

	if err := c.InitializeForeground(); err != nil {
		panic(err)
	}

	srv, err := server.New(c)
	if err != nil {
		panic(err)
	}

	return srv, nil
}

func createKube(sg *client.Client, version string) (*model.Kube, error) {
	cloudAccount := &model.CloudAccount{
		Name:        "test",
		Provider:    "aws",
		Credentials: map[string]string{"support": "1234"},
	}

	err := sg.CloudAccounts.Create(cloudAccount)

	if err != nil {
		return nil, err
	}

	kube := &model.Kube{
		CloudAccountName:  cloudAccount.Name,
		Name:              "test",
		MasterNodeSize:    "m4.large",
		KubernetesVersion: version,
		NodeSizes:         []string{"m4.large"},
		AWSConfig: &model.AWSKubeConfig{
			Region:           "us-east-1",
			AvailabilityZone: "us-east-1a",
		},
	}

	err = sg.Kubes.Create(kube)

	return kube, err
}
