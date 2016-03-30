package core

import (
	"os"

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
	etcdEndpoints []string
	k8sHost       string
	k8sUser       string
	k8sPass       string
	awsRegion     string

	AwsAZ string
)

func init() {
	etcdEndpoints = []string{os.Getenv("ETCD_ENDPOINT")}
	k8sHost = os.Getenv("K8S_HOST")
	k8sUser = os.Getenv("K8S_USER")
	k8sPass = os.Getenv("K8S_PASS")
	awsRegion = os.Getenv("AWS_REGION")

	AwsAZ = os.Getenv("AWS_AZ")
}

func New() *Core {
	c := Core{}
	c.DB = NewDB(etcdEndpoints)
	c.K8S = guber.NewClient(k8sHost, k8sUser, k8sPass)
	// NOTE / TODO AWS is configured through a file in ~
	c.EC2 = ec2.New(session.New(), &aws.Config{Region: aws.String(awsRegion)})
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
