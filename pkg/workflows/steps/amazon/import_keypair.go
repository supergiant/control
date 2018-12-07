package amazon

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/util"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const StepImportKeyPair = "awskeypairstep"

//KeyPairStep represents creation of keypair in aws
//since there is hard cap on keypairs per account supergiant will create one per clster
type KeyPairStep struct {
	GetEC2 GetEC2Fn
}

func NewImportKeyPairStep(fn GetEC2Fn) *KeyPairStep {
	return &KeyPairStep{
		GetEC2: fn,
	}
}

//InitImportKeyPair add the step to the registry
func InitImportKeyPair(fn GetEC2Fn) {
	steps.RegisterStep(StepImportKeyPair, NewImportKeyPairStep(fn))
}

//Verifies that a key exists,
func (s *KeyPairStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return ErrAuthorization
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
		PublicKeyMaterial: []byte(cfg.SshConfig.BootstrapPublicKey),
	}

	if cfg.LogBootstrapPrivateKey {
		key := strings.Replace(cfg.SshConfig.BootstrapPrivateKey, "\\n", "\n", -1)
		log.Infof("[%s] - bootstrap private key", s.Name())
		fmt.Fprintf(w, key)
	}

	output, err := EC2.ImportKeyPairWithContext(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "InvalidKeyPair.Duplicate") {
			cfg.AWSConfig.KeyPairName = bootstrapKeyPairName
			return errors.Wrap(ErrImportKeyPair, err.Error())
		}
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

	err = EC2.WaitUntilKeyPairExists(describeInput)

	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			logrus.Debugf("WaitUntilKeyPairExists caused %s", err.Message())
		}
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
