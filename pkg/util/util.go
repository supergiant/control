package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/node"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
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
func FillCloudAccountCredentials(ctx context.Context, cloudAccount *model.CloudAccount, config *steps.Config) error {
	config.ManifestConfig.ProviderString = string(cloudAccount.Provider)
	config.Provider = cloudAccount.Provider

	// Bind private key to config
	err := BindParams(cloudAccount.Credentials, &config.SshConfig)

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
	default:
		return sgerrors.ErrUnknownProvider
	}
}

func GetRandomNode(nodeMap map[string]*node.Node) *node.Node {
	for key := range nodeMap {
		return nodeMap[key]
	}

	return nil
}

func GetWriter(name string) (io.WriteCloser, error) {
	return os.OpenFile(path.Join("/tmp", name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func FindOutboundIP() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	publicIP := string(res)
	publicIP = strings.Replace(publicIP, "\n", "", -1)
	return publicIP, nil
}
