package workflows

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/clouds/awssdk"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/testutils/assert"
	"github.com/supergiant/supergiant/pkg/util"
	"github.com/supergiant/supergiant/pkg/workflows/steps"
	"github.com/supergiant/supergiant/pkg/workflows/steps/amazon"
)

const defaultETCDHost = "http://127.0.0.1:2379"

var defaultConfig clientv3.Config

func init() {
	assert.MustRunETCD(assert.DefaultETCDURL)
	defaultConfig = clientv3.Config{
		Endpoints: []string{defaultETCDHost},
	}
}

func TestAWSEC2Create(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	kv := storage.NewETCDRepository(defaultConfig)
	accounts := account.NewService(account.DefaultStoragePrefix, kv)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cloudAcc := &model.CloudAccount{
		Name:        "sgtest" + util.RandomString(16),
		Provider:    clouds.AWS,
		Credentials: map[string]string{},
	}
	err := accounts.Create(ctx, cloudAcc)
	require.NoError(t, err)

	key := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_KEY")
	region := "eu-west-1"

	sdk, err := awssdk.New(region, key, secret, "")
	require.NoError(t, err)

	out, err := sdk.EC2.DescribeImagesWithContext(ctx, &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("architecture"),
				Values: []*string{
					aws.String("x86_64"),
				},
			},
			{
				Name: aws.String("name"),
				Values: []*string{
					aws.String("*bionic*"),
				},
			},
		},
	})
	require.NoError(t, err)

	latestImage := out.Images[0]
	for _, ami := range out.Images {
		latestImageCreationTime, _ := time.Parse("2006-01-02T15:04:05.000Z", *latestImage.CreationDate)
		imageCreationTime, _ := time.Parse("2006-01-02T15:04:05.000Z", *ami.CreationDate)
		if imageCreationTime.After(latestImageCreationTime) {
			latestImage = ami
		}
	}
	cfg := &steps.Config{
		ClusterName:      "sgtest",
		CloudAccountName: cloudAcc.Name,
		IsMaster:         true,
		MasterNodes:      steps.MasterMap{},
		AWSConfig: steps.AWSConfig{
			EC2Config: steps.EC2Config{
				VolumeSize:   2,
				EbsOptimized: false,
				GPU:          false,
				InstanceType: "m4.large",
				ImageID:      *latestImage.ImageId,
			},
			Secret: secret,
			KeyID:  key,
			Region: region,
		},
	}
	createKeyPair := amazon.NewKeyPairStep(accounts)

	err = createKeyPair.Run(ctx, os.Stdout, cfg)
	if err != nil {
		err2 := createKeyPair.Rollback(ctx, os.Stdout, cfg)
		require.NoError(t, err2)
		require.NoError(t, err)
	}

	ec2Steps := make([]*amazon.StepCreateInstance, 0)

	for i := 0; i < 3; i++ {
		createEC2Step := amazon.NewCreateInstance()
		ec2Steps = append(ec2Steps, createEC2Step)

		err = createEC2Step.Run(ctx, os.Stdout, cfg)
		if err != nil {
			break
		}
		require.NoError(t, err)

		acc, err := accounts.Get(ctx, cfg.CloudAccountName)
		require.NoError(t, err)
		require.NotEmpty(t, acc.Credentials)
		require.NotEqual(t, "", acc.Credentials[model.KeyPrivateKey])
	}

	for _, v := range ec2Steps {
		err := v.Rollback(ctx, os.Stdout, cfg)
		//log error instead of failing test to make sure everything is deleted
		logrus.Error(err)
	}
}
