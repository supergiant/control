package amazon

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

const ImportKeyPairStepName = "awskeypairstep"

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
	steps.RegisterStep(ImportKeyPairStepName, NewImportKeyPairStep(fn))
}

func (s *KeyPairStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)

	EC2, err := s.GetEC2(cfg.AWSConfig)
	if err != nil {
		return ErrAuthorization
	}

	// Create key for user
	userKeyPairName := util.MakeKeyName(cfg.AWSConfig.KeyPairName, true)

	req := &ec2.ImportKeyPairInput{
		KeyName:           &userKeyPairName,
		PublicKeyMaterial: []byte(cfg.SshConfig.PublicKey),
	}

	log.Infof("[%s] - importing user certificate as keypair %s", s.Name(), userKeyPairName)
	output, err := EC2.ImportKeyPairWithContext(ctx, req)
	if err != nil {
		cfg.AWSConfig.KeyPairName = *output.KeyName
	}

	bootstrapKeyPairName := util.MakeKeyName(cfg.AWSConfig.KeyPairName, false)
	log.Infof("[%s] - importing cluster bootstrap certificate as keypair %s", s.Name(), bootstrapKeyPairName)
	req = &ec2.ImportKeyPairInput{
		KeyName:           &bootstrapKeyPairName,
		PublicKeyMaterial: []byte(cfg.SshConfig.BootstrapPublicKey),
	}

	output, err = EC2.ImportKeyPairWithContext(ctx, req)
	if err != nil {
		return errors.Wrap(ErrImportKeyPair, err.Error())
	}
	return nil
}

func (s *KeyPairStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	return nil
}

func (*KeyPairStep) Name() string {
	return ImportKeyPairStepName
}

func (*KeyPairStep) Description() string {
	return "If no keypair is present in config, creates a new keypair and writes private RSA key to the database and config"
}

func (*KeyPairStep) Depends() []string {
	return nil
}
