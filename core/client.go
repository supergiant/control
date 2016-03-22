package core

import (
	"guber"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Client struct {
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

func NewClient() *Client {
	c := Client{}
	c.DB = NewDB(etcdEndpoints)
	c.K8S = guber.NewClient(k8sHost, k8sUser, k8sPass)
	// NOTE / TODO AWS is configured through a file in ~
	c.EC2 = ec2.New(session.New(), &aws.Config{Region: aws.String(awsRegion)})
	return &c
}

// Top-level resources
//==============================================================================
func (c *Client) Apps() *AppResource {
	return &AppResource{c}
}

func (c *Client) ImageRepos() *ImageRepoResource {
	return &ImageRepoResource{c}
}

func (c *Client) Jobs() *JobResource {
	return &JobResource{c}
}
