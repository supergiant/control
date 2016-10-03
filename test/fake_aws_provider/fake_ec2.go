package fake_aws_provider

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2 struct {
	AcceptVpcPeeringConnectionRequestFn             func(*ec2.AcceptVpcPeeringConnectionInput) (*request.Request, *ec2.AcceptVpcPeeringConnectionOutput)
	AcceptVpcPeeringConnectionFn                    func(*ec2.AcceptVpcPeeringConnectionInput) (*ec2.AcceptVpcPeeringConnectionOutput, error)
	AllocateAddressRequestFn                        func(*ec2.AllocateAddressInput) (*request.Request, *ec2.AllocateAddressOutput)
	AllocateAddressFn                               func(*ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error)
	AllocateHostsRequestFn                          func(*ec2.AllocateHostsInput) (*request.Request, *ec2.AllocateHostsOutput)
	AllocateHostsFn                                 func(*ec2.AllocateHostsInput) (*ec2.AllocateHostsOutput, error)
	AssignPrivateIpAddressesRequestFn               func(*ec2.AssignPrivateIpAddressesInput) (*request.Request, *ec2.AssignPrivateIpAddressesOutput)
	AssignPrivateIpAddressesFn                      func(*ec2.AssignPrivateIpAddressesInput) (*ec2.AssignPrivateIpAddressesOutput, error)
	AssociateAddressRequestFn                       func(*ec2.AssociateAddressInput) (*request.Request, *ec2.AssociateAddressOutput)
	AssociateAddressFn                              func(*ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error)
	AssociateDhcpOptionsRequestFn                   func(*ec2.AssociateDhcpOptionsInput) (*request.Request, *ec2.AssociateDhcpOptionsOutput)
	AssociateDhcpOptionsFn                          func(*ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error)
	AssociateRouteTableRequestFn                    func(*ec2.AssociateRouteTableInput) (*request.Request, *ec2.AssociateRouteTableOutput)
	AssociateRouteTableFn                           func(*ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error)
	AttachClassicLinkVpcRequestFn                   func(*ec2.AttachClassicLinkVpcInput) (*request.Request, *ec2.AttachClassicLinkVpcOutput)
	AttachClassicLinkVpcFn                          func(*ec2.AttachClassicLinkVpcInput) (*ec2.AttachClassicLinkVpcOutput, error)
	AttachInternetGatewayRequestFn                  func(*ec2.AttachInternetGatewayInput) (*request.Request, *ec2.AttachInternetGatewayOutput)
	AttachInternetGatewayFn                         func(*ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error)
	AttachNetworkInterfaceRequestFn                 func(*ec2.AttachNetworkInterfaceInput) (*request.Request, *ec2.AttachNetworkInterfaceOutput)
	AttachNetworkInterfaceFn                        func(*ec2.AttachNetworkInterfaceInput) (*ec2.AttachNetworkInterfaceOutput, error)
	AttachVolumeRequestFn                           func(*ec2.AttachVolumeInput) (*request.Request, *ec2.VolumeAttachment)
	AttachVolumeFn                                  func(*ec2.AttachVolumeInput) (*ec2.VolumeAttachment, error)
	AttachVpnGatewayRequestFn                       func(*ec2.AttachVpnGatewayInput) (*request.Request, *ec2.AttachVpnGatewayOutput)
	AttachVpnGatewayFn                              func(*ec2.AttachVpnGatewayInput) (*ec2.AttachVpnGatewayOutput, error)
	AuthorizeSecurityGroupEgressRequestFn           func(*ec2.AuthorizeSecurityGroupEgressInput) (*request.Request, *ec2.AuthorizeSecurityGroupEgressOutput)
	AuthorizeSecurityGroupEgressFn                  func(*ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error)
	AuthorizeSecurityGroupIngressRequestFn          func(*ec2.AuthorizeSecurityGroupIngressInput) (*request.Request, *ec2.AuthorizeSecurityGroupIngressOutput)
	AuthorizeSecurityGroupIngressFn                 func(*ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	BundleInstanceRequestFn                         func(*ec2.BundleInstanceInput) (*request.Request, *ec2.BundleInstanceOutput)
	BundleInstanceFn                                func(*ec2.BundleInstanceInput) (*ec2.BundleInstanceOutput, error)
	CancelBundleTaskRequestFn                       func(*ec2.CancelBundleTaskInput) (*request.Request, *ec2.CancelBundleTaskOutput)
	CancelBundleTaskFn                              func(*ec2.CancelBundleTaskInput) (*ec2.CancelBundleTaskOutput, error)
	CancelConversionTaskRequestFn                   func(*ec2.CancelConversionTaskInput) (*request.Request, *ec2.CancelConversionTaskOutput)
	CancelConversionTaskFn                          func(*ec2.CancelConversionTaskInput) (*ec2.CancelConversionTaskOutput, error)
	CancelExportTaskRequestFn                       func(*ec2.CancelExportTaskInput) (*request.Request, *ec2.CancelExportTaskOutput)
	CancelExportTaskFn                              func(*ec2.CancelExportTaskInput) (*ec2.CancelExportTaskOutput, error)
	CancelImportTaskRequestFn                       func(*ec2.CancelImportTaskInput) (*request.Request, *ec2.CancelImportTaskOutput)
	CancelImportTaskFn                              func(*ec2.CancelImportTaskInput) (*ec2.CancelImportTaskOutput, error)
	CancelReservedInstancesListingRequestFn         func(*ec2.CancelReservedInstancesListingInput) (*request.Request, *ec2.CancelReservedInstancesListingOutput)
	CancelReservedInstancesListingFn                func(*ec2.CancelReservedInstancesListingInput) (*ec2.CancelReservedInstancesListingOutput, error)
	CancelSpotFleetRequestsRequestFn                func(*ec2.CancelSpotFleetRequestsInput) (*request.Request, *ec2.CancelSpotFleetRequestsOutput)
	CancelSpotFleetRequestsFn                       func(*ec2.CancelSpotFleetRequestsInput) (*ec2.CancelSpotFleetRequestsOutput, error)
	CancelSpotInstanceRequestsRequestFn             func(*ec2.CancelSpotInstanceRequestsInput) (*request.Request, *ec2.CancelSpotInstanceRequestsOutput)
	CancelSpotInstanceRequestsFn                    func(*ec2.CancelSpotInstanceRequestsInput) (*ec2.CancelSpotInstanceRequestsOutput, error)
	ConfirmProductInstanceRequestFn                 func(*ec2.ConfirmProductInstanceInput) (*request.Request, *ec2.ConfirmProductInstanceOutput)
	ConfirmProductInstanceFn                        func(*ec2.ConfirmProductInstanceInput) (*ec2.ConfirmProductInstanceOutput, error)
	CopyImageRequestFn                              func(*ec2.CopyImageInput) (*request.Request, *ec2.CopyImageOutput)
	CopyImageFn                                     func(*ec2.CopyImageInput) (*ec2.CopyImageOutput, error)
	CopySnapshotRequestFn                           func(*ec2.CopySnapshotInput) (*request.Request, *ec2.CopySnapshotOutput)
	CopySnapshotFn                                  func(*ec2.CopySnapshotInput) (*ec2.CopySnapshotOutput, error)
	CreateCustomerGatewayRequestFn                  func(*ec2.CreateCustomerGatewayInput) (*request.Request, *ec2.CreateCustomerGatewayOutput)
	CreateCustomerGatewayFn                         func(*ec2.CreateCustomerGatewayInput) (*ec2.CreateCustomerGatewayOutput, error)
	CreateDhcpOptionsRequestFn                      func(*ec2.CreateDhcpOptionsInput) (*request.Request, *ec2.CreateDhcpOptionsOutput)
	CreateDhcpOptionsFn                             func(*ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error)
	CreateFlowLogsRequestFn                         func(*ec2.CreateFlowLogsInput) (*request.Request, *ec2.CreateFlowLogsOutput)
	CreateFlowLogsFn                                func(*ec2.CreateFlowLogsInput) (*ec2.CreateFlowLogsOutput, error)
	CreateImageRequestFn                            func(*ec2.CreateImageInput) (*request.Request, *ec2.CreateImageOutput)
	CreateImageFn                                   func(*ec2.CreateImageInput) (*ec2.CreateImageOutput, error)
	CreateInstanceExportTaskRequestFn               func(*ec2.CreateInstanceExportTaskInput) (*request.Request, *ec2.CreateInstanceExportTaskOutput)
	CreateInstanceExportTaskFn                      func(*ec2.CreateInstanceExportTaskInput) (*ec2.CreateInstanceExportTaskOutput, error)
	CreateInternetGatewayRequestFn                  func(*ec2.CreateInternetGatewayInput) (*request.Request, *ec2.CreateInternetGatewayOutput)
	CreateInternetGatewayFn                         func(*ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error)
	CreateKeyPairRequestFn                          func(*ec2.CreateKeyPairInput) (*request.Request, *ec2.CreateKeyPairOutput)
	CreateKeyPairFn                                 func(*ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error)
	CreateNatGatewayRequestFn                       func(*ec2.CreateNatGatewayInput) (*request.Request, *ec2.CreateNatGatewayOutput)
	CreateNatGatewayFn                              func(*ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error)
	CreateNetworkAclRequestFn                       func(*ec2.CreateNetworkAclInput) (*request.Request, *ec2.CreateNetworkAclOutput)
	CreateNetworkAclFn                              func(*ec2.CreateNetworkAclInput) (*ec2.CreateNetworkAclOutput, error)
	CreateNetworkAclEntryRequestFn                  func(*ec2.CreateNetworkAclEntryInput) (*request.Request, *ec2.CreateNetworkAclEntryOutput)
	CreateNetworkAclEntryFn                         func(*ec2.CreateNetworkAclEntryInput) (*ec2.CreateNetworkAclEntryOutput, error)
	CreateNetworkInterfaceRequestFn                 func(*ec2.CreateNetworkInterfaceInput) (*request.Request, *ec2.CreateNetworkInterfaceOutput)
	CreateNetworkInterfaceFn                        func(*ec2.CreateNetworkInterfaceInput) (*ec2.CreateNetworkInterfaceOutput, error)
	CreatePlacementGroupRequestFn                   func(*ec2.CreatePlacementGroupInput) (*request.Request, *ec2.CreatePlacementGroupOutput)
	CreatePlacementGroupFn                          func(*ec2.CreatePlacementGroupInput) (*ec2.CreatePlacementGroupOutput, error)
	CreateReservedInstancesListingRequestFn         func(*ec2.CreateReservedInstancesListingInput) (*request.Request, *ec2.CreateReservedInstancesListingOutput)
	CreateReservedInstancesListingFn                func(*ec2.CreateReservedInstancesListingInput) (*ec2.CreateReservedInstancesListingOutput, error)
	CreateRouteRequestFn                            func(*ec2.CreateRouteInput) (*request.Request, *ec2.CreateRouteOutput)
	CreateRouteFn                                   func(*ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error)
	CreateRouteTableRequestFn                       func(*ec2.CreateRouteTableInput) (*request.Request, *ec2.CreateRouteTableOutput)
	CreateRouteTableFn                              func(*ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error)
	CreateSecurityGroupRequestFn                    func(*ec2.CreateSecurityGroupInput) (*request.Request, *ec2.CreateSecurityGroupOutput)
	CreateSecurityGroupFn                           func(*ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error)
	CreateSnapshotRequestFn                         func(*ec2.CreateSnapshotInput) (*request.Request, *ec2.Snapshot)
	CreateSnapshotFn                                func(*ec2.CreateSnapshotInput) (*ec2.Snapshot, error)
	CreateSpotDatafeedSubscriptionRequestFn         func(*ec2.CreateSpotDatafeedSubscriptionInput) (*request.Request, *ec2.CreateSpotDatafeedSubscriptionOutput)
	CreateSpotDatafeedSubscriptionFn                func(*ec2.CreateSpotDatafeedSubscriptionInput) (*ec2.CreateSpotDatafeedSubscriptionOutput, error)
	CreateSubnetRequestFn                           func(*ec2.CreateSubnetInput) (*request.Request, *ec2.CreateSubnetOutput)
	CreateSubnetFn                                  func(*ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error)
	CreateTagsRequestFn                             func(*ec2.CreateTagsInput) (*request.Request, *ec2.CreateTagsOutput)
	CreateTagsFn                                    func(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)
	CreateVolumeRequestFn                           func(*ec2.CreateVolumeInput) (*request.Request, *ec2.Volume)
	CreateVolumeFn                                  func(*ec2.CreateVolumeInput) (*ec2.Volume, error)
	CreateVpcRequestFn                              func(*ec2.CreateVpcInput) (*request.Request, *ec2.CreateVpcOutput)
	CreateVpcFn                                     func(*ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error)
	CreateVpcEndpointRequestFn                      func(*ec2.CreateVpcEndpointInput) (*request.Request, *ec2.CreateVpcEndpointOutput)
	CreateVpcEndpointFn                             func(*ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error)
	CreateVpcPeeringConnectionRequestFn             func(*ec2.CreateVpcPeeringConnectionInput) (*request.Request, *ec2.CreateVpcPeeringConnectionOutput)
	CreateVpcPeeringConnectionFn                    func(*ec2.CreateVpcPeeringConnectionInput) (*ec2.CreateVpcPeeringConnectionOutput, error)
	CreateVpnConnectionRequestFn                    func(*ec2.CreateVpnConnectionInput) (*request.Request, *ec2.CreateVpnConnectionOutput)
	CreateVpnConnectionFn                           func(*ec2.CreateVpnConnectionInput) (*ec2.CreateVpnConnectionOutput, error)
	CreateVpnConnectionRouteRequestFn               func(*ec2.CreateVpnConnectionRouteInput) (*request.Request, *ec2.CreateVpnConnectionRouteOutput)
	CreateVpnConnectionRouteFn                      func(*ec2.CreateVpnConnectionRouteInput) (*ec2.CreateVpnConnectionRouteOutput, error)
	CreateVpnGatewayRequestFn                       func(*ec2.CreateVpnGatewayInput) (*request.Request, *ec2.CreateVpnGatewayOutput)
	CreateVpnGatewayFn                              func(*ec2.CreateVpnGatewayInput) (*ec2.CreateVpnGatewayOutput, error)
	DeleteCustomerGatewayRequestFn                  func(*ec2.DeleteCustomerGatewayInput) (*request.Request, *ec2.DeleteCustomerGatewayOutput)
	DeleteCustomerGatewayFn                         func(*ec2.DeleteCustomerGatewayInput) (*ec2.DeleteCustomerGatewayOutput, error)
	DeleteDhcpOptionsRequestFn                      func(*ec2.DeleteDhcpOptionsInput) (*request.Request, *ec2.DeleteDhcpOptionsOutput)
	DeleteDhcpOptionsFn                             func(*ec2.DeleteDhcpOptionsInput) (*ec2.DeleteDhcpOptionsOutput, error)
	DeleteFlowLogsRequestFn                         func(*ec2.DeleteFlowLogsInput) (*request.Request, *ec2.DeleteFlowLogsOutput)
	DeleteFlowLogsFn                                func(*ec2.DeleteFlowLogsInput) (*ec2.DeleteFlowLogsOutput, error)
	DeleteInternetGatewayRequestFn                  func(*ec2.DeleteInternetGatewayInput) (*request.Request, *ec2.DeleteInternetGatewayOutput)
	DeleteInternetGatewayFn                         func(*ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error)
	DeleteKeyPairRequestFn                          func(*ec2.DeleteKeyPairInput) (*request.Request, *ec2.DeleteKeyPairOutput)
	DeleteKeyPairFn                                 func(*ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error)
	DeleteNatGatewayRequestFn                       func(*ec2.DeleteNatGatewayInput) (*request.Request, *ec2.DeleteNatGatewayOutput)
	DeleteNatGatewayFn                              func(*ec2.DeleteNatGatewayInput) (*ec2.DeleteNatGatewayOutput, error)
	DeleteNetworkAclRequestFn                       func(*ec2.DeleteNetworkAclInput) (*request.Request, *ec2.DeleteNetworkAclOutput)
	DeleteNetworkAclFn                              func(*ec2.DeleteNetworkAclInput) (*ec2.DeleteNetworkAclOutput, error)
	DeleteNetworkAclEntryRequestFn                  func(*ec2.DeleteNetworkAclEntryInput) (*request.Request, *ec2.DeleteNetworkAclEntryOutput)
	DeleteNetworkAclEntryFn                         func(*ec2.DeleteNetworkAclEntryInput) (*ec2.DeleteNetworkAclEntryOutput, error)
	DeleteNetworkInterfaceRequestFn                 func(*ec2.DeleteNetworkInterfaceInput) (*request.Request, *ec2.DeleteNetworkInterfaceOutput)
	DeleteNetworkInterfaceFn                        func(*ec2.DeleteNetworkInterfaceInput) (*ec2.DeleteNetworkInterfaceOutput, error)
	DeletePlacementGroupRequestFn                   func(*ec2.DeletePlacementGroupInput) (*request.Request, *ec2.DeletePlacementGroupOutput)
	DeletePlacementGroupFn                          func(*ec2.DeletePlacementGroupInput) (*ec2.DeletePlacementGroupOutput, error)
	DeleteRouteRequestFn                            func(*ec2.DeleteRouteInput) (*request.Request, *ec2.DeleteRouteOutput)
	DeleteRouteFn                                   func(*ec2.DeleteRouteInput) (*ec2.DeleteRouteOutput, error)
	DeleteRouteTableRequestFn                       func(*ec2.DeleteRouteTableInput) (*request.Request, *ec2.DeleteRouteTableOutput)
	DeleteRouteTableFn                              func(*ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error)
	DeleteSecurityGroupRequestFn                    func(*ec2.DeleteSecurityGroupInput) (*request.Request, *ec2.DeleteSecurityGroupOutput)
	DeleteSecurityGroupFn                           func(*ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error)
	DeleteSnapshotRequestFn                         func(*ec2.DeleteSnapshotInput) (*request.Request, *ec2.DeleteSnapshotOutput)
	DeleteSnapshotFn                                func(*ec2.DeleteSnapshotInput) (*ec2.DeleteSnapshotOutput, error)
	DeleteSpotDatafeedSubscriptionRequestFn         func(*ec2.DeleteSpotDatafeedSubscriptionInput) (*request.Request, *ec2.DeleteSpotDatafeedSubscriptionOutput)
	DeleteSpotDatafeedSubscriptionFn                func(*ec2.DeleteSpotDatafeedSubscriptionInput) (*ec2.DeleteSpotDatafeedSubscriptionOutput, error)
	DeleteSubnetRequestFn                           func(*ec2.DeleteSubnetInput) (*request.Request, *ec2.DeleteSubnetOutput)
	DeleteSubnetFn                                  func(*ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error)
	DeleteTagsRequestFn                             func(*ec2.DeleteTagsInput) (*request.Request, *ec2.DeleteTagsOutput)
	DeleteTagsFn                                    func(*ec2.DeleteTagsInput) (*ec2.DeleteTagsOutput, error)
	DeleteVolumeRequestFn                           func(*ec2.DeleteVolumeInput) (*request.Request, *ec2.DeleteVolumeOutput)
	DeleteVolumeFn                                  func(*ec2.DeleteVolumeInput) (*ec2.DeleteVolumeOutput, error)
	DeleteVpcRequestFn                              func(*ec2.DeleteVpcInput) (*request.Request, *ec2.DeleteVpcOutput)
	DeleteVpcFn                                     func(*ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error)
	DeleteVpcEndpointsRequestFn                     func(*ec2.DeleteVpcEndpointsInput) (*request.Request, *ec2.DeleteVpcEndpointsOutput)
	DeleteVpcEndpointsFn                            func(*ec2.DeleteVpcEndpointsInput) (*ec2.DeleteVpcEndpointsOutput, error)
	DeleteVpcPeeringConnectionRequestFn             func(*ec2.DeleteVpcPeeringConnectionInput) (*request.Request, *ec2.DeleteVpcPeeringConnectionOutput)
	DeleteVpcPeeringConnectionFn                    func(*ec2.DeleteVpcPeeringConnectionInput) (*ec2.DeleteVpcPeeringConnectionOutput, error)
	DeleteVpnConnectionRequestFn                    func(*ec2.DeleteVpnConnectionInput) (*request.Request, *ec2.DeleteVpnConnectionOutput)
	DeleteVpnConnectionFn                           func(*ec2.DeleteVpnConnectionInput) (*ec2.DeleteVpnConnectionOutput, error)
	DeleteVpnConnectionRouteRequestFn               func(*ec2.DeleteVpnConnectionRouteInput) (*request.Request, *ec2.DeleteVpnConnectionRouteOutput)
	DeleteVpnConnectionRouteFn                      func(*ec2.DeleteVpnConnectionRouteInput) (*ec2.DeleteVpnConnectionRouteOutput, error)
	DeleteVpnGatewayRequestFn                       func(*ec2.DeleteVpnGatewayInput) (*request.Request, *ec2.DeleteVpnGatewayOutput)
	DeleteVpnGatewayFn                              func(*ec2.DeleteVpnGatewayInput) (*ec2.DeleteVpnGatewayOutput, error)
	DeregisterImageRequestFn                        func(*ec2.DeregisterImageInput) (*request.Request, *ec2.DeregisterImageOutput)
	DeregisterImageFn                               func(*ec2.DeregisterImageInput) (*ec2.DeregisterImageOutput, error)
	DescribeAccountAttributesRequestFn              func(*ec2.DescribeAccountAttributesInput) (*request.Request, *ec2.DescribeAccountAttributesOutput)
	DescribeAccountAttributesFn                     func(*ec2.DescribeAccountAttributesInput) (*ec2.DescribeAccountAttributesOutput, error)
	DescribeAddressesRequestFn                      func(*ec2.DescribeAddressesInput) (*request.Request, *ec2.DescribeAddressesOutput)
	DescribeAddressesFn                             func(*ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error)
	DescribeAvailabilityZonesRequestFn              func(*ec2.DescribeAvailabilityZonesInput) (*request.Request, *ec2.DescribeAvailabilityZonesOutput)
	DescribeAvailabilityZonesFn                     func(*ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error)
	DescribeBundleTasksRequestFn                    func(*ec2.DescribeBundleTasksInput) (*request.Request, *ec2.DescribeBundleTasksOutput)
	DescribeBundleTasksFn                           func(*ec2.DescribeBundleTasksInput) (*ec2.DescribeBundleTasksOutput, error)
	DescribeClassicLinkInstancesRequestFn           func(*ec2.DescribeClassicLinkInstancesInput) (*request.Request, *ec2.DescribeClassicLinkInstancesOutput)
	DescribeClassicLinkInstancesFn                  func(*ec2.DescribeClassicLinkInstancesInput) (*ec2.DescribeClassicLinkInstancesOutput, error)
	DescribeConversionTasksRequestFn                func(*ec2.DescribeConversionTasksInput) (*request.Request, *ec2.DescribeConversionTasksOutput)
	DescribeConversionTasksFn                       func(*ec2.DescribeConversionTasksInput) (*ec2.DescribeConversionTasksOutput, error)
	DescribeCustomerGatewaysRequestFn               func(*ec2.DescribeCustomerGatewaysInput) (*request.Request, *ec2.DescribeCustomerGatewaysOutput)
	DescribeCustomerGatewaysFn                      func(*ec2.DescribeCustomerGatewaysInput) (*ec2.DescribeCustomerGatewaysOutput, error)
	DescribeDhcpOptionsRequestFn                    func(*ec2.DescribeDhcpOptionsInput) (*request.Request, *ec2.DescribeDhcpOptionsOutput)
	DescribeDhcpOptionsFn                           func(*ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error)
	DescribeExportTasksRequestFn                    func(*ec2.DescribeExportTasksInput) (*request.Request, *ec2.DescribeExportTasksOutput)
	DescribeExportTasksFn                           func(*ec2.DescribeExportTasksInput) (*ec2.DescribeExportTasksOutput, error)
	DescribeFlowLogsRequestFn                       func(*ec2.DescribeFlowLogsInput) (*request.Request, *ec2.DescribeFlowLogsOutput)
	DescribeFlowLogsFn                              func(*ec2.DescribeFlowLogsInput) (*ec2.DescribeFlowLogsOutput, error)
	DescribeHostsRequestFn                          func(*ec2.DescribeHostsInput) (*request.Request, *ec2.DescribeHostsOutput)
	DescribeHostsFn                                 func(*ec2.DescribeHostsInput) (*ec2.DescribeHostsOutput, error)
	DescribeIdFormatRequestFn                       func(*ec2.DescribeIdFormatInput) (*request.Request, *ec2.DescribeIdFormatOutput)
	DescribeIdFormatFn                              func(*ec2.DescribeIdFormatInput) (*ec2.DescribeIdFormatOutput, error)
	DescribeImageAttributeRequestFn                 func(*ec2.DescribeImageAttributeInput) (*request.Request, *ec2.DescribeImageAttributeOutput)
	DescribeImageAttributeFn                        func(*ec2.DescribeImageAttributeInput) (*ec2.DescribeImageAttributeOutput, error)
	DescribeImagesRequestFn                         func(*ec2.DescribeImagesInput) (*request.Request, *ec2.DescribeImagesOutput)
	DescribeImagesFn                                func(*ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error)
	DescribeImportImageTasksRequestFn               func(*ec2.DescribeImportImageTasksInput) (*request.Request, *ec2.DescribeImportImageTasksOutput)
	DescribeImportImageTasksFn                      func(*ec2.DescribeImportImageTasksInput) (*ec2.DescribeImportImageTasksOutput, error)
	DescribeImportSnapshotTasksRequestFn            func(*ec2.DescribeImportSnapshotTasksInput) (*request.Request, *ec2.DescribeImportSnapshotTasksOutput)
	DescribeImportSnapshotTasksFn                   func(*ec2.DescribeImportSnapshotTasksInput) (*ec2.DescribeImportSnapshotTasksOutput, error)
	DescribeInstanceAttributeRequestFn              func(*ec2.DescribeInstanceAttributeInput) (*request.Request, *ec2.DescribeInstanceAttributeOutput)
	DescribeInstanceAttributeFn                     func(*ec2.DescribeInstanceAttributeInput) (*ec2.DescribeInstanceAttributeOutput, error)
	DescribeInstanceStatusRequestFn                 func(*ec2.DescribeInstanceStatusInput) (*request.Request, *ec2.DescribeInstanceStatusOutput)
	DescribeInstanceStatusFn                        func(*ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error)
	DescribeInstanceStatusPagesFn                   func(*ec2.DescribeInstanceStatusInput, func(*ec2.DescribeInstanceStatusOutput, bool) bool) error
	DescribeInstancesRequestFn                      func(*ec2.DescribeInstancesInput) (*request.Request, *ec2.DescribeInstancesOutput)
	DescribeInstancesFn                             func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	DescribeInstancesPagesFn                        func(*ec2.DescribeInstancesInput, func(*ec2.DescribeInstancesOutput, bool) bool) error
	DescribeInternetGatewaysRequestFn               func(*ec2.DescribeInternetGatewaysInput) (*request.Request, *ec2.DescribeInternetGatewaysOutput)
	DescribeInternetGatewaysFn                      func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error)
	DescribeKeyPairsRequestFn                       func(*ec2.DescribeKeyPairsInput) (*request.Request, *ec2.DescribeKeyPairsOutput)
	DescribeKeyPairsFn                              func(*ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error)
	DescribeMovingAddressesRequestFn                func(*ec2.DescribeMovingAddressesInput) (*request.Request, *ec2.DescribeMovingAddressesOutput)
	DescribeMovingAddressesFn                       func(*ec2.DescribeMovingAddressesInput) (*ec2.DescribeMovingAddressesOutput, error)
	DescribeNatGatewaysRequestFn                    func(*ec2.DescribeNatGatewaysInput) (*request.Request, *ec2.DescribeNatGatewaysOutput)
	DescribeNatGatewaysFn                           func(*ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error)
	DescribeNetworkAclsRequestFn                    func(*ec2.DescribeNetworkAclsInput) (*request.Request, *ec2.DescribeNetworkAclsOutput)
	DescribeNetworkAclsFn                           func(*ec2.DescribeNetworkAclsInput) (*ec2.DescribeNetworkAclsOutput, error)
	DescribeNetworkInterfaceAttributeRequestFn      func(*ec2.DescribeNetworkInterfaceAttributeInput) (*request.Request, *ec2.DescribeNetworkInterfaceAttributeOutput)
	DescribeNetworkInterfaceAttributeFn             func(*ec2.DescribeNetworkInterfaceAttributeInput) (*ec2.DescribeNetworkInterfaceAttributeOutput, error)
	DescribeNetworkInterfacesRequestFn              func(*ec2.DescribeNetworkInterfacesInput) (*request.Request, *ec2.DescribeNetworkInterfacesOutput)
	DescribeNetworkInterfacesFn                     func(*ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error)
	DescribePlacementGroupsRequestFn                func(*ec2.DescribePlacementGroupsInput) (*request.Request, *ec2.DescribePlacementGroupsOutput)
	DescribePlacementGroupsFn                       func(*ec2.DescribePlacementGroupsInput) (*ec2.DescribePlacementGroupsOutput, error)
	DescribePrefixListsRequestFn                    func(*ec2.DescribePrefixListsInput) (*request.Request, *ec2.DescribePrefixListsOutput)
	DescribePrefixListsFn                           func(*ec2.DescribePrefixListsInput) (*ec2.DescribePrefixListsOutput, error)
	DescribeRegionsRequestFn                        func(*ec2.DescribeRegionsInput) (*request.Request, *ec2.DescribeRegionsOutput)
	DescribeRegionsFn                               func(*ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error)
	DescribeReservedInstancesRequestFn              func(*ec2.DescribeReservedInstancesInput) (*request.Request, *ec2.DescribeReservedInstancesOutput)
	DescribeReservedInstancesFn                     func(*ec2.DescribeReservedInstancesInput) (*ec2.DescribeReservedInstancesOutput, error)
	DescribeReservedInstancesListingsRequestFn      func(*ec2.DescribeReservedInstancesListingsInput) (*request.Request, *ec2.DescribeReservedInstancesListingsOutput)
	DescribeReservedInstancesListingsFn             func(*ec2.DescribeReservedInstancesListingsInput) (*ec2.DescribeReservedInstancesListingsOutput, error)
	DescribeReservedInstancesModificationsRequestFn func(*ec2.DescribeReservedInstancesModificationsInput) (*request.Request, *ec2.DescribeReservedInstancesModificationsOutput)
	DescribeReservedInstancesModificationsFn        func(*ec2.DescribeReservedInstancesModificationsInput) (*ec2.DescribeReservedInstancesModificationsOutput, error)
	DescribeReservedInstancesModificationsPagesFn   func(*ec2.DescribeReservedInstancesModificationsInput, func(*ec2.DescribeReservedInstancesModificationsOutput, bool) bool) error
	DescribeReservedInstancesOfferingsRequestFn     func(*ec2.DescribeReservedInstancesOfferingsInput) (*request.Request, *ec2.DescribeReservedInstancesOfferingsOutput)
	DescribeReservedInstancesOfferingsFn            func(*ec2.DescribeReservedInstancesOfferingsInput) (*ec2.DescribeReservedInstancesOfferingsOutput, error)
	DescribeReservedInstancesOfferingsPagesFn       func(*ec2.DescribeReservedInstancesOfferingsInput, func(*ec2.DescribeReservedInstancesOfferingsOutput, bool) bool) error
	DescribeRouteTablesRequestFn                    func(*ec2.DescribeRouteTablesInput) (*request.Request, *ec2.DescribeRouteTablesOutput)
	DescribeRouteTablesFn                           func(*ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error)
	DescribeScheduledInstanceAvailabilityRequestFn  func(*ec2.DescribeScheduledInstanceAvailabilityInput) (*request.Request, *ec2.DescribeScheduledInstanceAvailabilityOutput)
	DescribeScheduledInstanceAvailabilityFn         func(*ec2.DescribeScheduledInstanceAvailabilityInput) (*ec2.DescribeScheduledInstanceAvailabilityOutput, error)
	DescribeScheduledInstancesRequestFn             func(*ec2.DescribeScheduledInstancesInput) (*request.Request, *ec2.DescribeScheduledInstancesOutput)
	DescribeScheduledInstancesFn                    func(*ec2.DescribeScheduledInstancesInput) (*ec2.DescribeScheduledInstancesOutput, error)
	DescribeSecurityGroupsRequestFn                 func(*ec2.DescribeSecurityGroupsInput) (*request.Request, *ec2.DescribeSecurityGroupsOutput)
	DescribeSecurityGroupsFn                        func(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	DescribeSnapshotAttributeRequestFn              func(*ec2.DescribeSnapshotAttributeInput) (*request.Request, *ec2.DescribeSnapshotAttributeOutput)
	DescribeSnapshotAttributeFn                     func(*ec2.DescribeSnapshotAttributeInput) (*ec2.DescribeSnapshotAttributeOutput, error)
	DescribeSnapshotsRequestFn                      func(*ec2.DescribeSnapshotsInput) (*request.Request, *ec2.DescribeSnapshotsOutput)
	DescribeSnapshotsFn                             func(*ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error)
	DescribeSnapshotsPagesFn                        func(*ec2.DescribeSnapshotsInput, func(*ec2.DescribeSnapshotsOutput, bool) bool) error
	DescribeSpotDatafeedSubscriptionRequestFn       func(*ec2.DescribeSpotDatafeedSubscriptionInput) (*request.Request, *ec2.DescribeSpotDatafeedSubscriptionOutput)
	DescribeSpotDatafeedSubscriptionFn              func(*ec2.DescribeSpotDatafeedSubscriptionInput) (*ec2.DescribeSpotDatafeedSubscriptionOutput, error)
	DescribeSpotFleetInstancesRequestFn             func(*ec2.DescribeSpotFleetInstancesInput) (*request.Request, *ec2.DescribeSpotFleetInstancesOutput)
	DescribeSpotFleetInstancesFn                    func(*ec2.DescribeSpotFleetInstancesInput) (*ec2.DescribeSpotFleetInstancesOutput, error)
	DescribeSpotFleetRequestHistoryRequestFn        func(*ec2.DescribeSpotFleetRequestHistoryInput) (*request.Request, *ec2.DescribeSpotFleetRequestHistoryOutput)
	DescribeSpotFleetRequestHistoryFn               func(*ec2.DescribeSpotFleetRequestHistoryInput) (*ec2.DescribeSpotFleetRequestHistoryOutput, error)
	DescribeSpotFleetRequestsRequestFn              func(*ec2.DescribeSpotFleetRequestsInput) (*request.Request, *ec2.DescribeSpotFleetRequestsOutput)
	DescribeSpotFleetRequestsFn                     func(*ec2.DescribeSpotFleetRequestsInput) (*ec2.DescribeSpotFleetRequestsOutput, error)
	DescribeSpotInstanceRequestsRequestFn           func(*ec2.DescribeSpotInstanceRequestsInput) (*request.Request, *ec2.DescribeSpotInstanceRequestsOutput)
	DescribeSpotInstanceRequestsFn                  func(*ec2.DescribeSpotInstanceRequestsInput) (*ec2.DescribeSpotInstanceRequestsOutput, error)
	DescribeSpotPriceHistoryRequestFn               func(*ec2.DescribeSpotPriceHistoryInput) (*request.Request, *ec2.DescribeSpotPriceHistoryOutput)
	DescribeSpotPriceHistoryFn                      func(*ec2.DescribeSpotPriceHistoryInput) (*ec2.DescribeSpotPriceHistoryOutput, error)
	DescribeSpotPriceHistoryPagesFn                 func(*ec2.DescribeSpotPriceHistoryInput, func(*ec2.DescribeSpotPriceHistoryOutput, bool) bool) error
	DescribeSubnetsRequestFn                        func(*ec2.DescribeSubnetsInput) (*request.Request, *ec2.DescribeSubnetsOutput)
	DescribeSubnetsFn                               func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeTagsRequestFn                           func(*ec2.DescribeTagsInput) (*request.Request, *ec2.DescribeTagsOutput)
	DescribeTagsFn                                  func(*ec2.DescribeTagsInput) (*ec2.DescribeTagsOutput, error)
	DescribeTagsPagesFn                             func(*ec2.DescribeTagsInput, func(*ec2.DescribeTagsOutput, bool) bool) error
	DescribeVolumeAttributeRequestFn                func(*ec2.DescribeVolumeAttributeInput) (*request.Request, *ec2.DescribeVolumeAttributeOutput)
	DescribeVolumeAttributeFn                       func(*ec2.DescribeVolumeAttributeInput) (*ec2.DescribeVolumeAttributeOutput, error)
	DescribeVolumeStatusRequestFn                   func(*ec2.DescribeVolumeStatusInput) (*request.Request, *ec2.DescribeVolumeStatusOutput)
	DescribeVolumeStatusFn                          func(*ec2.DescribeVolumeStatusInput) (*ec2.DescribeVolumeStatusOutput, error)
	DescribeVolumeStatusPagesFn                     func(*ec2.DescribeVolumeStatusInput, func(*ec2.DescribeVolumeStatusOutput, bool) bool) error
	DescribeVolumesRequestFn                        func(*ec2.DescribeVolumesInput) (*request.Request, *ec2.DescribeVolumesOutput)
	DescribeVolumesFn                               func(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error)
	DescribeVolumesPagesFn                          func(*ec2.DescribeVolumesInput, func(*ec2.DescribeVolumesOutput, bool) bool) error
	DescribeVpcAttributeRequestFn                   func(*ec2.DescribeVpcAttributeInput) (*request.Request, *ec2.DescribeVpcAttributeOutput)
	DescribeVpcAttributeFn                          func(*ec2.DescribeVpcAttributeInput) (*ec2.DescribeVpcAttributeOutput, error)
	DescribeVpcClassicLinkRequestFn                 func(*ec2.DescribeVpcClassicLinkInput) (*request.Request, *ec2.DescribeVpcClassicLinkOutput)
	DescribeVpcClassicLinkFn                        func(*ec2.DescribeVpcClassicLinkInput) (*ec2.DescribeVpcClassicLinkOutput, error)
	DescribeVpcClassicLinkDnsSupportRequestFn       func(*ec2.DescribeVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.DescribeVpcClassicLinkDnsSupportOutput)
	DescribeVpcClassicLinkDnsSupportFn              func(*ec2.DescribeVpcClassicLinkDnsSupportInput) (*ec2.DescribeVpcClassicLinkDnsSupportOutput, error)
	DescribeVpcEndpointServicesRequestFn            func(*ec2.DescribeVpcEndpointServicesInput) (*request.Request, *ec2.DescribeVpcEndpointServicesOutput)
	DescribeVpcEndpointServicesFn                   func(*ec2.DescribeVpcEndpointServicesInput) (*ec2.DescribeVpcEndpointServicesOutput, error)
	DescribeVpcEndpointsRequestFn                   func(*ec2.DescribeVpcEndpointsInput) (*request.Request, *ec2.DescribeVpcEndpointsOutput)
	DescribeVpcEndpointsFn                          func(*ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error)
	DescribeVpcPeeringConnectionsRequestFn          func(*ec2.DescribeVpcPeeringConnectionsInput) (*request.Request, *ec2.DescribeVpcPeeringConnectionsOutput)
	DescribeVpcPeeringConnectionsFn                 func(*ec2.DescribeVpcPeeringConnectionsInput) (*ec2.DescribeVpcPeeringConnectionsOutput, error)
	DescribeVpcsRequestFn                           func(*ec2.DescribeVpcsInput) (*request.Request, *ec2.DescribeVpcsOutput)
	DescribeVpcsFn                                  func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeVpnConnectionsRequestFn                 func(*ec2.DescribeVpnConnectionsInput) (*request.Request, *ec2.DescribeVpnConnectionsOutput)
	DescribeVpnConnectionsFn                        func(*ec2.DescribeVpnConnectionsInput) (*ec2.DescribeVpnConnectionsOutput, error)
	DescribeVpnGatewaysRequestFn                    func(*ec2.DescribeVpnGatewaysInput) (*request.Request, *ec2.DescribeVpnGatewaysOutput)
	DescribeVpnGatewaysFn                           func(*ec2.DescribeVpnGatewaysInput) (*ec2.DescribeVpnGatewaysOutput, error)
	DetachClassicLinkVpcRequestFn                   func(*ec2.DetachClassicLinkVpcInput) (*request.Request, *ec2.DetachClassicLinkVpcOutput)
	DetachClassicLinkVpcFn                          func(*ec2.DetachClassicLinkVpcInput) (*ec2.DetachClassicLinkVpcOutput, error)
	DetachInternetGatewayRequestFn                  func(*ec2.DetachInternetGatewayInput) (*request.Request, *ec2.DetachInternetGatewayOutput)
	DetachInternetGatewayFn                         func(*ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error)
	DetachNetworkInterfaceRequestFn                 func(*ec2.DetachNetworkInterfaceInput) (*request.Request, *ec2.DetachNetworkInterfaceOutput)
	DetachNetworkInterfaceFn                        func(*ec2.DetachNetworkInterfaceInput) (*ec2.DetachNetworkInterfaceOutput, error)
	DetachVolumeRequestFn                           func(*ec2.DetachVolumeInput) (*request.Request, *ec2.VolumeAttachment)
	DetachVolumeFn                                  func(*ec2.DetachVolumeInput) (*ec2.VolumeAttachment, error)
	DetachVpnGatewayRequestFn                       func(*ec2.DetachVpnGatewayInput) (*request.Request, *ec2.DetachVpnGatewayOutput)
	DetachVpnGatewayFn                              func(*ec2.DetachVpnGatewayInput) (*ec2.DetachVpnGatewayOutput, error)
	DisableVgwRoutePropagationRequestFn             func(*ec2.DisableVgwRoutePropagationInput) (*request.Request, *ec2.DisableVgwRoutePropagationOutput)
	DisableVgwRoutePropagationFn                    func(*ec2.DisableVgwRoutePropagationInput) (*ec2.DisableVgwRoutePropagationOutput, error)
	DisableVpcClassicLinkRequestFn                  func(*ec2.DisableVpcClassicLinkInput) (*request.Request, *ec2.DisableVpcClassicLinkOutput)
	DisableVpcClassicLinkFn                         func(*ec2.DisableVpcClassicLinkInput) (*ec2.DisableVpcClassicLinkOutput, error)
	DisableVpcClassicLinkDnsSupportRequestFn        func(*ec2.DisableVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.DisableVpcClassicLinkDnsSupportOutput)
	DisableVpcClassicLinkDnsSupportFn               func(*ec2.DisableVpcClassicLinkDnsSupportInput) (*ec2.DisableVpcClassicLinkDnsSupportOutput, error)
	DisassociateAddressRequestFn                    func(*ec2.DisassociateAddressInput) (*request.Request, *ec2.DisassociateAddressOutput)
	DisassociateAddressFn                           func(*ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error)
	DisassociateRouteTableRequestFn                 func(*ec2.DisassociateRouteTableInput) (*request.Request, *ec2.DisassociateRouteTableOutput)
	DisassociateRouteTableFn                        func(*ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error)
	EnableVgwRoutePropagationRequestFn              func(*ec2.EnableVgwRoutePropagationInput) (*request.Request, *ec2.EnableVgwRoutePropagationOutput)
	EnableVgwRoutePropagationFn                     func(*ec2.EnableVgwRoutePropagationInput) (*ec2.EnableVgwRoutePropagationOutput, error)
	EnableVolumeIORequestFn                         func(*ec2.EnableVolumeIOInput) (*request.Request, *ec2.EnableVolumeIOOutput)
	EnableVolumeIOFn                                func(*ec2.EnableVolumeIOInput) (*ec2.EnableVolumeIOOutput, error)
	EnableVpcClassicLinkRequestFn                   func(*ec2.EnableVpcClassicLinkInput) (*request.Request, *ec2.EnableVpcClassicLinkOutput)
	EnableVpcClassicLinkFn                          func(*ec2.EnableVpcClassicLinkInput) (*ec2.EnableVpcClassicLinkOutput, error)
	EnableVpcClassicLinkDnsSupportRequestFn         func(*ec2.EnableVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.EnableVpcClassicLinkDnsSupportOutput)
	EnableVpcClassicLinkDnsSupportFn                func(*ec2.EnableVpcClassicLinkDnsSupportInput) (*ec2.EnableVpcClassicLinkDnsSupportOutput, error)
	GetConsoleOutputRequestFn                       func(*ec2.GetConsoleOutputInput) (*request.Request, *ec2.GetConsoleOutputOutput)
	GetConsoleOutputFn                              func(*ec2.GetConsoleOutputInput) (*ec2.GetConsoleOutputOutput, error)
	GetPasswordDataRequestFn                        func(*ec2.GetPasswordDataInput) (*request.Request, *ec2.GetPasswordDataOutput)
	GetPasswordDataFn                               func(*ec2.GetPasswordDataInput) (*ec2.GetPasswordDataOutput, error)
	ImportImageRequestFn                            func(*ec2.ImportImageInput) (*request.Request, *ec2.ImportImageOutput)
	ImportImageFn                                   func(*ec2.ImportImageInput) (*ec2.ImportImageOutput, error)
	ImportInstanceRequestFn                         func(*ec2.ImportInstanceInput) (*request.Request, *ec2.ImportInstanceOutput)
	ImportInstanceFn                                func(*ec2.ImportInstanceInput) (*ec2.ImportInstanceOutput, error)
	ImportKeyPairRequestFn                          func(*ec2.ImportKeyPairInput) (*request.Request, *ec2.ImportKeyPairOutput)
	ImportKeyPairFn                                 func(*ec2.ImportKeyPairInput) (*ec2.ImportKeyPairOutput, error)
	ImportSnapshotRequestFn                         func(*ec2.ImportSnapshotInput) (*request.Request, *ec2.ImportSnapshotOutput)
	ImportSnapshotFn                                func(*ec2.ImportSnapshotInput) (*ec2.ImportSnapshotOutput, error)
	ImportVolumeRequestFn                           func(*ec2.ImportVolumeInput) (*request.Request, *ec2.ImportVolumeOutput)
	ImportVolumeFn                                  func(*ec2.ImportVolumeInput) (*ec2.ImportVolumeOutput, error)
	ModifyHostsRequestFn                            func(*ec2.ModifyHostsInput) (*request.Request, *ec2.ModifyHostsOutput)
	ModifyHostsFn                                   func(*ec2.ModifyHostsInput) (*ec2.ModifyHostsOutput, error)
	ModifyIdFormatRequestFn                         func(*ec2.ModifyIdFormatInput) (*request.Request, *ec2.ModifyIdFormatOutput)
	ModifyIdFormatFn                                func(*ec2.ModifyIdFormatInput) (*ec2.ModifyIdFormatOutput, error)
	ModifyImageAttributeRequestFn                   func(*ec2.ModifyImageAttributeInput) (*request.Request, *ec2.ModifyImageAttributeOutput)
	ModifyImageAttributeFn                          func(*ec2.ModifyImageAttributeInput) (*ec2.ModifyImageAttributeOutput, error)
	ModifyInstanceAttributeRequestFn                func(*ec2.ModifyInstanceAttributeInput) (*request.Request, *ec2.ModifyInstanceAttributeOutput)
	ModifyInstanceAttributeFn                       func(*ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error)
	ModifyInstancePlacementRequestFn                func(*ec2.ModifyInstancePlacementInput) (*request.Request, *ec2.ModifyInstancePlacementOutput)
	ModifyInstancePlacementFn                       func(*ec2.ModifyInstancePlacementInput) (*ec2.ModifyInstancePlacementOutput, error)
	ModifyNetworkInterfaceAttributeRequestFn        func(*ec2.ModifyNetworkInterfaceAttributeInput) (*request.Request, *ec2.ModifyNetworkInterfaceAttributeOutput)
	ModifyNetworkInterfaceAttributeFn               func(*ec2.ModifyNetworkInterfaceAttributeInput) (*ec2.ModifyNetworkInterfaceAttributeOutput, error)
	ModifyReservedInstancesRequestFn                func(*ec2.ModifyReservedInstancesInput) (*request.Request, *ec2.ModifyReservedInstancesOutput)
	ModifyReservedInstancesFn                       func(*ec2.ModifyReservedInstancesInput) (*ec2.ModifyReservedInstancesOutput, error)
	ModifySnapshotAttributeRequestFn                func(*ec2.ModifySnapshotAttributeInput) (*request.Request, *ec2.ModifySnapshotAttributeOutput)
	ModifySnapshotAttributeFn                       func(*ec2.ModifySnapshotAttributeInput) (*ec2.ModifySnapshotAttributeOutput, error)
	ModifySpotFleetRequestRequestFn                 func(*ec2.ModifySpotFleetRequestInput) (*request.Request, *ec2.ModifySpotFleetRequestOutput)
	ModifySpotFleetRequestFn                        func(*ec2.ModifySpotFleetRequestInput) (*ec2.ModifySpotFleetRequestOutput, error)
	ModifySubnetAttributeRequestFn                  func(*ec2.ModifySubnetAttributeInput) (*request.Request, *ec2.ModifySubnetAttributeOutput)
	ModifySubnetAttributeFn                         func(*ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error)
	ModifyVolumeAttributeRequestFn                  func(*ec2.ModifyVolumeAttributeInput) (*request.Request, *ec2.ModifyVolumeAttributeOutput)
	ModifyVolumeAttributeFn                         func(*ec2.ModifyVolumeAttributeInput) (*ec2.ModifyVolumeAttributeOutput, error)
	ModifyVpcAttributeRequestFn                     func(*ec2.ModifyVpcAttributeInput) (*request.Request, *ec2.ModifyVpcAttributeOutput)
	ModifyVpcAttributeFn                            func(*ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error)
	ModifyVpcEndpointRequestFn                      func(*ec2.ModifyVpcEndpointInput) (*request.Request, *ec2.ModifyVpcEndpointOutput)
	ModifyVpcEndpointFn                             func(*ec2.ModifyVpcEndpointInput) (*ec2.ModifyVpcEndpointOutput, error)
	MonitorInstancesRequestFn                       func(*ec2.MonitorInstancesInput) (*request.Request, *ec2.MonitorInstancesOutput)
	MonitorInstancesFn                              func(*ec2.MonitorInstancesInput) (*ec2.MonitorInstancesOutput, error)
	MoveAddressToVpcRequestFn                       func(*ec2.MoveAddressToVpcInput) (*request.Request, *ec2.MoveAddressToVpcOutput)
	MoveAddressToVpcFn                              func(*ec2.MoveAddressToVpcInput) (*ec2.MoveAddressToVpcOutput, error)
	PurchaseReservedInstancesOfferingRequestFn      func(*ec2.PurchaseReservedInstancesOfferingInput) (*request.Request, *ec2.PurchaseReservedInstancesOfferingOutput)
	PurchaseReservedInstancesOfferingFn             func(*ec2.PurchaseReservedInstancesOfferingInput) (*ec2.PurchaseReservedInstancesOfferingOutput, error)
	PurchaseScheduledInstancesRequestFn             func(*ec2.PurchaseScheduledInstancesInput) (*request.Request, *ec2.PurchaseScheduledInstancesOutput)
	PurchaseScheduledInstancesFn                    func(*ec2.PurchaseScheduledInstancesInput) (*ec2.PurchaseScheduledInstancesOutput, error)
	RebootInstancesRequestFn                        func(*ec2.RebootInstancesInput) (*request.Request, *ec2.RebootInstancesOutput)
	RebootInstancesFn                               func(*ec2.RebootInstancesInput) (*ec2.RebootInstancesOutput, error)
	RegisterImageRequestFn                          func(*ec2.RegisterImageInput) (*request.Request, *ec2.RegisterImageOutput)
	RegisterImageFn                                 func(*ec2.RegisterImageInput) (*ec2.RegisterImageOutput, error)
	RejectVpcPeeringConnectionRequestFn             func(*ec2.RejectVpcPeeringConnectionInput) (*request.Request, *ec2.RejectVpcPeeringConnectionOutput)
	RejectVpcPeeringConnectionFn                    func(*ec2.RejectVpcPeeringConnectionInput) (*ec2.RejectVpcPeeringConnectionOutput, error)
	ReleaseAddressRequestFn                         func(*ec2.ReleaseAddressInput) (*request.Request, *ec2.ReleaseAddressOutput)
	ReleaseAddressFn                                func(*ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error)
	ReleaseHostsRequestFn                           func(*ec2.ReleaseHostsInput) (*request.Request, *ec2.ReleaseHostsOutput)
	ReleaseHostsFn                                  func(*ec2.ReleaseHostsInput) (*ec2.ReleaseHostsOutput, error)
	ReplaceNetworkAclAssociationRequestFn           func(*ec2.ReplaceNetworkAclAssociationInput) (*request.Request, *ec2.ReplaceNetworkAclAssociationOutput)
	ReplaceNetworkAclAssociationFn                  func(*ec2.ReplaceNetworkAclAssociationInput) (*ec2.ReplaceNetworkAclAssociationOutput, error)
	ReplaceNetworkAclEntryRequestFn                 func(*ec2.ReplaceNetworkAclEntryInput) (*request.Request, *ec2.ReplaceNetworkAclEntryOutput)
	ReplaceNetworkAclEntryFn                        func(*ec2.ReplaceNetworkAclEntryInput) (*ec2.ReplaceNetworkAclEntryOutput, error)
	ReplaceRouteRequestFn                           func(*ec2.ReplaceRouteInput) (*request.Request, *ec2.ReplaceRouteOutput)
	ReplaceRouteFn                                  func(*ec2.ReplaceRouteInput) (*ec2.ReplaceRouteOutput, error)
	ReplaceRouteTableAssociationRequestFn           func(*ec2.ReplaceRouteTableAssociationInput) (*request.Request, *ec2.ReplaceRouteTableAssociationOutput)
	ReplaceRouteTableAssociationFn                  func(*ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error)
	ReportInstanceStatusRequestFn                   func(*ec2.ReportInstanceStatusInput) (*request.Request, *ec2.ReportInstanceStatusOutput)
	ReportInstanceStatusFn                          func(*ec2.ReportInstanceStatusInput) (*ec2.ReportInstanceStatusOutput, error)
	RequestSpotFleetRequestFn                       func(*ec2.RequestSpotFleetInput) (*request.Request, *ec2.RequestSpotFleetOutput)
	RequestSpotFleetFn                              func(*ec2.RequestSpotFleetInput) (*ec2.RequestSpotFleetOutput, error)
	RequestSpotInstancesRequestFn                   func(*ec2.RequestSpotInstancesInput) (*request.Request, *ec2.RequestSpotInstancesOutput)
	RequestSpotInstancesFn                          func(*ec2.RequestSpotInstancesInput) (*ec2.RequestSpotInstancesOutput, error)
	ResetImageAttributeRequestFn                    func(*ec2.ResetImageAttributeInput) (*request.Request, *ec2.ResetImageAttributeOutput)
	ResetImageAttributeFn                           func(*ec2.ResetImageAttributeInput) (*ec2.ResetImageAttributeOutput, error)
	ResetInstanceAttributeRequestFn                 func(*ec2.ResetInstanceAttributeInput) (*request.Request, *ec2.ResetInstanceAttributeOutput)
	ResetInstanceAttributeFn                        func(*ec2.ResetInstanceAttributeInput) (*ec2.ResetInstanceAttributeOutput, error)
	ResetNetworkInterfaceAttributeRequestFn         func(*ec2.ResetNetworkInterfaceAttributeInput) (*request.Request, *ec2.ResetNetworkInterfaceAttributeOutput)
	ResetNetworkInterfaceAttributeFn                func(*ec2.ResetNetworkInterfaceAttributeInput) (*ec2.ResetNetworkInterfaceAttributeOutput, error)
	ResetSnapshotAttributeRequestFn                 func(*ec2.ResetSnapshotAttributeInput) (*request.Request, *ec2.ResetSnapshotAttributeOutput)
	ResetSnapshotAttributeFn                        func(*ec2.ResetSnapshotAttributeInput) (*ec2.ResetSnapshotAttributeOutput, error)
	RestoreAddressToClassicRequestFn                func(*ec2.RestoreAddressToClassicInput) (*request.Request, *ec2.RestoreAddressToClassicOutput)
	RestoreAddressToClassicFn                       func(*ec2.RestoreAddressToClassicInput) (*ec2.RestoreAddressToClassicOutput, error)
	RevokeSecurityGroupEgressRequestFn              func(*ec2.RevokeSecurityGroupEgressInput) (*request.Request, *ec2.RevokeSecurityGroupEgressOutput)
	RevokeSecurityGroupEgressFn                     func(*ec2.RevokeSecurityGroupEgressInput) (*ec2.RevokeSecurityGroupEgressOutput, error)
	RevokeSecurityGroupIngressRequestFn             func(*ec2.RevokeSecurityGroupIngressInput) (*request.Request, *ec2.RevokeSecurityGroupIngressOutput)
	RevokeSecurityGroupIngressFn                    func(*ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error)
	RunInstancesRequestFn                           func(*ec2.RunInstancesInput) (*request.Request, *ec2.Reservation)
	RunInstancesFn                                  func(*ec2.RunInstancesInput) (*ec2.Reservation, error)
	RunScheduledInstancesRequestFn                  func(*ec2.RunScheduledInstancesInput) (*request.Request, *ec2.RunScheduledInstancesOutput)
	RunScheduledInstancesFn                         func(*ec2.RunScheduledInstancesInput) (*ec2.RunScheduledInstancesOutput, error)
	StartInstancesRequestFn                         func(*ec2.StartInstancesInput) (*request.Request, *ec2.StartInstancesOutput)
	StartInstancesFn                                func(*ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error)
	StopInstancesRequestFn                          func(*ec2.StopInstancesInput) (*request.Request, *ec2.StopInstancesOutput)
	StopInstancesFn                                 func(*ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error)
	TerminateInstancesRequestFn                     func(*ec2.TerminateInstancesInput) (*request.Request, *ec2.TerminateInstancesOutput)
	TerminateInstancesFn                            func(*ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error)
	UnassignPrivateIpAddressesRequestFn             func(*ec2.UnassignPrivateIpAddressesInput) (*request.Request, *ec2.UnassignPrivateIpAddressesOutput)
	UnassignPrivateIpAddressesFn                    func(*ec2.UnassignPrivateIpAddressesInput) (*ec2.UnassignPrivateIpAddressesOutput, error)
	UnmonitorInstancesRequestFn                     func(*ec2.UnmonitorInstancesInput) (*request.Request, *ec2.UnmonitorInstancesOutput)
	UnmonitorInstancesFn                            func(*ec2.UnmonitorInstancesInput) (*ec2.UnmonitorInstancesOutput, error)
}

func (f *EC2) AcceptVpcPeeringConnectionRequest(input *ec2.AcceptVpcPeeringConnectionInput) (*request.Request, *ec2.AcceptVpcPeeringConnectionOutput) {
	if f.AcceptVpcPeeringConnectionRequestFn == nil {
		return nil, nil
	}
	return f.AcceptVpcPeeringConnectionRequestFn(input)
}

func (f *EC2) AcceptVpcPeeringConnection(input *ec2.AcceptVpcPeeringConnectionInput) (*ec2.AcceptVpcPeeringConnectionOutput, error) {
	if f.AcceptVpcPeeringConnectionFn == nil {
		return nil, nil
	}
	return f.AcceptVpcPeeringConnectionFn(input)
}

func (f *EC2) AllocateAddressRequest(input *ec2.AllocateAddressInput) (*request.Request, *ec2.AllocateAddressOutput) {
	if f.AllocateAddressRequestFn == nil {
		return nil, nil
	}
	return f.AllocateAddressRequestFn(input)
}

func (f *EC2) AllocateAddress(input *ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
	if f.AllocateAddressFn == nil {
		return nil, nil
	}
	return f.AllocateAddressFn(input)
}

func (f *EC2) AllocateHostsRequest(input *ec2.AllocateHostsInput) (*request.Request, *ec2.AllocateHostsOutput) {
	if f.AllocateHostsRequestFn == nil {
		return nil, nil
	}
	return f.AllocateHostsRequestFn(input)
}

func (f *EC2) AllocateHosts(input *ec2.AllocateHostsInput) (*ec2.AllocateHostsOutput, error) {
	if f.AllocateHostsFn == nil {
		return nil, nil
	}
	return f.AllocateHostsFn(input)
}

func (f *EC2) AssignPrivateIpAddressesRequest(input *ec2.AssignPrivateIpAddressesInput) (*request.Request, *ec2.AssignPrivateIpAddressesOutput) {
	if f.AssignPrivateIpAddressesRequestFn == nil {
		return nil, nil
	}
	return f.AssignPrivateIpAddressesRequestFn(input)
}

func (f *EC2) AssignPrivateIpAddresses(input *ec2.AssignPrivateIpAddressesInput) (*ec2.AssignPrivateIpAddressesOutput, error) {
	if f.AssignPrivateIpAddressesFn == nil {
		return nil, nil
	}
	return f.AssignPrivateIpAddressesFn(input)
}

func (f *EC2) AssociateAddressRequest(input *ec2.AssociateAddressInput) (*request.Request, *ec2.AssociateAddressOutput) {
	if f.AssociateAddressRequestFn == nil {
		return nil, nil
	}
	return f.AssociateAddressRequestFn(input)
}

func (f *EC2) AssociateAddress(input *ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error) {
	if f.AssociateAddressFn == nil {
		return nil, nil
	}
	return f.AssociateAddressFn(input)
}

func (f *EC2) AssociateDhcpOptionsRequest(input *ec2.AssociateDhcpOptionsInput) (*request.Request, *ec2.AssociateDhcpOptionsOutput) {
	if f.AssociateDhcpOptionsRequestFn == nil {
		return nil, nil
	}
	return f.AssociateDhcpOptionsRequestFn(input)
}

func (f *EC2) AssociateDhcpOptions(input *ec2.AssociateDhcpOptionsInput) (*ec2.AssociateDhcpOptionsOutput, error) {
	if f.AssociateDhcpOptionsFn == nil {
		return nil, nil
	}
	return f.AssociateDhcpOptionsFn(input)
}

func (f *EC2) AssociateRouteTableRequest(input *ec2.AssociateRouteTableInput) (*request.Request, *ec2.AssociateRouteTableOutput) {
	if f.AssociateRouteTableRequestFn == nil {
		return nil, nil
	}
	return f.AssociateRouteTableRequestFn(input)
}

func (f *EC2) AssociateRouteTable(input *ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
	if f.AssociateRouteTableFn == nil {
		return nil, nil
	}
	return f.AssociateRouteTableFn(input)
}

func (f *EC2) AttachClassicLinkVpcRequest(input *ec2.AttachClassicLinkVpcInput) (*request.Request, *ec2.AttachClassicLinkVpcOutput) {
	if f.AttachClassicLinkVpcRequestFn == nil {
		return nil, nil
	}
	return f.AttachClassicLinkVpcRequestFn(input)
}

func (f *EC2) AttachClassicLinkVpc(input *ec2.AttachClassicLinkVpcInput) (*ec2.AttachClassicLinkVpcOutput, error) {
	if f.AttachClassicLinkVpcFn == nil {
		return nil, nil
	}
	return f.AttachClassicLinkVpcFn(input)
}

func (f *EC2) AttachInternetGatewayRequest(input *ec2.AttachInternetGatewayInput) (*request.Request, *ec2.AttachInternetGatewayOutput) {
	if f.AttachInternetGatewayRequestFn == nil {
		return nil, nil
	}
	return f.AttachInternetGatewayRequestFn(input)
}

func (f *EC2) AttachInternetGateway(input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	if f.AttachInternetGatewayFn == nil {
		return nil, nil
	}
	return f.AttachInternetGatewayFn(input)
}

func (f *EC2) AttachNetworkInterfaceRequest(input *ec2.AttachNetworkInterfaceInput) (*request.Request, *ec2.AttachNetworkInterfaceOutput) {
	if f.AttachNetworkInterfaceRequestFn == nil {
		return nil, nil
	}
	return f.AttachNetworkInterfaceRequestFn(input)
}

func (f *EC2) AttachNetworkInterface(input *ec2.AttachNetworkInterfaceInput) (*ec2.AttachNetworkInterfaceOutput, error) {
	if f.AttachNetworkInterfaceFn == nil {
		return nil, nil
	}
	return f.AttachNetworkInterfaceFn(input)
}

func (f *EC2) AttachVolumeRequest(input *ec2.AttachVolumeInput) (*request.Request, *ec2.VolumeAttachment) {
	if f.AttachVolumeRequestFn == nil {
		return nil, nil
	}
	return f.AttachVolumeRequestFn(input)
}

func (f *EC2) AttachVolume(input *ec2.AttachVolumeInput) (*ec2.VolumeAttachment, error) {
	if f.AttachVolumeFn == nil {
		return nil, nil
	}
	return f.AttachVolumeFn(input)
}

func (f *EC2) AttachVpnGatewayRequest(input *ec2.AttachVpnGatewayInput) (*request.Request, *ec2.AttachVpnGatewayOutput) {
	if f.AttachVpnGatewayRequestFn == nil {
		return nil, nil
	}
	return f.AttachVpnGatewayRequestFn(input)
}

func (f *EC2) AttachVpnGateway(input *ec2.AttachVpnGatewayInput) (*ec2.AttachVpnGatewayOutput, error) {
	if f.AttachVpnGatewayFn == nil {
		return nil, nil
	}
	return f.AttachVpnGatewayFn(input)
}

func (f *EC2) AuthorizeSecurityGroupEgressRequest(input *ec2.AuthorizeSecurityGroupEgressInput) (*request.Request, *ec2.AuthorizeSecurityGroupEgressOutput) {
	if f.AuthorizeSecurityGroupEgressRequestFn == nil {
		return nil, nil
	}
	return f.AuthorizeSecurityGroupEgressRequestFn(input)
}

func (f *EC2) AuthorizeSecurityGroupEgress(input *ec2.AuthorizeSecurityGroupEgressInput) (*ec2.AuthorizeSecurityGroupEgressOutput, error) {
	if f.AuthorizeSecurityGroupEgressFn == nil {
		return nil, nil
	}
	return f.AuthorizeSecurityGroupEgressFn(input)
}

func (f *EC2) AuthorizeSecurityGroupIngressRequest(input *ec2.AuthorizeSecurityGroupIngressInput) (*request.Request, *ec2.AuthorizeSecurityGroupIngressOutput) {
	if f.AuthorizeSecurityGroupIngressRequestFn == nil {
		return nil, nil
	}
	return f.AuthorizeSecurityGroupIngressRequestFn(input)
}

func (f *EC2) AuthorizeSecurityGroupIngress(input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	if f.AuthorizeSecurityGroupIngressFn == nil {
		return nil, nil
	}
	return f.AuthorizeSecurityGroupIngressFn(input)
}

func (f *EC2) BundleInstanceRequest(input *ec2.BundleInstanceInput) (*request.Request, *ec2.BundleInstanceOutput) {
	if f.BundleInstanceRequestFn == nil {
		return nil, nil
	}
	return f.BundleInstanceRequestFn(input)
}

func (f *EC2) BundleInstance(input *ec2.BundleInstanceInput) (*ec2.BundleInstanceOutput, error) {
	if f.BundleInstanceFn == nil {
		return nil, nil
	}
	return f.BundleInstanceFn(input)
}

func (f *EC2) CancelBundleTaskRequest(input *ec2.CancelBundleTaskInput) (*request.Request, *ec2.CancelBundleTaskOutput) {
	if f.CancelBundleTaskRequestFn == nil {
		return nil, nil
	}
	return f.CancelBundleTaskRequestFn(input)
}

func (f *EC2) CancelBundleTask(input *ec2.CancelBundleTaskInput) (*ec2.CancelBundleTaskOutput, error) {
	if f.CancelBundleTaskFn == nil {
		return nil, nil
	}
	return f.CancelBundleTaskFn(input)
}

func (f *EC2) CancelConversionTaskRequest(input *ec2.CancelConversionTaskInput) (*request.Request, *ec2.CancelConversionTaskOutput) {
	if f.CancelConversionTaskRequestFn == nil {
		return nil, nil
	}
	return f.CancelConversionTaskRequestFn(input)
}

func (f *EC2) CancelConversionTask(input *ec2.CancelConversionTaskInput) (*ec2.CancelConversionTaskOutput, error) {
	if f.CancelConversionTaskFn == nil {
		return nil, nil
	}
	return f.CancelConversionTaskFn(input)
}

func (f *EC2) CancelExportTaskRequest(input *ec2.CancelExportTaskInput) (*request.Request, *ec2.CancelExportTaskOutput) {
	if f.CancelExportTaskRequestFn == nil {
		return nil, nil
	}
	return f.CancelExportTaskRequestFn(input)
}

func (f *EC2) CancelExportTask(input *ec2.CancelExportTaskInput) (*ec2.CancelExportTaskOutput, error) {
	if f.CancelExportTaskFn == nil {
		return nil, nil
	}
	return f.CancelExportTaskFn(input)
}

func (f *EC2) CancelImportTaskRequest(input *ec2.CancelImportTaskInput) (*request.Request, *ec2.CancelImportTaskOutput) {
	if f.CancelImportTaskRequestFn == nil {
		return nil, nil
	}
	return f.CancelImportTaskRequestFn(input)
}

func (f *EC2) CancelImportTask(input *ec2.CancelImportTaskInput) (*ec2.CancelImportTaskOutput, error) {
	if f.CancelImportTaskFn == nil {
		return nil, nil
	}
	return f.CancelImportTaskFn(input)
}

func (f *EC2) CancelReservedInstancesListingRequest(input *ec2.CancelReservedInstancesListingInput) (*request.Request, *ec2.CancelReservedInstancesListingOutput) {
	if f.CancelReservedInstancesListingRequestFn == nil {
		return nil, nil
	}
	return f.CancelReservedInstancesListingRequestFn(input)
}

func (f *EC2) CancelReservedInstancesListing(input *ec2.CancelReservedInstancesListingInput) (*ec2.CancelReservedInstancesListingOutput, error) {
	if f.CancelReservedInstancesListingFn == nil {
		return nil, nil
	}
	return f.CancelReservedInstancesListingFn(input)
}

func (f *EC2) CancelSpotFleetRequestsRequest(input *ec2.CancelSpotFleetRequestsInput) (*request.Request, *ec2.CancelSpotFleetRequestsOutput) {
	if f.CancelSpotFleetRequestsRequestFn == nil {
		return nil, nil
	}
	return f.CancelSpotFleetRequestsRequestFn(input)
}

func (f *EC2) CancelSpotFleetRequests(input *ec2.CancelSpotFleetRequestsInput) (*ec2.CancelSpotFleetRequestsOutput, error) {
	if f.CancelSpotFleetRequestsFn == nil {
		return nil, nil
	}
	return f.CancelSpotFleetRequestsFn(input)
}

func (f *EC2) CancelSpotInstanceRequestsRequest(input *ec2.CancelSpotInstanceRequestsInput) (*request.Request, *ec2.CancelSpotInstanceRequestsOutput) {
	if f.CancelSpotInstanceRequestsRequestFn == nil {
		return nil, nil
	}
	return f.CancelSpotInstanceRequestsRequestFn(input)
}

func (f *EC2) CancelSpotInstanceRequests(input *ec2.CancelSpotInstanceRequestsInput) (*ec2.CancelSpotInstanceRequestsOutput, error) {
	if f.CancelSpotInstanceRequestsFn == nil {
		return nil, nil
	}
	return f.CancelSpotInstanceRequestsFn(input)
}

func (f *EC2) ConfirmProductInstanceRequest(input *ec2.ConfirmProductInstanceInput) (*request.Request, *ec2.ConfirmProductInstanceOutput) {
	if f.ConfirmProductInstanceRequestFn == nil {
		return nil, nil
	}
	return f.ConfirmProductInstanceRequestFn(input)
}

func (f *EC2) ConfirmProductInstance(input *ec2.ConfirmProductInstanceInput) (*ec2.ConfirmProductInstanceOutput, error) {
	if f.ConfirmProductInstanceFn == nil {
		return nil, nil
	}
	return f.ConfirmProductInstanceFn(input)
}

func (f *EC2) CopyImageRequest(input *ec2.CopyImageInput) (*request.Request, *ec2.CopyImageOutput) {
	if f.CopyImageRequestFn == nil {
		return nil, nil
	}
	return f.CopyImageRequestFn(input)
}

func (f *EC2) CopyImage(input *ec2.CopyImageInput) (*ec2.CopyImageOutput, error) {
	if f.CopyImageFn == nil {
		return nil, nil
	}
	return f.CopyImageFn(input)
}

func (f *EC2) CopySnapshotRequest(input *ec2.CopySnapshotInput) (*request.Request, *ec2.CopySnapshotOutput) {
	if f.CopySnapshotRequestFn == nil {
		return nil, nil
	}
	return f.CopySnapshotRequestFn(input)
}

func (f *EC2) CopySnapshot(input *ec2.CopySnapshotInput) (*ec2.CopySnapshotOutput, error) {
	if f.CopySnapshotFn == nil {
		return nil, nil
	}
	return f.CopySnapshotFn(input)
}

func (f *EC2) CreateCustomerGatewayRequest(input *ec2.CreateCustomerGatewayInput) (*request.Request, *ec2.CreateCustomerGatewayOutput) {
	if f.CreateCustomerGatewayRequestFn == nil {
		return nil, nil
	}
	return f.CreateCustomerGatewayRequestFn(input)
}

func (f *EC2) CreateCustomerGateway(input *ec2.CreateCustomerGatewayInput) (*ec2.CreateCustomerGatewayOutput, error) {
	if f.CreateCustomerGatewayFn == nil {
		return nil, nil
	}
	return f.CreateCustomerGatewayFn(input)
}

func (f *EC2) CreateDhcpOptionsRequest(input *ec2.CreateDhcpOptionsInput) (*request.Request, *ec2.CreateDhcpOptionsOutput) {
	if f.CreateDhcpOptionsRequestFn == nil {
		return nil, nil
	}
	return f.CreateDhcpOptionsRequestFn(input)
}

func (f *EC2) CreateDhcpOptions(input *ec2.CreateDhcpOptionsInput) (*ec2.CreateDhcpOptionsOutput, error) {
	if f.CreateDhcpOptionsFn == nil {
		return nil, nil
	}
	return f.CreateDhcpOptionsFn(input)
}

func (f *EC2) CreateFlowLogsRequest(input *ec2.CreateFlowLogsInput) (*request.Request, *ec2.CreateFlowLogsOutput) {
	if f.CreateFlowLogsRequestFn == nil {
		return nil, nil
	}
	return f.CreateFlowLogsRequestFn(input)
}

func (f *EC2) CreateFlowLogs(input *ec2.CreateFlowLogsInput) (*ec2.CreateFlowLogsOutput, error) {
	if f.CreateFlowLogsFn == nil {
		return nil, nil
	}
	return f.CreateFlowLogsFn(input)
}

func (f *EC2) CreateImageRequest(input *ec2.CreateImageInput) (*request.Request, *ec2.CreateImageOutput) {
	if f.CreateImageRequestFn == nil {
		return nil, nil
	}
	return f.CreateImageRequestFn(input)
}

func (f *EC2) CreateImage(input *ec2.CreateImageInput) (*ec2.CreateImageOutput, error) {
	if f.CreateImageFn == nil {
		return nil, nil
	}
	return f.CreateImageFn(input)
}

func (f *EC2) CreateInstanceExportTaskRequest(input *ec2.CreateInstanceExportTaskInput) (*request.Request, *ec2.CreateInstanceExportTaskOutput) {
	if f.CreateInstanceExportTaskRequestFn == nil {
		return nil, nil
	}
	return f.CreateInstanceExportTaskRequestFn(input)
}

func (f *EC2) CreateInstanceExportTask(input *ec2.CreateInstanceExportTaskInput) (*ec2.CreateInstanceExportTaskOutput, error) {
	if f.CreateInstanceExportTaskFn == nil {
		return nil, nil
	}
	return f.CreateInstanceExportTaskFn(input)
}

func (f *EC2) CreateInternetGatewayRequest(input *ec2.CreateInternetGatewayInput) (*request.Request, *ec2.CreateInternetGatewayOutput) {
	if f.CreateInternetGatewayRequestFn == nil {
		return nil, nil
	}
	return f.CreateInternetGatewayRequestFn(input)
}

func (f *EC2) CreateInternetGateway(input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	if f.CreateInternetGatewayFn == nil {
		return nil, nil
	}
	return f.CreateInternetGatewayFn(input)
}

func (f *EC2) CreateKeyPairRequest(input *ec2.CreateKeyPairInput) (*request.Request, *ec2.CreateKeyPairOutput) {
	if f.CreateKeyPairRequestFn == nil {
		return nil, nil
	}
	return f.CreateKeyPairRequestFn(input)
}

func (f *EC2) CreateKeyPair(input *ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {
	if f.CreateKeyPairFn == nil {
		return nil, nil
	}
	return f.CreateKeyPairFn(input)
}

func (f *EC2) CreateNatGatewayRequest(input *ec2.CreateNatGatewayInput) (*request.Request, *ec2.CreateNatGatewayOutput) {
	if f.CreateNatGatewayRequestFn == nil {
		return nil, nil
	}
	return f.CreateNatGatewayRequestFn(input)
}

func (f *EC2) CreateNatGateway(input *ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
	if f.CreateNatGatewayFn == nil {
		return nil, nil
	}
	return f.CreateNatGatewayFn(input)
}

func (f *EC2) CreateNetworkAclRequest(input *ec2.CreateNetworkAclInput) (*request.Request, *ec2.CreateNetworkAclOutput) {
	if f.CreateNetworkAclRequestFn == nil {
		return nil, nil
	}
	return f.CreateNetworkAclRequestFn(input)
}

func (f *EC2) CreateNetworkAcl(input *ec2.CreateNetworkAclInput) (*ec2.CreateNetworkAclOutput, error) {
	if f.CreateNetworkAclFn == nil {
		return nil, nil
	}
	return f.CreateNetworkAclFn(input)
}

func (f *EC2) CreateNetworkAclEntryRequest(input *ec2.CreateNetworkAclEntryInput) (*request.Request, *ec2.CreateNetworkAclEntryOutput) {
	if f.CreateNetworkAclEntryRequestFn == nil {
		return nil, nil
	}
	return f.CreateNetworkAclEntryRequestFn(input)
}

func (f *EC2) CreateNetworkAclEntry(input *ec2.CreateNetworkAclEntryInput) (*ec2.CreateNetworkAclEntryOutput, error) {
	if f.CreateNetworkAclEntryFn == nil {
		return nil, nil
	}
	return f.CreateNetworkAclEntryFn(input)
}

func (f *EC2) CreateNetworkInterfaceRequest(input *ec2.CreateNetworkInterfaceInput) (*request.Request, *ec2.CreateNetworkInterfaceOutput) {
	if f.CreateNetworkInterfaceRequestFn == nil {
		return nil, nil
	}
	return f.CreateNetworkInterfaceRequestFn(input)
}

func (f *EC2) CreateNetworkInterface(input *ec2.CreateNetworkInterfaceInput) (*ec2.CreateNetworkInterfaceOutput, error) {
	if f.CreateNetworkInterfaceFn == nil {
		return nil, nil
	}
	return f.CreateNetworkInterfaceFn(input)
}

func (f *EC2) CreatePlacementGroupRequest(input *ec2.CreatePlacementGroupInput) (*request.Request, *ec2.CreatePlacementGroupOutput) {
	if f.CreatePlacementGroupRequestFn == nil {
		return nil, nil
	}
	return f.CreatePlacementGroupRequestFn(input)
}

func (f *EC2) CreatePlacementGroup(input *ec2.CreatePlacementGroupInput) (*ec2.CreatePlacementGroupOutput, error) {
	if f.CreatePlacementGroupFn == nil {
		return nil, nil
	}
	return f.CreatePlacementGroupFn(input)
}

func (f *EC2) CreateReservedInstancesListingRequest(input *ec2.CreateReservedInstancesListingInput) (*request.Request, *ec2.CreateReservedInstancesListingOutput) {
	if f.CreateReservedInstancesListingRequestFn == nil {
		return nil, nil
	}
	return f.CreateReservedInstancesListingRequestFn(input)
}

func (f *EC2) CreateReservedInstancesListing(input *ec2.CreateReservedInstancesListingInput) (*ec2.CreateReservedInstancesListingOutput, error) {
	if f.CreateReservedInstancesListingFn == nil {
		return nil, nil
	}
	return f.CreateReservedInstancesListingFn(input)
}

func (f *EC2) CreateRouteRequest(input *ec2.CreateRouteInput) (*request.Request, *ec2.CreateRouteOutput) {
	if f.CreateRouteRequestFn == nil {
		return nil, nil
	}
	return f.CreateRouteRequestFn(input)
}

func (f *EC2) CreateRoute(input *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
	if f.CreateRouteFn == nil {
		return nil, nil
	}
	return f.CreateRouteFn(input)
}

func (f *EC2) CreateRouteTableRequest(input *ec2.CreateRouteTableInput) (*request.Request, *ec2.CreateRouteTableOutput) {
	if f.CreateRouteTableRequestFn == nil {
		return nil, nil
	}
	return f.CreateRouteTableRequestFn(input)
}

func (f *EC2) CreateRouteTable(input *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
	if f.CreateRouteTableFn == nil {
		return nil, nil
	}
	return f.CreateRouteTableFn(input)
}

func (f *EC2) CreateSecurityGroupRequest(input *ec2.CreateSecurityGroupInput) (*request.Request, *ec2.CreateSecurityGroupOutput) {
	if f.CreateSecurityGroupRequestFn == nil {
		return nil, nil
	}
	return f.CreateSecurityGroupRequestFn(input)
}

func (f *EC2) CreateSecurityGroup(input *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	if f.CreateSecurityGroupFn == nil {
		return nil, nil
	}
	return f.CreateSecurityGroupFn(input)
}

func (f *EC2) CreateSnapshotRequest(input *ec2.CreateSnapshotInput) (*request.Request, *ec2.Snapshot) {
	if f.CreateSnapshotRequestFn == nil {
		return nil, nil
	}
	return f.CreateSnapshotRequestFn(input)
}

func (f *EC2) CreateSnapshot(input *ec2.CreateSnapshotInput) (*ec2.Snapshot, error) {
	if f.CreateSnapshotFn == nil {
		return nil, nil
	}
	return f.CreateSnapshotFn(input)
}

func (f *EC2) CreateSpotDatafeedSubscriptionRequest(input *ec2.CreateSpotDatafeedSubscriptionInput) (*request.Request, *ec2.CreateSpotDatafeedSubscriptionOutput) {
	if f.CreateSpotDatafeedSubscriptionRequestFn == nil {
		return nil, nil
	}
	return f.CreateSpotDatafeedSubscriptionRequestFn(input)
}

func (f *EC2) CreateSpotDatafeedSubscription(input *ec2.CreateSpotDatafeedSubscriptionInput) (*ec2.CreateSpotDatafeedSubscriptionOutput, error) {
	if f.CreateSpotDatafeedSubscriptionFn == nil {
		return nil, nil
	}
	return f.CreateSpotDatafeedSubscriptionFn(input)
}

func (f *EC2) CreateSubnetRequest(input *ec2.CreateSubnetInput) (*request.Request, *ec2.CreateSubnetOutput) {
	if f.CreateSubnetRequestFn == nil {
		return nil, nil
	}
	return f.CreateSubnetRequestFn(input)
}

func (f *EC2) CreateSubnet(input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	if f.CreateSubnetFn == nil {
		return nil, nil
	}
	return f.CreateSubnetFn(input)
}

func (f *EC2) CreateTagsRequest(input *ec2.CreateTagsInput) (*request.Request, *ec2.CreateTagsOutput) {
	if f.CreateTagsRequestFn == nil {
		return nil, nil
	}
	return f.CreateTagsRequestFn(input)
}

func (f *EC2) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	if f.CreateTagsFn == nil {
		return nil, nil
	}
	return f.CreateTagsFn(input)
}

func (f *EC2) CreateVolumeRequest(input *ec2.CreateVolumeInput) (*request.Request, *ec2.Volume) {
	if f.CreateVolumeRequestFn == nil {
		return nil, nil
	}
	return f.CreateVolumeRequestFn(input)
}

func (f *EC2) CreateVolume(input *ec2.CreateVolumeInput) (*ec2.Volume, error) {
	if f.CreateVolumeFn == nil {
		return nil, nil
	}
	return f.CreateVolumeFn(input)
}

func (f *EC2) CreateVpcRequest(input *ec2.CreateVpcInput) (*request.Request, *ec2.CreateVpcOutput) {
	if f.CreateVpcRequestFn == nil {
		return nil, nil
	}
	return f.CreateVpcRequestFn(input)
}

func (f *EC2) CreateVpc(input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	if f.CreateVpcFn == nil {
		return nil, nil
	}
	return f.CreateVpcFn(input)
}

func (f *EC2) CreateVpcEndpointRequest(input *ec2.CreateVpcEndpointInput) (*request.Request, *ec2.CreateVpcEndpointOutput) {
	if f.CreateVpcEndpointRequestFn == nil {
		return nil, nil
	}
	return f.CreateVpcEndpointRequestFn(input)
}

func (f *EC2) CreateVpcEndpoint(input *ec2.CreateVpcEndpointInput) (*ec2.CreateVpcEndpointOutput, error) {
	if f.CreateVpcEndpointFn == nil {
		return nil, nil
	}
	return f.CreateVpcEndpointFn(input)
}

func (f *EC2) CreateVpcPeeringConnectionRequest(input *ec2.CreateVpcPeeringConnectionInput) (*request.Request, *ec2.CreateVpcPeeringConnectionOutput) {
	if f.CreateVpcPeeringConnectionRequestFn == nil {
		return nil, nil
	}
	return f.CreateVpcPeeringConnectionRequestFn(input)
}

func (f *EC2) CreateVpcPeeringConnection(input *ec2.CreateVpcPeeringConnectionInput) (*ec2.CreateVpcPeeringConnectionOutput, error) {
	if f.CreateVpcPeeringConnectionFn == nil {
		return nil, nil
	}
	return f.CreateVpcPeeringConnectionFn(input)
}

func (f *EC2) CreateVpnConnectionRequest(input *ec2.CreateVpnConnectionInput) (*request.Request, *ec2.CreateVpnConnectionOutput) {
	if f.CreateVpnConnectionRequestFn == nil {
		return nil, nil
	}
	return f.CreateVpnConnectionRequestFn(input)
}

func (f *EC2) CreateVpnConnection(input *ec2.CreateVpnConnectionInput) (*ec2.CreateVpnConnectionOutput, error) {
	if f.CreateVpnConnectionFn == nil {
		return nil, nil
	}
	return f.CreateVpnConnectionFn(input)
}

func (f *EC2) CreateVpnConnectionRouteRequest(input *ec2.CreateVpnConnectionRouteInput) (*request.Request, *ec2.CreateVpnConnectionRouteOutput) {
	if f.CreateVpnConnectionRouteRequestFn == nil {
		return nil, nil
	}
	return f.CreateVpnConnectionRouteRequestFn(input)
}

func (f *EC2) CreateVpnConnectionRoute(input *ec2.CreateVpnConnectionRouteInput) (*ec2.CreateVpnConnectionRouteOutput, error) {
	if f.CreateVpnConnectionRouteFn == nil {
		return nil, nil
	}
	return f.CreateVpnConnectionRouteFn(input)
}

func (f *EC2) CreateVpnGatewayRequest(input *ec2.CreateVpnGatewayInput) (*request.Request, *ec2.CreateVpnGatewayOutput) {
	if f.CreateVpnGatewayRequestFn == nil {
		return nil, nil
	}
	return f.CreateVpnGatewayRequestFn(input)
}

func (f *EC2) CreateVpnGateway(input *ec2.CreateVpnGatewayInput) (*ec2.CreateVpnGatewayOutput, error) {
	if f.CreateVpnGatewayFn == nil {
		return nil, nil
	}
	return f.CreateVpnGatewayFn(input)
}

func (f *EC2) DeleteCustomerGatewayRequest(input *ec2.DeleteCustomerGatewayInput) (*request.Request, *ec2.DeleteCustomerGatewayOutput) {
	if f.DeleteCustomerGatewayRequestFn == nil {
		return nil, nil
	}
	return f.DeleteCustomerGatewayRequestFn(input)
}

func (f *EC2) DeleteCustomerGateway(input *ec2.DeleteCustomerGatewayInput) (*ec2.DeleteCustomerGatewayOutput, error) {
	if f.DeleteCustomerGatewayFn == nil {
		return nil, nil
	}
	return f.DeleteCustomerGatewayFn(input)
}

func (f *EC2) DeleteDhcpOptionsRequest(input *ec2.DeleteDhcpOptionsInput) (*request.Request, *ec2.DeleteDhcpOptionsOutput) {
	if f.DeleteDhcpOptionsRequestFn == nil {
		return nil, nil
	}
	return f.DeleteDhcpOptionsRequestFn(input)
}

func (f *EC2) DeleteDhcpOptions(input *ec2.DeleteDhcpOptionsInput) (*ec2.DeleteDhcpOptionsOutput, error) {
	if f.DeleteDhcpOptionsFn == nil {
		return nil, nil
	}
	return f.DeleteDhcpOptionsFn(input)
}

func (f *EC2) DeleteFlowLogsRequest(input *ec2.DeleteFlowLogsInput) (*request.Request, *ec2.DeleteFlowLogsOutput) {
	if f.DeleteFlowLogsRequestFn == nil {
		return nil, nil
	}
	return f.DeleteFlowLogsRequestFn(input)
}

func (f *EC2) DeleteFlowLogs(input *ec2.DeleteFlowLogsInput) (*ec2.DeleteFlowLogsOutput, error) {
	if f.DeleteFlowLogsFn == nil {
		return nil, nil
	}
	return f.DeleteFlowLogsFn(input)
}

func (f *EC2) DeleteInternetGatewayRequest(input *ec2.DeleteInternetGatewayInput) (*request.Request, *ec2.DeleteInternetGatewayOutput) {
	if f.DeleteInternetGatewayRequestFn == nil {
		return nil, nil
	}
	return f.DeleteInternetGatewayRequestFn(input)
}

func (f *EC2) DeleteInternetGateway(input *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
	if f.DeleteInternetGatewayFn == nil {
		return nil, nil
	}
	return f.DeleteInternetGatewayFn(input)
}

func (f *EC2) DeleteKeyPairRequest(input *ec2.DeleteKeyPairInput) (*request.Request, *ec2.DeleteKeyPairOutput) {
	if f.DeleteKeyPairRequestFn == nil {
		return nil, nil
	}
	return f.DeleteKeyPairRequestFn(input)
}

func (f *EC2) DeleteKeyPair(input *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
	if f.DeleteKeyPairFn == nil {
		return nil, nil
	}
	return f.DeleteKeyPairFn(input)
}

func (f *EC2) DeleteNatGatewayRequest(input *ec2.DeleteNatGatewayInput) (*request.Request, *ec2.DeleteNatGatewayOutput) {
	if f.DeleteNatGatewayRequestFn == nil {
		return nil, nil
	}
	return f.DeleteNatGatewayRequestFn(input)
}

func (f *EC2) DeleteNatGateway(input *ec2.DeleteNatGatewayInput) (*ec2.DeleteNatGatewayOutput, error) {
	if f.DeleteNatGatewayFn == nil {
		return nil, nil
	}
	return f.DeleteNatGatewayFn(input)
}

func (f *EC2) DeleteNetworkAclRequest(input *ec2.DeleteNetworkAclInput) (*request.Request, *ec2.DeleteNetworkAclOutput) {
	if f.DeleteNetworkAclRequestFn == nil {
		return nil, nil
	}
	return f.DeleteNetworkAclRequestFn(input)
}

func (f *EC2) DeleteNetworkAcl(input *ec2.DeleteNetworkAclInput) (*ec2.DeleteNetworkAclOutput, error) {
	if f.DeleteNetworkAclFn == nil {
		return nil, nil
	}
	return f.DeleteNetworkAclFn(input)
}

func (f *EC2) DeleteNetworkAclEntryRequest(input *ec2.DeleteNetworkAclEntryInput) (*request.Request, *ec2.DeleteNetworkAclEntryOutput) {
	if f.DeleteNetworkAclEntryRequestFn == nil {
		return nil, nil
	}
	return f.DeleteNetworkAclEntryRequestFn(input)
}

func (f *EC2) DeleteNetworkAclEntry(input *ec2.DeleteNetworkAclEntryInput) (*ec2.DeleteNetworkAclEntryOutput, error) {
	if f.DeleteNetworkAclEntryFn == nil {
		return nil, nil
	}
	return f.DeleteNetworkAclEntryFn(input)
}

func (f *EC2) DeleteNetworkInterfaceRequest(input *ec2.DeleteNetworkInterfaceInput) (*request.Request, *ec2.DeleteNetworkInterfaceOutput) {
	if f.DeleteNetworkInterfaceRequestFn == nil {
		return nil, nil
	}
	return f.DeleteNetworkInterfaceRequestFn(input)
}

func (f *EC2) DeleteNetworkInterface(input *ec2.DeleteNetworkInterfaceInput) (*ec2.DeleteNetworkInterfaceOutput, error) {
	if f.DeleteNetworkInterfaceFn == nil {
		return nil, nil
	}
	return f.DeleteNetworkInterfaceFn(input)
}

func (f *EC2) DeletePlacementGroupRequest(input *ec2.DeletePlacementGroupInput) (*request.Request, *ec2.DeletePlacementGroupOutput) {
	if f.DeletePlacementGroupRequestFn == nil {
		return nil, nil
	}
	return f.DeletePlacementGroupRequestFn(input)
}

func (f *EC2) DeletePlacementGroup(input *ec2.DeletePlacementGroupInput) (*ec2.DeletePlacementGroupOutput, error) {
	if f.DeletePlacementGroupFn == nil {
		return nil, nil
	}
	return f.DeletePlacementGroupFn(input)
}

func (f *EC2) DeleteRouteRequest(input *ec2.DeleteRouteInput) (*request.Request, *ec2.DeleteRouteOutput) {
	if f.DeleteRouteRequestFn == nil {
		return nil, nil
	}
	return f.DeleteRouteRequestFn(input)
}

func (f *EC2) DeleteRoute(input *ec2.DeleteRouteInput) (*ec2.DeleteRouteOutput, error) {
	if f.DeleteRouteFn == nil {
		return nil, nil
	}
	return f.DeleteRouteFn(input)
}

func (f *EC2) DeleteRouteTableRequest(input *ec2.DeleteRouteTableInput) (*request.Request, *ec2.DeleteRouteTableOutput) {
	if f.DeleteRouteTableRequestFn == nil {
		return nil, nil
	}
	return f.DeleteRouteTableRequestFn(input)
}

func (f *EC2) DeleteRouteTable(input *ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
	if f.DeleteRouteTableFn == nil {
		return nil, nil
	}
	return f.DeleteRouteTableFn(input)
}

func (f *EC2) DeleteSecurityGroupRequest(input *ec2.DeleteSecurityGroupInput) (*request.Request, *ec2.DeleteSecurityGroupOutput) {
	if f.DeleteSecurityGroupRequestFn == nil {
		return nil, nil
	}
	return f.DeleteSecurityGroupRequestFn(input)
}

func (f *EC2) DeleteSecurityGroup(input *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	if f.DeleteSecurityGroupFn == nil {
		return nil, nil
	}
	return f.DeleteSecurityGroupFn(input)
}

func (f *EC2) DeleteSnapshotRequest(input *ec2.DeleteSnapshotInput) (*request.Request, *ec2.DeleteSnapshotOutput) {
	if f.DeleteSnapshotRequestFn == nil {
		return nil, nil
	}
	return f.DeleteSnapshotRequestFn(input)
}

func (f *EC2) DeleteSnapshot(input *ec2.DeleteSnapshotInput) (*ec2.DeleteSnapshotOutput, error) {
	if f.DeleteSnapshotFn == nil {
		return nil, nil
	}
	return f.DeleteSnapshotFn(input)
}

func (f *EC2) DeleteSpotDatafeedSubscriptionRequest(input *ec2.DeleteSpotDatafeedSubscriptionInput) (*request.Request, *ec2.DeleteSpotDatafeedSubscriptionOutput) {
	if f.DeleteSpotDatafeedSubscriptionRequestFn == nil {
		return nil, nil
	}
	return f.DeleteSpotDatafeedSubscriptionRequestFn(input)
}

func (f *EC2) DeleteSpotDatafeedSubscription(input *ec2.DeleteSpotDatafeedSubscriptionInput) (*ec2.DeleteSpotDatafeedSubscriptionOutput, error) {
	if f.DeleteSpotDatafeedSubscriptionFn == nil {
		return nil, nil
	}
	return f.DeleteSpotDatafeedSubscriptionFn(input)
}

func (f *EC2) DeleteSubnetRequest(input *ec2.DeleteSubnetInput) (*request.Request, *ec2.DeleteSubnetOutput) {
	if f.DeleteSubnetRequestFn == nil {
		return nil, nil
	}
	return f.DeleteSubnetRequestFn(input)
}

func (f *EC2) DeleteSubnet(input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
	if f.DeleteSubnetFn == nil {
		return nil, nil
	}
	return f.DeleteSubnetFn(input)
}

func (f *EC2) DeleteTagsRequest(input *ec2.DeleteTagsInput) (*request.Request, *ec2.DeleteTagsOutput) {
	if f.DeleteTagsRequestFn == nil {
		return nil, nil
	}
	return f.DeleteTagsRequestFn(input)
}

func (f *EC2) DeleteTags(input *ec2.DeleteTagsInput) (*ec2.DeleteTagsOutput, error) {
	if f.DeleteTagsFn == nil {
		return nil, nil
	}
	return f.DeleteTagsFn(input)
}

func (f *EC2) DeleteVolumeRequest(input *ec2.DeleteVolumeInput) (*request.Request, *ec2.DeleteVolumeOutput) {
	if f.DeleteVolumeRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVolumeRequestFn(input)
}

func (f *EC2) DeleteVolume(input *ec2.DeleteVolumeInput) (*ec2.DeleteVolumeOutput, error) {
	if f.DeleteVolumeFn == nil {
		return nil, nil
	}
	return f.DeleteVolumeFn(input)
}

func (f *EC2) DeleteVpcRequest(input *ec2.DeleteVpcInput) (*request.Request, *ec2.DeleteVpcOutput) {
	if f.DeleteVpcRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVpcRequestFn(input)
}

func (f *EC2) DeleteVpc(input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
	if f.DeleteVpcFn == nil {
		return nil, nil
	}
	return f.DeleteVpcFn(input)
}

func (f *EC2) DeleteVpcEndpointsRequest(input *ec2.DeleteVpcEndpointsInput) (*request.Request, *ec2.DeleteVpcEndpointsOutput) {
	if f.DeleteVpcEndpointsRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVpcEndpointsRequestFn(input)
}

func (f *EC2) DeleteVpcEndpoints(input *ec2.DeleteVpcEndpointsInput) (*ec2.DeleteVpcEndpointsOutput, error) {
	if f.DeleteVpcEndpointsFn == nil {
		return nil, nil
	}
	return f.DeleteVpcEndpointsFn(input)
}

func (f *EC2) DeleteVpcPeeringConnectionRequest(input *ec2.DeleteVpcPeeringConnectionInput) (*request.Request, *ec2.DeleteVpcPeeringConnectionOutput) {
	if f.DeleteVpcPeeringConnectionRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVpcPeeringConnectionRequestFn(input)
}

func (f *EC2) DeleteVpcPeeringConnection(input *ec2.DeleteVpcPeeringConnectionInput) (*ec2.DeleteVpcPeeringConnectionOutput, error) {
	if f.DeleteVpcPeeringConnectionFn == nil {
		return nil, nil
	}
	return f.DeleteVpcPeeringConnectionFn(input)
}

func (f *EC2) DeleteVpnConnectionRequest(input *ec2.DeleteVpnConnectionInput) (*request.Request, *ec2.DeleteVpnConnectionOutput) {
	if f.DeleteVpnConnectionRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVpnConnectionRequestFn(input)
}

func (f *EC2) DeleteVpnConnection(input *ec2.DeleteVpnConnectionInput) (*ec2.DeleteVpnConnectionOutput, error) {
	if f.DeleteVpnConnectionFn == nil {
		return nil, nil
	}
	return f.DeleteVpnConnectionFn(input)
}

func (f *EC2) DeleteVpnConnectionRouteRequest(input *ec2.DeleteVpnConnectionRouteInput) (*request.Request, *ec2.DeleteVpnConnectionRouteOutput) {
	if f.DeleteVpnConnectionRouteRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVpnConnectionRouteRequestFn(input)
}

func (f *EC2) DeleteVpnConnectionRoute(input *ec2.DeleteVpnConnectionRouteInput) (*ec2.DeleteVpnConnectionRouteOutput, error) {
	if f.DeleteVpnConnectionRouteFn == nil {
		return nil, nil
	}
	return f.DeleteVpnConnectionRouteFn(input)
}

func (f *EC2) DeleteVpnGatewayRequest(input *ec2.DeleteVpnGatewayInput) (*request.Request, *ec2.DeleteVpnGatewayOutput) {
	if f.DeleteVpnGatewayRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVpnGatewayRequestFn(input)
}

func (f *EC2) DeleteVpnGateway(input *ec2.DeleteVpnGatewayInput) (*ec2.DeleteVpnGatewayOutput, error) {
	if f.DeleteVpnGatewayFn == nil {
		return nil, nil
	}
	return f.DeleteVpnGatewayFn(input)
}

func (f *EC2) DeregisterImageRequest(input *ec2.DeregisterImageInput) (*request.Request, *ec2.DeregisterImageOutput) {
	if f.DeregisterImageRequestFn == nil {
		return nil, nil
	}
	return f.DeregisterImageRequestFn(input)
}

func (f *EC2) DeregisterImage(input *ec2.DeregisterImageInput) (*ec2.DeregisterImageOutput, error) {
	if f.DeregisterImageFn == nil {
		return nil, nil
	}
	return f.DeregisterImageFn(input)
}

func (f *EC2) DescribeAccountAttributesRequest(input *ec2.DescribeAccountAttributesInput) (*request.Request, *ec2.DescribeAccountAttributesOutput) {
	if f.DescribeAccountAttributesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeAccountAttributesRequestFn(input)
}

func (f *EC2) DescribeAccountAttributes(input *ec2.DescribeAccountAttributesInput) (*ec2.DescribeAccountAttributesOutput, error) {
	if f.DescribeAccountAttributesFn == nil {
		return nil, nil
	}
	return f.DescribeAccountAttributesFn(input)
}

func (f *EC2) DescribeAddressesRequest(input *ec2.DescribeAddressesInput) (*request.Request, *ec2.DescribeAddressesOutput) {
	if f.DescribeAddressesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeAddressesRequestFn(input)
}

func (f *EC2) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	if f.DescribeAddressesFn == nil {
		return nil, nil
	}
	return f.DescribeAddressesFn(input)
}

func (f *EC2) DescribeAvailabilityZonesRequest(input *ec2.DescribeAvailabilityZonesInput) (*request.Request, *ec2.DescribeAvailabilityZonesOutput) {
	if f.DescribeAvailabilityZonesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeAvailabilityZonesRequestFn(input)
}

func (f *EC2) DescribeAvailabilityZones(input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	if f.DescribeAvailabilityZonesFn == nil {
		return nil, nil
	}
	return f.DescribeAvailabilityZonesFn(input)
}

func (f *EC2) DescribeBundleTasksRequest(input *ec2.DescribeBundleTasksInput) (*request.Request, *ec2.DescribeBundleTasksOutput) {
	if f.DescribeBundleTasksRequestFn == nil {
		return nil, nil
	}
	return f.DescribeBundleTasksRequestFn(input)
}

func (f *EC2) DescribeBundleTasks(input *ec2.DescribeBundleTasksInput) (*ec2.DescribeBundleTasksOutput, error) {
	if f.DescribeBundleTasksFn == nil {
		return nil, nil
	}
	return f.DescribeBundleTasksFn(input)
}

func (f *EC2) DescribeClassicLinkInstancesRequest(input *ec2.DescribeClassicLinkInstancesInput) (*request.Request, *ec2.DescribeClassicLinkInstancesOutput) {
	if f.DescribeClassicLinkInstancesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeClassicLinkInstancesRequestFn(input)
}

func (f *EC2) DescribeClassicLinkInstances(input *ec2.DescribeClassicLinkInstancesInput) (*ec2.DescribeClassicLinkInstancesOutput, error) {
	if f.DescribeClassicLinkInstancesFn == nil {
		return nil, nil
	}
	return f.DescribeClassicLinkInstancesFn(input)
}

func (f *EC2) DescribeConversionTasksRequest(input *ec2.DescribeConversionTasksInput) (*request.Request, *ec2.DescribeConversionTasksOutput) {
	if f.DescribeConversionTasksRequestFn == nil {
		return nil, nil
	}
	return f.DescribeConversionTasksRequestFn(input)
}

func (f *EC2) DescribeConversionTasks(input *ec2.DescribeConversionTasksInput) (*ec2.DescribeConversionTasksOutput, error) {
	if f.DescribeConversionTasksFn == nil {
		return nil, nil
	}
	return f.DescribeConversionTasksFn(input)
}

func (f *EC2) DescribeCustomerGatewaysRequest(input *ec2.DescribeCustomerGatewaysInput) (*request.Request, *ec2.DescribeCustomerGatewaysOutput) {
	if f.DescribeCustomerGatewaysRequestFn == nil {
		return nil, nil
	}
	return f.DescribeCustomerGatewaysRequestFn(input)
}

func (f *EC2) DescribeCustomerGateways(input *ec2.DescribeCustomerGatewaysInput) (*ec2.DescribeCustomerGatewaysOutput, error) {
	if f.DescribeCustomerGatewaysFn == nil {
		return nil, nil
	}
	return f.DescribeCustomerGatewaysFn(input)
}

func (f *EC2) DescribeDhcpOptionsRequest(input *ec2.DescribeDhcpOptionsInput) (*request.Request, *ec2.DescribeDhcpOptionsOutput) {
	if f.DescribeDhcpOptionsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeDhcpOptionsRequestFn(input)
}

func (f *EC2) DescribeDhcpOptions(input *ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
	if f.DescribeDhcpOptionsFn == nil {
		return nil, nil
	}
	return f.DescribeDhcpOptionsFn(input)
}

func (f *EC2) DescribeExportTasksRequest(input *ec2.DescribeExportTasksInput) (*request.Request, *ec2.DescribeExportTasksOutput) {
	if f.DescribeExportTasksRequestFn == nil {
		return nil, nil
	}
	return f.DescribeExportTasksRequestFn(input)
}

func (f *EC2) DescribeExportTasks(input *ec2.DescribeExportTasksInput) (*ec2.DescribeExportTasksOutput, error) {
	if f.DescribeExportTasksFn == nil {
		return nil, nil
	}
	return f.DescribeExportTasksFn(input)
}

func (f *EC2) DescribeFlowLogsRequest(input *ec2.DescribeFlowLogsInput) (*request.Request, *ec2.DescribeFlowLogsOutput) {
	if f.DescribeFlowLogsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeFlowLogsRequestFn(input)
}

func (f *EC2) DescribeFlowLogs(input *ec2.DescribeFlowLogsInput) (*ec2.DescribeFlowLogsOutput, error) {
	if f.DescribeFlowLogsFn == nil {
		return nil, nil
	}
	return f.DescribeFlowLogsFn(input)
}

func (f *EC2) DescribeHostsRequest(input *ec2.DescribeHostsInput) (*request.Request, *ec2.DescribeHostsOutput) {
	if f.DescribeHostsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeHostsRequestFn(input)
}

func (f *EC2) DescribeHosts(input *ec2.DescribeHostsInput) (*ec2.DescribeHostsOutput, error) {
	if f.DescribeHostsFn == nil {
		return nil, nil
	}
	return f.DescribeHostsFn(input)
}

func (f *EC2) DescribeIdFormatRequest(input *ec2.DescribeIdFormatInput) (*request.Request, *ec2.DescribeIdFormatOutput) {
	if f.DescribeIdFormatRequestFn == nil {
		return nil, nil
	}
	return f.DescribeIdFormatRequestFn(input)
}

func (f *EC2) DescribeIdFormat(input *ec2.DescribeIdFormatInput) (*ec2.DescribeIdFormatOutput, error) {
	if f.DescribeIdFormatFn == nil {
		return nil, nil
	}
	return f.DescribeIdFormatFn(input)
}

func (f *EC2) DescribeImageAttributeRequest(input *ec2.DescribeImageAttributeInput) (*request.Request, *ec2.DescribeImageAttributeOutput) {
	if f.DescribeImageAttributeRequestFn == nil {
		return nil, nil
	}
	return f.DescribeImageAttributeRequestFn(input)
}

func (f *EC2) DescribeImageAttribute(input *ec2.DescribeImageAttributeInput) (*ec2.DescribeImageAttributeOutput, error) {
	if f.DescribeImageAttributeFn == nil {
		return nil, nil
	}
	return f.DescribeImageAttributeFn(input)
}

func (f *EC2) DescribeImagesRequest(input *ec2.DescribeImagesInput) (*request.Request, *ec2.DescribeImagesOutput) {
	if f.DescribeImagesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeImagesRequestFn(input)
}

func (f *EC2) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	if f.DescribeImagesFn == nil {
		return nil, nil
	}
	return f.DescribeImagesFn(input)
}

func (f *EC2) DescribeImportImageTasksRequest(input *ec2.DescribeImportImageTasksInput) (*request.Request, *ec2.DescribeImportImageTasksOutput) {
	if f.DescribeImportImageTasksRequestFn == nil {
		return nil, nil
	}
	return f.DescribeImportImageTasksRequestFn(input)
}

func (f *EC2) DescribeImportImageTasks(input *ec2.DescribeImportImageTasksInput) (*ec2.DescribeImportImageTasksOutput, error) {
	if f.DescribeImportImageTasksFn == nil {
		return nil, nil
	}
	return f.DescribeImportImageTasksFn(input)
}

func (f *EC2) DescribeImportSnapshotTasksRequest(input *ec2.DescribeImportSnapshotTasksInput) (*request.Request, *ec2.DescribeImportSnapshotTasksOutput) {
	if f.DescribeImportSnapshotTasksRequestFn == nil {
		return nil, nil
	}
	return f.DescribeImportSnapshotTasksRequestFn(input)
}

func (f *EC2) DescribeImportSnapshotTasks(input *ec2.DescribeImportSnapshotTasksInput) (*ec2.DescribeImportSnapshotTasksOutput, error) {
	if f.DescribeImportSnapshotTasksFn == nil {
		return nil, nil
	}
	return f.DescribeImportSnapshotTasksFn(input)
}

func (f *EC2) DescribeInstanceAttributeRequest(input *ec2.DescribeInstanceAttributeInput) (*request.Request, *ec2.DescribeInstanceAttributeOutput) {
	if f.DescribeInstanceAttributeRequestFn == nil {
		return nil, nil
	}
	return f.DescribeInstanceAttributeRequestFn(input)
}

func (f *EC2) DescribeInstanceAttribute(input *ec2.DescribeInstanceAttributeInput) (*ec2.DescribeInstanceAttributeOutput, error) {
	if f.DescribeInstanceAttributeFn == nil {
		return nil, nil
	}
	return f.DescribeInstanceAttributeFn(input)
}

func (f *EC2) DescribeInstanceStatusRequest(input *ec2.DescribeInstanceStatusInput) (*request.Request, *ec2.DescribeInstanceStatusOutput) {
	if f.DescribeInstanceStatusRequestFn == nil {
		return nil, nil
	}
	return f.DescribeInstanceStatusRequestFn(input)
}

func (f *EC2) DescribeInstanceStatus(input *ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error) {
	if f.DescribeInstanceStatusFn == nil {
		return nil, nil
	}
	return f.DescribeInstanceStatusFn(input)
}

func (f *EC2) DescribeInstanceStatusPages(input *ec2.DescribeInstanceStatusInput, fn func(*ec2.DescribeInstanceStatusOutput, bool) bool) error {
	if f.DescribeInstanceStatusPagesFn == nil {
		return nil
	}
	return f.DescribeInstanceStatusPagesFn(input, fn)
}

func (f *EC2) DescribeInstancesRequest(input *ec2.DescribeInstancesInput) (*request.Request, *ec2.DescribeInstancesOutput) {
	if f.DescribeInstancesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeInstancesRequestFn(input)
}

func (f *EC2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if f.DescribeInstancesFn == nil {
		return nil, nil
	}
	return f.DescribeInstancesFn(input)
}

func (f *EC2) DescribeInstancesPages(input *ec2.DescribeInstancesInput, fn func(*ec2.DescribeInstancesOutput, bool) bool) error {
	if f.DescribeInstancesPagesFn == nil {
		return nil
	}
	return f.DescribeInstancesPagesFn(input, fn)
}

func (f *EC2) DescribeInternetGatewaysRequest(input *ec2.DescribeInternetGatewaysInput) (*request.Request, *ec2.DescribeInternetGatewaysOutput) {
	if f.DescribeInternetGatewaysRequestFn == nil {
		return nil, nil
	}
	return f.DescribeInternetGatewaysRequestFn(input)
}

func (f *EC2) DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	if f.DescribeInternetGatewaysFn == nil {
		return nil, nil
	}
	return f.DescribeInternetGatewaysFn(input)
}

func (f *EC2) DescribeKeyPairsRequest(input *ec2.DescribeKeyPairsInput) (*request.Request, *ec2.DescribeKeyPairsOutput) {
	if f.DescribeKeyPairsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeKeyPairsRequestFn(input)
}

func (f *EC2) DescribeKeyPairs(input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	if f.DescribeKeyPairsFn == nil {
		return nil, nil
	}
	return f.DescribeKeyPairsFn(input)
}

func (f *EC2) DescribeMovingAddressesRequest(input *ec2.DescribeMovingAddressesInput) (*request.Request, *ec2.DescribeMovingAddressesOutput) {
	if f.DescribeMovingAddressesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeMovingAddressesRequestFn(input)
}

func (f *EC2) DescribeMovingAddresses(input *ec2.DescribeMovingAddressesInput) (*ec2.DescribeMovingAddressesOutput, error) {
	if f.DescribeMovingAddressesFn == nil {
		return nil, nil
	}
	return f.DescribeMovingAddressesFn(input)
}

func (f *EC2) DescribeNatGatewaysRequest(input *ec2.DescribeNatGatewaysInput) (*request.Request, *ec2.DescribeNatGatewaysOutput) {
	if f.DescribeNatGatewaysRequestFn == nil {
		return nil, nil
	}
	return f.DescribeNatGatewaysRequestFn(input)
}

func (f *EC2) DescribeNatGateways(input *ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error) {
	if f.DescribeNatGatewaysFn == nil {
		return nil, nil
	}
	return f.DescribeNatGatewaysFn(input)
}

func (f *EC2) DescribeNetworkAclsRequest(input *ec2.DescribeNetworkAclsInput) (*request.Request, *ec2.DescribeNetworkAclsOutput) {
	if f.DescribeNetworkAclsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeNetworkAclsRequestFn(input)
}

func (f *EC2) DescribeNetworkAcls(input *ec2.DescribeNetworkAclsInput) (*ec2.DescribeNetworkAclsOutput, error) {
	if f.DescribeNetworkAclsFn == nil {
		return nil, nil
	}
	return f.DescribeNetworkAclsFn(input)
}

func (f *EC2) DescribeNetworkInterfaceAttributeRequest(input *ec2.DescribeNetworkInterfaceAttributeInput) (*request.Request, *ec2.DescribeNetworkInterfaceAttributeOutput) {
	if f.DescribeNetworkInterfaceAttributeRequestFn == nil {
		return nil, nil
	}
	return f.DescribeNetworkInterfaceAttributeRequestFn(input)
}

func (f *EC2) DescribeNetworkInterfaceAttribute(input *ec2.DescribeNetworkInterfaceAttributeInput) (*ec2.DescribeNetworkInterfaceAttributeOutput, error) {
	if f.DescribeNetworkInterfaceAttributeFn == nil {
		return nil, nil
	}
	return f.DescribeNetworkInterfaceAttributeFn(input)
}

func (f *EC2) DescribeNetworkInterfacesRequest(input *ec2.DescribeNetworkInterfacesInput) (*request.Request, *ec2.DescribeNetworkInterfacesOutput) {
	if f.DescribeNetworkInterfacesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeNetworkInterfacesRequestFn(input)
}

func (f *EC2) DescribeNetworkInterfaces(input *ec2.DescribeNetworkInterfacesInput) (*ec2.DescribeNetworkInterfacesOutput, error) {
	if f.DescribeNetworkInterfacesFn == nil {
		return nil, nil
	}
	return f.DescribeNetworkInterfacesFn(input)
}

func (f *EC2) DescribePlacementGroupsRequest(input *ec2.DescribePlacementGroupsInput) (*request.Request, *ec2.DescribePlacementGroupsOutput) {
	if f.DescribePlacementGroupsRequestFn == nil {
		return nil, nil
	}
	return f.DescribePlacementGroupsRequestFn(input)
}

func (f *EC2) DescribePlacementGroups(input *ec2.DescribePlacementGroupsInput) (*ec2.DescribePlacementGroupsOutput, error) {
	if f.DescribePlacementGroupsFn == nil {
		return nil, nil
	}
	return f.DescribePlacementGroupsFn(input)
}

func (f *EC2) DescribePrefixListsRequest(input *ec2.DescribePrefixListsInput) (*request.Request, *ec2.DescribePrefixListsOutput) {
	if f.DescribePrefixListsRequestFn == nil {
		return nil, nil
	}
	return f.DescribePrefixListsRequestFn(input)
}

func (f *EC2) DescribePrefixLists(input *ec2.DescribePrefixListsInput) (*ec2.DescribePrefixListsOutput, error) {
	if f.DescribePrefixListsFn == nil {
		return nil, nil
	}
	return f.DescribePrefixListsFn(input)
}

func (f *EC2) DescribeRegionsRequest(input *ec2.DescribeRegionsInput) (*request.Request, *ec2.DescribeRegionsOutput) {
	if f.DescribeRegionsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeRegionsRequestFn(input)
}

func (f *EC2) DescribeRegions(input *ec2.DescribeRegionsInput) (*ec2.DescribeRegionsOutput, error) {
	if f.DescribeRegionsFn == nil {
		return nil, nil
	}
	return f.DescribeRegionsFn(input)
}

func (f *EC2) DescribeReservedInstancesRequest(input *ec2.DescribeReservedInstancesInput) (*request.Request, *ec2.DescribeReservedInstancesOutput) {
	if f.DescribeReservedInstancesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesRequestFn(input)
}

func (f *EC2) DescribeReservedInstances(input *ec2.DescribeReservedInstancesInput) (*ec2.DescribeReservedInstancesOutput, error) {
	if f.DescribeReservedInstancesFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesFn(input)
}

func (f *EC2) DescribeReservedInstancesListingsRequest(input *ec2.DescribeReservedInstancesListingsInput) (*request.Request, *ec2.DescribeReservedInstancesListingsOutput) {
	if f.DescribeReservedInstancesListingsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesListingsRequestFn(input)
}

func (f *EC2) DescribeReservedInstancesListings(input *ec2.DescribeReservedInstancesListingsInput) (*ec2.DescribeReservedInstancesListingsOutput, error) {
	if f.DescribeReservedInstancesListingsFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesListingsFn(input)
}

func (f *EC2) DescribeReservedInstancesModificationsRequest(input *ec2.DescribeReservedInstancesModificationsInput) (*request.Request, *ec2.DescribeReservedInstancesModificationsOutput) {
	if f.DescribeReservedInstancesModificationsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesModificationsRequestFn(input)
}

func (f *EC2) DescribeReservedInstancesModifications(input *ec2.DescribeReservedInstancesModificationsInput) (*ec2.DescribeReservedInstancesModificationsOutput, error) {
	if f.DescribeReservedInstancesModificationsFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesModificationsFn(input)
}

func (f *EC2) DescribeReservedInstancesModificationsPages(input *ec2.DescribeReservedInstancesModificationsInput, fn func(*ec2.DescribeReservedInstancesModificationsOutput, bool) bool) error {
	if f.DescribeReservedInstancesModificationsPagesFn == nil {
		return nil
	}
	return f.DescribeReservedInstancesModificationsPagesFn(input, fn)
}

func (f *EC2) DescribeReservedInstancesOfferingsRequest(input *ec2.DescribeReservedInstancesOfferingsInput) (*request.Request, *ec2.DescribeReservedInstancesOfferingsOutput) {
	if f.DescribeReservedInstancesOfferingsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesOfferingsRequestFn(input)
}

func (f *EC2) DescribeReservedInstancesOfferings(input *ec2.DescribeReservedInstancesOfferingsInput) (*ec2.DescribeReservedInstancesOfferingsOutput, error) {
	if f.DescribeReservedInstancesOfferingsFn == nil {
		return nil, nil
	}
	return f.DescribeReservedInstancesOfferingsFn(input)
}

func (f *EC2) DescribeReservedInstancesOfferingsPages(input *ec2.DescribeReservedInstancesOfferingsInput, fn func(*ec2.DescribeReservedInstancesOfferingsOutput, bool) bool) error {
	if f.DescribeReservedInstancesOfferingsPagesFn == nil {
		return nil
	}
	return f.DescribeReservedInstancesOfferingsPagesFn(input, fn)
}

func (f *EC2) DescribeRouteTablesRequest(input *ec2.DescribeRouteTablesInput) (*request.Request, *ec2.DescribeRouteTablesOutput) {
	if f.DescribeRouteTablesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeRouteTablesRequestFn(input)
}

func (f *EC2) DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	if f.DescribeRouteTablesFn == nil {
		return nil, nil
	}
	return f.DescribeRouteTablesFn(input)
}

func (f *EC2) DescribeScheduledInstanceAvailabilityRequest(input *ec2.DescribeScheduledInstanceAvailabilityInput) (*request.Request, *ec2.DescribeScheduledInstanceAvailabilityOutput) {
	if f.DescribeScheduledInstanceAvailabilityRequestFn == nil {
		return nil, nil
	}
	return f.DescribeScheduledInstanceAvailabilityRequestFn(input)
}

func (f *EC2) DescribeScheduledInstanceAvailability(input *ec2.DescribeScheduledInstanceAvailabilityInput) (*ec2.DescribeScheduledInstanceAvailabilityOutput, error) {
	if f.DescribeScheduledInstanceAvailabilityFn == nil {
		return nil, nil
	}
	return f.DescribeScheduledInstanceAvailabilityFn(input)
}

func (f *EC2) DescribeScheduledInstancesRequest(input *ec2.DescribeScheduledInstancesInput) (*request.Request, *ec2.DescribeScheduledInstancesOutput) {
	if f.DescribeScheduledInstancesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeScheduledInstancesRequestFn(input)
}

func (f *EC2) DescribeScheduledInstances(input *ec2.DescribeScheduledInstancesInput) (*ec2.DescribeScheduledInstancesOutput, error) {
	if f.DescribeScheduledInstancesFn == nil {
		return nil, nil
	}
	return f.DescribeScheduledInstancesFn(input)
}

func (f *EC2) DescribeSecurityGroupsRequest(input *ec2.DescribeSecurityGroupsInput) (*request.Request, *ec2.DescribeSecurityGroupsOutput) {
	if f.DescribeSecurityGroupsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSecurityGroupsRequestFn(input)
}

func (f *EC2) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	if f.DescribeSecurityGroupsFn == nil {
		return nil, nil
	}
	return f.DescribeSecurityGroupsFn(input)
}

func (f *EC2) DescribeSnapshotAttributeRequest(input *ec2.DescribeSnapshotAttributeInput) (*request.Request, *ec2.DescribeSnapshotAttributeOutput) {
	if f.DescribeSnapshotAttributeRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSnapshotAttributeRequestFn(input)
}

func (f *EC2) DescribeSnapshotAttribute(input *ec2.DescribeSnapshotAttributeInput) (*ec2.DescribeSnapshotAttributeOutput, error) {
	if f.DescribeSnapshotAttributeFn == nil {
		return nil, nil
	}
	return f.DescribeSnapshotAttributeFn(input)
}

func (f *EC2) DescribeSnapshotsRequest(input *ec2.DescribeSnapshotsInput) (*request.Request, *ec2.DescribeSnapshotsOutput) {
	if f.DescribeSnapshotsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSnapshotsRequestFn(input)
}

func (f *EC2) DescribeSnapshots(input *ec2.DescribeSnapshotsInput) (*ec2.DescribeSnapshotsOutput, error) {
	if f.DescribeSnapshotsFn == nil {
		return nil, nil
	}
	return f.DescribeSnapshotsFn(input)
}

func (f *EC2) DescribeSnapshotsPages(input *ec2.DescribeSnapshotsInput, fn func(*ec2.DescribeSnapshotsOutput, bool) bool) error {
	if f.DescribeSnapshotsPagesFn == nil {
		return nil
	}
	return f.DescribeSnapshotsPagesFn(input, fn)
}

func (f *EC2) DescribeSpotDatafeedSubscriptionRequest(input *ec2.DescribeSpotDatafeedSubscriptionInput) (*request.Request, *ec2.DescribeSpotDatafeedSubscriptionOutput) {
	if f.DescribeSpotDatafeedSubscriptionRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSpotDatafeedSubscriptionRequestFn(input)
}

func (f *EC2) DescribeSpotDatafeedSubscription(input *ec2.DescribeSpotDatafeedSubscriptionInput) (*ec2.DescribeSpotDatafeedSubscriptionOutput, error) {
	if f.DescribeSpotDatafeedSubscriptionFn == nil {
		return nil, nil
	}
	return f.DescribeSpotDatafeedSubscriptionFn(input)
}

func (f *EC2) DescribeSpotFleetInstancesRequest(input *ec2.DescribeSpotFleetInstancesInput) (*request.Request, *ec2.DescribeSpotFleetInstancesOutput) {
	if f.DescribeSpotFleetInstancesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSpotFleetInstancesRequestFn(input)
}

func (f *EC2) DescribeSpotFleetInstances(input *ec2.DescribeSpotFleetInstancesInput) (*ec2.DescribeSpotFleetInstancesOutput, error) {
	if f.DescribeSpotFleetInstancesFn == nil {
		return nil, nil
	}
	return f.DescribeSpotFleetInstancesFn(input)
}

func (f *EC2) DescribeSpotFleetRequestHistoryRequest(input *ec2.DescribeSpotFleetRequestHistoryInput) (*request.Request, *ec2.DescribeSpotFleetRequestHistoryOutput) {
	if f.DescribeSpotFleetRequestHistoryRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSpotFleetRequestHistoryRequestFn(input)
}

func (f *EC2) DescribeSpotFleetRequestHistory(input *ec2.DescribeSpotFleetRequestHistoryInput) (*ec2.DescribeSpotFleetRequestHistoryOutput, error) {
	if f.DescribeSpotFleetRequestHistoryFn == nil {
		return nil, nil
	}
	return f.DescribeSpotFleetRequestHistoryFn(input)
}

func (f *EC2) DescribeSpotFleetRequestsRequest(input *ec2.DescribeSpotFleetRequestsInput) (*request.Request, *ec2.DescribeSpotFleetRequestsOutput) {
	if f.DescribeSpotFleetRequestsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSpotFleetRequestsRequestFn(input)
}

func (f *EC2) DescribeSpotFleetRequests(input *ec2.DescribeSpotFleetRequestsInput) (*ec2.DescribeSpotFleetRequestsOutput, error) {
	if f.DescribeSpotFleetRequestsFn == nil {
		return nil, nil
	}
	return f.DescribeSpotFleetRequestsFn(input)
}

func (f *EC2) DescribeSpotInstanceRequestsRequest(input *ec2.DescribeSpotInstanceRequestsInput) (*request.Request, *ec2.DescribeSpotInstanceRequestsOutput) {
	if f.DescribeSpotInstanceRequestsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSpotInstanceRequestsRequestFn(input)
}

func (f *EC2) DescribeSpotInstanceRequests(input *ec2.DescribeSpotInstanceRequestsInput) (*ec2.DescribeSpotInstanceRequestsOutput, error) {
	if f.DescribeSpotInstanceRequestsFn == nil {
		return nil, nil
	}
	return f.DescribeSpotInstanceRequestsFn(input)
}

func (f *EC2) DescribeSpotPriceHistoryRequest(input *ec2.DescribeSpotPriceHistoryInput) (*request.Request, *ec2.DescribeSpotPriceHistoryOutput) {
	if f.DescribeSpotPriceHistoryRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSpotPriceHistoryRequestFn(input)
}

func (f *EC2) DescribeSpotPriceHistory(input *ec2.DescribeSpotPriceHistoryInput) (*ec2.DescribeSpotPriceHistoryOutput, error) {
	if f.DescribeSpotPriceHistoryFn == nil {
		return nil, nil
	}
	return f.DescribeSpotPriceHistoryFn(input)
}

func (f *EC2) DescribeSpotPriceHistoryPages(input *ec2.DescribeSpotPriceHistoryInput, fn func(*ec2.DescribeSpotPriceHistoryOutput, bool) bool) error {
	if f.DescribeSpotPriceHistoryPagesFn == nil {
		return nil
	}
	return f.DescribeSpotPriceHistoryPagesFn(input, fn)
}

func (f *EC2) DescribeSubnetsRequest(input *ec2.DescribeSubnetsInput) (*request.Request, *ec2.DescribeSubnetsOutput) {
	if f.DescribeSubnetsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeSubnetsRequestFn(input)
}

func (f *EC2) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	if f.DescribeSubnetsFn == nil {
		return nil, nil
	}
	return f.DescribeSubnetsFn(input)
}

func (f *EC2) DescribeTagsRequest(input *ec2.DescribeTagsInput) (*request.Request, *ec2.DescribeTagsOutput) {
	if f.DescribeTagsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeTagsRequestFn(input)
}

func (f *EC2) DescribeTags(input *ec2.DescribeTagsInput) (*ec2.DescribeTagsOutput, error) {
	if f.DescribeTagsFn == nil {
		return nil, nil
	}
	return f.DescribeTagsFn(input)
}

func (f *EC2) DescribeTagsPages(input *ec2.DescribeTagsInput, fn func(*ec2.DescribeTagsOutput, bool) bool) error {
	if f.DescribeTagsPagesFn == nil {
		return nil
	}
	return f.DescribeTagsPagesFn(input, fn)
}

func (f *EC2) DescribeVolumeAttributeRequest(input *ec2.DescribeVolumeAttributeInput) (*request.Request, *ec2.DescribeVolumeAttributeOutput) {
	if f.DescribeVolumeAttributeRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVolumeAttributeRequestFn(input)
}

func (f *EC2) DescribeVolumeAttribute(input *ec2.DescribeVolumeAttributeInput) (*ec2.DescribeVolumeAttributeOutput, error) {
	if f.DescribeVolumeAttributeFn == nil {
		return nil, nil
	}
	return f.DescribeVolumeAttributeFn(input)
}

func (f *EC2) DescribeVolumeStatusRequest(input *ec2.DescribeVolumeStatusInput) (*request.Request, *ec2.DescribeVolumeStatusOutput) {
	if f.DescribeVolumeStatusRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVolumeStatusRequestFn(input)
}

func (f *EC2) DescribeVolumeStatus(input *ec2.DescribeVolumeStatusInput) (*ec2.DescribeVolumeStatusOutput, error) {
	if f.DescribeVolumeStatusFn == nil {
		return nil, nil
	}
	return f.DescribeVolumeStatusFn(input)
}

func (f *EC2) DescribeVolumeStatusPages(input *ec2.DescribeVolumeStatusInput, fn func(*ec2.DescribeVolumeStatusOutput, bool) bool) error {
	if f.DescribeVolumeStatusPagesFn == nil {
		return nil
	}
	return f.DescribeVolumeStatusPagesFn(input, fn)
}

func (f *EC2) DescribeVolumesRequest(input *ec2.DescribeVolumesInput) (*request.Request, *ec2.DescribeVolumesOutput) {
	if f.DescribeVolumesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVolumesRequestFn(input)
}

func (f *EC2) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	if f.DescribeVolumesFn == nil {
		return nil, nil
	}
	return f.DescribeVolumesFn(input)
}

func (f *EC2) DescribeVolumesPages(input *ec2.DescribeVolumesInput, fn func(*ec2.DescribeVolumesOutput, bool) bool) error {
	if f.DescribeVolumesPagesFn == nil {
		return nil
	}
	return f.DescribeVolumesPagesFn(input, fn)
}

func (f *EC2) DescribeVpcAttributeRequest(input *ec2.DescribeVpcAttributeInput) (*request.Request, *ec2.DescribeVpcAttributeOutput) {
	if f.DescribeVpcAttributeRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpcAttributeRequestFn(input)
}

func (f *EC2) DescribeVpcAttribute(input *ec2.DescribeVpcAttributeInput) (*ec2.DescribeVpcAttributeOutput, error) {
	if f.DescribeVpcAttributeFn == nil {
		return nil, nil
	}
	return f.DescribeVpcAttributeFn(input)
}

func (f *EC2) DescribeVpcClassicLinkRequest(input *ec2.DescribeVpcClassicLinkInput) (*request.Request, *ec2.DescribeVpcClassicLinkOutput) {
	if f.DescribeVpcClassicLinkRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpcClassicLinkRequestFn(input)
}

func (f *EC2) DescribeVpcClassicLink(input *ec2.DescribeVpcClassicLinkInput) (*ec2.DescribeVpcClassicLinkOutput, error) {
	if f.DescribeVpcClassicLinkFn == nil {
		return nil, nil
	}
	return f.DescribeVpcClassicLinkFn(input)
}

func (f *EC2) DescribeVpcClassicLinkDnsSupportRequest(input *ec2.DescribeVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.DescribeVpcClassicLinkDnsSupportOutput) {
	if f.DescribeVpcClassicLinkDnsSupportRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpcClassicLinkDnsSupportRequestFn(input)
}

func (f *EC2) DescribeVpcClassicLinkDnsSupport(input *ec2.DescribeVpcClassicLinkDnsSupportInput) (*ec2.DescribeVpcClassicLinkDnsSupportOutput, error) {
	if f.DescribeVpcClassicLinkDnsSupportFn == nil {
		return nil, nil
	}
	return f.DescribeVpcClassicLinkDnsSupportFn(input)
}

func (f *EC2) DescribeVpcEndpointServicesRequest(input *ec2.DescribeVpcEndpointServicesInput) (*request.Request, *ec2.DescribeVpcEndpointServicesOutput) {
	if f.DescribeVpcEndpointServicesRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpcEndpointServicesRequestFn(input)
}

func (f *EC2) DescribeVpcEndpointServices(input *ec2.DescribeVpcEndpointServicesInput) (*ec2.DescribeVpcEndpointServicesOutput, error) {
	if f.DescribeVpcEndpointServicesFn == nil {
		return nil, nil
	}
	return f.DescribeVpcEndpointServicesFn(input)
}

func (f *EC2) DescribeVpcEndpointsRequest(input *ec2.DescribeVpcEndpointsInput) (*request.Request, *ec2.DescribeVpcEndpointsOutput) {
	if f.DescribeVpcEndpointsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpcEndpointsRequestFn(input)
}

func (f *EC2) DescribeVpcEndpoints(input *ec2.DescribeVpcEndpointsInput) (*ec2.DescribeVpcEndpointsOutput, error) {
	if f.DescribeVpcEndpointsFn == nil {
		return nil, nil
	}
	return f.DescribeVpcEndpointsFn(input)
}

func (f *EC2) DescribeVpcPeeringConnectionsRequest(input *ec2.DescribeVpcPeeringConnectionsInput) (*request.Request, *ec2.DescribeVpcPeeringConnectionsOutput) {
	if f.DescribeVpcPeeringConnectionsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpcPeeringConnectionsRequestFn(input)
}

func (f *EC2) DescribeVpcPeeringConnections(input *ec2.DescribeVpcPeeringConnectionsInput) (*ec2.DescribeVpcPeeringConnectionsOutput, error) {
	if f.DescribeVpcPeeringConnectionsFn == nil {
		return nil, nil
	}
	return f.DescribeVpcPeeringConnectionsFn(input)
}

func (f *EC2) DescribeVpcsRequest(input *ec2.DescribeVpcsInput) (*request.Request, *ec2.DescribeVpcsOutput) {
	if f.DescribeVpcsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpcsRequestFn(input)
}

func (f *EC2) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	if f.DescribeVpcsFn == nil {
		return nil, nil
	}
	return f.DescribeVpcsFn(input)
}

func (f *EC2) DescribeVpnConnectionsRequest(input *ec2.DescribeVpnConnectionsInput) (*request.Request, *ec2.DescribeVpnConnectionsOutput) {
	if f.DescribeVpnConnectionsRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpnConnectionsRequestFn(input)
}

func (f *EC2) DescribeVpnConnections(input *ec2.DescribeVpnConnectionsInput) (*ec2.DescribeVpnConnectionsOutput, error) {
	if f.DescribeVpnConnectionsFn == nil {
		return nil, nil
	}
	return f.DescribeVpnConnectionsFn(input)
}

func (f *EC2) DescribeVpnGatewaysRequest(input *ec2.DescribeVpnGatewaysInput) (*request.Request, *ec2.DescribeVpnGatewaysOutput) {
	if f.DescribeVpnGatewaysRequestFn == nil {
		return nil, nil
	}
	return f.DescribeVpnGatewaysRequestFn(input)
}

func (f *EC2) DescribeVpnGateways(input *ec2.DescribeVpnGatewaysInput) (*ec2.DescribeVpnGatewaysOutput, error) {
	if f.DescribeVpnGatewaysFn == nil {
		return nil, nil
	}
	return f.DescribeVpnGatewaysFn(input)
}

func (f *EC2) DetachClassicLinkVpcRequest(input *ec2.DetachClassicLinkVpcInput) (*request.Request, *ec2.DetachClassicLinkVpcOutput) {
	if f.DetachClassicLinkVpcRequestFn == nil {
		return nil, nil
	}
	return f.DetachClassicLinkVpcRequestFn(input)
}

func (f *EC2) DetachClassicLinkVpc(input *ec2.DetachClassicLinkVpcInput) (*ec2.DetachClassicLinkVpcOutput, error) {
	if f.DetachClassicLinkVpcFn == nil {
		return nil, nil
	}
	return f.DetachClassicLinkVpcFn(input)
}

func (f *EC2) DetachInternetGatewayRequest(input *ec2.DetachInternetGatewayInput) (*request.Request, *ec2.DetachInternetGatewayOutput) {
	if f.DetachInternetGatewayRequestFn == nil {
		return nil, nil
	}
	return f.DetachInternetGatewayRequestFn(input)
}

func (f *EC2) DetachInternetGateway(input *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
	if f.DetachInternetGatewayFn == nil {
		return nil, nil
	}
	return f.DetachInternetGatewayFn(input)
}

func (f *EC2) DetachNetworkInterfaceRequest(input *ec2.DetachNetworkInterfaceInput) (*request.Request, *ec2.DetachNetworkInterfaceOutput) {
	if f.DetachNetworkInterfaceRequestFn == nil {
		return nil, nil
	}
	return f.DetachNetworkInterfaceRequestFn(input)
}

func (f *EC2) DetachNetworkInterface(input *ec2.DetachNetworkInterfaceInput) (*ec2.DetachNetworkInterfaceOutput, error) {
	if f.DetachNetworkInterfaceFn == nil {
		return nil, nil
	}
	return f.DetachNetworkInterfaceFn(input)
}

func (f *EC2) DetachVolumeRequest(input *ec2.DetachVolumeInput) (*request.Request, *ec2.VolumeAttachment) {
	if f.DetachVolumeRequestFn == nil {
		return nil, nil
	}
	return f.DetachVolumeRequestFn(input)
}

func (f *EC2) DetachVolume(input *ec2.DetachVolumeInput) (*ec2.VolumeAttachment, error) {
	if f.DetachVolumeFn == nil {
		return nil, nil
	}
	return f.DetachVolumeFn(input)
}

func (f *EC2) DetachVpnGatewayRequest(input *ec2.DetachVpnGatewayInput) (*request.Request, *ec2.DetachVpnGatewayOutput) {
	if f.DetachVpnGatewayRequestFn == nil {
		return nil, nil
	}
	return f.DetachVpnGatewayRequestFn(input)
}

func (f *EC2) DetachVpnGateway(input *ec2.DetachVpnGatewayInput) (*ec2.DetachVpnGatewayOutput, error) {
	if f.DetachVpnGatewayFn == nil {
		return nil, nil
	}
	return f.DetachVpnGatewayFn(input)
}

func (f *EC2) DisableVgwRoutePropagationRequest(input *ec2.DisableVgwRoutePropagationInput) (*request.Request, *ec2.DisableVgwRoutePropagationOutput) {
	if f.DisableVgwRoutePropagationRequestFn == nil {
		return nil, nil
	}
	return f.DisableVgwRoutePropagationRequestFn(input)
}

func (f *EC2) DisableVgwRoutePropagation(input *ec2.DisableVgwRoutePropagationInput) (*ec2.DisableVgwRoutePropagationOutput, error) {
	if f.DisableVgwRoutePropagationFn == nil {
		return nil, nil
	}
	return f.DisableVgwRoutePropagationFn(input)
}

func (f *EC2) DisableVpcClassicLinkRequest(input *ec2.DisableVpcClassicLinkInput) (*request.Request, *ec2.DisableVpcClassicLinkOutput) {
	if f.DisableVpcClassicLinkRequestFn == nil {
		return nil, nil
	}
	return f.DisableVpcClassicLinkRequestFn(input)
}

func (f *EC2) DisableVpcClassicLink(input *ec2.DisableVpcClassicLinkInput) (*ec2.DisableVpcClassicLinkOutput, error) {
	if f.DisableVpcClassicLinkFn == nil {
		return nil, nil
	}
	return f.DisableVpcClassicLinkFn(input)
}

func (f *EC2) DisableVpcClassicLinkDnsSupportRequest(input *ec2.DisableVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.DisableVpcClassicLinkDnsSupportOutput) {
	if f.DisableVpcClassicLinkDnsSupportRequestFn == nil {
		return nil, nil
	}
	return f.DisableVpcClassicLinkDnsSupportRequestFn(input)
}

func (f *EC2) DisableVpcClassicLinkDnsSupport(input *ec2.DisableVpcClassicLinkDnsSupportInput) (*ec2.DisableVpcClassicLinkDnsSupportOutput, error) {
	if f.DisableVpcClassicLinkDnsSupportFn == nil {
		return nil, nil
	}
	return f.DisableVpcClassicLinkDnsSupportFn(input)
}

func (f *EC2) DisassociateAddressRequest(input *ec2.DisassociateAddressInput) (*request.Request, *ec2.DisassociateAddressOutput) {
	if f.DisassociateAddressRequestFn == nil {
		return nil, nil
	}
	return f.DisassociateAddressRequestFn(input)
}

func (f *EC2) DisassociateAddress(input *ec2.DisassociateAddressInput) (*ec2.DisassociateAddressOutput, error) {
	if f.DisassociateAddressFn == nil {
		return nil, nil
	}
	return f.DisassociateAddressFn(input)
}

func (f *EC2) DisassociateRouteTableRequest(input *ec2.DisassociateRouteTableInput) (*request.Request, *ec2.DisassociateRouteTableOutput) {
	if f.DisassociateRouteTableRequestFn == nil {
		return nil, nil
	}
	return f.DisassociateRouteTableRequestFn(input)
}

func (f *EC2) DisassociateRouteTable(input *ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
	if f.DisassociateRouteTableFn == nil {
		return nil, nil
	}
	return f.DisassociateRouteTableFn(input)
}

func (f *EC2) EnableVgwRoutePropagationRequest(input *ec2.EnableVgwRoutePropagationInput) (*request.Request, *ec2.EnableVgwRoutePropagationOutput) {
	if f.EnableVgwRoutePropagationRequestFn == nil {
		return nil, nil
	}
	return f.EnableVgwRoutePropagationRequestFn(input)
}

func (f *EC2) EnableVgwRoutePropagation(input *ec2.EnableVgwRoutePropagationInput) (*ec2.EnableVgwRoutePropagationOutput, error) {
	if f.EnableVgwRoutePropagationFn == nil {
		return nil, nil
	}
	return f.EnableVgwRoutePropagationFn(input)
}

func (f *EC2) EnableVolumeIORequest(input *ec2.EnableVolumeIOInput) (*request.Request, *ec2.EnableVolumeIOOutput) {
	if f.EnableVolumeIORequestFn == nil {
		return nil, nil
	}
	return f.EnableVolumeIORequestFn(input)
}

func (f *EC2) EnableVolumeIO(input *ec2.EnableVolumeIOInput) (*ec2.EnableVolumeIOOutput, error) {
	if f.EnableVolumeIOFn == nil {
		return nil, nil
	}
	return f.EnableVolumeIOFn(input)
}

func (f *EC2) EnableVpcClassicLinkRequest(input *ec2.EnableVpcClassicLinkInput) (*request.Request, *ec2.EnableVpcClassicLinkOutput) {
	if f.EnableVpcClassicLinkRequestFn == nil {
		return nil, nil
	}
	return f.EnableVpcClassicLinkRequestFn(input)
}

func (f *EC2) EnableVpcClassicLink(input *ec2.EnableVpcClassicLinkInput) (*ec2.EnableVpcClassicLinkOutput, error) {
	if f.EnableVpcClassicLinkFn == nil {
		return nil, nil
	}
	return f.EnableVpcClassicLinkFn(input)
}

func (f *EC2) EnableVpcClassicLinkDnsSupportRequest(input *ec2.EnableVpcClassicLinkDnsSupportInput) (*request.Request, *ec2.EnableVpcClassicLinkDnsSupportOutput) {
	if f.EnableVpcClassicLinkDnsSupportRequestFn == nil {
		return nil, nil
	}
	return f.EnableVpcClassicLinkDnsSupportRequestFn(input)
}

func (f *EC2) EnableVpcClassicLinkDnsSupport(input *ec2.EnableVpcClassicLinkDnsSupportInput) (*ec2.EnableVpcClassicLinkDnsSupportOutput, error) {
	if f.EnableVpcClassicLinkDnsSupportFn == nil {
		return nil, nil
	}
	return f.EnableVpcClassicLinkDnsSupportFn(input)
}

func (f *EC2) GetConsoleOutputRequest(input *ec2.GetConsoleOutputInput) (*request.Request, *ec2.GetConsoleOutputOutput) {
	if f.GetConsoleOutputRequestFn == nil {
		return nil, nil
	}
	return f.GetConsoleOutputRequestFn(input)
}

func (f *EC2) GetConsoleOutput(input *ec2.GetConsoleOutputInput) (*ec2.GetConsoleOutputOutput, error) {
	if f.GetConsoleOutputFn == nil {
		return nil, nil
	}
	return f.GetConsoleOutputFn(input)
}

func (f *EC2) GetPasswordDataRequest(input *ec2.GetPasswordDataInput) (*request.Request, *ec2.GetPasswordDataOutput) {
	if f.GetPasswordDataRequestFn == nil {
		return nil, nil
	}
	return f.GetPasswordDataRequestFn(input)
}

func (f *EC2) GetPasswordData(input *ec2.GetPasswordDataInput) (*ec2.GetPasswordDataOutput, error) {
	if f.GetPasswordDataFn == nil {
		return nil, nil
	}
	return f.GetPasswordDataFn(input)
}

func (f *EC2) ImportImageRequest(input *ec2.ImportImageInput) (*request.Request, *ec2.ImportImageOutput) {
	if f.ImportImageRequestFn == nil {
		return nil, nil
	}
	return f.ImportImageRequestFn(input)
}

func (f *EC2) ImportImage(input *ec2.ImportImageInput) (*ec2.ImportImageOutput, error) {
	if f.ImportImageFn == nil {
		return nil, nil
	}
	return f.ImportImageFn(input)
}

func (f *EC2) ImportInstanceRequest(input *ec2.ImportInstanceInput) (*request.Request, *ec2.ImportInstanceOutput) {
	if f.ImportInstanceRequestFn == nil {
		return nil, nil
	}
	return f.ImportInstanceRequestFn(input)
}

func (f *EC2) ImportInstance(input *ec2.ImportInstanceInput) (*ec2.ImportInstanceOutput, error) {
	if f.ImportInstanceFn == nil {
		return nil, nil
	}
	return f.ImportInstanceFn(input)
}

func (f *EC2) ImportKeyPairRequest(input *ec2.ImportKeyPairInput) (*request.Request, *ec2.ImportKeyPairOutput) {
	if f.ImportKeyPairRequestFn == nil {
		return nil, nil
	}
	return f.ImportKeyPairRequestFn(input)
}

func (f *EC2) ImportKeyPair(input *ec2.ImportKeyPairInput) (*ec2.ImportKeyPairOutput, error) {
	if f.ImportKeyPairFn == nil {
		return nil, nil
	}
	return f.ImportKeyPairFn(input)
}

func (f *EC2) ImportSnapshotRequest(input *ec2.ImportSnapshotInput) (*request.Request, *ec2.ImportSnapshotOutput) {
	if f.ImportSnapshotRequestFn == nil {
		return nil, nil
	}
	return f.ImportSnapshotRequestFn(input)
}

func (f *EC2) ImportSnapshot(input *ec2.ImportSnapshotInput) (*ec2.ImportSnapshotOutput, error) {
	if f.ImportSnapshotFn == nil {
		return nil, nil
	}
	return f.ImportSnapshotFn(input)
}

func (f *EC2) ImportVolumeRequest(input *ec2.ImportVolumeInput) (*request.Request, *ec2.ImportVolumeOutput) {
	if f.ImportVolumeRequestFn == nil {
		return nil, nil
	}
	return f.ImportVolumeRequestFn(input)
}

func (f *EC2) ImportVolume(input *ec2.ImportVolumeInput) (*ec2.ImportVolumeOutput, error) {
	if f.ImportVolumeFn == nil {
		return nil, nil
	}
	return f.ImportVolumeFn(input)
}

func (f *EC2) ModifyHostsRequest(input *ec2.ModifyHostsInput) (*request.Request, *ec2.ModifyHostsOutput) {
	if f.ModifyHostsRequestFn == nil {
		return nil, nil
	}
	return f.ModifyHostsRequestFn(input)
}

func (f *EC2) ModifyHosts(input *ec2.ModifyHostsInput) (*ec2.ModifyHostsOutput, error) {
	if f.ModifyHostsFn == nil {
		return nil, nil
	}
	return f.ModifyHostsFn(input)
}

func (f *EC2) ModifyIdFormatRequest(input *ec2.ModifyIdFormatInput) (*request.Request, *ec2.ModifyIdFormatOutput) {
	if f.ModifyIdFormatRequestFn == nil {
		return nil, nil
	}
	return f.ModifyIdFormatRequestFn(input)
}

func (f *EC2) ModifyIdFormat(input *ec2.ModifyIdFormatInput) (*ec2.ModifyIdFormatOutput, error) {
	if f.ModifyIdFormatFn == nil {
		return nil, nil
	}
	return f.ModifyIdFormatFn(input)
}

func (f *EC2) ModifyImageAttributeRequest(input *ec2.ModifyImageAttributeInput) (*request.Request, *ec2.ModifyImageAttributeOutput) {
	if f.ModifyImageAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ModifyImageAttributeRequestFn(input)
}

func (f *EC2) ModifyImageAttribute(input *ec2.ModifyImageAttributeInput) (*ec2.ModifyImageAttributeOutput, error) {
	if f.ModifyImageAttributeFn == nil {
		return nil, nil
	}
	return f.ModifyImageAttributeFn(input)
}

func (f *EC2) ModifyInstanceAttributeRequest(input *ec2.ModifyInstanceAttributeInput) (*request.Request, *ec2.ModifyInstanceAttributeOutput) {
	if f.ModifyInstanceAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ModifyInstanceAttributeRequestFn(input)
}

func (f *EC2) ModifyInstanceAttribute(input *ec2.ModifyInstanceAttributeInput) (*ec2.ModifyInstanceAttributeOutput, error) {
	if f.ModifyInstanceAttributeFn == nil {
		return nil, nil
	}
	return f.ModifyInstanceAttributeFn(input)
}

func (f *EC2) ModifyInstancePlacementRequest(input *ec2.ModifyInstancePlacementInput) (*request.Request, *ec2.ModifyInstancePlacementOutput) {
	if f.ModifyInstancePlacementRequestFn == nil {
		return nil, nil
	}
	return f.ModifyInstancePlacementRequestFn(input)
}

func (f *EC2) ModifyInstancePlacement(input *ec2.ModifyInstancePlacementInput) (*ec2.ModifyInstancePlacementOutput, error) {
	if f.ModifyInstancePlacementFn == nil {
		return nil, nil
	}
	return f.ModifyInstancePlacementFn(input)
}

func (f *EC2) ModifyNetworkInterfaceAttributeRequest(input *ec2.ModifyNetworkInterfaceAttributeInput) (*request.Request, *ec2.ModifyNetworkInterfaceAttributeOutput) {
	if f.ModifyNetworkInterfaceAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ModifyNetworkInterfaceAttributeRequestFn(input)
}

func (f *EC2) ModifyNetworkInterfaceAttribute(input *ec2.ModifyNetworkInterfaceAttributeInput) (*ec2.ModifyNetworkInterfaceAttributeOutput, error) {
	if f.ModifyNetworkInterfaceAttributeFn == nil {
		return nil, nil
	}
	return f.ModifyNetworkInterfaceAttributeFn(input)
}

func (f *EC2) ModifyReservedInstancesRequest(input *ec2.ModifyReservedInstancesInput) (*request.Request, *ec2.ModifyReservedInstancesOutput) {
	if f.ModifyReservedInstancesRequestFn == nil {
		return nil, nil
	}
	return f.ModifyReservedInstancesRequestFn(input)
}

func (f *EC2) ModifyReservedInstances(input *ec2.ModifyReservedInstancesInput) (*ec2.ModifyReservedInstancesOutput, error) {
	if f.ModifyReservedInstancesFn == nil {
		return nil, nil
	}
	return f.ModifyReservedInstancesFn(input)
}

func (f *EC2) ModifySnapshotAttributeRequest(input *ec2.ModifySnapshotAttributeInput) (*request.Request, *ec2.ModifySnapshotAttributeOutput) {
	if f.ModifySnapshotAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ModifySnapshotAttributeRequestFn(input)
}

func (f *EC2) ModifySnapshotAttribute(input *ec2.ModifySnapshotAttributeInput) (*ec2.ModifySnapshotAttributeOutput, error) {
	if f.ModifySnapshotAttributeFn == nil {
		return nil, nil
	}
	return f.ModifySnapshotAttributeFn(input)
}

func (f *EC2) ModifySpotFleetRequestRequest(input *ec2.ModifySpotFleetRequestInput) (*request.Request, *ec2.ModifySpotFleetRequestOutput) {
	if f.ModifySpotFleetRequestRequestFn == nil {
		return nil, nil
	}
	return f.ModifySpotFleetRequestRequestFn(input)
}

func (f *EC2) ModifySpotFleetRequest(input *ec2.ModifySpotFleetRequestInput) (*ec2.ModifySpotFleetRequestOutput, error) {
	if f.ModifySpotFleetRequestFn == nil {
		return nil, nil
	}
	return f.ModifySpotFleetRequestFn(input)
}

func (f *EC2) ModifySubnetAttributeRequest(input *ec2.ModifySubnetAttributeInput) (*request.Request, *ec2.ModifySubnetAttributeOutput) {
	if f.ModifySubnetAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ModifySubnetAttributeRequestFn(input)
}

func (f *EC2) ModifySubnetAttribute(input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
	if f.ModifySubnetAttributeFn == nil {
		return nil, nil
	}
	return f.ModifySubnetAttributeFn(input)
}

func (f *EC2) ModifyVolumeAttributeRequest(input *ec2.ModifyVolumeAttributeInput) (*request.Request, *ec2.ModifyVolumeAttributeOutput) {
	if f.ModifyVolumeAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ModifyVolumeAttributeRequestFn(input)
}

func (f *EC2) ModifyVolumeAttribute(input *ec2.ModifyVolumeAttributeInput) (*ec2.ModifyVolumeAttributeOutput, error) {
	if f.ModifyVolumeAttributeFn == nil {
		return nil, nil
	}
	return f.ModifyVolumeAttributeFn(input)
}

func (f *EC2) ModifyVpcAttributeRequest(input *ec2.ModifyVpcAttributeInput) (*request.Request, *ec2.ModifyVpcAttributeOutput) {
	if f.ModifyVpcAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ModifyVpcAttributeRequestFn(input)
}

func (f *EC2) ModifyVpcAttribute(input *ec2.ModifyVpcAttributeInput) (*ec2.ModifyVpcAttributeOutput, error) {
	if f.ModifyVpcAttributeFn == nil {
		return nil, nil
	}
	return f.ModifyVpcAttributeFn(input)
}

func (f *EC2) ModifyVpcEndpointRequest(input *ec2.ModifyVpcEndpointInput) (*request.Request, *ec2.ModifyVpcEndpointOutput) {
	if f.ModifyVpcEndpointRequestFn == nil {
		return nil, nil
	}
	return f.ModifyVpcEndpointRequestFn(input)
}

func (f *EC2) ModifyVpcEndpoint(input *ec2.ModifyVpcEndpointInput) (*ec2.ModifyVpcEndpointOutput, error) {
	if f.ModifyVpcEndpointFn == nil {
		return nil, nil
	}
	return f.ModifyVpcEndpointFn(input)
}

func (f *EC2) MonitorInstancesRequest(input *ec2.MonitorInstancesInput) (*request.Request, *ec2.MonitorInstancesOutput) {
	if f.MonitorInstancesRequestFn == nil {
		return nil, nil
	}
	return f.MonitorInstancesRequestFn(input)
}

func (f *EC2) MonitorInstances(input *ec2.MonitorInstancesInput) (*ec2.MonitorInstancesOutput, error) {
	if f.MonitorInstancesFn == nil {
		return nil, nil
	}
	return f.MonitorInstancesFn(input)
}

func (f *EC2) MoveAddressToVpcRequest(input *ec2.MoveAddressToVpcInput) (*request.Request, *ec2.MoveAddressToVpcOutput) {
	if f.MoveAddressToVpcRequestFn == nil {
		return nil, nil
	}
	return f.MoveAddressToVpcRequestFn(input)
}

func (f *EC2) MoveAddressToVpc(input *ec2.MoveAddressToVpcInput) (*ec2.MoveAddressToVpcOutput, error) {
	if f.MoveAddressToVpcFn == nil {
		return nil, nil
	}
	return f.MoveAddressToVpcFn(input)
}

func (f *EC2) PurchaseReservedInstancesOfferingRequest(input *ec2.PurchaseReservedInstancesOfferingInput) (*request.Request, *ec2.PurchaseReservedInstancesOfferingOutput) {
	if f.PurchaseReservedInstancesOfferingRequestFn == nil {
		return nil, nil
	}
	return f.PurchaseReservedInstancesOfferingRequestFn(input)
}

func (f *EC2) PurchaseReservedInstancesOffering(input *ec2.PurchaseReservedInstancesOfferingInput) (*ec2.PurchaseReservedInstancesOfferingOutput, error) {
	if f.PurchaseReservedInstancesOfferingFn == nil {
		return nil, nil
	}
	return f.PurchaseReservedInstancesOfferingFn(input)
}

func (f *EC2) PurchaseScheduledInstancesRequest(input *ec2.PurchaseScheduledInstancesInput) (*request.Request, *ec2.PurchaseScheduledInstancesOutput) {
	if f.PurchaseScheduledInstancesRequestFn == nil {
		return nil, nil
	}
	return f.PurchaseScheduledInstancesRequestFn(input)
}

func (f *EC2) PurchaseScheduledInstances(input *ec2.PurchaseScheduledInstancesInput) (*ec2.PurchaseScheduledInstancesOutput, error) {
	if f.PurchaseScheduledInstancesFn == nil {
		return nil, nil
	}
	return f.PurchaseScheduledInstancesFn(input)
}

func (f *EC2) RebootInstancesRequest(input *ec2.RebootInstancesInput) (*request.Request, *ec2.RebootInstancesOutput) {
	if f.RebootInstancesRequestFn == nil {
		return nil, nil
	}
	return f.RebootInstancesRequestFn(input)
}

func (f *EC2) RebootInstances(input *ec2.RebootInstancesInput) (*ec2.RebootInstancesOutput, error) {
	if f.RebootInstancesFn == nil {
		return nil, nil
	}
	return f.RebootInstancesFn(input)
}

func (f *EC2) RegisterImageRequest(input *ec2.RegisterImageInput) (*request.Request, *ec2.RegisterImageOutput) {
	if f.RegisterImageRequestFn == nil {
		return nil, nil
	}
	return f.RegisterImageRequestFn(input)
}

func (f *EC2) RegisterImage(input *ec2.RegisterImageInput) (*ec2.RegisterImageOutput, error) {
	if f.RegisterImageFn == nil {
		return nil, nil
	}
	return f.RegisterImageFn(input)
}

func (f *EC2) RejectVpcPeeringConnectionRequest(input *ec2.RejectVpcPeeringConnectionInput) (*request.Request, *ec2.RejectVpcPeeringConnectionOutput) {
	if f.RejectVpcPeeringConnectionRequestFn == nil {
		return nil, nil
	}
	return f.RejectVpcPeeringConnectionRequestFn(input)
}

func (f *EC2) RejectVpcPeeringConnection(input *ec2.RejectVpcPeeringConnectionInput) (*ec2.RejectVpcPeeringConnectionOutput, error) {
	if f.RejectVpcPeeringConnectionFn == nil {
		return nil, nil
	}
	return f.RejectVpcPeeringConnectionFn(input)
}

func (f *EC2) ReleaseAddressRequest(input *ec2.ReleaseAddressInput) (*request.Request, *ec2.ReleaseAddressOutput) {
	if f.ReleaseAddressRequestFn == nil {
		return nil, nil
	}
	return f.ReleaseAddressRequestFn(input)
}

func (f *EC2) ReleaseAddress(input *ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
	if f.ReleaseAddressFn == nil {
		return nil, nil
	}
	return f.ReleaseAddressFn(input)
}

func (f *EC2) ReleaseHostsRequest(input *ec2.ReleaseHostsInput) (*request.Request, *ec2.ReleaseHostsOutput) {
	if f.ReleaseHostsRequestFn == nil {
		return nil, nil
	}
	return f.ReleaseHostsRequestFn(input)
}

func (f *EC2) ReleaseHosts(input *ec2.ReleaseHostsInput) (*ec2.ReleaseHostsOutput, error) {
	if f.ReleaseHostsFn == nil {
		return nil, nil
	}
	return f.ReleaseHostsFn(input)
}

func (f *EC2) ReplaceNetworkAclAssociationRequest(input *ec2.ReplaceNetworkAclAssociationInput) (*request.Request, *ec2.ReplaceNetworkAclAssociationOutput) {
	if f.ReplaceNetworkAclAssociationRequestFn == nil {
		return nil, nil
	}
	return f.ReplaceNetworkAclAssociationRequestFn(input)
}

func (f *EC2) ReplaceNetworkAclAssociation(input *ec2.ReplaceNetworkAclAssociationInput) (*ec2.ReplaceNetworkAclAssociationOutput, error) {
	if f.ReplaceNetworkAclAssociationFn == nil {
		return nil, nil
	}
	return f.ReplaceNetworkAclAssociationFn(input)
}

func (f *EC2) ReplaceNetworkAclEntryRequest(input *ec2.ReplaceNetworkAclEntryInput) (*request.Request, *ec2.ReplaceNetworkAclEntryOutput) {
	if f.ReplaceNetworkAclEntryRequestFn == nil {
		return nil, nil
	}
	return f.ReplaceNetworkAclEntryRequestFn(input)
}

func (f *EC2) ReplaceNetworkAclEntry(input *ec2.ReplaceNetworkAclEntryInput) (*ec2.ReplaceNetworkAclEntryOutput, error) {
	if f.ReplaceNetworkAclEntryFn == nil {
		return nil, nil
	}
	return f.ReplaceNetworkAclEntryFn(input)
}

func (f *EC2) ReplaceRouteRequest(input *ec2.ReplaceRouteInput) (*request.Request, *ec2.ReplaceRouteOutput) {
	if f.ReplaceRouteRequestFn == nil {
		return nil, nil
	}
	return f.ReplaceRouteRequestFn(input)
}

func (f *EC2) ReplaceRoute(input *ec2.ReplaceRouteInput) (*ec2.ReplaceRouteOutput, error) {
	if f.ReplaceRouteFn == nil {
		return nil, nil
	}
	return f.ReplaceRouteFn(input)
}

func (f *EC2) ReplaceRouteTableAssociationRequest(input *ec2.ReplaceRouteTableAssociationInput) (*request.Request, *ec2.ReplaceRouteTableAssociationOutput) {
	if f.ReplaceRouteTableAssociationRequestFn == nil {
		return nil, nil
	}
	return f.ReplaceRouteTableAssociationRequestFn(input)
}

func (f *EC2) ReplaceRouteTableAssociation(input *ec2.ReplaceRouteTableAssociationInput) (*ec2.ReplaceRouteTableAssociationOutput, error) {
	if f.ReplaceRouteTableAssociationFn == nil {
		return nil, nil
	}
	return f.ReplaceRouteTableAssociationFn(input)
}

func (f *EC2) ReportInstanceStatusRequest(input *ec2.ReportInstanceStatusInput) (*request.Request, *ec2.ReportInstanceStatusOutput) {
	if f.ReportInstanceStatusRequestFn == nil {
		return nil, nil
	}
	return f.ReportInstanceStatusRequestFn(input)
}

func (f *EC2) ReportInstanceStatus(input *ec2.ReportInstanceStatusInput) (*ec2.ReportInstanceStatusOutput, error) {
	if f.ReportInstanceStatusFn == nil {
		return nil, nil
	}
	return f.ReportInstanceStatusFn(input)
}

func (f *EC2) RequestSpotFleetRequest(input *ec2.RequestSpotFleetInput) (*request.Request, *ec2.RequestSpotFleetOutput) {
	if f.RequestSpotFleetRequestFn == nil {
		return nil, nil
	}
	return f.RequestSpotFleetRequestFn(input)
}

func (f *EC2) RequestSpotFleet(input *ec2.RequestSpotFleetInput) (*ec2.RequestSpotFleetOutput, error) {
	if f.RequestSpotFleetFn == nil {
		return nil, nil
	}
	return f.RequestSpotFleetFn(input)
}

func (f *EC2) RequestSpotInstancesRequest(input *ec2.RequestSpotInstancesInput) (*request.Request, *ec2.RequestSpotInstancesOutput) {
	if f.RequestSpotInstancesRequestFn == nil {
		return nil, nil
	}
	return f.RequestSpotInstancesRequestFn(input)
}

func (f *EC2) RequestSpotInstances(input *ec2.RequestSpotInstancesInput) (*ec2.RequestSpotInstancesOutput, error) {
	if f.RequestSpotInstancesFn == nil {
		return nil, nil
	}
	return f.RequestSpotInstancesFn(input)
}

func (f *EC2) ResetImageAttributeRequest(input *ec2.ResetImageAttributeInput) (*request.Request, *ec2.ResetImageAttributeOutput) {
	if f.ResetImageAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ResetImageAttributeRequestFn(input)
}

func (f *EC2) ResetImageAttribute(input *ec2.ResetImageAttributeInput) (*ec2.ResetImageAttributeOutput, error) {
	if f.ResetImageAttributeFn == nil {
		return nil, nil
	}
	return f.ResetImageAttributeFn(input)
}

func (f *EC2) ResetInstanceAttributeRequest(input *ec2.ResetInstanceAttributeInput) (*request.Request, *ec2.ResetInstanceAttributeOutput) {
	if f.ResetInstanceAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ResetInstanceAttributeRequestFn(input)
}

func (f *EC2) ResetInstanceAttribute(input *ec2.ResetInstanceAttributeInput) (*ec2.ResetInstanceAttributeOutput, error) {
	if f.ResetInstanceAttributeFn == nil {
		return nil, nil
	}
	return f.ResetInstanceAttributeFn(input)
}

func (f *EC2) ResetNetworkInterfaceAttributeRequest(input *ec2.ResetNetworkInterfaceAttributeInput) (*request.Request, *ec2.ResetNetworkInterfaceAttributeOutput) {
	if f.ResetNetworkInterfaceAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ResetNetworkInterfaceAttributeRequestFn(input)
}

func (f *EC2) ResetNetworkInterfaceAttribute(input *ec2.ResetNetworkInterfaceAttributeInput) (*ec2.ResetNetworkInterfaceAttributeOutput, error) {
	if f.ResetNetworkInterfaceAttributeFn == nil {
		return nil, nil
	}
	return f.ResetNetworkInterfaceAttributeFn(input)
}

func (f *EC2) ResetSnapshotAttributeRequest(input *ec2.ResetSnapshotAttributeInput) (*request.Request, *ec2.ResetSnapshotAttributeOutput) {
	if f.ResetSnapshotAttributeRequestFn == nil {
		return nil, nil
	}
	return f.ResetSnapshotAttributeRequestFn(input)
}

func (f *EC2) ResetSnapshotAttribute(input *ec2.ResetSnapshotAttributeInput) (*ec2.ResetSnapshotAttributeOutput, error) {
	if f.ResetSnapshotAttributeFn == nil {
		return nil, nil
	}
	return f.ResetSnapshotAttributeFn(input)
}

func (f *EC2) RestoreAddressToClassicRequest(input *ec2.RestoreAddressToClassicInput) (*request.Request, *ec2.RestoreAddressToClassicOutput) {
	if f.RestoreAddressToClassicRequestFn == nil {
		return nil, nil
	}
	return f.RestoreAddressToClassicRequestFn(input)
}

func (f *EC2) RestoreAddressToClassic(input *ec2.RestoreAddressToClassicInput) (*ec2.RestoreAddressToClassicOutput, error) {
	if f.RestoreAddressToClassicFn == nil {
		return nil, nil
	}
	return f.RestoreAddressToClassicFn(input)
}

func (f *EC2) RevokeSecurityGroupEgressRequest(input *ec2.RevokeSecurityGroupEgressInput) (*request.Request, *ec2.RevokeSecurityGroupEgressOutput) {
	if f.RevokeSecurityGroupEgressRequestFn == nil {
		return nil, nil
	}
	return f.RevokeSecurityGroupEgressRequestFn(input)
}

func (f *EC2) RevokeSecurityGroupEgress(input *ec2.RevokeSecurityGroupEgressInput) (*ec2.RevokeSecurityGroupEgressOutput, error) {
	if f.RevokeSecurityGroupEgressFn == nil {
		return nil, nil
	}
	return f.RevokeSecurityGroupEgressFn(input)
}

func (f *EC2) RevokeSecurityGroupIngressRequest(input *ec2.RevokeSecurityGroupIngressInput) (*request.Request, *ec2.RevokeSecurityGroupIngressOutput) {
	if f.RevokeSecurityGroupIngressRequestFn == nil {
		return nil, nil
	}
	return f.RevokeSecurityGroupIngressRequestFn(input)
}

func (f *EC2) RevokeSecurityGroupIngress(input *ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error) {
	if f.RevokeSecurityGroupIngressFn == nil {
		return nil, nil
	}
	return f.RevokeSecurityGroupIngressFn(input)
}

func (f *EC2) RunInstancesRequest(input *ec2.RunInstancesInput) (*request.Request, *ec2.Reservation) {
	if f.RunInstancesRequestFn == nil {
		return nil, nil
	}
	return f.RunInstancesRequestFn(input)
}

func (f *EC2) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if f.RunInstancesFn == nil {
		return nil, nil
	}
	return f.RunInstancesFn(input)
}

func (f *EC2) RunScheduledInstancesRequest(input *ec2.RunScheduledInstancesInput) (*request.Request, *ec2.RunScheduledInstancesOutput) {
	if f.RunScheduledInstancesRequestFn == nil {
		return nil, nil
	}
	return f.RunScheduledInstancesRequestFn(input)
}

func (f *EC2) RunScheduledInstances(input *ec2.RunScheduledInstancesInput) (*ec2.RunScheduledInstancesOutput, error) {
	if f.RunScheduledInstancesFn == nil {
		return nil, nil
	}
	return f.RunScheduledInstancesFn(input)
}

func (f *EC2) StartInstancesRequest(input *ec2.StartInstancesInput) (*request.Request, *ec2.StartInstancesOutput) {
	if f.StartInstancesRequestFn == nil {
		return nil, nil
	}
	return f.StartInstancesRequestFn(input)
}

func (f *EC2) StartInstances(input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	if f.StartInstancesFn == nil {
		return nil, nil
	}
	return f.StartInstancesFn(input)
}

func (f *EC2) StopInstancesRequest(input *ec2.StopInstancesInput) (*request.Request, *ec2.StopInstancesOutput) {
	if f.StopInstancesRequestFn == nil {
		return nil, nil
	}
	return f.StopInstancesRequestFn(input)
}

func (f *EC2) StopInstances(input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	if f.StopInstancesFn == nil {
		return nil, nil
	}
	return f.StopInstancesFn(input)
}

func (f *EC2) TerminateInstancesRequest(input *ec2.TerminateInstancesInput) (*request.Request, *ec2.TerminateInstancesOutput) {
	if f.TerminateInstancesRequestFn == nil {
		return nil, nil
	}
	return f.TerminateInstancesRequestFn(input)
}

func (f *EC2) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	if f.TerminateInstancesFn == nil {
		return nil, nil
	}
	return f.TerminateInstancesFn(input)
}

func (f *EC2) UnassignPrivateIpAddressesRequest(input *ec2.UnassignPrivateIpAddressesInput) (*request.Request, *ec2.UnassignPrivateIpAddressesOutput) {
	if f.UnassignPrivateIpAddressesRequestFn == nil {
		return nil, nil
	}
	return f.UnassignPrivateIpAddressesRequestFn(input)
}

func (f *EC2) UnassignPrivateIpAddresses(input *ec2.UnassignPrivateIpAddressesInput) (*ec2.UnassignPrivateIpAddressesOutput, error) {
	if f.UnassignPrivateIpAddressesFn == nil {
		return nil, nil
	}
	return f.UnassignPrivateIpAddressesFn(input)
}

func (f *EC2) UnmonitorInstancesRequest(input *ec2.UnmonitorInstancesInput) (*request.Request, *ec2.UnmonitorInstancesOutput) {
	if f.UnmonitorInstancesRequestFn == nil {
		return nil, nil
	}
	return f.UnmonitorInstancesRequestFn(input)
}

func (f *EC2) UnmonitorInstances(input *ec2.UnmonitorInstancesInput) (*ec2.UnmonitorInstancesOutput, error) {
	if f.UnmonitorInstancesFn == nil {
		return nil, nil
	}
	return f.UnmonitorInstancesFn(input)
}
