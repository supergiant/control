package provisioner

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"io/ioutil"
	"net/http"
)

type EtcdTokenGetter struct {
	discoveryUrl string
}

func NewEtcdTokenGetter() *EtcdTokenGetter {
	return &EtcdTokenGetter{
		discoveryUrl: "https://discovery.etcd.io/new?size=%d",
	}
}

func (e *EtcdTokenGetter) GetToken(ctx context.Context, num int) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(e.discoveryUrl, num), nil)
	req = req.WithContext(ctx)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Fill cloud account specific data gets data from the map and puts to particular cloud provider config
func FillNodeCloudSpecificData(provider clouds.Name, nodeProfile map[string]string, config *steps.Config) error {
	switch provider {
	case clouds.AWS:
		return util.BindParams(nodeProfile, &config.AWSConfig)
	case clouds.GCE:
		return util.BindParams(nodeProfile, &config.GCEConfig)
	case clouds.DigitalOcean:
		return util.BindParams(nodeProfile, &config.DigitalOceanConfig)
	case clouds.Packet:
		return util.BindParams(nodeProfile, &config.PacketConfig)
	case clouds.OpenStack:
		return util.BindParams(nodeProfile, &config.OSConfig)
	default:
		return sgerrors.ErrUnknownProvider
	}

	return nil
}

func nodesFromProfile(profile *profile.Profile) ([]*node.Node, []*node.Node) {
	masters := make([]*node.Node, 0, len(profile.MasterProfiles))
	nodes := make([]*node.Node, 0, len(profile.NodesProfiles))

	for _, p := range profile.MasterProfiles {
		n := &node.Node{
			Provider: profile.Provider,
			Region:   profile.Region,
		}

		util.BindParams(p, n)
		masters = append(masters, n)
	}

	for _, p := range profile.NodesProfiles {
		n := &node.Node{
			Provider: profile.Provider,
			Region:   profile.Region,
		}

		util.BindParams(p, n)
		nodes = append(nodes, n)
	}

	return masters, nodes
}
