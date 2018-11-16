package amazon

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"time"
)

const StepImportKeyPair = "awskeypairstep"

var (
	keyPairAttemptCound = 5
)

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

	bootstrapKeyPairName := util.MakeKeyName(cfg.ClusterName, false)
	log.Infof("[%s] - importing cluster bootstrap certificate as keypair %s",
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
			return nil
		}
		return errors.Wrap(ErrImportKeyPair, err.Error())
	}

	cfg.AWSConfig.KeyPairName = *output.KeyName

	delay := time.Second * 10
	// Wait until key pair become available
	for cnt := 0; cnt < keyPairAttemptCound; cnt++ {
		describeInput := &ec2.DescribeKeyPairsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("fingerprint"),
					Values: []*string{aws.String(*output.KeyFingerprint)},
				},
			},
		}

		_, err = EC2.DescribeKeyPairs(describeInput)

		if err != nil {
			logrus.Errorf("describe key pair %s error %v", bootstrapKeyPairName, err)
			time.Sleep(delay)
			delay = delay * 2
		} else {
			return nil
		}
	}

	return errors.Wrap(err, fmt.Sprintf("wait until key pair found %s",
		bootstrapKeyPairName))
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
