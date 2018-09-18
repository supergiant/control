package provisioner

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/ssh"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/profile"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
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
func FillNodeCloudSpecificData(provider clouds.Name, nodeProfile profile.NodeProfile, config *steps.Config) error {
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

func nodesFromProfile(config *steps.Config, profile *profile.Profile) (map[string]*node.Node, map[string]*node.Node) {
	masters := make(map[string]*node.Node)
	nodes := make(map[string]*node.Node)

	for _, p := range profile.MasterProfiles {
		n := &node.Node{
			Id:       util.MakeNodeName("todo-", true),
			Provider: profile.Provider,
			Region:   profile.Region,
		}

		util.BindParams(p, n)
		masters[n.Id] = n
	}

	for _, p := range profile.NodesProfiles {
		n := &node.Node{
			Id:       util.MakeNodeName(util.MakeNodeName("todo-", true), false),
			Provider: profile.Provider,
			Region:   profile.Region,
		}

		util.BindParams(p, n)
		nodes[n.Id] = n
	}

	return masters, nodes
}

func generateKeyPair(size int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)

	if err != nil {
		return "", "", err
	}

	privateKeyPem := encodePrivateKeyToPEM(privateKey)
	publicKey, err := generatePublicKey(&privateKey.PublicKey)

	if err != nil {
		return "", "", err
	}

	return string(privateKeyPem), string(publicKey), nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func generatePublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(publicKey)

	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}
