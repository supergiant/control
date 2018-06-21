package fake_aws_provider

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/efs/efsiface"
)

type EFS struct {
	efsiface.EFSAPI
	CreateFileSystemfn                         func(*efs.CreateFileSystemInput) (*efs.FileSystemDescription, error)
	CreateFileSystemRequestfn                  func(*efs.CreateFileSystemInput) (*request.Request, *efs.FileSystemDescription)
	CreateMountTargetfn                        func(*efs.CreateMountTargetInput) (*efs.MountTargetDescription, error)
	CreateMountTargetRequestfn                 func(*efs.CreateMountTargetInput) (*request.Request, *efs.MountTargetDescription)
	CreateTagsfn                               func(*efs.CreateTagsInput) (*efs.CreateTagsOutput, error)
	CreateTagsRequestfn                        func(*efs.CreateTagsInput) (*request.Request, *efs.CreateTagsOutput)
	DeleteFileSystemfn                         func(*efs.DeleteFileSystemInput) (*efs.DeleteFileSystemOutput, error)
	DeleteFileSystemRequestfn                  func(*efs.DeleteFileSystemInput) (*request.Request, *efs.DeleteFileSystemOutput)
	DeleteMountTargetfn                        func(*efs.DeleteMountTargetInput) (*efs.DeleteMountTargetOutput, error)
	DeleteMountTargetRequestfn                 func(*efs.DeleteMountTargetInput) (*request.Request, *efs.DeleteMountTargetOutput)
	DeleteTagsfn                               func(*efs.DeleteTagsInput) (*efs.DeleteTagsOutput, error)
	DeleteTagsRequestfn                        func(*efs.DeleteTagsInput) (*request.Request, *efs.DeleteTagsOutput)
	DescribeFileSystemsfn                      func(*efs.DescribeFileSystemsInput) (*efs.DescribeFileSystemsOutput, error)
	DescribeFileSystemsRequestfn               func(*efs.DescribeFileSystemsInput) (*request.Request, *efs.DescribeFileSystemsOutput)
	DescribeMountTargetSecurityGroupsfn        func(*efs.DescribeMountTargetSecurityGroupsInput) (*efs.DescribeMountTargetSecurityGroupsOutput, error)
	DescribeMountTargetSecurityGroupsRequestfn func(*efs.DescribeMountTargetSecurityGroupsInput) (*request.Request, *efs.DescribeMountTargetSecurityGroupsOutput)
	DescribeMountTargetsfn                     func(*efs.DescribeMountTargetsInput) (*efs.DescribeMountTargetsOutput, error)
	DescribeMountTargetsRequestfn              func(*efs.DescribeMountTargetsInput) (*request.Request, *efs.DescribeMountTargetsOutput)
	DescribeTagsfn                             func(*efs.DescribeTagsInput) (*efs.DescribeTagsOutput, error)
	DescribeTagsRequestfn                      func(*efs.DescribeTagsInput) (*request.Request, *efs.DescribeTagsOutput)
	ModifyMountTargetSecurityGroupsfn          func(*efs.ModifyMountTargetSecurityGroupsInput) (*efs.ModifyMountTargetSecurityGroupsOutput, error)
	ModifyMountTargetSecurityGroupsRequestfn   func(*efs.ModifyMountTargetSecurityGroupsInput) (*request.Request, *efs.ModifyMountTargetSecurityGroupsOutput)
}

func (f *EFS) CreateFileSystem(input *efs.CreateFileSystemInput) (*efs.FileSystemDescription, error) {
	if f.CreateFileSystemfn == nil {
		return nil, nil
	}
	return f.CreateFileSystemfn(input)
}
func (f *EFS) CreateFileSystemRequest(input *efs.CreateFileSystemInput) (*request.Request, *efs.FileSystemDescription) {
	if f.CreateFileSystemRequestfn == nil {
		return nil, nil
	}
	return f.CreateFileSystemRequestfn(input)
}

func (f *EFS) CreateMountTarget(input *efs.CreateMountTargetInput) (*efs.MountTargetDescription, error) {
	if f.CreateMountTargetfn == nil {
		return nil, nil
	}
	return f.CreateMountTargetfn(input)
}
func (f *EFS) CreateMountTargetRequest(input *efs.CreateMountTargetInput) (*request.Request, *efs.MountTargetDescription) {
	if f.CreateMountTargetRequestfn == nil {
		return nil, nil
	}
	return f.CreateMountTargetRequestfn(input)
}

func (f *EFS) CreateTags(input *efs.CreateTagsInput) (*efs.CreateTagsOutput, error) {
	if f.CreateTagsfn == nil {
		return nil, nil
	}
	return f.CreateTagsfn(input)
}
func (f *EFS) CreateTagsRequest(input *efs.CreateTagsInput) (*request.Request, *efs.CreateTagsOutput) {
	if f.CreateTagsRequestfn == nil {
		return nil, nil
	}
	return f.CreateTagsRequestfn(input)
}

func (f *EFS) DeleteFileSystem(input *efs.DeleteFileSystemInput) (*efs.DeleteFileSystemOutput, error) {
	if f.DeleteFileSystemfn == nil {
		return nil, nil
	}
	return f.DeleteFileSystemfn(input)
}
func (f *EFS) DeleteFileSystemRequest(input *efs.DeleteFileSystemInput) (*request.Request, *efs.DeleteFileSystemOutput) {
	if f.DeleteFileSystemRequestfn == nil {
		return nil, nil
	}
	return f.DeleteFileSystemRequestfn(input)
}

func (f *EFS) DeleteMountTarget(input *efs.DeleteMountTargetInput) (*efs.DeleteMountTargetOutput, error) {
	if f.DeleteMountTargetfn == nil {
		return nil, nil
	}
	return f.DeleteMountTargetfn(input)
}
func (f *EFS) DeleteMountTargetRequest(input *efs.DeleteMountTargetInput) (*request.Request, *efs.DeleteMountTargetOutput) {
	if f.DeleteMountTargetRequestfn == nil {
		return nil, nil
	}
	return f.DeleteMountTargetRequestfn(input)
}

func (f *EFS) DeleteTags(input *efs.DeleteTagsInput) (*efs.DeleteTagsOutput, error) {
	if f.DeleteTagsfn == nil {
		return nil, nil
	}
	return f.DeleteTagsfn(input)
}
func (f *EFS) DeleteTagsRequest(input *efs.DeleteTagsInput) (*request.Request, *efs.DeleteTagsOutput) {
	if f.DeleteTagsRequestfn == nil {
		return nil, nil
	}
	return f.DeleteTagsRequestfn(input)
}

func (f *EFS) DescribeFileSystems(input *efs.DescribeFileSystemsInput) (*efs.DescribeFileSystemsOutput, error) {
	if f.DescribeFileSystemsfn == nil {
		return nil, nil
	}
	return f.DescribeFileSystemsfn(input)
}
func (f *EFS) DescribeFileSystemsRequest(input *efs.DescribeFileSystemsInput) (*request.Request, *efs.DescribeFileSystemsOutput) {
	if f.DescribeFileSystemsRequestfn == nil {
		return nil, nil
	}
	return f.DescribeFileSystemsRequestfn(input)
}

func (f *EFS) DescribeMountTargetSecurityGroups(input *efs.DescribeMountTargetSecurityGroupsInput) (*efs.DescribeMountTargetSecurityGroupsOutput, error) {
	if f.DescribeMountTargetSecurityGroupsfn == nil {
		return nil, nil
	}
	return f.DescribeMountTargetSecurityGroupsfn(input)
}
func (f *EFS) DescribeMountTargetSecurityGroupsRequest(input *efs.DescribeMountTargetSecurityGroupsInput) (*request.Request, *efs.DescribeMountTargetSecurityGroupsOutput) {
	if f.DescribeMountTargetSecurityGroupsRequestfn == nil {
		return nil, nil
	}
	return f.DescribeMountTargetSecurityGroupsRequestfn(input)
}

func (f *EFS) DescribeMountTargets(input *efs.DescribeMountTargetsInput) (*efs.DescribeMountTargetsOutput, error) {
	if f.DescribeMountTargetsfn == nil {
		return nil, nil
	}
	return f.DescribeMountTargetsfn(input)
}
func (f *EFS) DescribeMountTargetsRequest(input *efs.DescribeMountTargetsInput) (*request.Request, *efs.DescribeMountTargetsOutput) {
	if f.DescribeMountTargetsRequestfn == nil {
		return nil, nil
	}
	return f.DescribeMountTargetsRequestfn(input)
}

func (f *EFS) DescribeTags(input *efs.DescribeTagsInput) (*efs.DescribeTagsOutput, error) {
	if f.DescribeTagsfn == nil {
		return nil, nil
	}
	return f.DescribeTagsfn(input)
}
func (f *EFS) DescribeTagsRequest(input *efs.DescribeTagsInput) (*request.Request, *efs.DescribeTagsOutput) {
	if f.DescribeTagsRequestfn == nil {
		return nil, nil
	}
	return f.DescribeTagsRequestfn(input)
}

func (f *EFS) ModifyMountTargetSecurityGroups(input *efs.ModifyMountTargetSecurityGroupsInput) (*efs.ModifyMountTargetSecurityGroupsOutput, error) {
	if f.ModifyMountTargetSecurityGroupsfn == nil {
		return nil, nil
	}
	return f.ModifyMountTargetSecurityGroupsfn(input)
}
func (f *EFS) ModifyMountTargetSecurityGroupsRequest(input *efs.ModifyMountTargetSecurityGroupsInput) (*request.Request, *efs.ModifyMountTargetSecurityGroupsOutput) {
	if f.ModifyMountTargetSecurityGroupsRequestfn == nil {
		return nil, nil
	}
	return f.ModifyMountTargetSecurityGroupsRequestfn(input)
}
