package util

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

// RandomString generates random string with reservoir sampling algorithm https://en.wikipedia.org/wiki/Reservoir_sampling
func RandomString(n int) string {
	buffer := make([]byte, n)
	copy(buffer, letterBytes)

	for i := n; i < len(letterBytes); i++ {
		rndIndex := src.Int63() % int64(i)

		if rndIndex < int64(n) {
			buffer[rndIndex] = letterBytes[i]
		}
	}

	return string(buffer)
}

func MakeNodeName(clusterName string, nodeId string, isMaster bool) string {
	if isMaster {
		return fmt.Sprintf("%s-%s-%s", clusterName, "master", nodeId[:4])
	}

	return fmt.Sprintf("%s-%s-%s", clusterName, "node", nodeId[:4])
}

// bind params uses json serializing and reflect package that is underneath
// to avoid direct access to map for getting appropriate field values.
func BindParams(params map[string]string, object interface{}) error {
	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, object)

	if err != nil {
		return err
	}

	return nil
}

func MakeRole(isMaster bool) string {
	if isMaster {
		return "master"
	} else {
		return "node"
	}
}

func GetLogger(w io.Writer) (log *logrus.Logger) {
	log = logrus.New()
	log.Out = w
	log.SetLevel(logrus.StandardLogger().Level)
	return
}

func MakeFileName(taskID string) string {
	return fmt.Sprintf("%s.log", taskID)
}

func MakeKeyName(name string, isUser bool) string {
	if name == "" {
		name = RandomString(12)
	}

	if isUser {
		return fmt.Sprintf("%s-user", name)
	}

	return fmt.Sprintf("%s-provision", name)
}

// TODO(stgleb): move getting cloud account outside of this function
// Gets cloud account from storage and fills config object with those credentials
func FillCloudAccountCredentials(cloudAccount *model.CloudAccount, config *steps.Config) error {
	config.Provider = cloudAccount.Provider

	// Bind private key to config
	err := BindParams(cloudAccount.Credentials, &config.Kube.SSHConfig)

	if err != nil {
		return err
	}

	// TODO(stgleb):  Add support for other cloud providers
	switch cloudAccount.Provider {
	case clouds.AWS:
		return BindParams(cloudAccount.Credentials, &config.AWSConfig)
	case clouds.DigitalOcean:
		return BindParams(cloudAccount.Credentials, &config.DigitalOceanConfig)
	case clouds.GCE:
		return BindParams(cloudAccount.Credentials, &config.GCEConfig)
	case clouds.Azure:
		return BindParams(cloudAccount.Credentials, &config.AzureConfig)
	default:
		return sgerrors.ErrUnknownProvider
	}
}

func GetRandomNode(nodeMap map[string]*model.Machine) *model.Machine {
	for key := range nodeMap {
		return nodeMap[key]
	}

	return nil
}

func GetWriter(name string) (io.WriteCloser, error) {
	return os.OpenFile(path.Join("/tmp", name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func LoadCloudSpecificDataFromKube(k *model.Kube, config *steps.Config) error {
	if k == nil {
		return sgerrors.ErrNilEntity
	}
	config.Kube = *k

	// TODO: Is it ok?
	if k.CloudSpec == nil {
		return nil
	}

	switch config.Provider {
	case clouds.AWS:
		// Load AZ -> subnet mapping for cluster
		config.AWSConfig.Subnets = k.Subnets
		config.AWSConfig.Region = k.Region
		config.AWSConfig.AvailabilityZone = k.CloudSpec[clouds.AwsAZ]
		config.AWSConfig.VPCCIDR = k.CloudSpec[clouds.AwsVpcCIDR]
		config.AWSConfig.VPCID = k.CloudSpec[clouds.AwsVpcID]
		config.AWSConfig.KeyPairName = k.CloudSpec[clouds.AwsKeyPairName]
		config.AWSConfig.MastersSecurityGroupID = k.CloudSpec[clouds.AwsMastersSecGroupID]
		config.AWSConfig.NodesSecurityGroupID = k.CloudSpec[clouds.AwsNodesSecgroupID]
		config.AWSConfig.RouteTableID = k.CloudSpec[clouds.AwsRouteTableID]
		config.AWSConfig.InternetGatewayID = k.CloudSpec[clouds.AwsInternetGateWayID]
		config.AWSConfig.MastersInstanceProfile = k.CloudSpec[clouds.AwsMasterInstanceProfile]
		config.AWSConfig.NodesInstanceProfile = k.CloudSpec[clouds.AwsNodeInstanceProfile]
		config.AWSConfig.ImageID = k.CloudSpec[clouds.AwsImageID]
		config.Kube.SSHConfig.BootstrapPrivateKey = k.CloudSpec[clouds.AwsSshBootstrapPrivateKey]
		config.Kube.SSHConfig.PublicKey = k.CloudSpec[clouds.AwsUserProvidedSshPublicKey]
		config.AWSConfig.ExternalLoadBalancerName = k.ExternalLoadBalancerName
		config.AWSConfig.InternalLoadBalancerName = k.InternalLoadBalancerName
	case clouds.GCE:
		config.GCEConfig.Region = k.Region

	case clouds.DigitalOcean:

	case clouds.Azure:
		config.AzureConfig.Location = k.Region

	default:
		return errors.Wrapf(sgerrors.ErrUnsupportedProvider, "Load cloud specific data from kube %s", k.ID)
	}

	return nil
}
