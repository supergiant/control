package etcd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/docker"
)

const StepName = "etcd"

type Step struct {
	script *template.Template
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func New(tpl *template.Template) *Step {
	return &Step{
		script: tpl,
	}
}

func (s *Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	config.EtcdConfig.Name = config.Node.ID
	config.EtcdConfig.AdvertiseHost = config.Node.PrivateIp
	ctx2, _ := context.WithTimeout(ctx, config.EtcdConfig.Timeout)

	vars := struct {
		ETCDConfig               steps.EtcdConfig
		InitialClusterIPs        string
		InitialAdvertisePeerURLs string
		AdvertiseURLs            string
		NodePrivateIP            string
	}{
		ETCDConfig:    config.EtcdConfig,
		NodePrivateIP: config.Node.PrivateIp,
	}

	masters := config.GetMasters()
	//one master configuration
	if len(masters) == 1 {
		vars.InitialClusterIPs = fmt.Sprintf("%s=http://%s:%s", config.Node.ID, config.Node.PrivateIp, config.EtcdConfig.ManagementPort)
		vars.InitialAdvertisePeerURLs = fmt.Sprintf("http://%s:%s", config.Node.PrivateIp, config.EtcdConfig.ManagementPort)
	} else {
		initialClusterIPs := make([]string, 0)
		advertisePeers := make([]string, 0)

		for _, master := range masters {
			initialClusterIPs = append(initialClusterIPs, fmt.Sprintf("%s=http://%s:%s", master.ID, master.PrivateIp, config.EtcdConfig.ManagementPort))
			advertisePeers = append(advertisePeers, fmt.Sprintf("http://%s:%s", master.PrivateIp, config.EtcdConfig.ManagementPort))
		}

		vars.InitialClusterIPs = strings.Join(initialClusterIPs, ",")
		vars.InitialAdvertisePeerURLs = strings.Join(advertisePeers, ",")
		vars.AdvertiseURLs = strings.Replace(vars.InitialAdvertisePeerURLs, ":2380", ":2379", -1)
	}

	err := steps.RunTemplate(ctx2, s.script, config.Runner, out, vars)

	if err != nil {
		return errors.Wrap(err, "install etcd step")
	}

	return nil
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return ""
}

func (s *Step) Depends() []string {
	return []string{docker.StepName}
}
