package amazon

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/workflows/steps"
)

const (
	roleMaster = "master"

	assumePolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": { "Service": "ec2.amazonaws.com"},
      "Action": "sts:AssumeRole"
    }
  ]
}`

	// https://github.com/kubernetes/cloud-provider-aws#iam-policy
	masterIAMPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:DescribeLaunchConfigurations",
        "autoscaling:DescribeTags",
        "ec2:DescribeInstances",
        "ec2:DescribeRegions",
        "ec2:DescribeRouteTables",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSubnets",
        "ec2:DescribeVolumes",
        "ec2:CreateSecurityGroup",
        "ec2:CreateTags",
        "ec2:CreateVolume",
        "ec2:ModifyInstanceAttribute",
        "ec2:ModifyVolume",
        "ec2:AttachVolume",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CreateRoute",
        "ec2:DeleteRoute",
        "ec2:DeleteSecurityGroup",
        "ec2:DeleteVolume",
        "ec2:DetachVolume",
        "ec2:RevokeSecurityGroupIngress",
        "ec2:DescribeVpcs",
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:AttachLoadBalancerToSubnets",
        "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
        "elasticloadbalancing:CreateLoadBalancer",
        "elasticloadbalancing:CreateLoadBalancerPolicy",
        "elasticloadbalancing:CreateLoadBalancerListeners",
        "elasticloadbalancing:ConfigureHealthCheck",
        "elasticloadbalancing:DeleteLoadBalancer",
        "elasticloadbalancing:DeleteLoadBalancerListeners",
        "elasticloadbalancing:DescribeLoadBalancers",
        "elasticloadbalancing:DescribeLoadBalancerAttributes",
        "elasticloadbalancing:DetachLoadBalancerFromSubnets",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:ModifyLoadBalancerAttributes",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:CreateListener",
        "elasticloadbalancing:CreateTargetGroup",
        "elasticloadbalancing:DeleteListener",
        "elasticloadbalancing:DeleteTargetGroup",
        "elasticloadbalancing:DescribeListeners",
        "elasticloadbalancing:DescribeLoadBalancerPolicies",
        "elasticloadbalancing:DescribeTargetGroups",
        "elasticloadbalancing:DescribeTargetHealth",
        "elasticloadbalancing:ModifyListener",
        "elasticloadbalancing:ModifyTargetGroup",
        "elasticloadbalancing:RegisterTargets",
        "elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
        "iam:CreateServiceLinkedRole",
        "kms:DescribeKey"
      ],
      "Resource": [
        "*"
      ]
    },
  ]
}`

	nodeIAMPolicy = `{
      "Version": "2012-10-17",
      "Statement": [
          {
              "Effect": "Allow",
              "Action": [
                  "ec2:DescribeInstances",
                  "ec2:DescribeRegions",
                  "ecr:GetAuthorizationToken",
                  "ecr:BatchCheckLayerAvailability",
                  "ecr:GetDownloadUrlForLayer",
                  "ecr:GetRepositoryPolicy",
                  "ecr:DescribeRepositories",
                  "ecr:ListImages",
                  "ecr:BatchGetImage"
              ],
              "Resource": "*"
          } 
      ]
}`
)

var (
	ErrEmptyResponse = errors.New("empty response")
)

const (
	StepNameCreateInstanceProfiles = "aws_create_instance_profiles"
)

type StepCreateInstanceProfiles struct {
	GetIAM GetIAMFn
}

func InitCreateInstanceProfiles(iamfn GetIAMFn) {
	steps.RegisterStep(StepNameCreateInstanceProfiles, NewCreateInstanceProfiles(iamfn))
}

func NewCreateInstanceProfiles(iamfn GetIAMFn) *StepCreateInstanceProfiles {
	return &StepCreateInstanceProfiles{
		GetIAM: iamfn,
	}
}

func (s StepCreateInstanceProfiles) Run(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	iamS, err := s.GetIAM(cfg.AWSConfig)
	if err != nil {
		return errors.Wrapf(err, "%s: failed to authorize in AWS: %v", s.Name(), err)
	}

	// TODO: use a separate config: aws.Nodes/aws.Masters?
	cfg.AWSConfig.MastersInstanceProfile, err = ensureIAMProfile(ctx, iamS, cfg.ClusterID, string(model.RoleMaster))
	if err != nil {
		return errors.Wrapf(err, "%s: failed to authorize in AWS: %v", s.Name(), err)
	}
	logrus.Infof("%s: set up %s instance profile", s.Name(), cfg.AWSConfig.MastersInstanceProfile)

	cfg.AWSConfig.NodesInstanceProfile, err = ensureIAMProfile(ctx, iamS, cfg.ClusterID, string(model.RoleNode))
	if err != nil {
		return errors.Wrapf(err, "%s: failed to authorize in AWS: %v", s.Name(), err)
	}
	logrus.Infof("%s: set up %s instance profile", s.Name(), cfg.AWSConfig.NodesInstanceProfile)

	return nil
}

func (s StepCreateInstanceProfiles) Rollback(ctx context.Context, w io.Writer, cfg *steps.Config) error {
	// TODO: implement instance profile removal
	return nil
}

func (StepCreateInstanceProfiles) Name() string {
	return StepNameCreateInstanceProfiles
}

func (StepCreateInstanceProfiles) Description() string {
	return "Create EC2 Instance master/node profiles"
}

func (StepCreateInstanceProfiles) Depends() []string {
	return nil
}

func ensureIAMProfile(ctx context.Context, iamS iamiface.IAMAPI, prefix, role string) (string, error) {
	var err error
	name := buildIAMName(prefix, role)

	if err = createIAMRolePolicy(ctx, iamS, name, policyFor(role)); err != nil {
		return "", errors.Wrapf(err, "ensure %s policy exists", name)
	}
	if err = createIAMRole(ctx, iamS, name, assumePolicy); err != nil {
		return "", errors.Wrapf(err, "ensure %s role exists", name)
	}
	if err = createIAMInstanceProfile(ctx, iamS, name); err != nil {
		return "", errors.Wrapf(err, "ensure %s instance profile exists", name)
	}

	return name, nil
}

func createIAMInstanceProfile(ctx context.Context, iamS iamiface.IAMAPI, name string) error {
	getInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	var instanceProfile *iam.InstanceProfile

	resp, err := iamS.GetInstanceProfileWithContext(ctx, getInput)
	if err != nil {
		if !isNotFoundErr(err) {
			return err
		}

		input := &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(name),
			Path:                aws.String("/"),
		}
		createResp, err := iamS.CreateInstanceProfileWithContext(ctx, input)
		if err != nil {
			return err
		}
		if createResp == nil || createResp.InstanceProfile == nil {
			return errors.Wrap(ErrEmptyResponse, "create instance profile")
		}
		instanceProfile = createResp.InstanceProfile

	} else {
		if resp == nil || resp.InstanceProfile == nil {
			return errors.Wrap(ErrEmptyResponse, "get instance profile")
		}
		instanceProfile = resp.InstanceProfile
	}

	if len(instanceProfile.Roles) == 0 {
		addInput := &iam.AddRoleToInstanceProfileInput{
			RoleName:            aws.String(name),
			InstanceProfileName: aws.String(name),
		}
		if _, err = iamS.AddRoleToInstanceProfileWithContext(ctx, addInput); err != nil {
			return err
		}
	}

	return nil
}

func createIAMRole(ctx context.Context, iamS iamiface.IAMAPI, name string, policy string) error {
	getInput := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}
	_, err := iamS.GetRoleWithContext(ctx, getInput)
	if err == nil {
		return nil
	}
	if !isNotFoundErr(err) {
		return err
	}
	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		Path:                     aws.String("/"),
		AssumeRolePolicyDocument: aws.String(policy),
	}
	_, err = iamS.CreateRoleWithContext(ctx, input)
	return err
}

func createIAMRolePolicy(ctx context.Context, iamS iamiface.IAMAPI, name string, policy string) error {
	getInput := &iam.GetRolePolicyInput{
		RoleName:   aws.String(name),
		PolicyName: aws.String(name),
	}
	_, err := iamS.GetRolePolicyWithContext(ctx, getInput)
	if err == nil {
		return nil
	}
	if !isNotFoundErr(err) {
		return err
	}
	putRoleInput := &iam.PutRolePolicyInput{
		RoleName:       aws.String(name),
		PolicyName:     aws.String(name),
		PolicyDocument: aws.String(policy),
	}
	_, err = iamS.PutRolePolicyWithContext(ctx, putRoleInput)
	return err
}

func isNotFoundErr(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == iam.ErrCodeNoSuchEntityException {
			return true
		}
	}
	return false
}

func policyFor(role string) string {
	if role == roleMaster {
		return masterIAMPolicy
	}
	return nodeIAMPolicy
}

func buildIAMName(prefix, role string) string {
	// TODO: use cluster specific names after adding roles removal.
	return strings.Join([]string{"kubernetes", role}, "-")
}
