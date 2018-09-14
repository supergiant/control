package amazon

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const StepName = "awskeypairstep"

//KeyPairStep represents creation of keypair in aws
//since there is hard cap on keypairs per account supergiant will create one per clster
type KeyPairStep struct {
}

func NewKeyPairStep() *KeyPairStep {
	return &KeyPairStep{}
}

//InitCreateKeyPair add the step to the registry
func InitCreateKeyPair() {
	steps.RegisterStep(StepName, NewKeyPairStep())
}

func (s *KeyPairStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	log.Infof("[%s] - started!", s.Name())

	sdk, err := GetSDK(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	// Create key for user
	userKeyPairName := util.MakeKeyName(cfg.AWSConfig.KeyPairName, true)

	req := &ec2.ImportKeyPairInput{
		KeyName:           &userKeyPairName,
		PublicKeyMaterial: []byte(cfg.SshConfig.PublicKey),
	}

	output, err := sdk.EC2.ImportKeyPairWithContext(ctx, req)

	keyPairName := util.MakeKeyName(cfg.AWSConfig.KeyPairName, false)

	if err != nil {
		cfg.AWSConfig.KeyPairName = *output.KeyName
	}

	req = &ec2.ImportKeyPairInput{
		KeyName:           &keyPairName,
		PublicKeyMaterial: []byte(cfg.SshConfig.BootstrapPublicKey),
	}

	output, err = sdk.EC2.ImportKeyPairWithContext(ctx, req)

	if err != nil {
		return errors.Wrap(err, "create provision key pair")
	}

	log.Infof("[%s] - success!", s.Name())
	return nil
}

func (s *KeyPairStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	sdk, err := GetSDK(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	_, err = sdk.EC2.DeleteKeyPairWithContext(ctx, &ec2.DeleteKeyPairInput{
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
