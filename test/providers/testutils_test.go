// +build integration

package providers

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/aws"
	"github.com/supergiant/supergiant/pkg/provider/digitalocean"
	"github.com/supergiant/supergiant/pkg/provider/gce"
	"github.com/supergiant/supergiant/pkg/provider/kubernetes"
	"github.com/supergiant/supergiant/pkg/provider/openstack"
	"github.com/supergiant/supergiant/pkg/provider/packet"
	"github.com/supergiant/supergiant/pkg/server"
)

func getPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func newServer() (*server.Server, error) {
	port, err := getPort()

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
	c.Log = logrus.New()

	c.Settings = settings
	// See relevant NOTE in core.go
	c.AWSProvider = func(creds map[string]string) core.Provider {
		return &aws.Provider{
			Core: c,
			EC2:  aws.EC2,
			IAM:  aws.IAM,
			ELB:  aws.ELB,
			S3:   aws.S3,
			EFS:  aws.EFS,
		}
	}
	c.DOProvider = func(creds map[string]string) core.Provider {
		return &digitalocean.Provider{
			Core:   c,
			Client: digitalocean.Client,
		}
	}
	c.OSProvider = func(creds map[string]string) core.Provider {
		return &openstack.Provider{
			Core:   c,
			Client: openstack.Client,
		}
	}
	c.GCEProvider = func(creds map[string]string) core.Provider {
		return &gce.Provider{
			Core:   c,
			Client: gce.Client,
		}
	}
	c.PACKProvider = func(creds map[string]string) core.Provider {
		return &packet.Provider{
			Core:   c,
			Client: packet.Client,
		}
	}

	c.K8SProvider = &kubernetes.Provider{Core: c}

	if err := c.InitializeForeground(); err != nil {
		return nil, err
	}

	srv, err := server.New(c)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func createCloudAccount(client *client.Client, credentials map[string]string, provider string) (*model.CloudAccount, error) {
	cloudAccount := &model.CloudAccount{
		Name:        "test",
		Provider:    provider,
		Credentials: credentials,
	}

	err := client.CloudAccounts.Create(cloudAccount)

	if err != nil {
		return nil, err
	}

	return cloudAccount, nil
}

func createAdmin(c *core.Core) *model.User {
	admin := &model.User{
		Username: "bossman",
		Password: "password",
		Role:     "admin",
	}
	c.Users.Create(admin)
	return admin
}

func createKubeAWS(sg *client.Client, cloudAccount *model.CloudAccount, kubeName, awsRegion, awsAZ, pubKey, version string) (*model.Kube, error) {
	kube := &model.Kube{
		CloudAccountName:  cloudAccount.Name,
		CloudAccount:      cloudAccount,
		Name:              kubeName,
		MasterNodeSize:    "m4.large",
		KubernetesVersion: version,
		SSHPubKey:         pubKey,
		NodeSizes:         []string{"m4.large"},
		AWSConfig: &model.AWSKubeConfig{
			Region:           awsRegion,
			AvailabilityZone: awsAZ,
		},
	}

	err := sg.Kubes.Create(kube)

	return kube, err
}

func createKubeDO(sg *client.Client, cloudAccount *model.CloudAccount, kubeName, region, keyFingerPrint, version string) (*model.Kube, error) {
	kube := &model.Kube{
		CloudAccountName:  cloudAccount.Name,
		CloudAccount:      cloudAccount,
		Name:              kubeName,
		MasterNodeSize:    "1gb",
		KubernetesVersion: version,
		NodeSizes:         []string{"2gb"},
		DigitalOceanConfig: &model.DOKubeConfig{
			Region:            region,
			SSHKeyFingerprint: []string{keyFingerPrint},
		},
	}

	err := sg.Kubes.Create(kube)

	return kube, err
}

func createKubeGCE(sg *client.Client, cloudAccount *model.CloudAccount, kubeName, zone, pubKey, version string) (*model.Kube, error) {
	kube := &model.Kube{
		CloudAccountName:  cloudAccount.Name,
		CloudAccount:      cloudAccount,
		Name:              kubeName,
		MasterNodeSize:    "n1-standard-1",
		KubernetesVersion: version,
		SSHPubKey:         pubKey,

		NodeSizes: []string{"n1-standard-1"},
		GCEConfig: &model.GCEKubeConfig{
			SSHPubKey:         pubKey,
			Zone:              zone,
			KubernetesVersion: version,
		},
	}

	err := sg.Kubes.Create(kube)

	return kube, err
}
