package util

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/clouds/digitaloceansdk"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/sgerrors"
	"github.com/supergiant/control/pkg/workflows/steps"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type CloudAccountValidator interface {
	ValidateCredentials(cloudAccount *model.CloudAccount) error
}

type CloudAccountValidatorImpl struct {
	digitalOcean func(map[string]string) error
	aws          func(map[string]string) error
	gce          func(map[string]string) error
	azure        func(map[string]string) error
	openstack    func(map[string]string) error
}

func NewCloudAccountValidator() *CloudAccountValidatorImpl {
	return &CloudAccountValidatorImpl{
		digitalOcean: validateDigitalOceanCredentials,
		aws:          validateAWSCredentials,
		gce:          validateGCECredentials,
		azure:        validateAzureCredentials,
		openstack:    validateOpenstackCredentials,
	}
}

func (v *CloudAccountValidatorImpl) ValidateCredentials(cloudAccount *model.CloudAccount) error {
	switch cloudAccount.Provider {
	case clouds.DigitalOcean:
		return v.digitalOcean(cloudAccount.Credentials)
	case clouds.AWS:
		return v.aws(cloudAccount.Credentials)
	case clouds.GCE:
		return v.gce(cloudAccount.Credentials)
	case clouds.Azure:
		return v.azure(cloudAccount.Credentials)
	case clouds.OpenStack:
		return v.openstack(cloudAccount.Credentials)
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

func validateGCECredentials(creds map[string]string) error {
	clientScopes := []string{
		compute.ComputeScope,
		compute.CloudPlatformScope,
		dns.NdevClouddnsReadwriteScope,
		compute.DevstorageFullControlScope,
	}

	conf := jwt.Config{
		Email:      creds[clouds.GCEClientEmail],
		PrivateKey: []byte(creds[clouds.GCEPrivateKey]),
		Scopes:     clientScopes,
		TokenURL:   creds[clouds.GCETokenURI],
	}

	client := conf.Client(context.Background())

	computeService, err := compute.New(client)
	if err != nil {
		logrus.Errorf("Error creating compute object %v", err)
		return err
	}

	// find the ubuntu image.
	_, err = computeService.Images.GetFromFamily(
		"ubuntu-os-cloud", "ubuntu-1804-lts").Do()
	if err != nil {
		logrus.Errorf("Error getting image %v", err)
		return err
	}

	return nil
}

func validateAzureCredentials(creds map[string]string) error {
	if creds == nil {
		return ErrInvalidCredentials
	}

	for _, k := range []string{
		clouds.AzureTenantID,
		clouds.AzureSubscriptionID,
		clouds.AzureClientID,
		clouds.AzureClientSecret,
	} {
		creds[k] = strings.TrimSpace(creds[k])
		if creds[k] == "" {
			return errors.Wrapf(ErrInvalidCredentials, "azure: %s should be provided", k)
		}
	}

	return nil
}

func validateOpenstackCredentials(creds map[string]string) error {
	if creds == nil {
		return ErrInvalidCredentials
	}

	for _, k := range []string{
		clouds.OpenstackAuthUrl,
		clouds.OpenstackDomainId,
		clouds.OpenstackPassword,
		clouds.OpenstackTenantName,
		clouds.OpenstackUserName,
	} {
		creds[k] = strings.TrimSpace(creds[k])
		if creds[k] == "" {
			return errors.Wrapf(ErrInvalidCredentials, "%s: %s should be provided",
				clouds.OpenStack, k)
		}
	}

	return nil
}
