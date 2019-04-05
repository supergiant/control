package azure

import "github.com/supergiant/control/pkg/workflows/steps"

func Init() {
	steps.RegisterStep(GetAuthorizerStepName, NewGetAuthorizerStepStep())
	steps.RegisterStep(CreateGroupStepName, NewCreateGroupStep())
	steps.RegisterStep(CreateLBStepName, NewCreateLBStep(NewSDK()))
	steps.RegisterStep(CreateVNetAndSubnetsStepName, NewCreateVirtualNetworkStep())
	steps.RegisterStep(CreateSecurityGroupStepName, NewCreateSecurityGroupStep())
	steps.RegisterStep(CreateVMStepName, NewCreateVMStep(NewSDK()))
	steps.RegisterStep(DeleteVMStepName, NewDeleteVMStep(NewSDK()))
	steps.RegisterStep(DeleteClusterStepName, NewDeleteClusterStep(NewSDK()))
}
