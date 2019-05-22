package configmap

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	capacityapi "github.com/supergiant/capacity/pkg/api"
	"k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/kubeconfig"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	// userdataPrefix is used to set script interpreter and log the output.
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts
	userdataPrefix = `#!/bin/bash
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1
`
)

// TODO: take it from the capacity repo
const (
	CapacityConfigMapName      = "capacity"
	CapacityConfigMapNamespace = "kube-system"
	CapacityConfigMapKey       = "kubescaler.conf"

	KeyID          = "awsKeyID"
	SecretKey      = "awsSecretKey"
	Region         = "awsRegion"
	KeyName        = "awsKeyName"
	ImageID        = "awsImageID"
	IAMRole        = "awsIAMRole"
	SecurityGroups = "awsSecurityGroups"
	SubnetID       = "awsSubnetID"
	VolType        = "awsVolType"
	VolSize        = "awsVolSize"
	Tags           = "awsTags"

	DefaultVolSize = "80"
	DefaultVolType = "gp2"

	KubeAddr = "KubeAddr"
)

const StepName = "configmap"

type Step struct {
	timeout      time.Duration
	attemptCount int
}

func New() *Step {
	t := &Step{
		timeout:      time.Minute * 1,
		attemptCount: 5,
	}

	return t
}

func Init() {
	steps.RegisterStep(StepName, New())
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	k8sClient, err := buildKubeClient(config)
	if err != nil {
		return errors.Wrap(err, "build kubernetes client")
	}

	userdata, err := parse(config.ConfigMap.Data, map[string]string{
		KubeAddr: config.ExternalDNSName,
	})
	if err != nil {
		return errors.Wrap(err, "parse userdata template")
	}

	capcfg := capacityapi.Config{
		ClusterName:  config.ClusterName,
		Userdata:     base64.StdEncoding.EncodeToString(userdata),
		ProviderName: "aws",
		Provider: map[string]string{
			KeyID:          config.AWSConfig.KeyID,
			SecretKey:      config.AWSConfig.Secret,
			Region:         config.AWSConfig.Region,
			IAMRole:        config.AWSConfig.NodesInstanceProfile,
			ImageID:        config.AWSConfig.ImageID,
			KeyName:        config.AWSConfig.KeyPairName,
			SecurityGroups: config.AWSConfig.NodesSecurityGroupID,
			VolSize:        DefaultVolSize,
			VolType:        DefaultVolType,
			Tags:           fmt.Sprintf("%s=%s", clouds.TagClusterID, config.ClusterID),
			SubnetID:       getSubnet(config.AWSConfig.Subnets),
			// TODO: add AZ option
		},
		WorkersCountMin: 2,
		WorkersCountMax: 10,
		Paused:          aws.Bool(true),
	}

	capcfgraw, err := json.Marshal(capcfg)
	if err != nil {
		return err
	}

	_, err = k8sClient.ConfigMaps(CapacityConfigMapNamespace).Create(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: CapacityConfigMapName,
		},
		Data: map[string]string{
			CapacityConfigMapKey: string(capcfgraw),
		},
	})

	if !k8serrors.IsAlreadyExists(err) && err != nil {
		return errors.Wrapf(err, "create %s/%s config map", CapacityConfigMapNamespace, CapacityConfigMapName)
	}

	return nil
}

func (s *Step) Rollback(ctx context.Context, out io.Writer, config *steps.Config) error {
	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "create configmap for capacity service"
}

func (s *Step) Depends() []string {
	return nil
}

func buildKubeClient(config *steps.Config) (clientcorev1.CoreV1Interface, error) {
	config.Kube.Auth.AdminCert = config.CertificatesConfig.AdminCert
	config.Kube.Auth.AdminKey = config.CertificatesConfig.AdminKey
	config.Kube.Auth.CACert = config.CertificatesConfig.CACert
	config.Kube.ExternalDNSName = config.ExternalDNSName

	return kubeconfig.CoreV1Client(&config.Kube)
}

func parse(in string, data interface{}) ([]byte, error) {
	tpl, err := template.New("userdata").Parse(in)
	if err != nil {
		return nil, err
	}
	w := &bytes.Buffer{}
	if err = tpl.Execute(w, data); err != nil {
		return nil, err
	}

	out := bytes.TrimSpace(w.Bytes())
	if !bytes.HasPrefix(out, []byte("#!")) {
		// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts
		out = append([]byte(userdataPrefix), out...)
	}
	return out, nil
}

func getSubnet(azSubnets map[string]string) string {
	for _, subnet := range azSubnets {
		return subnet
	}
	return ""
}
