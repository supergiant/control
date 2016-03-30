package core

import (
	"github.com/supergiant/guber"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Core struct {
	DB  *DB
	K8S *guber.Client
	EC2 *ec2.EC2
}

var (
	EtcdEndpoints []string
	K8sHost       string
	K8sUser       string
	K8sPass       string
	AwsRegion     string

	AwsAZ string
)

func New() *Core {
	c := Core{}
	c.DB = NewDB(EtcdEndpoints)
	c.K8S = guber.NewClient(K8sHost, K8sUser, K8sPass)
	// NOTE / TODO AWS is configured through a file in ~
	c.EC2 = ec2.New(session.New(), &aws.Config{Region: aws.String(AwsRegion)})
	return &c
}

// Top-level resources
//==============================================================================
func (c *Core) Apps() *AppCollection {
	return &AppCollection{c}
}

func (c *Core) ImageRepos() *ImageRepoCollection {
	return &ImageRepoCollection{c}
}

func (c *Core) Tasks() *TaskCollection {
	return &TaskCollection{c}
}
