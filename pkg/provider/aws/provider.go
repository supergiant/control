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
	"github.com/supergiant/supergiant/pkg/kubernetes"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
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

// CreateVolume creates a kubernetes Volume.
func (p *Provider) CreateVolume(m *model.Volume, action *core.Action) error {
	return p.createVolume(m, nil)
}

// KubernetesVolumeDefinition defines object layout of a AWS volume.
func (p *Provider) KubernetesVolumeDefinition(m *model.Volume) *kubernetes.Volume {
	return &kubernetes.Volume{
		Name: m.Name,
		AwsElasticBlockStore: &kubernetes.AwsElasticBlockStore{
			VolumeID: m.ProviderID,
			FSType:   "ext4",
		},
	}
}

// ResizeVolume resizes a AWS volume.
func (p *Provider) ResizeVolume(m *model.Volume, action *core.Action) error {
	return p.resizeVolume(m, action)
}

// WaitForVolumeAvailable waits for AWS volume to be available.
func (p *Provider) WaitForVolumeAvailable(m *model.Volume, action *core.Action) error {
	return p.waitForAvailable(m)
}

// DeleteVolume deletes a aws volume.
func (p *Provider) DeleteVolume(m *model.Volume, action *core.Action) error {
	return p.deleteVolume(m)
}

// CreateEntrypoint creates a AWS LoadBalancer
func (p *Provider) CreateEntrypoint(m *model.Entrypoint, action *core.Action) error {
	return p.createELB(m)
}

// DeleteEntrypoint deletes a aws loadbalancer.
func (p *Provider) DeleteEntrypoint(m *model.Entrypoint, action *core.Action) error {
	return p.deleteELB(m)
}

// CreateEntrypointListener creates a listener for a aws loadbalancer.
func (p *Provider) CreateEntrypointListener(m *model.EntrypointListener, action *core.Action) error {
	input := &elb.CreateLoadBalancerListenersInput{
		LoadBalancerName: aws.String(m.Entrypoint.ProviderID),
		Listeners: []*elb.Listener{
			{
				LoadBalancerPort: aws.Int64(m.EntrypointPort),
				Protocol:         aws.String(m.EntrypointProtocol),
				InstancePort:     aws.Int64(m.NodePort),
				InstanceProtocol: aws.String(m.NodeProtocol),
			},
		},
	}
	_, err := p.ELB(m.Entrypoint.Kube).CreateLoadBalancerListeners(input)
	return err
}

// DeleteEntrypointListener deletes a listener form an aws loadbalancer.
func (p *Provider) DeleteEntrypointListener(m *model.EntrypointListener, action *core.Action) error {
	input := &elb.DeleteLoadBalancerListenersInput{
		LoadBalancerName: aws.String(m.Entrypoint.ProviderID),
		LoadBalancerPorts: []*int64{
			aws.Int64(m.EntrypointPort),
		},
	}
	_, err := p.ELB(m.Entrypoint.Kube).DeleteLoadBalancerListeners(input)
	if isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods                                                            //
////////////////////////////////////////////////////////////////////////////////

// EC2 client
func EC2(kube *model.Kube) ec2iface.EC2API {
	return ec2.New(globalAWSSession, awsConfig(kube))
}

// S3 client
func S3(kube *model.Kube) s3iface.S3API {
	return s3.New(globalAWSSession, awsConfig(kube))
}

// ELB client
func ELB(kube *model.Kube) elbiface.ELBAPI {
	return elb.New(globalAWSSession, awsConfig(kube))
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

//func (p *Provider) createNode(m *model.Node) error {
//	server, err := p.createServer(m)
//	if err != nil {
//		return err
//	}
//	p.setAttrsFromServer(m, server)
//	if err := p.Core.DB.Save(m); err != nil {
//		return err
//	}
//	for _, entrypoint := range m.Kube.Entrypoints {
//		if err := p.registerNodes(entrypoint, m); err != nil {
//			return err
//		}
//	}
//	return nil
//}

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

func (p *Provider) createELB(m *model.Entrypoint) error {

	var subnets []*string
	for _, subnet := range m.Kube.AWSConfig.PublicSubnetIPRange {
		if subnet["subnet_id"] != "" {
			subnets = append(subnets, aws.String(subnet["subnet_id"]))
		}
	}

	params := &elb.CreateLoadBalancerInput{
		Listeners: []*elb.Listener{ // NOTE we must provide at least 1 listener, it is currently arbitrary
			{
				InstancePort:     aws.Int64(420),
				LoadBalancerPort: aws.Int64(420),
				Protocol:         aws.String("TCP"),
			},
		},
		LoadBalancerName: aws.String(m.ProviderID),
		Scheme:           aws.String("internet-facing"),
		SecurityGroups: []*string{
			aws.String(m.Kube.AWSConfig.ELBSecurityGroupID),
		},
		Subnets: subnets,
	}
	resp, err := p.ELB(m.Kube).CreateLoadBalancer(params)
	if err != nil {
		return err
	}

	// Save Address
	m.Address = *resp.DNSName
	err = p.Core.DB.Save(m)
	if err != nil {
		return err
	}

	err = p.registerNodes(m, m.Kube.Nodes...)
	if err != nil {
		return err
	}

	// Configure health check
	healthParams := &elb.ConfigureHealthCheckInput{
		LoadBalancerName: aws.String(m.ProviderID),
		HealthCheck: &elb.HealthCheck{
			Target:             aws.String("HTTPS:10250/healthz"),
			HealthyThreshold:   aws.Int64(2),
			UnhealthyThreshold: aws.Int64(10),
			Interval:           aws.Int64(30),
			Timeout:            aws.Int64(5),
		},
	}
	_, err = p.ELB(m.Kube).ConfigureHealthCheck(healthParams)
	return err
}

func (p *Provider) registerNodes(m *model.Entrypoint, nodes ...*model.Node) error {
	var elbInstances []*elb.Instance
	for _, node := range nodes {
		elbInstances = append(elbInstances, &elb.Instance{
			InstanceId: aws.String(node.ProviderID),
		})
	}
	input := &elb.RegisterInstancesWithLoadBalancerInput{
		LoadBalancerName: aws.String(m.ProviderID),
		Instances:        elbInstances,
	}
	_, err := p.ELB(m.Kube).RegisterInstancesWithLoadBalancer(input)
	return err
}

func (p *Provider) deleteELB(m *model.Entrypoint) error {
	// Delete ELB
	params := &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(m.ProviderID),
	}
	_, err := p.ELB(m.Kube).DeleteLoadBalancer(params)
	if isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

func (p *Provider) createVolume(volume *model.Volume, snapshotID *string) error {
	volInput := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(volume.Kube.AWSConfig.AvailabilityZone),
		VolumeType:       aws.String(volume.Type),
		Size:             aws.Int64(int64(volume.Size)),
		SnapshotId:       snapshotID,
	}
	awsVol, err := p.EC2(volume.Kube).CreateVolume(volInput)
	if err != nil {
		return err
	}

	volume.ProviderID = *awsVol.VolumeId
	volume.Size = int(*awsVol.Size)
	err = p.Core.DB.Save(volume)
	if err != nil {
		return err
	}

	tagsInput := &ec2.CreateTagsInput{
		Resources: []*string{
			awsVol.VolumeId,
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(volume.Name),
			},
		},
	}
	_, err = p.EC2(volume.Kube).CreateTags(tagsInput)
	return err
}

func (p *Provider) resizeVolume(m *model.Volume, action *core.Action) error {
	snapshot, err := p.createSnapshot(m, action)
	if err != nil {
		return err
	}
	if err := p.deleteVolume(m); err != nil {
		return err
	}
	if err := p.createVolume(m, snapshot.SnapshotId); err != nil {
		return err
	}
	if err := p.deleteSnapshot(m, snapshot); err != nil {
		p.Core.Log.Errorf("Error deleting snapshot %s: %s", *snapshot.SnapshotId, err.Error())
	}
	return nil
}

func (p *Provider) deleteVolume(volume *model.Volume) error {
	if volume.ProviderID == "" {
		return nil
	}
	if err := p.waitForAvailable(volume); err != nil {
		return err
	}
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(volume.ProviderID),
	}
	if _, err := p.EC2(volume.Kube).DeleteVolume(input); isErrAndNotAWSNotFound(err) {
		return err
	}
	return nil
}

func (p *Provider) waitForAvailable(volume *model.Volume) error {
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(volume.ProviderID),
				},
			},
		},
	}

	desc := fmt.Sprintf("EBS volume %s to be available or deleted", volume.Name)
	return util.WaitFor(desc, 5*time.Minute, 10*time.Second, func() (bool, error) {
		resp, err := p.EC2(volume.Kube).DescribeVolumes(input)
		if err != nil {
			return false, err
		}
		if len(resp.Volumes) == 0 {
			return true, nil
		}
		state := *resp.Volumes[0].State
		return state == "available" || state == "deleted", nil
	})
}

func (p *Provider) createSnapshot(volume *model.Volume, action *core.Action) (*ec2.Snapshot, error) {
	input := &ec2.CreateSnapshotInput{
		Description: aws.String(fmt.Sprintf("%s-%s", volume.Name, time.Now())),
		VolumeId:    aws.String(volume.ProviderID),
	}
	snapshot, err := p.EC2(volume.Kube).CreateSnapshot(input)
	if err != nil {
		return nil, err
	}
	getInput := &ec2.DescribeSnapshotsInput{
		SnapshotIds: []*string{snapshot.SnapshotId},
	}

	desc := fmt.Sprintf("Snapshot %s to complete", volume.Name)
	waitErr := action.CancellableWaitFor(desc, 4*time.Hour, 15*time.Second, func() (bool, error) {
		resp, err := p.EC2(volume.Kube).DescribeSnapshots(getInput)
		if err != nil {
			return false, err
		}
		if len(resp.Snapshots) == 0 {
			return true, nil
		}
		state := *resp.Snapshots[0].State
		return state == "completed", nil
	})
	if waitErr != nil {
		return nil, waitErr
	}

	return snapshot, nil
}

func (p *Provider) deleteSnapshot(volume *model.Volume, snapshot *ec2.Snapshot) error {
	input := &ec2.DeleteSnapshotInput{
		SnapshotId: snapshot.SnapshotId,
	}
	if _, err := p.EC2(volume.Kube).DeleteSnapshot(input); isErrAndNotAWSNotFound(err) {
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
