package fake_aws_provider

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

type IAM struct {
	iamiface.IAMAPI
	AddClientIDToOpenIDConnectProviderRequestFn      func(*iam.AddClientIDToOpenIDConnectProviderInput) (*request.Request, *iam.AddClientIDToOpenIDConnectProviderOutput)
	AddClientIDToOpenIDConnectProviderFn             func(*iam.AddClientIDToOpenIDConnectProviderInput) (*iam.AddClientIDToOpenIDConnectProviderOutput, error)
	AddRoleToInstanceProfileRequestFn                func(*iam.AddRoleToInstanceProfileInput) (*request.Request, *iam.AddRoleToInstanceProfileOutput)
	AddRoleToInstanceProfileFn                       func(*iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error)
	AddUserToGroupRequestFn                          func(*iam.AddUserToGroupInput) (*request.Request, *iam.AddUserToGroupOutput)
	AddUserToGroupFn                                 func(*iam.AddUserToGroupInput) (*iam.AddUserToGroupOutput, error)
	AttachGroupPolicyRequestFn                       func(*iam.AttachGroupPolicyInput) (*request.Request, *iam.AttachGroupPolicyOutput)
	AttachGroupPolicyFn                              func(*iam.AttachGroupPolicyInput) (*iam.AttachGroupPolicyOutput, error)
	AttachRolePolicyRequestFn                        func(*iam.AttachRolePolicyInput) (*request.Request, *iam.AttachRolePolicyOutput)
	AttachRolePolicyFn                               func(*iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error)
	AttachUserPolicyRequestFn                        func(*iam.AttachUserPolicyInput) (*request.Request, *iam.AttachUserPolicyOutput)
	AttachUserPolicyFn                               func(*iam.AttachUserPolicyInput) (*iam.AttachUserPolicyOutput, error)
	ChangePasswordRequestFn                          func(*iam.ChangePasswordInput) (*request.Request, *iam.ChangePasswordOutput)
	ChangePasswordFn                                 func(*iam.ChangePasswordInput) (*iam.ChangePasswordOutput, error)
	CreateAccessKeyRequestFn                         func(*iam.CreateAccessKeyInput) (*request.Request, *iam.CreateAccessKeyOutput)
	CreateAccessKeyFn                                func(*iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error)
	CreateAccountAliasRequestFn                      func(*iam.CreateAccountAliasInput) (*request.Request, *iam.CreateAccountAliasOutput)
	CreateAccountAliasFn                             func(*iam.CreateAccountAliasInput) (*iam.CreateAccountAliasOutput, error)
	CreateGroupRequestFn                             func(*iam.CreateGroupInput) (*request.Request, *iam.CreateGroupOutput)
	CreateGroupFn                                    func(*iam.CreateGroupInput) (*iam.CreateGroupOutput, error)
	CreateInstanceProfileRequestFn                   func(*iam.CreateInstanceProfileInput) (*request.Request, *iam.CreateInstanceProfileOutput)
	CreateInstanceProfileFn                          func(*iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error)
	CreateLoginProfileRequestFn                      func(*iam.CreateLoginProfileInput) (*request.Request, *iam.CreateLoginProfileOutput)
	CreateLoginProfileFn                             func(*iam.CreateLoginProfileInput) (*iam.CreateLoginProfileOutput, error)
	CreateOpenIDConnectProviderRequestFn             func(*iam.CreateOpenIDConnectProviderInput) (*request.Request, *iam.CreateOpenIDConnectProviderOutput)
	CreateOpenIDConnectProviderFn                    func(*iam.CreateOpenIDConnectProviderInput) (*iam.CreateOpenIDConnectProviderOutput, error)
	CreatePolicyRequestFn                            func(*iam.CreatePolicyInput) (*request.Request, *iam.CreatePolicyOutput)
	CreatePolicyFn                                   func(*iam.CreatePolicyInput) (*iam.CreatePolicyOutput, error)
	CreatePolicyVersionRequestFn                     func(*iam.CreatePolicyVersionInput) (*request.Request, *iam.CreatePolicyVersionOutput)
	CreatePolicyVersionFn                            func(*iam.CreatePolicyVersionInput) (*iam.CreatePolicyVersionOutput, error)
	CreateRoleRequestFn                              func(*iam.CreateRoleInput) (*request.Request, *iam.CreateRoleOutput)
	CreateRoleFn                                     func(*iam.CreateRoleInput) (*iam.CreateRoleOutput, error)
	CreateSAMLProviderRequestFn                      func(*iam.CreateSAMLProviderInput) (*request.Request, *iam.CreateSAMLProviderOutput)
	CreateSAMLProviderFn                             func(*iam.CreateSAMLProviderInput) (*iam.CreateSAMLProviderOutput, error)
	CreateUserRequestFn                              func(*iam.CreateUserInput) (*request.Request, *iam.CreateUserOutput)
	CreateUserFn                                     func(*iam.CreateUserInput) (*iam.CreateUserOutput, error)
	CreateVirtualMFADeviceRequestFn                  func(*iam.CreateVirtualMFADeviceInput) (*request.Request, *iam.CreateVirtualMFADeviceOutput)
	CreateVirtualMFADeviceFn                         func(*iam.CreateVirtualMFADeviceInput) (*iam.CreateVirtualMFADeviceOutput, error)
	DeactivateMFADeviceRequestFn                     func(*iam.DeactivateMFADeviceInput) (*request.Request, *iam.DeactivateMFADeviceOutput)
	DeactivateMFADeviceFn                            func(*iam.DeactivateMFADeviceInput) (*iam.DeactivateMFADeviceOutput, error)
	DeleteAccessKeyRequestFn                         func(*iam.DeleteAccessKeyInput) (*request.Request, *iam.DeleteAccessKeyOutput)
	DeleteAccessKeyFn                                func(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error)
	DeleteAccountAliasRequestFn                      func(*iam.DeleteAccountAliasInput) (*request.Request, *iam.DeleteAccountAliasOutput)
	DeleteAccountAliasFn                             func(*iam.DeleteAccountAliasInput) (*iam.DeleteAccountAliasOutput, error)
	DeleteAccountPasswordPolicyRequestFn             func(*iam.DeleteAccountPasswordPolicyInput) (*request.Request, *iam.DeleteAccountPasswordPolicyOutput)
	DeleteAccountPasswordPolicyFn                    func(*iam.DeleteAccountPasswordPolicyInput) (*iam.DeleteAccountPasswordPolicyOutput, error)
	DeleteGroupRequestFn                             func(*iam.DeleteGroupInput) (*request.Request, *iam.DeleteGroupOutput)
	DeleteGroupFn                                    func(*iam.DeleteGroupInput) (*iam.DeleteGroupOutput, error)
	DeleteGroupPolicyRequestFn                       func(*iam.DeleteGroupPolicyInput) (*request.Request, *iam.DeleteGroupPolicyOutput)
	DeleteGroupPolicyFn                              func(*iam.DeleteGroupPolicyInput) (*iam.DeleteGroupPolicyOutput, error)
	DeleteInstanceProfileRequestFn                   func(*iam.DeleteInstanceProfileInput) (*request.Request, *iam.DeleteInstanceProfileOutput)
	DeleteInstanceProfileFn                          func(*iam.DeleteInstanceProfileInput) (*iam.DeleteInstanceProfileOutput, error)
	DeleteLoginProfileRequestFn                      func(*iam.DeleteLoginProfileInput) (*request.Request, *iam.DeleteLoginProfileOutput)
	DeleteLoginProfileFn                             func(*iam.DeleteLoginProfileInput) (*iam.DeleteLoginProfileOutput, error)
	DeleteOpenIDConnectProviderRequestFn             func(*iam.DeleteOpenIDConnectProviderInput) (*request.Request, *iam.DeleteOpenIDConnectProviderOutput)
	DeleteOpenIDConnectProviderFn                    func(*iam.DeleteOpenIDConnectProviderInput) (*iam.DeleteOpenIDConnectProviderOutput, error)
	DeletePolicyRequestFn                            func(*iam.DeletePolicyInput) (*request.Request, *iam.DeletePolicyOutput)
	DeletePolicyFn                                   func(*iam.DeletePolicyInput) (*iam.DeletePolicyOutput, error)
	DeletePolicyVersionRequestFn                     func(*iam.DeletePolicyVersionInput) (*request.Request, *iam.DeletePolicyVersionOutput)
	DeletePolicyVersionFn                            func(*iam.DeletePolicyVersionInput) (*iam.DeletePolicyVersionOutput, error)
	DeleteRoleRequestFn                              func(*iam.DeleteRoleInput) (*request.Request, *iam.DeleteRoleOutput)
	DeleteRoleFn                                     func(*iam.DeleteRoleInput) (*iam.DeleteRoleOutput, error)
	DeleteRolePolicyRequestFn                        func(*iam.DeleteRolePolicyInput) (*request.Request, *iam.DeleteRolePolicyOutput)
	DeleteRolePolicyFn                               func(*iam.DeleteRolePolicyInput) (*iam.DeleteRolePolicyOutput, error)
	DeleteSAMLProviderRequestFn                      func(*iam.DeleteSAMLProviderInput) (*request.Request, *iam.DeleteSAMLProviderOutput)
	DeleteSAMLProviderFn                             func(*iam.DeleteSAMLProviderInput) (*iam.DeleteSAMLProviderOutput, error)
	DeleteSSHPublicKeyRequestFn                      func(*iam.DeleteSSHPublicKeyInput) (*request.Request, *iam.DeleteSSHPublicKeyOutput)
	DeleteSSHPublicKeyFn                             func(*iam.DeleteSSHPublicKeyInput) (*iam.DeleteSSHPublicKeyOutput, error)
	DeleteServerCertificateRequestFn                 func(*iam.DeleteServerCertificateInput) (*request.Request, *iam.DeleteServerCertificateOutput)
	DeleteServerCertificateFn                        func(*iam.DeleteServerCertificateInput) (*iam.DeleteServerCertificateOutput, error)
	DeleteSigningCertificateRequestFn                func(*iam.DeleteSigningCertificateInput) (*request.Request, *iam.DeleteSigningCertificateOutput)
	DeleteSigningCertificateFn                       func(*iam.DeleteSigningCertificateInput) (*iam.DeleteSigningCertificateOutput, error)
	DeleteUserRequestFn                              func(*iam.DeleteUserInput) (*request.Request, *iam.DeleteUserOutput)
	DeleteUserFn                                     func(*iam.DeleteUserInput) (*iam.DeleteUserOutput, error)
	DeleteUserPolicyRequestFn                        func(*iam.DeleteUserPolicyInput) (*request.Request, *iam.DeleteUserPolicyOutput)
	DeleteUserPolicyFn                               func(*iam.DeleteUserPolicyInput) (*iam.DeleteUserPolicyOutput, error)
	DeleteVirtualMFADeviceRequestFn                  func(*iam.DeleteVirtualMFADeviceInput) (*request.Request, *iam.DeleteVirtualMFADeviceOutput)
	DeleteVirtualMFADeviceFn                         func(*iam.DeleteVirtualMFADeviceInput) (*iam.DeleteVirtualMFADeviceOutput, error)
	DetachGroupPolicyRequestFn                       func(*iam.DetachGroupPolicyInput) (*request.Request, *iam.DetachGroupPolicyOutput)
	DetachGroupPolicyFn                              func(*iam.DetachGroupPolicyInput) (*iam.DetachGroupPolicyOutput, error)
	DetachRolePolicyRequestFn                        func(*iam.DetachRolePolicyInput) (*request.Request, *iam.DetachRolePolicyOutput)
	DetachRolePolicyFn                               func(*iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error)
	DetachUserPolicyRequestFn                        func(*iam.DetachUserPolicyInput) (*request.Request, *iam.DetachUserPolicyOutput)
	DetachUserPolicyFn                               func(*iam.DetachUserPolicyInput) (*iam.DetachUserPolicyOutput, error)
	EnableMFADeviceRequestFn                         func(*iam.EnableMFADeviceInput) (*request.Request, *iam.EnableMFADeviceOutput)
	EnableMFADeviceFn                                func(*iam.EnableMFADeviceInput) (*iam.EnableMFADeviceOutput, error)
	GenerateCredentialReportRequestFn                func(*iam.GenerateCredentialReportInput) (*request.Request, *iam.GenerateCredentialReportOutput)
	GenerateCredentialReportFn                       func(*iam.GenerateCredentialReportInput) (*iam.GenerateCredentialReportOutput, error)
	GetAccessKeyLastUsedRequestFn                    func(*iam.GetAccessKeyLastUsedInput) (*request.Request, *iam.GetAccessKeyLastUsedOutput)
	GetAccessKeyLastUsedFn                           func(*iam.GetAccessKeyLastUsedInput) (*iam.GetAccessKeyLastUsedOutput, error)
	GetAccountAuthorizationDetailsRequestFn          func(*iam.GetAccountAuthorizationDetailsInput) (*request.Request, *iam.GetAccountAuthorizationDetailsOutput)
	GetAccountAuthorizationDetailsFn                 func(*iam.GetAccountAuthorizationDetailsInput) (*iam.GetAccountAuthorizationDetailsOutput, error)
	GetAccountAuthorizationDetailsPagesFn            func(*iam.GetAccountAuthorizationDetailsInput, func(*iam.GetAccountAuthorizationDetailsOutput, bool) bool) error
	GetAccountPasswordPolicyRequestFn                func(*iam.GetAccountPasswordPolicyInput) (*request.Request, *iam.GetAccountPasswordPolicyOutput)
	GetAccountPasswordPolicyFn                       func(*iam.GetAccountPasswordPolicyInput) (*iam.GetAccountPasswordPolicyOutput, error)
	GetAccountSummaryRequestFn                       func(*iam.GetAccountSummaryInput) (*request.Request, *iam.GetAccountSummaryOutput)
	GetAccountSummaryFn                              func(*iam.GetAccountSummaryInput) (*iam.GetAccountSummaryOutput, error)
	GetContextKeysForCustomPolicyRequestFn           func(*iam.GetContextKeysForCustomPolicyInput) (*request.Request, *iam.GetContextKeysForPolicyResponse)
	GetContextKeysForCustomPolicyFn                  func(*iam.GetContextKeysForCustomPolicyInput) (*iam.GetContextKeysForPolicyResponse, error)
	GetContextKeysForPrincipalPolicyRequestFn        func(*iam.GetContextKeysForPrincipalPolicyInput) (*request.Request, *iam.GetContextKeysForPolicyResponse)
	GetContextKeysForPrincipalPolicyFn               func(*iam.GetContextKeysForPrincipalPolicyInput) (*iam.GetContextKeysForPolicyResponse, error)
	GetCredentialReportRequestFn                     func(*iam.GetCredentialReportInput) (*request.Request, *iam.GetCredentialReportOutput)
	GetCredentialReportFn                            func(*iam.GetCredentialReportInput) (*iam.GetCredentialReportOutput, error)
	GetGroupRequestFn                                func(*iam.GetGroupInput) (*request.Request, *iam.GetGroupOutput)
	GetGroupFn                                       func(*iam.GetGroupInput) (*iam.GetGroupOutput, error)
	GetGroupPagesFn                                  func(*iam.GetGroupInput, func(*iam.GetGroupOutput, bool) bool) error
	GetGroupPolicyRequestFn                          func(*iam.GetGroupPolicyInput) (*request.Request, *iam.GetGroupPolicyOutput)
	GetGroupPolicyFn                                 func(*iam.GetGroupPolicyInput) (*iam.GetGroupPolicyOutput, error)
	GetInstanceProfileRequestFn                      func(*iam.GetInstanceProfileInput) (*request.Request, *iam.GetInstanceProfileOutput)
	GetInstanceProfileFn                             func(*iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error)
	GetLoginProfileRequestFn                         func(*iam.GetLoginProfileInput) (*request.Request, *iam.GetLoginProfileOutput)
	GetLoginProfileFn                                func(*iam.GetLoginProfileInput) (*iam.GetLoginProfileOutput, error)
	GetOpenIDConnectProviderRequestFn                func(*iam.GetOpenIDConnectProviderInput) (*request.Request, *iam.GetOpenIDConnectProviderOutput)
	GetOpenIDConnectProviderFn                       func(*iam.GetOpenIDConnectProviderInput) (*iam.GetOpenIDConnectProviderOutput, error)
	GetPolicyRequestFn                               func(*iam.GetPolicyInput) (*request.Request, *iam.GetPolicyOutput)
	GetPolicyFn                                      func(*iam.GetPolicyInput) (*iam.GetPolicyOutput, error)
	GetPolicyVersionRequestFn                        func(*iam.GetPolicyVersionInput) (*request.Request, *iam.GetPolicyVersionOutput)
	GetPolicyVersionFn                               func(*iam.GetPolicyVersionInput) (*iam.GetPolicyVersionOutput, error)
	GetRoleRequestFn                                 func(*iam.GetRoleInput) (*request.Request, *iam.GetRoleOutput)
	GetRoleFn                                        func(*iam.GetRoleInput) (*iam.GetRoleOutput, error)
	GetRolePolicyRequestFn                           func(*iam.GetRolePolicyInput) (*request.Request, *iam.GetRolePolicyOutput)
	GetRolePolicyFn                                  func(*iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error)
	GetSAMLProviderRequestFn                         func(*iam.GetSAMLProviderInput) (*request.Request, *iam.GetSAMLProviderOutput)
	GetSAMLProviderFn                                func(*iam.GetSAMLProviderInput) (*iam.GetSAMLProviderOutput, error)
	GetSSHPublicKeyRequestFn                         func(*iam.GetSSHPublicKeyInput) (*request.Request, *iam.GetSSHPublicKeyOutput)
	GetSSHPublicKeyFn                                func(*iam.GetSSHPublicKeyInput) (*iam.GetSSHPublicKeyOutput, error)
	GetServerCertificateRequestFn                    func(*iam.GetServerCertificateInput) (*request.Request, *iam.GetServerCertificateOutput)
	GetServerCertificateFn                           func(*iam.GetServerCertificateInput) (*iam.GetServerCertificateOutput, error)
	GetUserRequestFn                                 func(*iam.GetUserInput) (*request.Request, *iam.GetUserOutput)
	GetUserFn                                        func(*iam.GetUserInput) (*iam.GetUserOutput, error)
	GetUserPolicyRequestFn                           func(*iam.GetUserPolicyInput) (*request.Request, *iam.GetUserPolicyOutput)
	GetUserPolicyFn                                  func(*iam.GetUserPolicyInput) (*iam.GetUserPolicyOutput, error)
	ListAccessKeysRequestFn                          func(*iam.ListAccessKeysInput) (*request.Request, *iam.ListAccessKeysOutput)
	ListAccessKeysFn                                 func(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error)
	ListAccessKeysPagesFn                            func(*iam.ListAccessKeysInput, func(*iam.ListAccessKeysOutput, bool) bool) error
	ListAccountAliasesRequestFn                      func(*iam.ListAccountAliasesInput) (*request.Request, *iam.ListAccountAliasesOutput)
	ListAccountAliasesFn                             func(*iam.ListAccountAliasesInput) (*iam.ListAccountAliasesOutput, error)
	ListAccountAliasesPagesFn                        func(*iam.ListAccountAliasesInput, func(*iam.ListAccountAliasesOutput, bool) bool) error
	ListAttachedGroupPoliciesRequestFn               func(*iam.ListAttachedGroupPoliciesInput) (*request.Request, *iam.ListAttachedGroupPoliciesOutput)
	ListAttachedGroupPoliciesFn                      func(*iam.ListAttachedGroupPoliciesInput) (*iam.ListAttachedGroupPoliciesOutput, error)
	ListAttachedGroupPoliciesPagesFn                 func(*iam.ListAttachedGroupPoliciesInput, func(*iam.ListAttachedGroupPoliciesOutput, bool) bool) error
	ListAttachedRolePoliciesRequestFn                func(*iam.ListAttachedRolePoliciesInput) (*request.Request, *iam.ListAttachedRolePoliciesOutput)
	ListAttachedRolePoliciesFn                       func(*iam.ListAttachedRolePoliciesInput) (*iam.ListAttachedRolePoliciesOutput, error)
	ListAttachedRolePoliciesPagesFn                  func(*iam.ListAttachedRolePoliciesInput, func(*iam.ListAttachedRolePoliciesOutput, bool) bool) error
	ListAttachedUserPoliciesRequestFn                func(*iam.ListAttachedUserPoliciesInput) (*request.Request, *iam.ListAttachedUserPoliciesOutput)
	ListAttachedUserPoliciesFn                       func(*iam.ListAttachedUserPoliciesInput) (*iam.ListAttachedUserPoliciesOutput, error)
	ListAttachedUserPoliciesPagesFn                  func(*iam.ListAttachedUserPoliciesInput, func(*iam.ListAttachedUserPoliciesOutput, bool) bool) error
	ListEntitiesForPolicyRequestFn                   func(*iam.ListEntitiesForPolicyInput) (*request.Request, *iam.ListEntitiesForPolicyOutput)
	ListEntitiesForPolicyFn                          func(*iam.ListEntitiesForPolicyInput) (*iam.ListEntitiesForPolicyOutput, error)
	ListEntitiesForPolicyPagesFn                     func(*iam.ListEntitiesForPolicyInput, func(*iam.ListEntitiesForPolicyOutput, bool) bool) error
	ListGroupPoliciesRequestFn                       func(*iam.ListGroupPoliciesInput) (*request.Request, *iam.ListGroupPoliciesOutput)
	ListGroupPoliciesFn                              func(*iam.ListGroupPoliciesInput) (*iam.ListGroupPoliciesOutput, error)
	ListGroupPoliciesPagesFn                         func(*iam.ListGroupPoliciesInput, func(*iam.ListGroupPoliciesOutput, bool) bool) error
	ListGroupsRequestFn                              func(*iam.ListGroupsInput) (*request.Request, *iam.ListGroupsOutput)
	ListGroupsFn                                     func(*iam.ListGroupsInput) (*iam.ListGroupsOutput, error)
	ListGroupsPagesFn                                func(*iam.ListGroupsInput, func(*iam.ListGroupsOutput, bool) bool) error
	ListGroupsForUserRequestFn                       func(*iam.ListGroupsForUserInput) (*request.Request, *iam.ListGroupsForUserOutput)
	ListGroupsForUserFn                              func(*iam.ListGroupsForUserInput) (*iam.ListGroupsForUserOutput, error)
	ListGroupsForUserPagesFn                         func(*iam.ListGroupsForUserInput, func(*iam.ListGroupsForUserOutput, bool) bool) error
	ListInstanceProfilesRequestFn                    func(*iam.ListInstanceProfilesInput) (*request.Request, *iam.ListInstanceProfilesOutput)
	ListInstanceProfilesFn                           func(*iam.ListInstanceProfilesInput) (*iam.ListInstanceProfilesOutput, error)
	ListInstanceProfilesPagesFn                      func(*iam.ListInstanceProfilesInput, func(*iam.ListInstanceProfilesOutput, bool) bool) error
	ListInstanceProfilesForRoleRequestFn             func(*iam.ListInstanceProfilesForRoleInput) (*request.Request, *iam.ListInstanceProfilesForRoleOutput)
	ListInstanceProfilesForRoleFn                    func(*iam.ListInstanceProfilesForRoleInput) (*iam.ListInstanceProfilesForRoleOutput, error)
	ListInstanceProfilesForRolePagesFn               func(*iam.ListInstanceProfilesForRoleInput, func(*iam.ListInstanceProfilesForRoleOutput, bool) bool) error
	ListMFADevicesRequestFn                          func(*iam.ListMFADevicesInput) (*request.Request, *iam.ListMFADevicesOutput)
	ListMFADevicesFn                                 func(*iam.ListMFADevicesInput) (*iam.ListMFADevicesOutput, error)
	ListMFADevicesPagesFn                            func(*iam.ListMFADevicesInput, func(*iam.ListMFADevicesOutput, bool) bool) error
	ListOpenIDConnectProvidersRequestFn              func(*iam.ListOpenIDConnectProvidersInput) (*request.Request, *iam.ListOpenIDConnectProvidersOutput)
	ListOpenIDConnectProvidersFn                     func(*iam.ListOpenIDConnectProvidersInput) (*iam.ListOpenIDConnectProvidersOutput, error)
	ListPoliciesRequestFn                            func(*iam.ListPoliciesInput) (*request.Request, *iam.ListPoliciesOutput)
	ListPoliciesFn                                   func(*iam.ListPoliciesInput) (*iam.ListPoliciesOutput, error)
	ListPoliciesPagesFn                              func(*iam.ListPoliciesInput, func(*iam.ListPoliciesOutput, bool) bool) error
	ListPolicyVersionsRequestFn                      func(*iam.ListPolicyVersionsInput) (*request.Request, *iam.ListPolicyVersionsOutput)
	ListPolicyVersionsFn                             func(*iam.ListPolicyVersionsInput) (*iam.ListPolicyVersionsOutput, error)
	ListRolePoliciesRequestFn                        func(*iam.ListRolePoliciesInput) (*request.Request, *iam.ListRolePoliciesOutput)
	ListRolePoliciesFn                               func(*iam.ListRolePoliciesInput) (*iam.ListRolePoliciesOutput, error)
	ListRolePoliciesPagesFn                          func(*iam.ListRolePoliciesInput, func(*iam.ListRolePoliciesOutput, bool) bool) error
	ListRolesRequestFn                               func(*iam.ListRolesInput) (*request.Request, *iam.ListRolesOutput)
	ListRolesFn                                      func(*iam.ListRolesInput) (*iam.ListRolesOutput, error)
	ListRolesPagesFn                                 func(*iam.ListRolesInput, func(*iam.ListRolesOutput, bool) bool) error
	ListSAMLProvidersRequestFn                       func(*iam.ListSAMLProvidersInput) (*request.Request, *iam.ListSAMLProvidersOutput)
	ListSAMLProvidersFn                              func(*iam.ListSAMLProvidersInput) (*iam.ListSAMLProvidersOutput, error)
	ListSSHPublicKeysRequestFn                       func(*iam.ListSSHPublicKeysInput) (*request.Request, *iam.ListSSHPublicKeysOutput)
	ListSSHPublicKeysFn                              func(*iam.ListSSHPublicKeysInput) (*iam.ListSSHPublicKeysOutput, error)
	ListServerCertificatesRequestFn                  func(*iam.ListServerCertificatesInput) (*request.Request, *iam.ListServerCertificatesOutput)
	ListServerCertificatesFn                         func(*iam.ListServerCertificatesInput) (*iam.ListServerCertificatesOutput, error)
	ListServerCertificatesPagesFn                    func(*iam.ListServerCertificatesInput, func(*iam.ListServerCertificatesOutput, bool) bool) error
	ListSigningCertificatesRequestFn                 func(*iam.ListSigningCertificatesInput) (*request.Request, *iam.ListSigningCertificatesOutput)
	ListSigningCertificatesFn                        func(*iam.ListSigningCertificatesInput) (*iam.ListSigningCertificatesOutput, error)
	ListSigningCertificatesPagesFn                   func(*iam.ListSigningCertificatesInput, func(*iam.ListSigningCertificatesOutput, bool) bool) error
	ListUserPoliciesRequestFn                        func(*iam.ListUserPoliciesInput) (*request.Request, *iam.ListUserPoliciesOutput)
	ListUserPoliciesFn                               func(*iam.ListUserPoliciesInput) (*iam.ListUserPoliciesOutput, error)
	ListUserPoliciesPagesFn                          func(*iam.ListUserPoliciesInput, func(*iam.ListUserPoliciesOutput, bool) bool) error
	ListUsersRequestFn                               func(*iam.ListUsersInput) (*request.Request, *iam.ListUsersOutput)
	ListUsersFn                                      func(*iam.ListUsersInput) (*iam.ListUsersOutput, error)
	ListUsersPagesFn                                 func(*iam.ListUsersInput, func(*iam.ListUsersOutput, bool) bool) error
	ListVirtualMFADevicesRequestFn                   func(*iam.ListVirtualMFADevicesInput) (*request.Request, *iam.ListVirtualMFADevicesOutput)
	ListVirtualMFADevicesFn                          func(*iam.ListVirtualMFADevicesInput) (*iam.ListVirtualMFADevicesOutput, error)
	ListVirtualMFADevicesPagesFn                     func(*iam.ListVirtualMFADevicesInput, func(*iam.ListVirtualMFADevicesOutput, bool) bool) error
	PutGroupPolicyRequestFn                          func(*iam.PutGroupPolicyInput) (*request.Request, *iam.PutGroupPolicyOutput)
	PutGroupPolicyFn                                 func(*iam.PutGroupPolicyInput) (*iam.PutGroupPolicyOutput, error)
	PutRolePolicyRequestFn                           func(*iam.PutRolePolicyInput) (*request.Request, *iam.PutRolePolicyOutput)
	PutRolePolicyFn                                  func(*iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error)
	PutUserPolicyRequestFn                           func(*iam.PutUserPolicyInput) (*request.Request, *iam.PutUserPolicyOutput)
	PutUserPolicyFn                                  func(*iam.PutUserPolicyInput) (*iam.PutUserPolicyOutput, error)
	RemoveClientIDFromOpenIDConnectProviderRequestFn func(*iam.RemoveClientIDFromOpenIDConnectProviderInput) (*request.Request, *iam.RemoveClientIDFromOpenIDConnectProviderOutput)
	RemoveClientIDFromOpenIDConnectProviderFn        func(*iam.RemoveClientIDFromOpenIDConnectProviderInput) (*iam.RemoveClientIDFromOpenIDConnectProviderOutput, error)
	RemoveRoleFromInstanceProfileRequestFn           func(*iam.RemoveRoleFromInstanceProfileInput) (*request.Request, *iam.RemoveRoleFromInstanceProfileOutput)
	RemoveRoleFromInstanceProfileFn                  func(*iam.RemoveRoleFromInstanceProfileInput) (*iam.RemoveRoleFromInstanceProfileOutput, error)
	RemoveUserFromGroupRequestFn                     func(*iam.RemoveUserFromGroupInput) (*request.Request, *iam.RemoveUserFromGroupOutput)
	RemoveUserFromGroupFn                            func(*iam.RemoveUserFromGroupInput) (*iam.RemoveUserFromGroupOutput, error)
	ResyncMFADeviceRequestFn                         func(*iam.ResyncMFADeviceInput) (*request.Request, *iam.ResyncMFADeviceOutput)
	ResyncMFADeviceFn                                func(*iam.ResyncMFADeviceInput) (*iam.ResyncMFADeviceOutput, error)
	SetDefaultPolicyVersionRequestFn                 func(*iam.SetDefaultPolicyVersionInput) (*request.Request, *iam.SetDefaultPolicyVersionOutput)
	SetDefaultPolicyVersionFn                        func(*iam.SetDefaultPolicyVersionInput) (*iam.SetDefaultPolicyVersionOutput, error)
	SimulateCustomPolicyRequestFn                    func(*iam.SimulateCustomPolicyInput) (*request.Request, *iam.SimulatePolicyResponse)
	SimulateCustomPolicyFn                           func(*iam.SimulateCustomPolicyInput) (*iam.SimulatePolicyResponse, error)
	SimulatePrincipalPolicyRequestFn                 func(*iam.SimulatePrincipalPolicyInput) (*request.Request, *iam.SimulatePolicyResponse)
	SimulatePrincipalPolicyFn                        func(*iam.SimulatePrincipalPolicyInput) (*iam.SimulatePolicyResponse, error)
	UpdateAccessKeyRequestFn                         func(*iam.UpdateAccessKeyInput) (*request.Request, *iam.UpdateAccessKeyOutput)
	UpdateAccessKeyFn                                func(*iam.UpdateAccessKeyInput) (*iam.UpdateAccessKeyOutput, error)
	UpdateAccountPasswordPolicyRequestFn             func(*iam.UpdateAccountPasswordPolicyInput) (*request.Request, *iam.UpdateAccountPasswordPolicyOutput)
	UpdateAccountPasswordPolicyFn                    func(*iam.UpdateAccountPasswordPolicyInput) (*iam.UpdateAccountPasswordPolicyOutput, error)
	UpdateAssumeRolePolicyRequestFn                  func(*iam.UpdateAssumeRolePolicyInput) (*request.Request, *iam.UpdateAssumeRolePolicyOutput)
	UpdateAssumeRolePolicyFn                         func(*iam.UpdateAssumeRolePolicyInput) (*iam.UpdateAssumeRolePolicyOutput, error)
	UpdateGroupRequestFn                             func(*iam.UpdateGroupInput) (*request.Request, *iam.UpdateGroupOutput)
	UpdateGroupFn                                    func(*iam.UpdateGroupInput) (*iam.UpdateGroupOutput, error)
	UpdateLoginProfileRequestFn                      func(*iam.UpdateLoginProfileInput) (*request.Request, *iam.UpdateLoginProfileOutput)
	UpdateLoginProfileFn                             func(*iam.UpdateLoginProfileInput) (*iam.UpdateLoginProfileOutput, error)
	UpdateOpenIDConnectProviderThumbprintRequestFn   func(*iam.UpdateOpenIDConnectProviderThumbprintInput) (*request.Request, *iam.UpdateOpenIDConnectProviderThumbprintOutput)
	UpdateOpenIDConnectProviderThumbprintFn          func(*iam.UpdateOpenIDConnectProviderThumbprintInput) (*iam.UpdateOpenIDConnectProviderThumbprintOutput, error)
	UpdateSAMLProviderRequestFn                      func(*iam.UpdateSAMLProviderInput) (*request.Request, *iam.UpdateSAMLProviderOutput)
	UpdateSAMLProviderFn                             func(*iam.UpdateSAMLProviderInput) (*iam.UpdateSAMLProviderOutput, error)
	UpdateSSHPublicKeyRequestFn                      func(*iam.UpdateSSHPublicKeyInput) (*request.Request, *iam.UpdateSSHPublicKeyOutput)
	UpdateSSHPublicKeyFn                             func(*iam.UpdateSSHPublicKeyInput) (*iam.UpdateSSHPublicKeyOutput, error)
	UpdateServerCertificateRequestFn                 func(*iam.UpdateServerCertificateInput) (*request.Request, *iam.UpdateServerCertificateOutput)
	UpdateServerCertificateFn                        func(*iam.UpdateServerCertificateInput) (*iam.UpdateServerCertificateOutput, error)
	UpdateSigningCertificateRequestFn                func(*iam.UpdateSigningCertificateInput) (*request.Request, *iam.UpdateSigningCertificateOutput)
	UpdateSigningCertificateFn                       func(*iam.UpdateSigningCertificateInput) (*iam.UpdateSigningCertificateOutput, error)
	UpdateUserRequestFn                              func(*iam.UpdateUserInput) (*request.Request, *iam.UpdateUserOutput)
	UpdateUserFn                                     func(*iam.UpdateUserInput) (*iam.UpdateUserOutput, error)
	UploadSSHPublicKeyRequestFn                      func(*iam.UploadSSHPublicKeyInput) (*request.Request, *iam.UploadSSHPublicKeyOutput)
	UploadSSHPublicKeyFn                             func(*iam.UploadSSHPublicKeyInput) (*iam.UploadSSHPublicKeyOutput, error)
	UploadServerCertificateRequestFn                 func(*iam.UploadServerCertificateInput) (*request.Request, *iam.UploadServerCertificateOutput)
	UploadServerCertificateFn                        func(*iam.UploadServerCertificateInput) (*iam.UploadServerCertificateOutput, error)
	UploadSigningCertificateRequestFn                func(*iam.UploadSigningCertificateInput) (*request.Request, *iam.UploadSigningCertificateOutput)
	UploadSigningCertificateFn                       func(*iam.UploadSigningCertificateInput) (*iam.UploadSigningCertificateOutput, error)
}

func (f *IAM) AddClientIDToOpenIDConnectProviderRequest(input *iam.AddClientIDToOpenIDConnectProviderInput) (*request.Request, *iam.AddClientIDToOpenIDConnectProviderOutput) {
	if f.AddClientIDToOpenIDConnectProviderRequestFn == nil {
		return nil, nil
	}
	return f.AddClientIDToOpenIDConnectProviderRequestFn(input)
}

func (f *IAM) AddClientIDToOpenIDConnectProvider(input *iam.AddClientIDToOpenIDConnectProviderInput) (*iam.AddClientIDToOpenIDConnectProviderOutput, error) {
	if f.AddClientIDToOpenIDConnectProviderFn == nil {
		return nil, nil
	}
	return f.AddClientIDToOpenIDConnectProviderFn(input)
}

func (f *IAM) AddRoleToInstanceProfileRequest(input *iam.AddRoleToInstanceProfileInput) (*request.Request, *iam.AddRoleToInstanceProfileOutput) {
	if f.AddRoleToInstanceProfileRequestFn == nil {
		return nil, nil
	}
	return f.AddRoleToInstanceProfileRequestFn(input)
}

func (f *IAM) AddRoleToInstanceProfile(input *iam.AddRoleToInstanceProfileInput) (*iam.AddRoleToInstanceProfileOutput, error) {
	if f.AddRoleToInstanceProfileFn == nil {
		return nil, nil
	}
	return f.AddRoleToInstanceProfileFn(input)
}

func (f *IAM) AddUserToGroupRequest(input *iam.AddUserToGroupInput) (*request.Request, *iam.AddUserToGroupOutput) {
	if f.AddUserToGroupRequestFn == nil {
		return nil, nil
	}
	return f.AddUserToGroupRequestFn(input)
}

func (f *IAM) AddUserToGroup(input *iam.AddUserToGroupInput) (*iam.AddUserToGroupOutput, error) {
	if f.AddUserToGroupFn == nil {
		return nil, nil
	}
	return f.AddUserToGroupFn(input)
}

func (f *IAM) AttachGroupPolicyRequest(input *iam.AttachGroupPolicyInput) (*request.Request, *iam.AttachGroupPolicyOutput) {
	if f.AttachGroupPolicyRequestFn == nil {
		return nil, nil
	}
	return f.AttachGroupPolicyRequestFn(input)
}

func (f *IAM) AttachGroupPolicy(input *iam.AttachGroupPolicyInput) (*iam.AttachGroupPolicyOutput, error) {
	if f.AttachGroupPolicyFn == nil {
		return nil, nil
	}
	return f.AttachGroupPolicyFn(input)
}

func (f *IAM) AttachRolePolicyRequest(input *iam.AttachRolePolicyInput) (*request.Request, *iam.AttachRolePolicyOutput) {
	if f.AttachRolePolicyRequestFn == nil {
		return nil, nil
	}
	return f.AttachRolePolicyRequestFn(input)
}

func (f *IAM) AttachRolePolicy(input *iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
	if f.AttachRolePolicyFn == nil {
		return nil, nil
	}
	return f.AttachRolePolicyFn(input)
}

func (f *IAM) AttachUserPolicyRequest(input *iam.AttachUserPolicyInput) (*request.Request, *iam.AttachUserPolicyOutput) {
	if f.AttachUserPolicyRequestFn == nil {
		return nil, nil
	}
	return f.AttachUserPolicyRequestFn(input)
}

func (f *IAM) AttachUserPolicy(input *iam.AttachUserPolicyInput) (*iam.AttachUserPolicyOutput, error) {
	if f.AttachUserPolicyFn == nil {
		return nil, nil
	}
	return f.AttachUserPolicyFn(input)
}

func (f *IAM) ChangePasswordRequest(input *iam.ChangePasswordInput) (*request.Request, *iam.ChangePasswordOutput) {
	if f.ChangePasswordRequestFn == nil {
		return nil, nil
	}
	return f.ChangePasswordRequestFn(input)
}

func (f *IAM) ChangePassword(input *iam.ChangePasswordInput) (*iam.ChangePasswordOutput, error) {
	if f.ChangePasswordFn == nil {
		return nil, nil
	}
	return f.ChangePasswordFn(input)
}

func (f *IAM) CreateAccessKeyRequest(input *iam.CreateAccessKeyInput) (*request.Request, *iam.CreateAccessKeyOutput) {
	if f.CreateAccessKeyRequestFn == nil {
		return nil, nil
	}
	return f.CreateAccessKeyRequestFn(input)
}

func (f *IAM) CreateAccessKey(input *iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	if f.CreateAccessKeyFn == nil {
		return nil, nil
	}
	return f.CreateAccessKeyFn(input)
}

func (f *IAM) CreateAccountAliasRequest(input *iam.CreateAccountAliasInput) (*request.Request, *iam.CreateAccountAliasOutput) {
	if f.CreateAccountAliasRequestFn == nil {
		return nil, nil
	}
	return f.CreateAccountAliasRequestFn(input)
}

func (f *IAM) CreateAccountAlias(input *iam.CreateAccountAliasInput) (*iam.CreateAccountAliasOutput, error) {
	if f.CreateAccountAliasFn == nil {
		return nil, nil
	}
	return f.CreateAccountAliasFn(input)
}

func (f *IAM) CreateGroupRequest(input *iam.CreateGroupInput) (*request.Request, *iam.CreateGroupOutput) {
	if f.CreateGroupRequestFn == nil {
		return nil, nil
	}
	return f.CreateGroupRequestFn(input)
}

func (f *IAM) CreateGroup(input *iam.CreateGroupInput) (*iam.CreateGroupOutput, error) {
	if f.CreateGroupFn == nil {
		return nil, nil
	}
	return f.CreateGroupFn(input)
}

func (f *IAM) CreateInstanceProfileRequest(input *iam.CreateInstanceProfileInput) (*request.Request, *iam.CreateInstanceProfileOutput) {
	if f.CreateInstanceProfileRequestFn == nil {
		return nil, nil
	}
	return f.CreateInstanceProfileRequestFn(input)
}

func (f *IAM) CreateInstanceProfile(input *iam.CreateInstanceProfileInput) (*iam.CreateInstanceProfileOutput, error) {
	if f.CreateInstanceProfileFn == nil {
		return nil, nil
	}
	return f.CreateInstanceProfileFn(input)
}

func (f *IAM) CreateLoginProfileRequest(input *iam.CreateLoginProfileInput) (*request.Request, *iam.CreateLoginProfileOutput) {
	if f.CreateLoginProfileRequestFn == nil {
		return nil, nil
	}
	return f.CreateLoginProfileRequestFn(input)
}

func (f *IAM) CreateLoginProfile(input *iam.CreateLoginProfileInput) (*iam.CreateLoginProfileOutput, error) {
	if f.CreateLoginProfileFn == nil {
		return nil, nil
	}
	return f.CreateLoginProfileFn(input)
}

func (f *IAM) CreateOpenIDConnectProviderRequest(input *iam.CreateOpenIDConnectProviderInput) (*request.Request, *iam.CreateOpenIDConnectProviderOutput) {
	if f.CreateOpenIDConnectProviderRequestFn == nil {
		return nil, nil
	}
	return f.CreateOpenIDConnectProviderRequestFn(input)
}

func (f *IAM) CreateOpenIDConnectProvider(input *iam.CreateOpenIDConnectProviderInput) (*iam.CreateOpenIDConnectProviderOutput, error) {
	if f.CreateOpenIDConnectProviderFn == nil {
		return nil, nil
	}
	return f.CreateOpenIDConnectProviderFn(input)
}

func (f *IAM) CreatePolicyRequest(input *iam.CreatePolicyInput) (*request.Request, *iam.CreatePolicyOutput) {
	if f.CreatePolicyRequestFn == nil {
		return nil, nil
	}
	return f.CreatePolicyRequestFn(input)
}

func (f *IAM) CreatePolicy(input *iam.CreatePolicyInput) (*iam.CreatePolicyOutput, error) {
	if f.CreatePolicyFn == nil {
		return nil, nil
	}
	return f.CreatePolicyFn(input)
}

func (f *IAM) CreatePolicyVersionRequest(input *iam.CreatePolicyVersionInput) (*request.Request, *iam.CreatePolicyVersionOutput) {
	if f.CreatePolicyVersionRequestFn == nil {
		return nil, nil
	}
	return f.CreatePolicyVersionRequestFn(input)
}

func (f *IAM) CreatePolicyVersion(input *iam.CreatePolicyVersionInput) (*iam.CreatePolicyVersionOutput, error) {
	if f.CreatePolicyVersionFn == nil {
		return nil, nil
	}
	return f.CreatePolicyVersionFn(input)
}

func (f *IAM) CreateRoleRequest(input *iam.CreateRoleInput) (*request.Request, *iam.CreateRoleOutput) {
	if f.CreateRoleRequestFn == nil {
		return nil, nil
	}
	return f.CreateRoleRequestFn(input)
}

func (f *IAM) CreateRole(input *iam.CreateRoleInput) (*iam.CreateRoleOutput, error) {
	if f.CreateRoleFn == nil {
		return nil, nil
	}
	return f.CreateRoleFn(input)
}

func (f *IAM) CreateSAMLProviderRequest(input *iam.CreateSAMLProviderInput) (*request.Request, *iam.CreateSAMLProviderOutput) {
	if f.CreateSAMLProviderRequestFn == nil {
		return nil, nil
	}
	return f.CreateSAMLProviderRequestFn(input)
}

func (f *IAM) CreateSAMLProvider(input *iam.CreateSAMLProviderInput) (*iam.CreateSAMLProviderOutput, error) {
	if f.CreateSAMLProviderFn == nil {
		return nil, nil
	}
	return f.CreateSAMLProviderFn(input)
}

func (f *IAM) CreateUserRequest(input *iam.CreateUserInput) (*request.Request, *iam.CreateUserOutput) {
	if f.CreateUserRequestFn == nil {
		return nil, nil
	}
	return f.CreateUserRequestFn(input)
}

func (f *IAM) CreateUser(input *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
	if f.CreateUserFn == nil {
		return nil, nil
	}
	return f.CreateUserFn(input)
}

func (f *IAM) CreateVirtualMFADeviceRequest(input *iam.CreateVirtualMFADeviceInput) (*request.Request, *iam.CreateVirtualMFADeviceOutput) {
	if f.CreateVirtualMFADeviceRequestFn == nil {
		return nil, nil
	}
	return f.CreateVirtualMFADeviceRequestFn(input)
}

func (f *IAM) CreateVirtualMFADevice(input *iam.CreateVirtualMFADeviceInput) (*iam.CreateVirtualMFADeviceOutput, error) {
	if f.CreateVirtualMFADeviceFn == nil {
		return nil, nil
	}
	return f.CreateVirtualMFADeviceFn(input)
}

func (f *IAM) DeactivateMFADeviceRequest(input *iam.DeactivateMFADeviceInput) (*request.Request, *iam.DeactivateMFADeviceOutput) {
	if f.DeactivateMFADeviceRequestFn == nil {
		return nil, nil
	}
	return f.DeactivateMFADeviceRequestFn(input)
}

func (f *IAM) DeactivateMFADevice(input *iam.DeactivateMFADeviceInput) (*iam.DeactivateMFADeviceOutput, error) {
	if f.DeactivateMFADeviceFn == nil {
		return nil, nil
	}
	return f.DeactivateMFADeviceFn(input)
}

func (f *IAM) DeleteAccessKeyRequest(input *iam.DeleteAccessKeyInput) (*request.Request, *iam.DeleteAccessKeyOutput) {
	if f.DeleteAccessKeyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteAccessKeyRequestFn(input)
}

func (f *IAM) DeleteAccessKey(input *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	if f.DeleteAccessKeyFn == nil {
		return nil, nil
	}
	return f.DeleteAccessKeyFn(input)
}

func (f *IAM) DeleteAccountAliasRequest(input *iam.DeleteAccountAliasInput) (*request.Request, *iam.DeleteAccountAliasOutput) {
	if f.DeleteAccountAliasRequestFn == nil {
		return nil, nil
	}
	return f.DeleteAccountAliasRequestFn(input)
}

func (f *IAM) DeleteAccountAlias(input *iam.DeleteAccountAliasInput) (*iam.DeleteAccountAliasOutput, error) {
	if f.DeleteAccountAliasFn == nil {
		return nil, nil
	}
	return f.DeleteAccountAliasFn(input)
}

func (f *IAM) DeleteAccountPasswordPolicyRequest(input *iam.DeleteAccountPasswordPolicyInput) (*request.Request, *iam.DeleteAccountPasswordPolicyOutput) {
	if f.DeleteAccountPasswordPolicyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteAccountPasswordPolicyRequestFn(input)
}

func (f *IAM) DeleteAccountPasswordPolicy(input *iam.DeleteAccountPasswordPolicyInput) (*iam.DeleteAccountPasswordPolicyOutput, error) {
	if f.DeleteAccountPasswordPolicyFn == nil {
		return nil, nil
	}
	return f.DeleteAccountPasswordPolicyFn(input)
}

func (f *IAM) DeleteGroupRequest(input *iam.DeleteGroupInput) (*request.Request, *iam.DeleteGroupOutput) {
	if f.DeleteGroupRequestFn == nil {
		return nil, nil
	}
	return f.DeleteGroupRequestFn(input)
}

func (f *IAM) DeleteGroup(input *iam.DeleteGroupInput) (*iam.DeleteGroupOutput, error) {
	if f.DeleteGroupFn == nil {
		return nil, nil
	}
	return f.DeleteGroupFn(input)
}

func (f *IAM) DeleteGroupPolicyRequest(input *iam.DeleteGroupPolicyInput) (*request.Request, *iam.DeleteGroupPolicyOutput) {
	if f.DeleteGroupPolicyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteGroupPolicyRequestFn(input)
}

func (f *IAM) DeleteGroupPolicy(input *iam.DeleteGroupPolicyInput) (*iam.DeleteGroupPolicyOutput, error) {
	if f.DeleteGroupPolicyFn == nil {
		return nil, nil
	}
	return f.DeleteGroupPolicyFn(input)
}

func (f *IAM) DeleteInstanceProfileRequest(input *iam.DeleteInstanceProfileInput) (*request.Request, *iam.DeleteInstanceProfileOutput) {
	if f.DeleteInstanceProfileRequestFn == nil {
		return nil, nil
	}
	return f.DeleteInstanceProfileRequestFn(input)
}

func (f *IAM) DeleteInstanceProfile(input *iam.DeleteInstanceProfileInput) (*iam.DeleteInstanceProfileOutput, error) {
	if f.DeleteInstanceProfileFn == nil {
		return nil, nil
	}
	return f.DeleteInstanceProfileFn(input)
}

func (f *IAM) DeleteLoginProfileRequest(input *iam.DeleteLoginProfileInput) (*request.Request, *iam.DeleteLoginProfileOutput) {
	if f.DeleteLoginProfileRequestFn == nil {
		return nil, nil
	}
	return f.DeleteLoginProfileRequestFn(input)
}

func (f *IAM) DeleteLoginProfile(input *iam.DeleteLoginProfileInput) (*iam.DeleteLoginProfileOutput, error) {
	if f.DeleteLoginProfileFn == nil {
		return nil, nil
	}
	return f.DeleteLoginProfileFn(input)
}

func (f *IAM) DeleteOpenIDConnectProviderRequest(input *iam.DeleteOpenIDConnectProviderInput) (*request.Request, *iam.DeleteOpenIDConnectProviderOutput) {
	if f.DeleteOpenIDConnectProviderRequestFn == nil {
		return nil, nil
	}
	return f.DeleteOpenIDConnectProviderRequestFn(input)
}

func (f *IAM) DeleteOpenIDConnectProvider(input *iam.DeleteOpenIDConnectProviderInput) (*iam.DeleteOpenIDConnectProviderOutput, error) {
	if f.DeleteOpenIDConnectProviderFn == nil {
		return nil, nil
	}
	return f.DeleteOpenIDConnectProviderFn(input)
}

func (f *IAM) DeletePolicyRequest(input *iam.DeletePolicyInput) (*request.Request, *iam.DeletePolicyOutput) {
	if f.DeletePolicyRequestFn == nil {
		return nil, nil
	}
	return f.DeletePolicyRequestFn(input)
}

func (f *IAM) DeletePolicy(input *iam.DeletePolicyInput) (*iam.DeletePolicyOutput, error) {
	if f.DeletePolicyFn == nil {
		return nil, nil
	}
	return f.DeletePolicyFn(input)
}

func (f *IAM) DeletePolicyVersionRequest(input *iam.DeletePolicyVersionInput) (*request.Request, *iam.DeletePolicyVersionOutput) {
	if f.DeletePolicyVersionRequestFn == nil {
		return nil, nil
	}
	return f.DeletePolicyVersionRequestFn(input)
}

func (f *IAM) DeletePolicyVersion(input *iam.DeletePolicyVersionInput) (*iam.DeletePolicyVersionOutput, error) {
	if f.DeletePolicyVersionFn == nil {
		return nil, nil
	}
	return f.DeletePolicyVersionFn(input)
}

func (f *IAM) DeleteRoleRequest(input *iam.DeleteRoleInput) (*request.Request, *iam.DeleteRoleOutput) {
	if f.DeleteRoleRequestFn == nil {
		return nil, nil
	}
	return f.DeleteRoleRequestFn(input)
}

func (f *IAM) DeleteRole(input *iam.DeleteRoleInput) (*iam.DeleteRoleOutput, error) {
	if f.DeleteRoleFn == nil {
		return nil, nil
	}
	return f.DeleteRoleFn(input)
}

func (f *IAM) DeleteRolePolicyRequest(input *iam.DeleteRolePolicyInput) (*request.Request, *iam.DeleteRolePolicyOutput) {
	if f.DeleteRolePolicyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteRolePolicyRequestFn(input)
}

func (f *IAM) DeleteRolePolicy(input *iam.DeleteRolePolicyInput) (*iam.DeleteRolePolicyOutput, error) {
	if f.DeleteRolePolicyFn == nil {
		return nil, nil
	}
	return f.DeleteRolePolicyFn(input)
}

func (f *IAM) DeleteSAMLProviderRequest(input *iam.DeleteSAMLProviderInput) (*request.Request, *iam.DeleteSAMLProviderOutput) {
	if f.DeleteSAMLProviderRequestFn == nil {
		return nil, nil
	}
	return f.DeleteSAMLProviderRequestFn(input)
}

func (f *IAM) DeleteSAMLProvider(input *iam.DeleteSAMLProviderInput) (*iam.DeleteSAMLProviderOutput, error) {
	if f.DeleteSAMLProviderFn == nil {
		return nil, nil
	}
	return f.DeleteSAMLProviderFn(input)
}

func (f *IAM) DeleteSSHPublicKeyRequest(input *iam.DeleteSSHPublicKeyInput) (*request.Request, *iam.DeleteSSHPublicKeyOutput) {
	if f.DeleteSSHPublicKeyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteSSHPublicKeyRequestFn(input)
}

func (f *IAM) DeleteSSHPublicKey(input *iam.DeleteSSHPublicKeyInput) (*iam.DeleteSSHPublicKeyOutput, error) {
	if f.DeleteSSHPublicKeyFn == nil {
		return nil, nil
	}
	return f.DeleteSSHPublicKeyFn(input)
}

func (f *IAM) DeleteServerCertificateRequest(input *iam.DeleteServerCertificateInput) (*request.Request, *iam.DeleteServerCertificateOutput) {
	if f.DeleteServerCertificateRequestFn == nil {
		return nil, nil
	}
	return f.DeleteServerCertificateRequestFn(input)
}

func (f *IAM) DeleteServerCertificate(input *iam.DeleteServerCertificateInput) (*iam.DeleteServerCertificateOutput, error) {
	if f.DeleteServerCertificateFn == nil {
		return nil, nil
	}
	return f.DeleteServerCertificateFn(input)
}

func (f *IAM) DeleteSigningCertificateRequest(input *iam.DeleteSigningCertificateInput) (*request.Request, *iam.DeleteSigningCertificateOutput) {
	if f.DeleteSigningCertificateRequestFn == nil {
		return nil, nil
	}
	return f.DeleteSigningCertificateRequestFn(input)
}

func (f *IAM) DeleteSigningCertificate(input *iam.DeleteSigningCertificateInput) (*iam.DeleteSigningCertificateOutput, error) {
	if f.DeleteSigningCertificateFn == nil {
		return nil, nil
	}
	return f.DeleteSigningCertificateFn(input)
}

func (f *IAM) DeleteUserRequest(input *iam.DeleteUserInput) (*request.Request, *iam.DeleteUserOutput) {
	if f.DeleteUserRequestFn == nil {
		return nil, nil
	}
	return f.DeleteUserRequestFn(input)
}

func (f *IAM) DeleteUser(input *iam.DeleteUserInput) (*iam.DeleteUserOutput, error) {
	if f.DeleteUserFn == nil {
		return nil, nil
	}
	return f.DeleteUserFn(input)
}

func (f *IAM) DeleteUserPolicyRequest(input *iam.DeleteUserPolicyInput) (*request.Request, *iam.DeleteUserPolicyOutput) {
	if f.DeleteUserPolicyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteUserPolicyRequestFn(input)
}

func (f *IAM) DeleteUserPolicy(input *iam.DeleteUserPolicyInput) (*iam.DeleteUserPolicyOutput, error) {
	if f.DeleteUserPolicyFn == nil {
		return nil, nil
	}
	return f.DeleteUserPolicyFn(input)
}

func (f *IAM) DeleteVirtualMFADeviceRequest(input *iam.DeleteVirtualMFADeviceInput) (*request.Request, *iam.DeleteVirtualMFADeviceOutput) {
	if f.DeleteVirtualMFADeviceRequestFn == nil {
		return nil, nil
	}
	return f.DeleteVirtualMFADeviceRequestFn(input)
}

func (f *IAM) DeleteVirtualMFADevice(input *iam.DeleteVirtualMFADeviceInput) (*iam.DeleteVirtualMFADeviceOutput, error) {
	if f.DeleteVirtualMFADeviceFn == nil {
		return nil, nil
	}
	return f.DeleteVirtualMFADeviceFn(input)
}

func (f *IAM) DetachGroupPolicyRequest(input *iam.DetachGroupPolicyInput) (*request.Request, *iam.DetachGroupPolicyOutput) {
	if f.DetachGroupPolicyRequestFn == nil {
		return nil, nil
	}
	return f.DetachGroupPolicyRequestFn(input)
}

func (f *IAM) DetachGroupPolicy(input *iam.DetachGroupPolicyInput) (*iam.DetachGroupPolicyOutput, error) {
	if f.DetachGroupPolicyFn == nil {
		return nil, nil
	}
	return f.DetachGroupPolicyFn(input)
}

func (f *IAM) DetachRolePolicyRequest(input *iam.DetachRolePolicyInput) (*request.Request, *iam.DetachRolePolicyOutput) {
	if f.DetachRolePolicyRequestFn == nil {
		return nil, nil
	}
	return f.DetachRolePolicyRequestFn(input)
}

func (f *IAM) DetachRolePolicy(input *iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error) {
	if f.DetachRolePolicyFn == nil {
		return nil, nil
	}
	return f.DetachRolePolicyFn(input)
}

func (f *IAM) DetachUserPolicyRequest(input *iam.DetachUserPolicyInput) (*request.Request, *iam.DetachUserPolicyOutput) {
	if f.DetachUserPolicyRequestFn == nil {
		return nil, nil
	}
	return f.DetachUserPolicyRequestFn(input)
}

func (f *IAM) DetachUserPolicy(input *iam.DetachUserPolicyInput) (*iam.DetachUserPolicyOutput, error) {
	if f.DetachUserPolicyFn == nil {
		return nil, nil
	}
	return f.DetachUserPolicyFn(input)
}

func (f *IAM) EnableMFADeviceRequest(input *iam.EnableMFADeviceInput) (*request.Request, *iam.EnableMFADeviceOutput) {
	if f.EnableMFADeviceRequestFn == nil {
		return nil, nil
	}
	return f.EnableMFADeviceRequestFn(input)
}

func (f *IAM) EnableMFADevice(input *iam.EnableMFADeviceInput) (*iam.EnableMFADeviceOutput, error) {
	if f.EnableMFADeviceFn == nil {
		return nil, nil
	}
	return f.EnableMFADeviceFn(input)
}

func (f *IAM) GenerateCredentialReportRequest(input *iam.GenerateCredentialReportInput) (*request.Request, *iam.GenerateCredentialReportOutput) {
	if f.GenerateCredentialReportRequestFn == nil {
		return nil, nil
	}
	return f.GenerateCredentialReportRequestFn(input)
}

func (f *IAM) GenerateCredentialReport(input *iam.GenerateCredentialReportInput) (*iam.GenerateCredentialReportOutput, error) {
	if f.GenerateCredentialReportFn == nil {
		return nil, nil
	}
	return f.GenerateCredentialReportFn(input)
}

func (f *IAM) GetAccessKeyLastUsedRequest(input *iam.GetAccessKeyLastUsedInput) (*request.Request, *iam.GetAccessKeyLastUsedOutput) {
	if f.GetAccessKeyLastUsedRequestFn == nil {
		return nil, nil
	}
	return f.GetAccessKeyLastUsedRequestFn(input)
}

func (f *IAM) GetAccessKeyLastUsed(input *iam.GetAccessKeyLastUsedInput) (*iam.GetAccessKeyLastUsedOutput, error) {
	if f.GetAccessKeyLastUsedFn == nil {
		return nil, nil
	}
	return f.GetAccessKeyLastUsedFn(input)
}

func (f *IAM) GetAccountAuthorizationDetailsRequest(input *iam.GetAccountAuthorizationDetailsInput) (*request.Request, *iam.GetAccountAuthorizationDetailsOutput) {
	if f.GetAccountAuthorizationDetailsRequestFn == nil {
		return nil, nil
	}
	return f.GetAccountAuthorizationDetailsRequestFn(input)
}

func (f *IAM) GetAccountAuthorizationDetails(input *iam.GetAccountAuthorizationDetailsInput) (*iam.GetAccountAuthorizationDetailsOutput, error) {
	if f.GetAccountAuthorizationDetailsFn == nil {
		return nil, nil
	}
	return f.GetAccountAuthorizationDetailsFn(input)
}

func (f *IAM) GetAccountAuthorizationDetailsPages(input *iam.GetAccountAuthorizationDetailsInput, fn func(*iam.GetAccountAuthorizationDetailsOutput, bool) bool) error {
	if f.GetAccountAuthorizationDetailsPagesFn == nil {
		return nil
	}
	return f.GetAccountAuthorizationDetailsPagesFn(input, fn)
}

func (f *IAM) GetAccountPasswordPolicyRequest(input *iam.GetAccountPasswordPolicyInput) (*request.Request, *iam.GetAccountPasswordPolicyOutput) {
	if f.GetAccountPasswordPolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetAccountPasswordPolicyRequestFn(input)
}

func (f *IAM) GetAccountPasswordPolicy(input *iam.GetAccountPasswordPolicyInput) (*iam.GetAccountPasswordPolicyOutput, error) {
	if f.GetAccountPasswordPolicyFn == nil {
		return nil, nil
	}
	return f.GetAccountPasswordPolicyFn(input)
}

func (f *IAM) GetAccountSummaryRequest(input *iam.GetAccountSummaryInput) (*request.Request, *iam.GetAccountSummaryOutput) {
	if f.GetAccountSummaryRequestFn == nil {
		return nil, nil
	}
	return f.GetAccountSummaryRequestFn(input)
}

func (f *IAM) GetAccountSummary(input *iam.GetAccountSummaryInput) (*iam.GetAccountSummaryOutput, error) {
	if f.GetAccountSummaryFn == nil {
		return nil, nil
	}
	return f.GetAccountSummaryFn(input)
}

func (f *IAM) GetContextKeysForCustomPolicyRequest(input *iam.GetContextKeysForCustomPolicyInput) (*request.Request, *iam.GetContextKeysForPolicyResponse) {
	if f.GetContextKeysForCustomPolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetContextKeysForCustomPolicyRequestFn(input)
}

func (f *IAM) GetContextKeysForCustomPolicy(input *iam.GetContextKeysForCustomPolicyInput) (*iam.GetContextKeysForPolicyResponse, error) {
	if f.GetContextKeysForCustomPolicyFn == nil {
		return nil, nil
	}
	return f.GetContextKeysForCustomPolicyFn(input)
}

func (f *IAM) GetContextKeysForPrincipalPolicyRequest(input *iam.GetContextKeysForPrincipalPolicyInput) (*request.Request, *iam.GetContextKeysForPolicyResponse) {
	if f.GetContextKeysForPrincipalPolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetContextKeysForPrincipalPolicyRequestFn(input)
}

func (f *IAM) GetContextKeysForPrincipalPolicy(input *iam.GetContextKeysForPrincipalPolicyInput) (*iam.GetContextKeysForPolicyResponse, error) {
	if f.GetContextKeysForPrincipalPolicyFn == nil {
		return nil, nil
	}
	return f.GetContextKeysForPrincipalPolicyFn(input)
}

func (f *IAM) GetCredentialReportRequest(input *iam.GetCredentialReportInput) (*request.Request, *iam.GetCredentialReportOutput) {
	if f.GetCredentialReportRequestFn == nil {
		return nil, nil
	}
	return f.GetCredentialReportRequestFn(input)
}

func (f *IAM) GetCredentialReport(input *iam.GetCredentialReportInput) (*iam.GetCredentialReportOutput, error) {
	if f.GetCredentialReportFn == nil {
		return nil, nil
	}
	return f.GetCredentialReportFn(input)
}

func (f *IAM) GetGroupRequest(input *iam.GetGroupInput) (*request.Request, *iam.GetGroupOutput) {
	if f.GetGroupRequestFn == nil {
		return nil, nil
	}
	return f.GetGroupRequestFn(input)
}

func (f *IAM) GetGroup(input *iam.GetGroupInput) (*iam.GetGroupOutput, error) {
	if f.GetGroupFn == nil {
		return nil, nil
	}
	return f.GetGroupFn(input)
}

func (f *IAM) GetGroupPages(input *iam.GetGroupInput, fn func(*iam.GetGroupOutput, bool) bool) error {
	if f.GetGroupPagesFn == nil {
		return nil
	}
	return f.GetGroupPagesFn(input, fn)
}

func (f *IAM) GetGroupPolicyRequest(input *iam.GetGroupPolicyInput) (*request.Request, *iam.GetGroupPolicyOutput) {
	if f.GetGroupPolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetGroupPolicyRequestFn(input)
}

func (f *IAM) GetGroupPolicy(input *iam.GetGroupPolicyInput) (*iam.GetGroupPolicyOutput, error) {
	if f.GetGroupPolicyFn == nil {
		return nil, nil
	}
	return f.GetGroupPolicyFn(input)
}

func (f *IAM) GetInstanceProfileRequest(input *iam.GetInstanceProfileInput) (*request.Request, *iam.GetInstanceProfileOutput) {
	if f.GetInstanceProfileRequestFn == nil {
		return nil, nil
	}
	return f.GetInstanceProfileRequestFn(input)
}

func (f *IAM) GetInstanceProfile(input *iam.GetInstanceProfileInput) (*iam.GetInstanceProfileOutput, error) {
	if f.GetInstanceProfileFn == nil {
		return nil, nil
	}
	return f.GetInstanceProfileFn(input)
}

func (f *IAM) GetLoginProfileRequest(input *iam.GetLoginProfileInput) (*request.Request, *iam.GetLoginProfileOutput) {
	if f.GetLoginProfileRequestFn == nil {
		return nil, nil
	}
	return f.GetLoginProfileRequestFn(input)
}

func (f *IAM) GetLoginProfile(input *iam.GetLoginProfileInput) (*iam.GetLoginProfileOutput, error) {
	if f.GetLoginProfileFn == nil {
		return nil, nil
	}
	return f.GetLoginProfileFn(input)
}

func (f *IAM) GetOpenIDConnectProviderRequest(input *iam.GetOpenIDConnectProviderInput) (*request.Request, *iam.GetOpenIDConnectProviderOutput) {
	if f.GetOpenIDConnectProviderRequestFn == nil {
		return nil, nil
	}
	return f.GetOpenIDConnectProviderRequestFn(input)
}

func (f *IAM) GetOpenIDConnectProvider(input *iam.GetOpenIDConnectProviderInput) (*iam.GetOpenIDConnectProviderOutput, error) {
	if f.GetOpenIDConnectProviderFn == nil {
		return nil, nil
	}
	return f.GetOpenIDConnectProviderFn(input)
}

func (f *IAM) GetPolicyRequest(input *iam.GetPolicyInput) (*request.Request, *iam.GetPolicyOutput) {
	if f.GetPolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetPolicyRequestFn(input)
}

func (f *IAM) GetPolicy(input *iam.GetPolicyInput) (*iam.GetPolicyOutput, error) {
	if f.GetPolicyFn == nil {
		return nil, nil
	}
	return f.GetPolicyFn(input)
}

func (f *IAM) GetPolicyVersionRequest(input *iam.GetPolicyVersionInput) (*request.Request, *iam.GetPolicyVersionOutput) {
	if f.GetPolicyVersionRequestFn == nil {
		return nil, nil
	}
	return f.GetPolicyVersionRequestFn(input)
}

func (f *IAM) GetPolicyVersion(input *iam.GetPolicyVersionInput) (*iam.GetPolicyVersionOutput, error) {
	if f.GetPolicyVersionFn == nil {
		return nil, nil
	}
	return f.GetPolicyVersionFn(input)
}

func (f *IAM) GetRoleRequest(input *iam.GetRoleInput) (*request.Request, *iam.GetRoleOutput) {
	if f.GetRoleRequestFn == nil {
		return nil, nil
	}
	return f.GetRoleRequestFn(input)
}

func (f *IAM) GetRole(input *iam.GetRoleInput) (*iam.GetRoleOutput, error) {
	if f.GetRoleFn == nil {
		return nil, nil
	}
	return f.GetRoleFn(input)
}

func (f *IAM) GetRolePolicyRequest(input *iam.GetRolePolicyInput) (*request.Request, *iam.GetRolePolicyOutput) {
	if f.GetRolePolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetRolePolicyRequestFn(input)
}

func (f *IAM) GetRolePolicy(input *iam.GetRolePolicyInput) (*iam.GetRolePolicyOutput, error) {
	if f.GetRolePolicyFn == nil {
		return nil, nil
	}
	return f.GetRolePolicyFn(input)
}

func (f *IAM) GetSAMLProviderRequest(input *iam.GetSAMLProviderInput) (*request.Request, *iam.GetSAMLProviderOutput) {
	if f.GetSAMLProviderRequestFn == nil {
		return nil, nil
	}
	return f.GetSAMLProviderRequestFn(input)
}

func (f *IAM) GetSAMLProvider(input *iam.GetSAMLProviderInput) (*iam.GetSAMLProviderOutput, error) {
	if f.GetSAMLProviderFn == nil {
		return nil, nil
	}
	return f.GetSAMLProviderFn(input)
}

func (f *IAM) GetSSHPublicKeyRequest(input *iam.GetSSHPublicKeyInput) (*request.Request, *iam.GetSSHPublicKeyOutput) {
	if f.GetSSHPublicKeyRequestFn == nil {
		return nil, nil
	}
	return f.GetSSHPublicKeyRequestFn(input)
}

func (f *IAM) GetSSHPublicKey(input *iam.GetSSHPublicKeyInput) (*iam.GetSSHPublicKeyOutput, error) {
	if f.GetSSHPublicKeyFn == nil {
		return nil, nil
	}
	return f.GetSSHPublicKeyFn(input)
}

func (f *IAM) GetServerCertificateRequest(input *iam.GetServerCertificateInput) (*request.Request, *iam.GetServerCertificateOutput) {
	if f.GetServerCertificateRequestFn == nil {
		return nil, nil
	}
	return f.GetServerCertificateRequestFn(input)
}

func (f *IAM) GetServerCertificate(input *iam.GetServerCertificateInput) (*iam.GetServerCertificateOutput, error) {
	if f.GetServerCertificateFn == nil {
		return nil, nil
	}
	return f.GetServerCertificateFn(input)
}

func (f *IAM) GetUserRequest(input *iam.GetUserInput) (*request.Request, *iam.GetUserOutput) {
	if f.GetUserRequestFn == nil {
		return nil, nil
	}
	return f.GetUserRequestFn(input)
}

func (f *IAM) GetUser(input *iam.GetUserInput) (*iam.GetUserOutput, error) {
	if f.GetUserFn == nil {
		return nil, nil
	}
	return f.GetUserFn(input)
}

func (f *IAM) GetUserPolicyRequest(input *iam.GetUserPolicyInput) (*request.Request, *iam.GetUserPolicyOutput) {
	if f.GetUserPolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetUserPolicyRequestFn(input)
}

func (f *IAM) GetUserPolicy(input *iam.GetUserPolicyInput) (*iam.GetUserPolicyOutput, error) {
	if f.GetUserPolicyFn == nil {
		return nil, nil
	}
	return f.GetUserPolicyFn(input)
}

func (f *IAM) ListAccessKeysRequest(input *iam.ListAccessKeysInput) (*request.Request, *iam.ListAccessKeysOutput) {
	if f.ListAccessKeysRequestFn == nil {
		return nil, nil
	}
	return f.ListAccessKeysRequestFn(input)
}

func (f *IAM) ListAccessKeys(input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	if f.ListAccessKeysFn == nil {
		return nil, nil
	}
	return f.ListAccessKeysFn(input)
}

func (f *IAM) ListAccessKeysPages(input *iam.ListAccessKeysInput, fn func(*iam.ListAccessKeysOutput, bool) bool) error {
	if f.ListAccessKeysPagesFn == nil {
		return nil
	}
	return f.ListAccessKeysPagesFn(input, fn)
}

func (f *IAM) ListAccountAliasesRequest(input *iam.ListAccountAliasesInput) (*request.Request, *iam.ListAccountAliasesOutput) {
	if f.ListAccountAliasesRequestFn == nil {
		return nil, nil
	}
	return f.ListAccountAliasesRequestFn(input)
}

func (f *IAM) ListAccountAliases(input *iam.ListAccountAliasesInput) (*iam.ListAccountAliasesOutput, error) {
	if f.ListAccountAliasesFn == nil {
		return nil, nil
	}
	return f.ListAccountAliasesFn(input)
}

func (f *IAM) ListAccountAliasesPages(input *iam.ListAccountAliasesInput, fn func(*iam.ListAccountAliasesOutput, bool) bool) error {
	if f.ListAccountAliasesPagesFn == nil {
		return nil
	}
	return f.ListAccountAliasesPagesFn(input, fn)
}

func (f *IAM) ListAttachedGroupPoliciesRequest(input *iam.ListAttachedGroupPoliciesInput) (*request.Request, *iam.ListAttachedGroupPoliciesOutput) {
	if f.ListAttachedGroupPoliciesRequestFn == nil {
		return nil, nil
	}
	return f.ListAttachedGroupPoliciesRequestFn(input)
}

func (f *IAM) ListAttachedGroupPolicies(input *iam.ListAttachedGroupPoliciesInput) (*iam.ListAttachedGroupPoliciesOutput, error) {
	if f.ListAttachedGroupPoliciesFn == nil {
		return nil, nil
	}
	return f.ListAttachedGroupPoliciesFn(input)
}

func (f *IAM) ListAttachedGroupPoliciesPages(input *iam.ListAttachedGroupPoliciesInput, fn func(*iam.ListAttachedGroupPoliciesOutput, bool) bool) error {
	if f.ListAttachedGroupPoliciesPagesFn == nil {
		return nil
	}
	return f.ListAttachedGroupPoliciesPagesFn(input, fn)
}

func (f *IAM) ListAttachedRolePoliciesRequest(input *iam.ListAttachedRolePoliciesInput) (*request.Request, *iam.ListAttachedRolePoliciesOutput) {
	if f.ListAttachedRolePoliciesRequestFn == nil {
		return nil, nil
	}
	return f.ListAttachedRolePoliciesRequestFn(input)
}

func (f *IAM) ListAttachedRolePolicies(input *iam.ListAttachedRolePoliciesInput) (*iam.ListAttachedRolePoliciesOutput, error) {
	if f.ListAttachedRolePoliciesFn == nil {
		return nil, nil
	}
	return f.ListAttachedRolePoliciesFn(input)
}

func (f *IAM) ListAttachedRolePoliciesPages(input *iam.ListAttachedRolePoliciesInput, fn func(*iam.ListAttachedRolePoliciesOutput, bool) bool) error {
	if f.ListAttachedRolePoliciesPagesFn == nil {
		return nil
	}
	return f.ListAttachedRolePoliciesPagesFn(input, fn)
}

func (f *IAM) ListAttachedUserPoliciesRequest(input *iam.ListAttachedUserPoliciesInput) (*request.Request, *iam.ListAttachedUserPoliciesOutput) {
	if f.ListAttachedUserPoliciesRequestFn == nil {
		return nil, nil
	}
	return f.ListAttachedUserPoliciesRequestFn(input)
}

func (f *IAM) ListAttachedUserPolicies(input *iam.ListAttachedUserPoliciesInput) (*iam.ListAttachedUserPoliciesOutput, error) {
	if f.ListAttachedUserPoliciesFn == nil {
		return nil, nil
	}
	return f.ListAttachedUserPoliciesFn(input)
}

func (f *IAM) ListAttachedUserPoliciesPages(input *iam.ListAttachedUserPoliciesInput, fn func(*iam.ListAttachedUserPoliciesOutput, bool) bool) error {
	if f.ListAttachedUserPoliciesPagesFn == nil {
		return nil
	}
	return f.ListAttachedUserPoliciesPagesFn(input, fn)
}

func (f *IAM) ListEntitiesForPolicyRequest(input *iam.ListEntitiesForPolicyInput) (*request.Request, *iam.ListEntitiesForPolicyOutput) {
	if f.ListEntitiesForPolicyRequestFn == nil {
		return nil, nil
	}
	return f.ListEntitiesForPolicyRequestFn(input)
}

func (f *IAM) ListEntitiesForPolicy(input *iam.ListEntitiesForPolicyInput) (*iam.ListEntitiesForPolicyOutput, error) {
	if f.ListEntitiesForPolicyFn == nil {
		return nil, nil
	}
	return f.ListEntitiesForPolicyFn(input)
}

func (f *IAM) ListEntitiesForPolicyPages(input *iam.ListEntitiesForPolicyInput, fn func(*iam.ListEntitiesForPolicyOutput, bool) bool) error {
	if f.ListEntitiesForPolicyPagesFn == nil {
		return nil
	}
	return f.ListEntitiesForPolicyPagesFn(input, fn)
}

func (f *IAM) ListGroupPoliciesRequest(input *iam.ListGroupPoliciesInput) (*request.Request, *iam.ListGroupPoliciesOutput) {
	if f.ListGroupPoliciesRequestFn == nil {
		return nil, nil
	}
	return f.ListGroupPoliciesRequestFn(input)
}

func (f *IAM) ListGroupPolicies(input *iam.ListGroupPoliciesInput) (*iam.ListGroupPoliciesOutput, error) {
	if f.ListGroupPoliciesFn == nil {
		return nil, nil
	}
	return f.ListGroupPoliciesFn(input)
}

func (f *IAM) ListGroupPoliciesPages(input *iam.ListGroupPoliciesInput, fn func(*iam.ListGroupPoliciesOutput, bool) bool) error {
	if f.ListGroupPoliciesPagesFn == nil {
		return nil
	}
	return f.ListGroupPoliciesPagesFn(input, fn)
}

func (f *IAM) ListGroupsRequest(input *iam.ListGroupsInput) (*request.Request, *iam.ListGroupsOutput) {
	if f.ListGroupsRequestFn == nil {
		return nil, nil
	}
	return f.ListGroupsRequestFn(input)
}

func (f *IAM) ListGroups(input *iam.ListGroupsInput) (*iam.ListGroupsOutput, error) {
	if f.ListGroupsFn == nil {
		return nil, nil
	}
	return f.ListGroupsFn(input)
}

func (f *IAM) ListGroupsPages(input *iam.ListGroupsInput, fn func(*iam.ListGroupsOutput, bool) bool) error {
	if f.ListGroupsPagesFn == nil {
		return nil
	}
	return f.ListGroupsPagesFn(input, fn)
}

func (f *IAM) ListGroupsForUserRequest(input *iam.ListGroupsForUserInput) (*request.Request, *iam.ListGroupsForUserOutput) {
	if f.ListGroupsForUserRequestFn == nil {
		return nil, nil
	}
	return f.ListGroupsForUserRequestFn(input)
}

func (f *IAM) ListGroupsForUser(input *iam.ListGroupsForUserInput) (*iam.ListGroupsForUserOutput, error) {
	if f.ListGroupsForUserFn == nil {
		return nil, nil
	}
	return f.ListGroupsForUserFn(input)
}

func (f *IAM) ListGroupsForUserPages(input *iam.ListGroupsForUserInput, fn func(*iam.ListGroupsForUserOutput, bool) bool) error {
	if f.ListGroupsForUserPagesFn == nil {
		return nil
	}
	return f.ListGroupsForUserPagesFn(input, fn)
}

func (f *IAM) ListInstanceProfilesRequest(input *iam.ListInstanceProfilesInput) (*request.Request, *iam.ListInstanceProfilesOutput) {
	if f.ListInstanceProfilesRequestFn == nil {
		return nil, nil
	}
	return f.ListInstanceProfilesRequestFn(input)
}

func (f *IAM) ListInstanceProfiles(input *iam.ListInstanceProfilesInput) (*iam.ListInstanceProfilesOutput, error) {
	if f.ListInstanceProfilesFn == nil {
		return nil, nil
	}
	return f.ListInstanceProfilesFn(input)
}

func (f *IAM) ListInstanceProfilesPages(input *iam.ListInstanceProfilesInput, fn func(*iam.ListInstanceProfilesOutput, bool) bool) error {
	if f.ListInstanceProfilesPagesFn == nil {
		return nil
	}
	return f.ListInstanceProfilesPagesFn(input, fn)
}

func (f *IAM) ListInstanceProfilesForRoleRequest(input *iam.ListInstanceProfilesForRoleInput) (*request.Request, *iam.ListInstanceProfilesForRoleOutput) {
	if f.ListInstanceProfilesForRoleRequestFn == nil {
		return nil, nil
	}
	return f.ListInstanceProfilesForRoleRequestFn(input)
}

func (f *IAM) ListInstanceProfilesForRole(input *iam.ListInstanceProfilesForRoleInput) (*iam.ListInstanceProfilesForRoleOutput, error) {
	if f.ListInstanceProfilesForRoleFn == nil {
		return nil, nil
	}
	return f.ListInstanceProfilesForRoleFn(input)
}

func (f *IAM) ListInstanceProfilesForRolePages(input *iam.ListInstanceProfilesForRoleInput, fn func(*iam.ListInstanceProfilesForRoleOutput, bool) bool) error {
	if f.ListInstanceProfilesForRolePagesFn == nil {
		return nil
	}
	return f.ListInstanceProfilesForRolePagesFn(input, fn)
}

func (f *IAM) ListMFADevicesRequest(input *iam.ListMFADevicesInput) (*request.Request, *iam.ListMFADevicesOutput) {
	if f.ListMFADevicesRequestFn == nil {
		return nil, nil
	}
	return f.ListMFADevicesRequestFn(input)
}

func (f *IAM) ListMFADevices(input *iam.ListMFADevicesInput) (*iam.ListMFADevicesOutput, error) {
	if f.ListMFADevicesFn == nil {
		return nil, nil
	}
	return f.ListMFADevicesFn(input)
}

func (f *IAM) ListMFADevicesPages(input *iam.ListMFADevicesInput, fn func(*iam.ListMFADevicesOutput, bool) bool) error {
	if f.ListMFADevicesPagesFn == nil {
		return nil
	}
	return f.ListMFADevicesPagesFn(input, fn)
}

func (f *IAM) ListOpenIDConnectProvidersRequest(input *iam.ListOpenIDConnectProvidersInput) (*request.Request, *iam.ListOpenIDConnectProvidersOutput) {
	if f.ListOpenIDConnectProvidersRequestFn == nil {
		return nil, nil
	}
	return f.ListOpenIDConnectProvidersRequestFn(input)
}

func (f *IAM) ListOpenIDConnectProviders(input *iam.ListOpenIDConnectProvidersInput) (*iam.ListOpenIDConnectProvidersOutput, error) {
	if f.ListOpenIDConnectProvidersFn == nil {
		return nil, nil
	}
	return f.ListOpenIDConnectProvidersFn(input)
}

func (f *IAM) ListPoliciesRequest(input *iam.ListPoliciesInput) (*request.Request, *iam.ListPoliciesOutput) {
	if f.ListPoliciesRequestFn == nil {
		return nil, nil
	}
	return f.ListPoliciesRequestFn(input)
}

func (f *IAM) ListPolicies(input *iam.ListPoliciesInput) (*iam.ListPoliciesOutput, error) {
	if f.ListPoliciesFn == nil {
		return nil, nil
	}
	return f.ListPoliciesFn(input)
}

func (f *IAM) ListPoliciesPages(input *iam.ListPoliciesInput, fn func(*iam.ListPoliciesOutput, bool) bool) error {
	if f.ListPoliciesPagesFn == nil {
		return nil
	}
	return f.ListPoliciesPagesFn(input, fn)
}

func (f *IAM) ListPolicyVersionsRequest(input *iam.ListPolicyVersionsInput) (*request.Request, *iam.ListPolicyVersionsOutput) {
	if f.ListPolicyVersionsRequestFn == nil {
		return nil, nil
	}
	return f.ListPolicyVersionsRequestFn(input)
}

func (f *IAM) ListPolicyVersions(input *iam.ListPolicyVersionsInput) (*iam.ListPolicyVersionsOutput, error) {
	if f.ListPolicyVersionsFn == nil {
		return nil, nil
	}
	return f.ListPolicyVersionsFn(input)
}

func (f *IAM) ListRolePoliciesRequest(input *iam.ListRolePoliciesInput) (*request.Request, *iam.ListRolePoliciesOutput) {
	if f.ListRolePoliciesRequestFn == nil {
		return nil, nil
	}
	return f.ListRolePoliciesRequestFn(input)
}

func (f *IAM) ListRolePolicies(input *iam.ListRolePoliciesInput) (*iam.ListRolePoliciesOutput, error) {
	if f.ListRolePoliciesFn == nil {
		return nil, nil
	}
	return f.ListRolePoliciesFn(input)
}

func (f *IAM) ListRolePoliciesPages(input *iam.ListRolePoliciesInput, fn func(*iam.ListRolePoliciesOutput, bool) bool) error {
	if f.ListRolePoliciesPagesFn == nil {
		return nil
	}
	return f.ListRolePoliciesPagesFn(input, fn)
}

func (f *IAM) ListRolesRequest(input *iam.ListRolesInput) (*request.Request, *iam.ListRolesOutput) {
	if f.ListRolesRequestFn == nil {
		return nil, nil
	}
	return f.ListRolesRequestFn(input)
}

func (f *IAM) ListRoles(input *iam.ListRolesInput) (*iam.ListRolesOutput, error) {
	if f.ListRolesFn == nil {
		return nil, nil
	}
	return f.ListRolesFn(input)
}

func (f *IAM) ListRolesPages(input *iam.ListRolesInput, fn func(*iam.ListRolesOutput, bool) bool) error {
	if f.ListRolesPagesFn == nil {
		return nil
	}
	return f.ListRolesPagesFn(input, fn)
}

func (f *IAM) ListSAMLProvidersRequest(input *iam.ListSAMLProvidersInput) (*request.Request, *iam.ListSAMLProvidersOutput) {
	if f.ListSAMLProvidersRequestFn == nil {
		return nil, nil
	}
	return f.ListSAMLProvidersRequestFn(input)
}

func (f *IAM) ListSAMLProviders(input *iam.ListSAMLProvidersInput) (*iam.ListSAMLProvidersOutput, error) {
	if f.ListSAMLProvidersFn == nil {
		return nil, nil
	}
	return f.ListSAMLProvidersFn(input)
}

func (f *IAM) ListSSHPublicKeysRequest(input *iam.ListSSHPublicKeysInput) (*request.Request, *iam.ListSSHPublicKeysOutput) {
	if f.ListSSHPublicKeysRequestFn == nil {
		return nil, nil
	}
	return f.ListSSHPublicKeysRequestFn(input)
}

func (f *IAM) ListSSHPublicKeys(input *iam.ListSSHPublicKeysInput) (*iam.ListSSHPublicKeysOutput, error) {
	if f.ListSSHPublicKeysFn == nil {
		return nil, nil
	}
	return f.ListSSHPublicKeysFn(input)
}

func (f *IAM) ListServerCertificatesRequest(input *iam.ListServerCertificatesInput) (*request.Request, *iam.ListServerCertificatesOutput) {
	if f.ListServerCertificatesRequestFn == nil {
		return nil, nil
	}
	return f.ListServerCertificatesRequestFn(input)
}

func (f *IAM) ListServerCertificates(input *iam.ListServerCertificatesInput) (*iam.ListServerCertificatesOutput, error) {
	if f.ListServerCertificatesFn == nil {
		return nil, nil
	}
	return f.ListServerCertificatesFn(input)
}

func (f *IAM) ListServerCertificatesPages(input *iam.ListServerCertificatesInput, fn func(*iam.ListServerCertificatesOutput, bool) bool) error {
	if f.ListServerCertificatesPagesFn == nil {
		return nil
	}
	return f.ListServerCertificatesPagesFn(input, fn)
}

func (f *IAM) ListSigningCertificatesRequest(input *iam.ListSigningCertificatesInput) (*request.Request, *iam.ListSigningCertificatesOutput) {
	if f.ListSigningCertificatesRequestFn == nil {
		return nil, nil
	}
	return f.ListSigningCertificatesRequestFn(input)
}

func (f *IAM) ListSigningCertificates(input *iam.ListSigningCertificatesInput) (*iam.ListSigningCertificatesOutput, error) {
	if f.ListSigningCertificatesFn == nil {
		return nil, nil
	}
	return f.ListSigningCertificatesFn(input)
}

func (f *IAM) ListSigningCertificatesPages(input *iam.ListSigningCertificatesInput, fn func(*iam.ListSigningCertificatesOutput, bool) bool) error {
	if f.ListSigningCertificatesPagesFn == nil {
		return nil
	}
	return f.ListSigningCertificatesPagesFn(input, fn)
}

func (f *IAM) ListUserPoliciesRequest(input *iam.ListUserPoliciesInput) (*request.Request, *iam.ListUserPoliciesOutput) {
	if f.ListUserPoliciesRequestFn == nil {
		return nil, nil
	}
	return f.ListUserPoliciesRequestFn(input)
}

func (f *IAM) ListUserPolicies(input *iam.ListUserPoliciesInput) (*iam.ListUserPoliciesOutput, error) {
	if f.ListUserPoliciesFn == nil {
		return nil, nil
	}
	return f.ListUserPoliciesFn(input)
}

func (f *IAM) ListUserPoliciesPages(input *iam.ListUserPoliciesInput, fn func(*iam.ListUserPoliciesOutput, bool) bool) error {
	if f.ListUserPoliciesPagesFn == nil {
		return nil
	}
	return f.ListUserPoliciesPagesFn(input, fn)
}

func (f *IAM) ListUsersRequest(input *iam.ListUsersInput) (*request.Request, *iam.ListUsersOutput) {
	if f.ListUsersRequestFn == nil {
		return nil, nil
	}
	return f.ListUsersRequestFn(input)
}

func (f *IAM) ListUsers(input *iam.ListUsersInput) (*iam.ListUsersOutput, error) {
	if f.ListUsersFn == nil {
		return nil, nil
	}
	return f.ListUsersFn(input)
}

func (f *IAM) ListUsersPages(input *iam.ListUsersInput, fn func(*iam.ListUsersOutput, bool) bool) error {
	if f.ListUsersPagesFn == nil {
		return nil
	}
	return f.ListUsersPagesFn(input, fn)
}

func (f *IAM) ListVirtualMFADevicesRequest(input *iam.ListVirtualMFADevicesInput) (*request.Request, *iam.ListVirtualMFADevicesOutput) {
	if f.ListVirtualMFADevicesRequestFn == nil {
		return nil, nil
	}
	return f.ListVirtualMFADevicesRequestFn(input)
}

func (f *IAM) ListVirtualMFADevices(input *iam.ListVirtualMFADevicesInput) (*iam.ListVirtualMFADevicesOutput, error) {
	if f.ListVirtualMFADevicesFn == nil {
		return nil, nil
	}
	return f.ListVirtualMFADevicesFn(input)
}

func (f *IAM) ListVirtualMFADevicesPages(input *iam.ListVirtualMFADevicesInput, fn func(*iam.ListVirtualMFADevicesOutput, bool) bool) error {
	if f.ListVirtualMFADevicesPagesFn == nil {
		return nil
	}
	return f.ListVirtualMFADevicesPagesFn(input, fn)
}

func (f *IAM) PutGroupPolicyRequest(input *iam.PutGroupPolicyInput) (*request.Request, *iam.PutGroupPolicyOutput) {
	if f.PutGroupPolicyRequestFn == nil {
		return nil, nil
	}
	return f.PutGroupPolicyRequestFn(input)
}

func (f *IAM) PutGroupPolicy(input *iam.PutGroupPolicyInput) (*iam.PutGroupPolicyOutput, error) {
	if f.PutGroupPolicyFn == nil {
		return nil, nil
	}
	return f.PutGroupPolicyFn(input)
}

func (f *IAM) PutRolePolicyRequest(input *iam.PutRolePolicyInput) (*request.Request, *iam.PutRolePolicyOutput) {
	if f.PutRolePolicyRequestFn == nil {
		return nil, nil
	}
	return f.PutRolePolicyRequestFn(input)
}

func (f *IAM) PutRolePolicy(input *iam.PutRolePolicyInput) (*iam.PutRolePolicyOutput, error) {
	if f.PutRolePolicyFn == nil {
		return nil, nil
	}
	return f.PutRolePolicyFn(input)
}

func (f *IAM) PutUserPolicyRequest(input *iam.PutUserPolicyInput) (*request.Request, *iam.PutUserPolicyOutput) {
	if f.PutUserPolicyRequestFn == nil {
		return nil, nil
	}
	return f.PutUserPolicyRequestFn(input)
}

func (f *IAM) PutUserPolicy(input *iam.PutUserPolicyInput) (*iam.PutUserPolicyOutput, error) {
	if f.PutUserPolicyFn == nil {
		return nil, nil
	}
	return f.PutUserPolicyFn(input)
}

func (f *IAM) RemoveClientIDFromOpenIDConnectProviderRequest(input *iam.RemoveClientIDFromOpenIDConnectProviderInput) (*request.Request, *iam.RemoveClientIDFromOpenIDConnectProviderOutput) {
	if f.RemoveClientIDFromOpenIDConnectProviderRequestFn == nil {
		return nil, nil
	}
	return f.RemoveClientIDFromOpenIDConnectProviderRequestFn(input)
}

func (f *IAM) RemoveClientIDFromOpenIDConnectProvider(input *iam.RemoveClientIDFromOpenIDConnectProviderInput) (*iam.RemoveClientIDFromOpenIDConnectProviderOutput, error) {
	if f.RemoveClientIDFromOpenIDConnectProviderFn == nil {
		return nil, nil
	}
	return f.RemoveClientIDFromOpenIDConnectProviderFn(input)
}

func (f *IAM) RemoveRoleFromInstanceProfileRequest(input *iam.RemoveRoleFromInstanceProfileInput) (*request.Request, *iam.RemoveRoleFromInstanceProfileOutput) {
	if f.RemoveRoleFromInstanceProfileRequestFn == nil {
		return nil, nil
	}
	return f.RemoveRoleFromInstanceProfileRequestFn(input)
}

func (f *IAM) RemoveRoleFromInstanceProfile(input *iam.RemoveRoleFromInstanceProfileInput) (*iam.RemoveRoleFromInstanceProfileOutput, error) {
	if f.RemoveRoleFromInstanceProfileFn == nil {
		return nil, nil
	}
	return f.RemoveRoleFromInstanceProfileFn(input)
}

func (f *IAM) RemoveUserFromGroupRequest(input *iam.RemoveUserFromGroupInput) (*request.Request, *iam.RemoveUserFromGroupOutput) {
	if f.RemoveUserFromGroupRequestFn == nil {
		return nil, nil
	}
	return f.RemoveUserFromGroupRequestFn(input)
}

func (f *IAM) RemoveUserFromGroup(input *iam.RemoveUserFromGroupInput) (*iam.RemoveUserFromGroupOutput, error) {
	if f.RemoveUserFromGroupFn == nil {
		return nil, nil
	}
	return f.RemoveUserFromGroupFn(input)
}

func (f *IAM) ResyncMFADeviceRequest(input *iam.ResyncMFADeviceInput) (*request.Request, *iam.ResyncMFADeviceOutput) {
	if f.ResyncMFADeviceRequestFn == nil {
		return nil, nil
	}
	return f.ResyncMFADeviceRequestFn(input)
}

func (f *IAM) ResyncMFADevice(input *iam.ResyncMFADeviceInput) (*iam.ResyncMFADeviceOutput, error) {
	if f.ResyncMFADeviceFn == nil {
		return nil, nil
	}
	return f.ResyncMFADeviceFn(input)
}

func (f *IAM) SetDefaultPolicyVersionRequest(input *iam.SetDefaultPolicyVersionInput) (*request.Request, *iam.SetDefaultPolicyVersionOutput) {
	if f.SetDefaultPolicyVersionRequestFn == nil {
		return nil, nil
	}
	return f.SetDefaultPolicyVersionRequestFn(input)
}

func (f *IAM) SetDefaultPolicyVersion(input *iam.SetDefaultPolicyVersionInput) (*iam.SetDefaultPolicyVersionOutput, error) {
	if f.SetDefaultPolicyVersionFn == nil {
		return nil, nil
	}
	return f.SetDefaultPolicyVersionFn(input)
}

func (f *IAM) SimulateCustomPolicyRequest(input *iam.SimulateCustomPolicyInput) (*request.Request, *iam.SimulatePolicyResponse) {
	if f.SimulateCustomPolicyRequestFn == nil {
		return nil, nil
	}
	return f.SimulateCustomPolicyRequestFn(input)
}

func (f *IAM) SimulateCustomPolicy(input *iam.SimulateCustomPolicyInput) (*iam.SimulatePolicyResponse, error) {
	if f.SimulateCustomPolicyFn == nil {
		return nil, nil
	}
	return f.SimulateCustomPolicyFn(input)
}

func (f *IAM) SimulatePrincipalPolicyRequest(input *iam.SimulatePrincipalPolicyInput) (*request.Request, *iam.SimulatePolicyResponse) {
	if f.SimulatePrincipalPolicyRequestFn == nil {
		return nil, nil
	}
	return f.SimulatePrincipalPolicyRequestFn(input)
}

func (f *IAM) SimulatePrincipalPolicy(input *iam.SimulatePrincipalPolicyInput) (*iam.SimulatePolicyResponse, error) {
	if f.SimulatePrincipalPolicyFn == nil {
		return nil, nil
	}
	return f.SimulatePrincipalPolicyFn(input)
}

func (f *IAM) UpdateAccessKeyRequest(input *iam.UpdateAccessKeyInput) (*request.Request, *iam.UpdateAccessKeyOutput) {
	if f.UpdateAccessKeyRequestFn == nil {
		return nil, nil
	}
	return f.UpdateAccessKeyRequestFn(input)
}

func (f *IAM) UpdateAccessKey(input *iam.UpdateAccessKeyInput) (*iam.UpdateAccessKeyOutput, error) {
	if f.UpdateAccessKeyFn == nil {
		return nil, nil
	}
	return f.UpdateAccessKeyFn(input)
}

func (f *IAM) UpdateAccountPasswordPolicyRequest(input *iam.UpdateAccountPasswordPolicyInput) (*request.Request, *iam.UpdateAccountPasswordPolicyOutput) {
	if f.UpdateAccountPasswordPolicyRequestFn == nil {
		return nil, nil
	}
	return f.UpdateAccountPasswordPolicyRequestFn(input)
}

func (f *IAM) UpdateAccountPasswordPolicy(input *iam.UpdateAccountPasswordPolicyInput) (*iam.UpdateAccountPasswordPolicyOutput, error) {
	if f.UpdateAccountPasswordPolicyFn == nil {
		return nil, nil
	}
	return f.UpdateAccountPasswordPolicyFn(input)
}

func (f *IAM) UpdateAssumeRolePolicyRequest(input *iam.UpdateAssumeRolePolicyInput) (*request.Request, *iam.UpdateAssumeRolePolicyOutput) {
	if f.UpdateAssumeRolePolicyRequestFn == nil {
		return nil, nil
	}
	return f.UpdateAssumeRolePolicyRequestFn(input)
}

func (f *IAM) UpdateAssumeRolePolicy(input *iam.UpdateAssumeRolePolicyInput) (*iam.UpdateAssumeRolePolicyOutput, error) {
	if f.UpdateAssumeRolePolicyFn == nil {
		return nil, nil
	}
	return f.UpdateAssumeRolePolicyFn(input)
}

func (f *IAM) UpdateGroupRequest(input *iam.UpdateGroupInput) (*request.Request, *iam.UpdateGroupOutput) {
	if f.UpdateGroupRequestFn == nil {
		return nil, nil
	}
	return f.UpdateGroupRequestFn(input)
}

func (f *IAM) UpdateGroup(input *iam.UpdateGroupInput) (*iam.UpdateGroupOutput, error) {
	if f.UpdateGroupFn == nil {
		return nil, nil
	}
	return f.UpdateGroupFn(input)
}

func (f *IAM) UpdateLoginProfileRequest(input *iam.UpdateLoginProfileInput) (*request.Request, *iam.UpdateLoginProfileOutput) {
	if f.UpdateLoginProfileRequestFn == nil {
		return nil, nil
	}
	return f.UpdateLoginProfileRequestFn(input)
}

func (f *IAM) UpdateLoginProfile(input *iam.UpdateLoginProfileInput) (*iam.UpdateLoginProfileOutput, error) {
	if f.UpdateLoginProfileFn == nil {
		return nil, nil
	}
	return f.UpdateLoginProfileFn(input)
}

func (f *IAM) UpdateOpenIDConnectProviderThumbprintRequest(input *iam.UpdateOpenIDConnectProviderThumbprintInput) (*request.Request, *iam.UpdateOpenIDConnectProviderThumbprintOutput) {
	if f.UpdateOpenIDConnectProviderThumbprintRequestFn == nil {
		return nil, nil
	}
	return f.UpdateOpenIDConnectProviderThumbprintRequestFn(input)
}

func (f *IAM) UpdateOpenIDConnectProviderThumbprint(input *iam.UpdateOpenIDConnectProviderThumbprintInput) (*iam.UpdateOpenIDConnectProviderThumbprintOutput, error) {
	if f.UpdateOpenIDConnectProviderThumbprintFn == nil {
		return nil, nil
	}
	return f.UpdateOpenIDConnectProviderThumbprintFn(input)
}

func (f *IAM) UpdateSAMLProviderRequest(input *iam.UpdateSAMLProviderInput) (*request.Request, *iam.UpdateSAMLProviderOutput) {
	if f.UpdateSAMLProviderRequestFn == nil {
		return nil, nil
	}
	return f.UpdateSAMLProviderRequestFn(input)
}

func (f *IAM) UpdateSAMLProvider(input *iam.UpdateSAMLProviderInput) (*iam.UpdateSAMLProviderOutput, error) {
	if f.UpdateSAMLProviderFn == nil {
		return nil, nil
	}
	return f.UpdateSAMLProviderFn(input)
}

func (f *IAM) UpdateSSHPublicKeyRequest(input *iam.UpdateSSHPublicKeyInput) (*request.Request, *iam.UpdateSSHPublicKeyOutput) {
	if f.UpdateSSHPublicKeyRequestFn == nil {
		return nil, nil
	}
	return f.UpdateSSHPublicKeyRequestFn(input)
}

func (f *IAM) UpdateSSHPublicKey(input *iam.UpdateSSHPublicKeyInput) (*iam.UpdateSSHPublicKeyOutput, error) {
	if f.UpdateSSHPublicKeyFn == nil {
		return nil, nil
	}
	return f.UpdateSSHPublicKeyFn(input)
}

func (f *IAM) UpdateServerCertificateRequest(input *iam.UpdateServerCertificateInput) (*request.Request, *iam.UpdateServerCertificateOutput) {
	if f.UpdateServerCertificateRequestFn == nil {
		return nil, nil
	}
	return f.UpdateServerCertificateRequestFn(input)
}

func (f *IAM) UpdateServerCertificate(input *iam.UpdateServerCertificateInput) (*iam.UpdateServerCertificateOutput, error) {
	if f.UpdateServerCertificateFn == nil {
		return nil, nil
	}
	return f.UpdateServerCertificateFn(input)
}

func (f *IAM) UpdateSigningCertificateRequest(input *iam.UpdateSigningCertificateInput) (*request.Request, *iam.UpdateSigningCertificateOutput) {
	if f.UpdateSigningCertificateRequestFn == nil {
		return nil, nil
	}
	return f.UpdateSigningCertificateRequestFn(input)
}

func (f *IAM) UpdateSigningCertificate(input *iam.UpdateSigningCertificateInput) (*iam.UpdateSigningCertificateOutput, error) {
	if f.UpdateSigningCertificateFn == nil {
		return nil, nil
	}
	return f.UpdateSigningCertificateFn(input)
}

func (f *IAM) UpdateUserRequest(input *iam.UpdateUserInput) (*request.Request, *iam.UpdateUserOutput) {
	if f.UpdateUserRequestFn == nil {
		return nil, nil
	}
	return f.UpdateUserRequestFn(input)
}

func (f *IAM) UpdateUser(input *iam.UpdateUserInput) (*iam.UpdateUserOutput, error) {
	if f.UpdateUserFn == nil {
		return nil, nil
	}
	return f.UpdateUserFn(input)
}

func (f *IAM) UploadSSHPublicKeyRequest(input *iam.UploadSSHPublicKeyInput) (*request.Request, *iam.UploadSSHPublicKeyOutput) {
	if f.UploadSSHPublicKeyRequestFn == nil {
		return nil, nil
	}
	return f.UploadSSHPublicKeyRequestFn(input)
}

func (f *IAM) UploadSSHPublicKey(input *iam.UploadSSHPublicKeyInput) (*iam.UploadSSHPublicKeyOutput, error) {
	if f.UploadSSHPublicKeyFn == nil {
		return nil, nil
	}
	return f.UploadSSHPublicKeyFn(input)
}

func (f *IAM) UploadServerCertificateRequest(input *iam.UploadServerCertificateInput) (*request.Request, *iam.UploadServerCertificateOutput) {
	if f.UploadServerCertificateRequestFn == nil {
		return nil, nil
	}
	return f.UploadServerCertificateRequestFn(input)
}

func (f *IAM) UploadServerCertificate(input *iam.UploadServerCertificateInput) (*iam.UploadServerCertificateOutput, error) {
	if f.UploadServerCertificateFn == nil {
		return nil, nil
	}
	return f.UploadServerCertificateFn(input)
}

func (f *IAM) UploadSigningCertificateRequest(input *iam.UploadSigningCertificateInput) (*request.Request, *iam.UploadSigningCertificateOutput) {
	if f.UploadSigningCertificateRequestFn == nil {
		return nil, nil
	}
	return f.UploadSigningCertificateRequestFn(input)
}

func (f *IAM) UploadSigningCertificate(input *iam.UploadSigningCertificateInput) (*iam.UploadSigningCertificateOutput, error) {
	if f.UploadSigningCertificateFn == nil {
		return nil, nil
	}
	return f.UploadSigningCertificateFn(input)
}
