package provisioner

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows"
	"github.com/supergiant/control/pkg/workflows/steps"
)

type RateLimiter struct {
	bucket *time.Ticker
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		bucket: time.NewTicker(interval),
	}
}

// Take either returns giving calling code ability to execute or blocks until
// bucket is full again
func (r *RateLimiter) Take() {
	<-r.bucket.C
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
	case clouds.Azure:
		return util.BindParams(nodeProfile, &config.AzureConfig)
	default:
		return sgerrors.ErrUnknownProvider
	}

	return nil
}

func MergeConfig(source *steps.Config, destination *steps.Config) error {
	switch source.Provider {
	case clouds.AWS:
		data, err := json.Marshal(&source.AWSConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.AWSConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	case clouds.GCE:
		data, err := json.Marshal(&source.GCEConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.GCEConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	case clouds.DigitalOcean:
		data, err := json.Marshal(&source.DigitalOceanConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.DigitalOceanConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	case clouds.Packet:
		return nil
	case clouds.OpenStack:
		return nil
	case clouds.Azure:
		data, err := json.Marshal(&source.AzureConfig)

		if err != nil {
			return errors.Wrapf(err, "merge config marshall config1")
		}

		err = json.Unmarshal(data, &destination.AzureConfig)

		if err != nil {
			return errors.Wrapf(err, "Merge config")
		}
	default:
		return sgerrors.ErrUnknownProvider
	}

	// These items must be shared among configs
	destination.ExternalDNSName = source.ExternalDNSName
	destination.InternalDNSName = source.InternalDNSName
	destination.SetKubeStateChan(source.KubeStateChan())
	destination.SetNodeChan(source.NodeChan())
	destination.SetConfigChan(source.ConfigChan())
	destination.Masters = source.Masters
	destination.Nodes = source.Nodes
	logrus.Info("Copy masters %v", destination.Masters)
	logrus.Info("Copy nodes %v", destination.Nodes)

	return nil
}

func nodesFromProfile(clusterName string, masterTasks, nodeTasks []*workflows.Task, profile *profile.Profile) (map[string]*model.Machine, map[string]*model.Machine) {
	masters := make(map[string]*model.Machine)
	nodes := make(map[string]*model.Machine)

	for index, p := range profile.MasterProfiles {
		taskId := masterTasks[index].ID
		name := util.MakeNodeName(clusterName, taskId, true)

		// TODO(stgleb): check if we can lowercase node names for all nodes
		if profile.Provider == clouds.GCE {
			name = strings.ToLower(name)
		}
		n := &model.Machine{
			TaskID:   taskId,
			Name:     name,
			Provider: profile.Provider,
			Region:   profile.Region,
			State:    model.MachineStatePlanned,
		}

		util.BindParams(p, n)
		masters[n.Name] = n
	}

	for index, p := range profile.NodesProfiles {
		taskId := nodeTasks[index].ID
		name := util.MakeNodeName(clusterName, taskId[:4], false)

		// TODO(stgleb): check if we can lowercase node names for all nodes
		if profile.Provider == clouds.GCE {
			name = strings.ToLower(name)
		}
		n := &model.Machine{
			TaskID:   taskId,
			Name:     name,
			Provider: profile.Provider,
			Region:   profile.Region,
			State:    model.MachineStatePlanned,
		}

		util.BindParams(p, n)
		nodes[n.Name] = n
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

func grabTaskIds(taskMap map[string][]*workflows.Task) map[string][]string {
	taskIds := make(map[string][]string, 0)

	for taskSet := range taskMap {
		tasks := make([]string, 0, len(taskMap[taskSet]))

		for _, task := range taskMap[taskSet] {
			tasks = append(tasks, task.ID)
		}

		taskIds[taskSet] = tasks
	}

	return taskIds
}
