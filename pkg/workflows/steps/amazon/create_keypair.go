package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "awskeypairstep"

//KeyPairStep represents creation of keypair in aws
//since there is hard cap on keypairs per account supergiant will create one per clster
type KeyPairStep struct {
	GetEC2 GetEC2Fn
}

func NewKeyPairStep(fn GetEC2Fn) *KeyPairStep {
	return &KeyPairStep{
		GetEC2: fn,
	}
}

//InitCreateKeyPair add the step to the registry
func InitCreateKeyPair(fn GetEC2Fn) {
	steps.RegisterStep(StepName, NewKeyPairStep(fn))
}

func (s *KeyPairStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	// Create key for user
	userKeyPairName := util.MakeKeyName(cfg.AWSConfig.KeyPairName, true)

	req := &ec2.ImportKeyPairInput{
		KeyName:           &userKeyPairName,
		PublicKeyMaterial: []byte(cfg.SshConfig.PublicKey),
	}

	output, err := EC2.ImportKeyPairWithContext(ctx, req)

	keyPairName := util.MakeKeyName(cfg.AWSConfig.KeyPairName, false)

	if err != nil {
		cfg.AWSConfig.KeyPairName = *output.KeyName
	}

	req = &ec2.ImportKeyPairInput{
		KeyName:           &keyPairName,
		PublicKeyMaterial: []byte(cfg.SshConfig.BootstrapPublicKey),
	}

	output, err = EC2.ImportKeyPairWithContext(ctx, req)

	if err != nil {
		return errors.Wrap(err, "create provision key pair")
	}

	log.Infof("[%s] - success!", s.Name())
	return nil
}

func (s *KeyPairStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	_, err = EC2.DeleteKeyPairWithContext(ctx, &ec2.DeleteKeyPairInput{
		KeyName: aws.String(cfg.AWSConfig.KeyPairName),
	})
	return err
}

func (*KeyPairStep) Name() string {
	return StepName
}

func (*KeyPairStep) Description() string {
	return "If no keypair is present in config, creates a new keypair and writes private RSA key to the database and config"
}

func (*KeyPairStep) Depends() []string {
	return nil
}
