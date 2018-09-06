package amazon

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

//KeyPairStep represents creation of keypair in aws
//since there is hard cap on keypairs per account supergiant will create one per clster
type KeyPairStep struct {
	accounts accountWrapper
}

func NewKeyPairStep(wrapper accountWrapper) *KeyPairStep {
	return &KeyPairStep{
		accounts: wrapper,
	}
}

//InitCreateKeyPair add the step to the registry
func InitCreateKeyPair(wrapper accountWrapper) {
	steps.RegisterStep(StepNameKeyPair, NewKeyPairStep(wrapper))
}

const StepNameKeyPair = "keypairstep"

type accountWrapper interface {
	Get(context.Context, string) (*model.CloudAccount, error)
	Update(context.Context, *model.CloudAccount) error
}

func (s *KeyPairStep) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	log := util.GetLogger(w)
	log.Infof("[%s] - started!", s.Name())

	account, err := s.accounts.Get(ctx, cfg.CloudAccountName)
	if err != nil {
		return errors.Wrap(err, "aws: no cloud account found")
	}

	sdk, err := GetSDK(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	//If a user chooses to use pre-existing ec2 keypair it should be in the database already
	if cfg.AWSConfig.KeyPairName != "" {
		err := s.GetKeyFromAccount(cfg, account, log)
		if err != nil {
			return err
		}
	} else {
		//create new keypair with the same name as cloud account
		cfg.AWSConfig.KeyPairName = cfg.CloudAccountName

		out, err := sdk.EC2.CreateKeyPairWithContext(ctx, &ec2.CreateKeyPairInput{KeyName: &cfg.AWSConfig.KeyPairName})
		if err != nil {
			if strings.Contains(err.Error(), "InvalidKeyPair.Duplicate") {
				err := s.GetKeyFromAccount(cfg, account, log)
				if err != nil {
					return err
				} else {
					return errors.Wrap(err, "aws: failed to create key pair")
				}
			}

			return errors.Wrap(err, "aws: failed to create key pair")
		}

		if out.KeyMaterial == nil || out.KeyFingerprint == nil {
			return errors.New("aws: faield to obtain keypair data")
		}

		account.Credentials[clouds.CredsPrivateKey] = *out.KeyMaterial

		if err := s.accounts.Update(ctx, account); err != nil {
			return err
		}

		log.Debugf("[%s] created new RSA key for keypair %s", StepNameKeyPair, cfg.AWSConfig.KeyPairName)
		log.Debugln(*out.KeyMaterial)
	}
	log.Infof("[%s] - success!", s.Name())
	return nil
}

func (s *KeyPairStep) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	sdk, err := GetSDK(cfg.AWSConfig)
	if err != nil {
		return errors.New("aws: authorization")
	}

	cfg.SshConfig.PrivateKey = ""

	_, err = sdk.EC2.DeleteKeyPairWithContext(ctx, &ec2.DeleteKeyPairInput{
		KeyName: aws.String(cfg.AWSConfig.KeyPairName),
	})
	return err
}

func (s *KeyPairStep) GetKeyFromAccount(cfg *steps.Config, account *model.CloudAccount, log *logrus.Logger) error {
	key, ok := account.Credentials[clouds.CredsPrivateKey]
	if !ok || key == "" {
		log.Errorf("[%s] - no ssh key present in database, aborting", s.Name())
		return errors.New("aws: no ssh key found")
	}
	cfg.SshConfig.PrivateKey = key
	return nil
}

func (*KeyPairStep) Name() string {
	return StepNameKeyPair
}

func (*KeyPairStep) Description() string {
	return "If no keypair is present in config, creates a new keypair and writes private RSA key to the database and config"
}

func (*KeyPairStep) Depends() []string {
	return nil
}
