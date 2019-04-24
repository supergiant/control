package configmap

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/control/pkg/model"
	"io"
	"k8s.io/api/core/v1"
	"text/template"
	"time"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/supergiant/control/pkg/kubeconfig"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepName = "configmap"

type Step struct {
	timeout time.Duration
	attemptCount int
	script *template.Template
}

func New(script *template.Template) *Step {
	t := &Step{
		timeout: time.Minute * 1,
		attemptCount: 5,
		script: script,
	}

	return t
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	var err error

	config.Kube.Auth.AdminCert = config.CertificatesConfig.AdminCert
	config.Kube.Auth.AdminKey = config.CertificatesConfig.AdminKey
	config.Kube.Auth.CACert = config.CertificatesConfig.CACert
	config.Kube.ExternalDNSName = config.ExternalDNSName

	master := config.GetMaster()
	config.Kube.Masters = map[string]*model.Machine{
		master.Name: master,
	}

	cfg, err := kubeconfig.NewConfigFor(&config.Kube)

	if err != nil {
		return errors.Wrap(err, "configmap create kubeconfig from kube")
	}

	k8sClient, err := clientcorev1.NewForConfig(cfg)

	if err != nil {
		return errors.Wrapf(err, "configmap build kubernetes client")
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: config.ConfigMap.Namespace,
			// TODO(stgleb): extract to constant
			Name: "capacity",
		},
		Data: map[string]string{
			"script": config.ConfigMap.Data,
		},
	}

	timeout := s.timeout

	for i := 0; i < s.attemptCount; i++ {
		_, err = k8sClient.ConfigMaps(config.ConfigMap.Namespace).Create(configMap)

		if err == nil {
			break
		}

		logrus.Debugf("create config map error %v", err)
		time.Sleep(s.timeout)
		timeout *= 2
	}

	if err != nil {
		return errors.Wrapf(err, "create config map")
	}

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
