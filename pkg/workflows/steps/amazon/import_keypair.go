package amazon

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepImportKeyPair = "aws_import_keypair_step"

type keyImporter interface {
	ImportKeyPairWithContext(aws.Context, *ec2.ImportKeyPairInput, ...request.Option) (*ec2.ImportKeyPairOutput, error)
	WaitUntilKeyPairExists(*ec2.DescribeKeyPairsInput) error
}

// KeyPairStep represents creation of keypair in aws
// since there is hard cap on keypairs per account supergiant will create one per cluster
type KeyPairStep struct {
	GetEC2 GetEC2Fn
	getSvc func(steps.AWSConfig) (keyImporter, error)
}

//InitImportKeyPair add the step to the registry
func InitImportKeyPair(fn GetEC2Fn) {
	steps.RegisterStep(StepImportKeyPair, NewImportKeyPairStep(fn))
}

func NewImportKeyPairStep(fn GetEC2Fn) *KeyPairStep {
	return &KeyPairStep{
		getSvc: func(cfg steps.AWSConfig) (keyImporter, error) {
			EC2, err := fn(cfg)

			if err != nil {
				return nil, ErrAuthorization
			}

			return EC2, nil
		},
	}
}

//Verifies that a key exists,
func (s *KeyPairStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	svc, err := s.getSvc(cfg.AWSConfig)

	if err != nil {
		logrus.Errorf("Getting service caused %v", err)
		return errors.Wrapf(err, "%s caused error when getting service",
			StepImportKeyPair)
	}

	if len(cfg.ClusterID) < 4 {
		return errors.New("Cluster ID is too short")
	}

	// NOTE(stgleb): Add unique part to key pair name that allows to
	// create cluster with the same name and avoid name collision of key pairs.
	bootstrapKeyPairName := util.MakeKeyName(fmt.Sprintf("%s-%s",
		cfg.ClusterName,
		cfg.ClusterID[:4]),
		false)
	log.Infof("[%s] - importing cluster bootstrap key as keypair %s",
		s.Name(), bootstrapKeyPairName)
	req := &ec2.ImportKeyPairInput{
		KeyName:           &bootstrapKeyPairName,
		PublicKeyMaterial: []byte(cfg.Kube.SSHConfig.BootstrapPublicKey),
	}

	output, err := svc.ImportKeyPairWithContext(ctx, req)
	if err != nil {
		return errors.Wrap(ErrImportKeyPair, err.Error())
	}

	cfg.AWSConfig.KeyPairName = *output.KeyName

	describeInput := &ec2.DescribeKeyPairsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("fingerprint"),
				Values: []*string{aws.String(*output.KeyFingerprint)},
			},
		},
	}

	err = svc.WaitUntilKeyPairExists(describeInput)

	if err != nil {
		logrus.Debugf("WaitUntilKeyPairExists caused %s", err.Error())
		return errors.Wrap(err, fmt.Sprintf("wait until key pair found %s",
			bootstrapKeyPairName))
	}

	return nil
}

func (s *KeyPairStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func (*KeyPairStep) Name() string {
	return StepImportKeyPair
}

func (*KeyPairStep) Description() string {
	return "If no keypair is present in config, creates a new keypair"
}

func (*KeyPairStep) Depends() []string {
	return nil
}
