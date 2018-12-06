package clusterservices

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	tm "github.com/supergiant/control/pkg/templatemanager"
	"github.com/supergiant/control/pkg/workflows/steps"
	"github.com/supergiant/control/pkg/workflows/steps/dryrun"
)

const StepName = "clusterservices"

type Step struct {
	script *template.Template
}

func (s *Step) Rollback(context.Context, io.Writer, *steps.Config) error {
	return nil
}

func Init() {
	tpl, err := tm.GetTemplate(StepName)

	if err != nil {
		panic(fmt.Sprintf("template %s not found", StepName))
	}

	steps.RegisterStep(StepName, New(tpl))
}

func New(script *template.Template) *Step {
	t := &Step{
		script: script,
	}

	return t
}

func (s Step) Run(ctx context.Context, out io.Writer, config *steps.Config) error {
	capacityCfg, err := buildCapacityConfig(config)
	if err != nil {
		return errors.Wrap(err, "install capacity")
	}

	// TODO: install it with the helm client?
	if err := steps.RunTemplate(ctx, s.script, config.Runner, out, capacityCfg, config.DryRun); err != nil {
		return errors.Wrap(err, "install capacity")
	}

	return nil
}

func (s *Step) Name() string {
	return StepName
}

func (s *Step) Description() string {
	return "Step for deploying supergiant services"
}

func (s *Step) Depends() []string {
	return nil
}

type CapacityConfig struct {
	Version          string
	KubescalerConfig string
	Userdata         string
}

func buildCapacityConfig(cfg *steps.Config) (*CapacityConfig, error) {
	kubescalerCfg, err := kubescalerConfigFrom(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "setup kubescaler config")
	}
	userdata, err := dryrun.NodeScript(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "get node setup script")
	}
	
	return &CapacityConfig{
		Version:          "0.1.0",
		KubescalerConfig: kubescalerCfg,
		Userdata:         userdata,
	}, nil
}

type Config struct {
	ClusterName     string            `json:"clusterName"`
	ProviderName    string            `json:"providerName"`
	Provider        map[string]string `json:"provider"`
	Paused          *bool             `json:"paused,omitempty"`
	ScanInterval    string            `json:"scanInterval"`
	WorkersCountMin int               `json:"workersCountMin"`
	WorkersCountMax int               `json:"workersCountMax"`
	MachineTypes    []string          `json:"machineTypes"`
}

// TODO: retrieve all parameters from steps.Config
func kubescalerConfigFrom(cfg *steps.Config) (string, error) {
	// TODO: capacity doesn't recognize availability zones
	var subnetID string
	for _, id := range cfg.AWSConfig.Subnets {
		if id != "" {
			subnetID = id
			break
		}
	}

	raw, err := json.Marshal(Config{
		ClusterName:  cfg.ClusterName,
		ProviderName: "aws",
		Provider: map[string]string{
			"awsIAMRole":        cfg.AWSConfig.NodesInstanceProfile,
			"awsImageID":        cfg.AWSConfig.ImageID,
			"awsKeyID":          cfg.AWSConfig.KeyID,
			"awsKeyName":        cfg.AWSConfig.KeyPairName,
			"awsRegion":         cfg.AWSConfig.Region,
			"awsSecretKey":      cfg.AWSConfig.Secret,
			"awsSecurityGroups": cfg.AWSConfig.NodesSecurityGroupID,
			"awsSubnetID":       subnetID,
			"awsVolSize":        cfg.AWSConfig.VolumeSize,
			"awsVolType":        "gp2",
		},
		Paused:          aws.Bool(true),
		WorkersCountMin: 1,
		WorkersCountMax: 1,
	})
	return string(raw), err
}
