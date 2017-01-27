package aws

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

// TODO this and the similar concept in Kubes should be moved to core, not global vars
var globalAWSSession = session.New()

// Provider AWS provider object
type Provider struct {
	Core *core.Core
	EC2  func(*model.Kube) ec2iface.EC2API
	S3   func(*model.Kube) s3iface.S3API
	IAM  func(*model.Kube) iamiface.IAMAPI
	ELB  func(*model.Kube) elbiface.ELBAPI
}

// ValidateAccount validates that the AWS credentials entered work.
func (p *Provider) ValidateAccount(m *model.CloudAccount) error {
	// Doesn't really matter what we do here, as long as it works
	mockKube := &model.Kube{
		CloudAccount: m,
		AWSConfig: &model.AWSKubeConfig{
			Region: "us-east-1",
		},
	}
	_, err := p.EC2(mockKube).DescribeKeyPairs(new(ec2.DescribeKeyPairsInput))
	return err
}

// DeleteNode deletes a Kubernetes minion.
func (p *Provider) DeleteNode(m *model.Node, action *core.Action) error {
	return p.deleteServer(m)
}

func (p *Provider) CreateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.CreateLoadBalancer(m, action)
}

func (p *Provider) UpdateLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.UpdateLoadBalancer(m, action)
}

func (p *Provider) DeleteLoadBalancer(m *model.LoadBalancer, action *core.Action) error {
	return p.Core.K8SProvider.DeleteLoadBalancer(m, action)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

// EC2 client
func EC2(kube *model.Kube) ec2iface.EC2API {
	return ec2.New(globalAWSSession, awsConfig(kube))
}

// ELB client
func ELB(kube *model.Kube) elbiface.ELBAPI {
	return elb.New(globalAWSSession, awsConfig(kube))
}

// S3 client
func S3(kube *model.Kube) s3iface.S3API {
	return s3.New(globalAWSSession, awsConfig(kube))
}

// IAM client
func IAM(kube *model.Kube) iamiface.IAMAPI {
	return iam.New(globalAWSSession, awsConfig(kube))
}

func awsConfig(kube *model.Kube) *aws.Config {
	c := kube.CloudAccount.Credentials
	creds := credentials.NewStaticCredentials(c["access_key"], c["secret_key"], "")
	creds.Get()
	return aws.NewConfig().WithRegion(kube.AWSConfig.Region).WithCredentials(creds)
}

//------------------------------------------------------------------------------

func (p *Provider) setAttrsFromServer(m *model.Node, server *ec2.Instance) {
	m.ProviderID = *server.InstanceId
	m.Name = *server.PrivateDnsName
	m.Size = *server.InstanceType
	m.ProviderCreationTimestamp = *server.LaunchTime
}

func (p *Provider) deleteServer(m *model.Node) error {

	// TODO move out of here
	if m.Kube == nil {
		p.Core.Log.Warnf("Deleting Node %d without deleting server because Kube is nil", *m.ID)
		return nil
	}

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(m.ProviderID)},
	}
	_, err := p.EC2(m.Kube).TerminateInstances(input)
	if isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------

// is it NOT Not Found
func isErrAndNotAWSNotFound(err error) bool {
	return err != nil && !regexp.MustCompile(`([Nn]ot *[Ff]ound|404)`).MatchString(err.Error())
}

func createIAMRole(iamS iamiface.IAMAPI, name string, policy string) error {
	getInput := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}
	_, err := iamS.GetRole(getInput)
	if err == nil {
		return nil
	} else if isErrAndNotAWSNotFound(err) {
		return err
	}
	input := &iam.CreateRoleInput{
		RoleName: aws.String(name),
		Path:     aws.String("/"),
		AssumeRolePolicyDocument: aws.String(policy),
	}
	_, err = iamS.CreateRole(input)
	return err
}

func createIAMRolePolicy(iamS iamiface.IAMAPI, name string, policy string) error {
	getInput := &iam.GetRolePolicyInput{
		RoleName:   aws.String(name),
		PolicyName: aws.String(name),
	}
	_, err := iamS.GetRolePolicy(getInput)
	if err == nil {
		return nil
	} else if isErrAndNotAWSNotFound(err) {
		return err
	}

	putRoleInput := &iam.PutRolePolicyInput{
		RoleName:       aws.String(name),
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(policy),
	}
	_, err = iamS.PutRolePolicy(putRoleInput)
	return err
}

func createIAMInstanceProfile(iamS iamiface.IAMAPI, name string) error {
	getInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	var instanceProfile *iam.InstanceProfile

	resp, err := iamS.GetInstanceProfile(getInput)
	if err != nil {
		if isErrAndNotAWSNotFound(err) {
			return err
		}

		// Create
		input := &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(name),
			Path:                aws.String("/"),
		}
		createResp, createErr := iamS.CreateInstanceProfile(input)
		if createErr != nil {
			return createErr
		}
		instanceProfile = createResp.InstanceProfile

	} else {
		instanceProfile = resp.InstanceProfile
	}

	if len(instanceProfile.Roles) == 0 {
		addInput := &iam.AddRoleToInstanceProfileInput{
			RoleName:            aws.String(name),
			InstanceProfileName: aws.String(name),
		}
		if _, err = iamS.AddRoleToInstanceProfile(addInput); err != nil {
			return err
		}
	}

	return nil
}

func tagAWSResource(ec2S ec2iface.EC2API, idstr string, tags map[string]string) error {
	var ec2Tags []*ec2.Tag
	for key, val := range tags {
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(val),
		})
	}
	input := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(idstr)},
		Tags:      ec2Tags,
	}
	_, err := ec2S.CreateTags(input)
	return err
}

func getAMI(ec2S ec2iface.EC2API) (string, error) {
	images, err := ec2S.DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("architecture"),
				Values: []*string{
					aws.String("x86_64"),
				},
			},
			&ec2.Filter{
				Name: aws.String("owner-id"),
				Values: []*string{
					aws.String("595879546273"),
				},
			},
			&ec2.Filter{
				Name: aws.String("name"),
				Values: []*string{
					aws.String("*stable*"),
				},
			},
			&ec2.Filter{
				Name: aws.String("virtualization-type"),
				Values: []*string{
					aws.String("hvm"),
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	var latestImage *ec2.Image
	for _, image := range images.Images {
		// latest year
		if latestImage == nil {
			latestImage = image
			continue
		}

		latestImageCreationTime, err := time.Parse("2006-01-02T15:04:05.000Z", *latestImage.CreationDate)
		if err != nil {
			panic(err)
		}
		imageCreationTime, err := time.Parse("2006-01-02T15:04:05.000Z", *image.CreationDate)
		if err != nil {
			panic(err)
		}

		if imageCreationTime.After(latestImageCreationTime) {
			latestImage = image
		}
	}
	return *latestImage.ImageId, nil
}

func etcdToken(num string) (string, error) {
	resp, err := http.Get("https://discovery.etcd.io/new?size=" + num + "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
