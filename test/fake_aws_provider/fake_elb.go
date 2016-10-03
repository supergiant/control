package fake_aws_provider

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elb"
)

type ELB struct {
	AddTagsRequestFn                                 func(*elb.AddTagsInput) (*request.Request, *elb.AddTagsOutput)
	AddTagsFn                                        func(*elb.AddTagsInput) (*elb.AddTagsOutput, error)
	ApplySecurityGroupsToLoadBalancerRequestFn       func(*elb.ApplySecurityGroupsToLoadBalancerInput) (*request.Request, *elb.ApplySecurityGroupsToLoadBalancerOutput)
	ApplySecurityGroupsToLoadBalancerFn              func(*elb.ApplySecurityGroupsToLoadBalancerInput) (*elb.ApplySecurityGroupsToLoadBalancerOutput, error)
	AttachLoadBalancerToSubnetsRequestFn             func(*elb.AttachLoadBalancerToSubnetsInput) (*request.Request, *elb.AttachLoadBalancerToSubnetsOutput)
	AttachLoadBalancerToSubnetsFn                    func(*elb.AttachLoadBalancerToSubnetsInput) (*elb.AttachLoadBalancerToSubnetsOutput, error)
	ConfigureHealthCheckRequestFn                    func(*elb.ConfigureHealthCheckInput) (*request.Request, *elb.ConfigureHealthCheckOutput)
	ConfigureHealthCheckFn                           func(*elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error)
	CreateAppCookieStickinessPolicyRequestFn         func(*elb.CreateAppCookieStickinessPolicyInput) (*request.Request, *elb.CreateAppCookieStickinessPolicyOutput)
	CreateAppCookieStickinessPolicyFn                func(*elb.CreateAppCookieStickinessPolicyInput) (*elb.CreateAppCookieStickinessPolicyOutput, error)
	CreateLBCookieStickinessPolicyRequestFn          func(*elb.CreateLBCookieStickinessPolicyInput) (*request.Request, *elb.CreateLBCookieStickinessPolicyOutput)
	CreateLBCookieStickinessPolicyFn                 func(*elb.CreateLBCookieStickinessPolicyInput) (*elb.CreateLBCookieStickinessPolicyOutput, error)
	CreateLoadBalancerRequestFn                      func(*elb.CreateLoadBalancerInput) (*request.Request, *elb.CreateLoadBalancerOutput)
	CreateLoadBalancerFn                             func(*elb.CreateLoadBalancerInput) (*elb.CreateLoadBalancerOutput, error)
	CreateLoadBalancerListenersRequestFn             func(*elb.CreateLoadBalancerListenersInput) (*request.Request, *elb.CreateLoadBalancerListenersOutput)
	CreateLoadBalancerListenersFn                    func(*elb.CreateLoadBalancerListenersInput) (*elb.CreateLoadBalancerListenersOutput, error)
	CreateLoadBalancerPolicyRequestFn                func(*elb.CreateLoadBalancerPolicyInput) (*request.Request, *elb.CreateLoadBalancerPolicyOutput)
	CreateLoadBalancerPolicyFn                       func(*elb.CreateLoadBalancerPolicyInput) (*elb.CreateLoadBalancerPolicyOutput, error)
	DeleteLoadBalancerRequestFn                      func(*elb.DeleteLoadBalancerInput) (*request.Request, *elb.DeleteLoadBalancerOutput)
	DeleteLoadBalancerFn                             func(*elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error)
	DeleteLoadBalancerListenersRequestFn             func(*elb.DeleteLoadBalancerListenersInput) (*request.Request, *elb.DeleteLoadBalancerListenersOutput)
	DeleteLoadBalancerListenersFn                    func(*elb.DeleteLoadBalancerListenersInput) (*elb.DeleteLoadBalancerListenersOutput, error)
	DeleteLoadBalancerPolicyRequestFn                func(*elb.DeleteLoadBalancerPolicyInput) (*request.Request, *elb.DeleteLoadBalancerPolicyOutput)
	DeleteLoadBalancerPolicyFn                       func(*elb.DeleteLoadBalancerPolicyInput) (*elb.DeleteLoadBalancerPolicyOutput, error)
	DeregisterInstancesFromLoadBalancerRequestFn     func(*elb.DeregisterInstancesFromLoadBalancerInput) (*request.Request, *elb.DeregisterInstancesFromLoadBalancerOutput)
	DeregisterInstancesFromLoadBalancerFn            func(*elb.DeregisterInstancesFromLoadBalancerInput) (*elb.DeregisterInstancesFromLoadBalancerOutput, error)
	DescribeInstanceHealthRequestFn                  func(*elb.DescribeInstanceHealthInput) (*request.Request, *elb.DescribeInstanceHealthOutput)
	DescribeInstanceHealthFn                         func(*elb.DescribeInstanceHealthInput) (*elb.DescribeInstanceHealthOutput, error)
	DescribeLoadBalancerAttributesRequestFn          func(*elb.DescribeLoadBalancerAttributesInput) (*request.Request, *elb.DescribeLoadBalancerAttributesOutput)
	DescribeLoadBalancerAttributesFn                 func(*elb.DescribeLoadBalancerAttributesInput) (*elb.DescribeLoadBalancerAttributesOutput, error)
	DescribeLoadBalancerPoliciesRequestFn            func(*elb.DescribeLoadBalancerPoliciesInput) (*request.Request, *elb.DescribeLoadBalancerPoliciesOutput)
	DescribeLoadBalancerPoliciesFn                   func(*elb.DescribeLoadBalancerPoliciesInput) (*elb.DescribeLoadBalancerPoliciesOutput, error)
	DescribeLoadBalancerPolicyTypesRequestFn         func(*elb.DescribeLoadBalancerPolicyTypesInput) (*request.Request, *elb.DescribeLoadBalancerPolicyTypesOutput)
	DescribeLoadBalancerPolicyTypesFn                func(*elb.DescribeLoadBalancerPolicyTypesInput) (*elb.DescribeLoadBalancerPolicyTypesOutput, error)
	DescribeLoadBalancersRequestFn                   func(*elb.DescribeLoadBalancersInput) (*request.Request, *elb.DescribeLoadBalancersOutput)
	DescribeLoadBalancersFn                          func(*elb.DescribeLoadBalancersInput) (*elb.DescribeLoadBalancersOutput, error)
	DescribeLoadBalancersPagesFn                     func(*elb.DescribeLoadBalancersInput, func(*elb.DescribeLoadBalancersOutput, bool) bool) error
	DescribeTagsRequestFn                            func(*elb.DescribeTagsInput) (*request.Request, *elb.DescribeTagsOutput)
	DescribeTagsFn                                   func(*elb.DescribeTagsInput) (*elb.DescribeTagsOutput, error)
	DetachLoadBalancerFromSubnetsRequestFn           func(*elb.DetachLoadBalancerFromSubnetsInput) (*request.Request, *elb.DetachLoadBalancerFromSubnetsOutput)
	DetachLoadBalancerFromSubnetsFn                  func(*elb.DetachLoadBalancerFromSubnetsInput) (*elb.DetachLoadBalancerFromSubnetsOutput, error)
	DisableAvailabilityZonesForLoadBalancerRequestFn func(*elb.DisableAvailabilityZonesForLoadBalancerInput) (*request.Request, *elb.DisableAvailabilityZonesForLoadBalancerOutput)
	DisableAvailabilityZonesForLoadBalancerFn        func(*elb.DisableAvailabilityZonesForLoadBalancerInput) (*elb.DisableAvailabilityZonesForLoadBalancerOutput, error)
	EnableAvailabilityZonesForLoadBalancerRequestFn  func(*elb.EnableAvailabilityZonesForLoadBalancerInput) (*request.Request, *elb.EnableAvailabilityZonesForLoadBalancerOutput)
	EnableAvailabilityZonesForLoadBalancerFn         func(*elb.EnableAvailabilityZonesForLoadBalancerInput) (*elb.EnableAvailabilityZonesForLoadBalancerOutput, error)
	ModifyLoadBalancerAttributesRequestFn            func(*elb.ModifyLoadBalancerAttributesInput) (*request.Request, *elb.ModifyLoadBalancerAttributesOutput)
	ModifyLoadBalancerAttributesFn                   func(*elb.ModifyLoadBalancerAttributesInput) (*elb.ModifyLoadBalancerAttributesOutput, error)
	RegisterInstancesWithLoadBalancerRequestFn       func(*elb.RegisterInstancesWithLoadBalancerInput) (*request.Request, *elb.RegisterInstancesWithLoadBalancerOutput)
	RegisterInstancesWithLoadBalancerFn              func(*elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error)
	RemoveTagsRequestFn                              func(*elb.RemoveTagsInput) (*request.Request, *elb.RemoveTagsOutput)
	RemoveTagsFn                                     func(*elb.RemoveTagsInput) (*elb.RemoveTagsOutput, error)
	SetLoadBalancerListenerSSLCertificateRequestFn   func(*elb.SetLoadBalancerListenerSSLCertificateInput) (*request.Request, *elb.SetLoadBalancerListenerSSLCertificateOutput)
	SetLoadBalancerListenerSSLCertificateFn          func(*elb.SetLoadBalancerListenerSSLCertificateInput) (*elb.SetLoadBalancerListenerSSLCertificateOutput, error)
	SetLoadBalancerPoliciesForBackendServerRequestFn func(*elb.SetLoadBalancerPoliciesForBackendServerInput) (*request.Request, *elb.SetLoadBalancerPoliciesForBackendServerOutput)
	SetLoadBalancerPoliciesForBackendServerFn        func(*elb.SetLoadBalancerPoliciesForBackendServerInput) (*elb.SetLoadBalancerPoliciesForBackendServerOutput, error)
	SetLoadBalancerPoliciesOfListenerRequestFn       func(*elb.SetLoadBalancerPoliciesOfListenerInput) (*request.Request, *elb.SetLoadBalancerPoliciesOfListenerOutput)
	SetLoadBalancerPoliciesOfListenerFn              func(*elb.SetLoadBalancerPoliciesOfListenerInput) (*elb.SetLoadBalancerPoliciesOfListenerOutput, error)
}

func (f *ELB) AddTagsRequest(input *elb.AddTagsInput) (*request.Request, *elb.AddTagsOutput) {
	if f.AddTagsRequestFn == nil {
		return nil, nil
	}
	return f.AddTagsRequestFn(input)
}

func (f *ELB) AddTags(input *elb.AddTagsInput) (*elb.AddTagsOutput, error) {
	if f.AddTagsFn == nil {
		return nil, nil
	}
	return f.AddTagsFn(input)
}

func (f *ELB) ApplySecurityGroupsToLoadBalancerRequest(input *elb.ApplySecurityGroupsToLoadBalancerInput) (*request.Request, *elb.ApplySecurityGroupsToLoadBalancerOutput) {
	if f.ApplySecurityGroupsToLoadBalancerRequestFn == nil {
		return nil, nil
	}
	return f.ApplySecurityGroupsToLoadBalancerRequestFn(input)
}

func (f *ELB) ApplySecurityGroupsToLoadBalancer(input *elb.ApplySecurityGroupsToLoadBalancerInput) (*elb.ApplySecurityGroupsToLoadBalancerOutput, error) {
	if f.ApplySecurityGroupsToLoadBalancerFn == nil {
		return nil, nil
	}
	return f.ApplySecurityGroupsToLoadBalancerFn(input)
}

func (f *ELB) AttachLoadBalancerToSubnetsRequest(input *elb.AttachLoadBalancerToSubnetsInput) (*request.Request, *elb.AttachLoadBalancerToSubnetsOutput) {
	if f.AttachLoadBalancerToSubnetsRequestFn == nil {
		return nil, nil
	}
	return f.AttachLoadBalancerToSubnetsRequestFn(input)
}

func (f *ELB) AttachLoadBalancerToSubnets(input *elb.AttachLoadBalancerToSubnetsInput) (*elb.AttachLoadBalancerToSubnetsOutput, error) {
	if f.AttachLoadBalancerToSubnetsFn == nil {
		return nil, nil
	}
	return f.AttachLoadBalancerToSubnetsFn(input)
}

func (f *ELB) ConfigureHealthCheckRequest(input *elb.ConfigureHealthCheckInput) (*request.Request, *elb.ConfigureHealthCheckOutput) {
	if f.ConfigureHealthCheckRequestFn == nil {
		return nil, nil
	}
	return f.ConfigureHealthCheckRequestFn(input)
}

func (f *ELB) ConfigureHealthCheck(input *elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error) {
	if f.ConfigureHealthCheckFn == nil {
		return nil, nil
	}
	return f.ConfigureHealthCheckFn(input)
}

func (f *ELB) CreateAppCookieStickinessPolicyRequest(input *elb.CreateAppCookieStickinessPolicyInput) (*request.Request, *elb.CreateAppCookieStickinessPolicyOutput) {
	if f.CreateAppCookieStickinessPolicyRequestFn == nil {
		return nil, nil
	}
	return f.CreateAppCookieStickinessPolicyRequestFn(input)
}

func (f *ELB) CreateAppCookieStickinessPolicy(input *elb.CreateAppCookieStickinessPolicyInput) (*elb.CreateAppCookieStickinessPolicyOutput, error) {
	if f.CreateAppCookieStickinessPolicyFn == nil {
		return nil, nil
	}
	return f.CreateAppCookieStickinessPolicyFn(input)
}

func (f *ELB) CreateLBCookieStickinessPolicyRequest(input *elb.CreateLBCookieStickinessPolicyInput) (*request.Request, *elb.CreateLBCookieStickinessPolicyOutput) {
	if f.CreateLBCookieStickinessPolicyRequestFn == nil {
		return nil, nil
	}
	return f.CreateLBCookieStickinessPolicyRequestFn(input)
}

func (f *ELB) CreateLBCookieStickinessPolicy(input *elb.CreateLBCookieStickinessPolicyInput) (*elb.CreateLBCookieStickinessPolicyOutput, error) {
	if f.CreateLBCookieStickinessPolicyFn == nil {
		return nil, nil
	}
	return f.CreateLBCookieStickinessPolicyFn(input)
}

func (f *ELB) CreateLoadBalancerRequest(input *elb.CreateLoadBalancerInput) (*request.Request, *elb.CreateLoadBalancerOutput) {
	if f.CreateLoadBalancerRequestFn == nil {
		return nil, nil
	}
	return f.CreateLoadBalancerRequestFn(input)
}

func (f *ELB) CreateLoadBalancer(input *elb.CreateLoadBalancerInput) (*elb.CreateLoadBalancerOutput, error) {
	if f.CreateLoadBalancerFn == nil {
		return nil, nil
	}
	return f.CreateLoadBalancerFn(input)
}

func (f *ELB) CreateLoadBalancerListenersRequest(input *elb.CreateLoadBalancerListenersInput) (*request.Request, *elb.CreateLoadBalancerListenersOutput) {
	if f.CreateLoadBalancerListenersRequestFn == nil {
		return nil, nil
	}
	return f.CreateLoadBalancerListenersRequestFn(input)
}

func (f *ELB) CreateLoadBalancerListeners(input *elb.CreateLoadBalancerListenersInput) (*elb.CreateLoadBalancerListenersOutput, error) {
	if f.CreateLoadBalancerListenersFn == nil {
		return nil, nil
	}
	return f.CreateLoadBalancerListenersFn(input)
}

func (f *ELB) CreateLoadBalancerPolicyRequest(input *elb.CreateLoadBalancerPolicyInput) (*request.Request, *elb.CreateLoadBalancerPolicyOutput) {
	if f.CreateLoadBalancerPolicyRequestFn == nil {
		return nil, nil
	}
	return f.CreateLoadBalancerPolicyRequestFn(input)
}

func (f *ELB) CreateLoadBalancerPolicy(input *elb.CreateLoadBalancerPolicyInput) (*elb.CreateLoadBalancerPolicyOutput, error) {
	if f.CreateLoadBalancerPolicyFn == nil {
		return nil, nil
	}
	return f.CreateLoadBalancerPolicyFn(input)
}

func (f *ELB) DeleteLoadBalancerRequest(input *elb.DeleteLoadBalancerInput) (*request.Request, *elb.DeleteLoadBalancerOutput) {
	if f.DeleteLoadBalancerRequestFn == nil {
		return nil, nil
	}
	return f.DeleteLoadBalancerRequestFn(input)
}

func (f *ELB) DeleteLoadBalancer(input *elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error) {
	if f.DeleteLoadBalancerFn == nil {
		return nil, nil
	}
	return f.DeleteLoadBalancerFn(input)
}

func (f *ELB) DeleteLoadBalancerListenersRequest(input *elb.DeleteLoadBalancerListenersInput) (*request.Request, *elb.DeleteLoadBalancerListenersOutput) {
	if f.DeleteLoadBalancerListenersRequestFn == nil {
		return nil, nil
	}
	return f.DeleteLoadBalancerListenersRequestFn(input)
}

func (f *ELB) DeleteLoadBalancerListeners(input *elb.DeleteLoadBalancerListenersInput) (*elb.DeleteLoadBalancerListenersOutput, error) {
	if f.DeleteLoadBalancerListenersFn == nil {
		return nil, nil
	}
	return f.DeleteLoadBalancerListenersFn(input)
}

func (f *ELB) DeleteLoadBalancerPolicyRequest(input *elb.DeleteLoadBalancerPolicyInput) (*request.Request, *elb.DeleteLoadBalancerPolicyOutput) {
	if f.DeleteLoadBalancerPolicyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteLoadBalancerPolicyRequestFn(input)
}

func (f *ELB) DeleteLoadBalancerPolicy(input *elb.DeleteLoadBalancerPolicyInput) (*elb.DeleteLoadBalancerPolicyOutput, error) {
	if f.DeleteLoadBalancerPolicyFn == nil {
		return nil, nil
	}
	return f.DeleteLoadBalancerPolicyFn(input)
}

func (f *ELB) DeregisterInstancesFromLoadBalancerRequest(input *elb.DeregisterInstancesFromLoadBalancerInput) (*request.Request, *elb.DeregisterInstancesFromLoadBalancerOutput) {
	if f.DeregisterInstancesFromLoadBalancerRequestFn == nil {
		return nil, nil
	}
	return f.DeregisterInstancesFromLoadBalancerRequestFn(input)
}

func (f *ELB) DeregisterInstancesFromLoadBalancer(input *elb.DeregisterInstancesFromLoadBalancerInput) (*elb.DeregisterInstancesFromLoadBalancerOutput, error) {
	if f.DeregisterInstancesFromLoadBalancerFn == nil {
		return nil, nil
	}
	return f.DeregisterInstancesFromLoadBalancerFn(input)
}

func (f *ELB) DescribeInstanceHealthRequest(input *elb.DescribeInstanceHealthInput) (*request.Request, *elb.DescribeInstanceHealthOutput) {
	if f.DescribeInstanceHealthRequestFn == nil {
		return nil, nil
	}
	return f.DescribeInstanceHealthRequestFn(input)
}

func (f *ELB) DescribeInstanceHealth(input *elb.DescribeInstanceHealthInput) (*elb.DescribeInstanceHealthOutput, error) {
	if f.DescribeInstanceHealthFn == nil {
		return nil, nil
	}
	return f.DescribeInstanceHealthFn(input)
}

func (f *ELB) DescribeLoadBalancerAttributesRequest(input *elb.DescribeLoadBalancerAttributesInput) (*request.Request, *elb.DescribeLoadBalancerAttributesOutput) {
	if f.DescribeLoadBalancerAttributesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancerAttributesRequestFn(input)
}

func (f *ELB) DescribeLoadBalancerAttributes(input *elb.DescribeLoadBalancerAttributesInput) (*elb.DescribeLoadBalancerAttributesOutput, error) {
	if f.DescribeLoadBalancerAttributesFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancerAttributesFn(input)
}

func (f *ELB) DescribeLoadBalancerPoliciesRequest(input *elb.DescribeLoadBalancerPoliciesInput) (*request.Request, *elb.DescribeLoadBalancerPoliciesOutput) {
	if f.DescribeLoadBalancerPoliciesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancerPoliciesRequestFn(input)
}

func (f *ELB) DescribeLoadBalancerPolicies(input *elb.DescribeLoadBalancerPoliciesInput) (*elb.DescribeLoadBalancerPoliciesOutput, error) {
	if f.DescribeLoadBalancerPoliciesFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancerPoliciesFn(input)
}

func (f *ELB) DescribeLoadBalancerPolicyTypesRequest(input *elb.DescribeLoadBalancerPolicyTypesInput) (*request.Request, *elb.DescribeLoadBalancerPolicyTypesOutput) {
	if f.DescribeLoadBalancerPolicyTypesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancerPolicyTypesRequestFn(input)
}

func (f *ELB) DescribeLoadBalancerPolicyTypes(input *elb.DescribeLoadBalancerPolicyTypesInput) (*elb.DescribeLoadBalancerPolicyTypesOutput, error) {
	if f.DescribeLoadBalancerPolicyTypesFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancerPolicyTypesFn(input)
}

func (f *ELB) DescribeLoadBalancersRequest(input *elb.DescribeLoadBalancersInput) (*request.Request, *elb.DescribeLoadBalancersOutput) {
	if f.DescribeLoadBalancersRequestFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancersRequestFn(input)
}

func (f *ELB) DescribeLoadBalancers(input *elb.DescribeLoadBalancersInput) (*elb.DescribeLoadBalancersOutput, error) {
	if f.DescribeLoadBalancersFn == nil {
		return nil, nil
	}
	return f.DescribeLoadBalancersFn(input)
}

func (f *ELB) DescribeLoadBalancersPages(input *elb.DescribeLoadBalancersInput, fn func(*elb.DescribeLoadBalancersOutput, bool) bool) error {
	if f.DescribeLoadBalancersPagesFn == nil {
		return nil
	}
	return f.DescribeLoadBalancersPagesFn(input, fn)
}

func (f *ELB) DescribeTagsRequest(input *elb.DescribeTagsInput) (*request.Request, *elb.DescribeTagsOutput) {
	if f.DescribeTagsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeTagsRequestFn(input)
}

func (f *ELB) DescribeTags(input *elb.DescribeTagsInput) (*elb.DescribeTagsOutput, error) {
	if f.DescribeTagsFn == nil {
		return nil, nil
	}
	return f.DescribeTagsFn(input)
}

func (f *ELB) DetachLoadBalancerFromSubnetsRequest(input *elb.DetachLoadBalancerFromSubnetsInput) (*request.Request, *elb.DetachLoadBalancerFromSubnetsOutput) {
	if f.DetachLoadBalancerFromSubnetsRequestFn == nil {
		return nil, nil
	}
	return f.DetachLoadBalancerFromSubnetsRequestFn(input)
}

func (f *ELB) DetachLoadBalancerFromSubnets(input *elb.DetachLoadBalancerFromSubnetsInput) (*elb.DetachLoadBalancerFromSubnetsOutput, error) {
	if f.DetachLoadBalancerFromSubnetsFn == nil {
		return nil, nil
	}
	return f.DetachLoadBalancerFromSubnetsFn(input)
}

func (f *ELB) DisableAvailabilityZonesForLoadBalancerRequest(input *elb.DisableAvailabilityZonesForLoadBalancerInput) (*request.Request, *elb.DisableAvailabilityZonesForLoadBalancerOutput) {
	if f.DisableAvailabilityZonesForLoadBalancerRequestFn == nil {
		return nil, nil
	}
	return f.DisableAvailabilityZonesForLoadBalancerRequestFn(input)
}

func (f *ELB) DisableAvailabilityZonesForLoadBalancer(input *elb.DisableAvailabilityZonesForLoadBalancerInput) (*elb.DisableAvailabilityZonesForLoadBalancerOutput, error) {
	if f.DisableAvailabilityZonesForLoadBalancerFn == nil {
		return nil, nil
	}
	return f.DisableAvailabilityZonesForLoadBalancerFn(input)
}

func (f *ELB) EnableAvailabilityZonesForLoadBalancerRequest(input *elb.EnableAvailabilityZonesForLoadBalancerInput) (*request.Request, *elb.EnableAvailabilityZonesForLoadBalancerOutput) {
	if f.EnableAvailabilityZonesForLoadBalancerRequestFn == nil {
		return nil, nil
	}
	return f.EnableAvailabilityZonesForLoadBalancerRequestFn(input)
}

func (f *ELB) EnableAvailabilityZonesForLoadBalancer(input *elb.EnableAvailabilityZonesForLoadBalancerInput) (*elb.EnableAvailabilityZonesForLoadBalancerOutput, error) {
	if f.EnableAvailabilityZonesForLoadBalancerFn == nil {
		return nil, nil
	}
	return f.EnableAvailabilityZonesForLoadBalancerFn(input)
}

func (f *ELB) ModifyLoadBalancerAttributesRequest(input *elb.ModifyLoadBalancerAttributesInput) (*request.Request, *elb.ModifyLoadBalancerAttributesOutput) {
	if f.ModifyLoadBalancerAttributesRequestFn == nil {
		return nil, nil
	}
	return f.ModifyLoadBalancerAttributesRequestFn(input)
}

func (f *ELB) ModifyLoadBalancerAttributes(input *elb.ModifyLoadBalancerAttributesInput) (*elb.ModifyLoadBalancerAttributesOutput, error) {
	if f.ModifyLoadBalancerAttributesFn == nil {
		return nil, nil
	}
	return f.ModifyLoadBalancerAttributesFn(input)
}

func (f *ELB) RegisterInstancesWithLoadBalancerRequest(input *elb.RegisterInstancesWithLoadBalancerInput) (*request.Request, *elb.RegisterInstancesWithLoadBalancerOutput) {
	if f.RegisterInstancesWithLoadBalancerRequestFn == nil {
		return nil, nil
	}
	return f.RegisterInstancesWithLoadBalancerRequestFn(input)
}

func (f *ELB) RegisterInstancesWithLoadBalancer(input *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
	if f.RegisterInstancesWithLoadBalancerFn == nil {
		return nil, nil
	}
	return f.RegisterInstancesWithLoadBalancerFn(input)
}

func (f *ELB) RemoveTagsRequest(input *elb.RemoveTagsInput) (*request.Request, *elb.RemoveTagsOutput) {
	if f.RemoveTagsRequestFn == nil {
		return nil, nil
	}
	return f.RemoveTagsRequestFn(input)
}

func (f *ELB) RemoveTags(input *elb.RemoveTagsInput) (*elb.RemoveTagsOutput, error) {
	if f.RemoveTagsFn == nil {
		return nil, nil
	}
	return f.RemoveTagsFn(input)
}

func (f *ELB) SetLoadBalancerListenerSSLCertificateRequest(input *elb.SetLoadBalancerListenerSSLCertificateInput) (*request.Request, *elb.SetLoadBalancerListenerSSLCertificateOutput) {
	if f.SetLoadBalancerListenerSSLCertificateRequestFn == nil {
		return nil, nil
	}
	return f.SetLoadBalancerListenerSSLCertificateRequestFn(input)
}

func (f *ELB) SetLoadBalancerListenerSSLCertificate(input *elb.SetLoadBalancerListenerSSLCertificateInput) (*elb.SetLoadBalancerListenerSSLCertificateOutput, error) {
	if f.SetLoadBalancerListenerSSLCertificateFn == nil {
		return nil, nil
	}
	return f.SetLoadBalancerListenerSSLCertificateFn(input)
}

func (f *ELB) SetLoadBalancerPoliciesForBackendServerRequest(input *elb.SetLoadBalancerPoliciesForBackendServerInput) (*request.Request, *elb.SetLoadBalancerPoliciesForBackendServerOutput) {
	if f.SetLoadBalancerPoliciesForBackendServerRequestFn == nil {
		return nil, nil
	}
	return f.SetLoadBalancerPoliciesForBackendServerRequestFn(input)
}

func (f *ELB) SetLoadBalancerPoliciesForBackendServer(input *elb.SetLoadBalancerPoliciesForBackendServerInput) (*elb.SetLoadBalancerPoliciesForBackendServerOutput, error) {
	if f.SetLoadBalancerPoliciesForBackendServerFn == nil {
		return nil, nil
	}
	return f.SetLoadBalancerPoliciesForBackendServerFn(input)
}

func (f *ELB) SetLoadBalancerPoliciesOfListenerRequest(input *elb.SetLoadBalancerPoliciesOfListenerInput) (*request.Request, *elb.SetLoadBalancerPoliciesOfListenerOutput) {
	if f.SetLoadBalancerPoliciesOfListenerRequestFn == nil {
		return nil, nil
	}
	return f.SetLoadBalancerPoliciesOfListenerRequestFn(input)
}

func (f *ELB) SetLoadBalancerPoliciesOfListener(input *elb.SetLoadBalancerPoliciesOfListenerInput) (*elb.SetLoadBalancerPoliciesOfListenerOutput, error) {
	if f.SetLoadBalancerPoliciesOfListenerFn == nil {
		return nil, nil
	}
	return f.SetLoadBalancerPoliciesOfListenerFn(input)
}
