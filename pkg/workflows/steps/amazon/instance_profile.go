package amazon

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/pkg/errors"
)

const (
	roleMaster = "master"
	roleNode   = "node"

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

func ensureIAMProfile(iamS iamiface.IAMAPI, prefix string, isMaster bool) (string, error) {
	var err error
	role := toRole(isMaster)
	name := buildIAMName(prefix, role)

	if err = createIAMRolePolicy(iamS, name, policyFor(role)); err != nil {
		return "", errors.Wrapf(err, "ensure %s policy exists", name)
	}
	if err = createIAMRole(iamS, name, assumePolicy); err != nil {
		return "", errors.Wrapf(err, "ensure %s role exists", name)
	}
	if err = createIAMInstanceProfile(iamS, name); err != nil {
		return "", errors.Wrapf(err, "ensure %s instance profile exists", name)
	}

	return name, nil
}

func createIAMInstanceProfile(iamS iamiface.IAMAPI, name string) error {
	getInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	var instanceProfile *iam.InstanceProfile

	resp, err := iamS.GetInstanceProfile(getInput)
	if err != nil {
		if !isAlreadyExistErr(err) {
			return err
		}

		input := &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(name),
			Path:                aws.String("/"),
		}
		createResp, err := iamS.CreateInstanceProfile(input)
		if err != nil {
			return err
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

func createIAMRole(iamS iamiface.IAMAPI, name string, policy string) error {
	getInput := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}
	_, err := iamS.GetRole(getInput)
	if err == nil {
		return nil
	} else if !isAlreadyExistErr(err) {
		return err
	}
	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		Path:                     aws.String("/"),
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
	} else if !isAlreadyExistErr(err) {
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

func isAlreadyExistErr(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == iam.ErrCodeEntityAlreadyExistsException {
			return true
		}
	}
	return false
}

// TODO: use node.Role
func toRole(isMaster bool) string {
	if isMaster {
		return roleMaster
	}
	return roleNode
}

func policyFor(role string) string {
	if role == roleMaster {
		return masterIAMPolicy
	}
	return nodeIAMPolicy
}

func buildIAMName(prefix, role string) string {
	// TODO: use cluster specific names, add roles removal.
	return strings.Join([]string{"kubernetes", role}, "-")
}
