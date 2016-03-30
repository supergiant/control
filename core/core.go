package core

import (
	"github.com/supergiant/guber"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
)

type Core struct {
	DB          *DB
	K8S         *guber.Client
	EC2         *ec2.EC2
	ELB         *elb.ELB
	AutoScaling *autoscaling.AutoScaling
}

var (
	EtcdEndpoints []string
	K8sHost       string
	K8sUser       string
	K8sPass       string
	AwsRegion     string

	AwsAZ       string
	AwsSgID     string
	AwsSubnetID string
)

func New() *Core {
	c := Core{}
	c.DB = NewDB(EtcdEndpoints)
	c.K8S = guber.NewClient(K8sHost, K8sUser, K8sPass)
	// NOTE / TODO AWS is configured through a file in ~
	awsConf := &aws.Config{Region: aws.String(AwsRegion)}
	c.EC2 = ec2.New(session.New(), awsConf)
	c.ELB = elb.New(session.New(), awsConf)
	c.AutoScaling = autoscaling.New(session.New(), awsConf)
	return &c
}

// Top-level resources
//==============================================================================
func (c *Core) Apps() *AppCollection {
	return &AppCollection{c}
}

func (c *Core) Entrypoints() *EntrypointCollection {
	return &EntrypointCollection{c}
}

func (c *Core) ImageRepos() *ImageRepoCollection {
	return &ImageRepoCollection{c}
}

func (c *Core) Tasks() *TaskCollection {
	return &TaskCollection{c}
}
