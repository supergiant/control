package util

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"

	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/sgerrors"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
)

type CloudAccountValidator interface {
	ValidateCredentials(cloudAccount *model.CloudAccount) error
}

type CloudAccountValidatorImpl struct {
	digitalOcean func(map[string]string) error
	aws          func(map[string]string) error
}

func NewCloudAccountValidator() *CloudAccountValidatorImpl {
	return &CloudAccountValidatorImpl{
		digitalOcean: validateDigitalOceanCredentials,
		aws:          validateAWSCredentials,
	}
}

func (validator *CloudAccountValidatorImpl) ValidateCredentials(cloudAccount *model.CloudAccount) error {
	switch cloudAccount.Provider {
	case clouds.DigitalOcean:
		return validator.digitalOcean(cloudAccount.Credentials)
	case clouds.AWS:
		return validator.aws(cloudAccount.Credentials)
	}

	return sgerrors.ErrUnsupportedProvider
}

func validateDigitalOceanCredentials(creds map[string]string) error {
	config := &steps.DOConfig{}
	err := BindParams(creds, config)

	if err != nil {
		return err
	}

	ts := &digitaloceansdk.TokenSource{
		AccessToken: config.AccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)
	client := godo.NewClient(oauthClient)

	_, _, err = client.Droplets.List(context.Background(), new(godo.ListOptions))

	return err
}

func validateAWSCredentials(creds map[string]string) error {
	config := &steps.AWSConfig{}
	err := BindParams(creds, config)

	if err != nil {
		return err
	}

	awsCfg := aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(config.KeyID, config.Secret, ""),
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: awsCfg,
	})

	if err != nil {
		return err
	}

	ec2Client := ec2.New(sess)

	_, err = ec2Client.DescribeKeyPairs(new(ec2.DescribeKeyPairsInput))
	return err
}
