package fake_aws_provider

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3 struct {
	s3iface.S3API
	AbortMultipartUploadRequestFn    func(*s3.AbortMultipartUploadInput) (*request.Request, *s3.AbortMultipartUploadOutput)
	AbortMultipartUploadFn           func(*s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error)
	CompleteMultipartUploadRequestFn func(*s3.CompleteMultipartUploadInput) (*request.Request, *s3.CompleteMultipartUploadOutput)
	CompleteMultipartUploadFn        func(*s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)
	CopyObjectRequestFn              func(*s3.CopyObjectInput) (*request.Request, *s3.CopyObjectOutput)
	CopyObjectFn                     func(*s3.CopyObjectInput) (*s3.CopyObjectOutput, error)
	CreateBucketRequestFn            func(*s3.CreateBucketInput) (*request.Request, *s3.CreateBucketOutput)
	CreateBucketFn                   func(*s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	CreateMultipartUploadRequestFn   func(*s3.CreateMultipartUploadInput) (*request.Request, *s3.CreateMultipartUploadOutput)
	CreateMultipartUploadFn          func(*s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)
	DeleteBucketRequestFn            func(*s3.DeleteBucketInput) (*request.Request, *s3.DeleteBucketOutput)
	DeleteBucketFn                   func(*s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error)
	DeleteBucketCorsRequestFn        func(*s3.DeleteBucketCorsInput) (*request.Request, *s3.DeleteBucketCorsOutput)
	DeleteBucketCorsFn               func(*s3.DeleteBucketCorsInput) (*s3.DeleteBucketCorsOutput, error)
	DeleteBucketLifecycleRequestFn   func(*s3.DeleteBucketLifecycleInput) (*request.Request, *s3.DeleteBucketLifecycleOutput)
	DeleteBucketLifecycleFn          func(*s3.DeleteBucketLifecycleInput) (*s3.DeleteBucketLifecycleOutput, error)
	DeleteBucketPolicyRequestFn      func(*s3.DeleteBucketPolicyInput) (*request.Request, *s3.DeleteBucketPolicyOutput)
	DeleteBucketPolicyFn             func(*s3.DeleteBucketPolicyInput) (*s3.DeleteBucketPolicyOutput, error)
	DeleteBucketReplicationRequestFn func(*s3.DeleteBucketReplicationInput) (*request.Request, *s3.DeleteBucketReplicationOutput)
	DeleteBucketReplicationFn        func(*s3.DeleteBucketReplicationInput) (*s3.DeleteBucketReplicationOutput, error)
	DeleteBucketTaggingRequestFn     func(*s3.DeleteBucketTaggingInput) (*request.Request, *s3.DeleteBucketTaggingOutput)
	DeleteBucketTaggingFn            func(*s3.DeleteBucketTaggingInput) (*s3.DeleteBucketTaggingOutput, error)
	DeleteBucketWebsiteRequestFn     func(*s3.DeleteBucketWebsiteInput) (*request.Request, *s3.DeleteBucketWebsiteOutput)
	DeleteBucketWebsiteFn            func(*s3.DeleteBucketWebsiteInput) (*s3.DeleteBucketWebsiteOutput, error)
	DeleteObjectRequestFn            func(*s3.DeleteObjectInput) (*request.Request, *s3.DeleteObjectOutput)
	DeleteObjectFn                   func(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	DeleteObjectsRequestFn           func(*s3.DeleteObjectsInput) (*request.Request, *s3.DeleteObjectsOutput)
	DeleteObjectsFn                  func(*s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)
	//	GetBucketAccelerateConfigurationRequestFn   func(*s3.GetBucketAccelerateConfigurationInput) (*request.Request, *s3.GetBucketAccelerateConfigurationOutput)
	//	GetBucketAccelerateConfigurationFn          func(*s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error)
	GetBucketAclRequestFn                       func(*s3.GetBucketAclInput) (*request.Request, *s3.GetBucketAclOutput)
	GetBucketAclFn                              func(*s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error)
	GetBucketCorsRequestFn                      func(*s3.GetBucketCorsInput) (*request.Request, *s3.GetBucketCorsOutput)
	GetBucketCorsFn                             func(*s3.GetBucketCorsInput) (*s3.GetBucketCorsOutput, error)
	GetBucketLifecycleRequestFn                 func(*s3.GetBucketLifecycleInput) (*request.Request, *s3.GetBucketLifecycleOutput)
	GetBucketLifecycleFn                        func(*s3.GetBucketLifecycleInput) (*s3.GetBucketLifecycleOutput, error)
	GetBucketLifecycleConfigurationRequestFn    func(*s3.GetBucketLifecycleConfigurationInput) (*request.Request, *s3.GetBucketLifecycleConfigurationOutput)
	GetBucketLifecycleConfigurationFn           func(*s3.GetBucketLifecycleConfigurationInput) (*s3.GetBucketLifecycleConfigurationOutput, error)
	GetBucketLocationRequestFn                  func(*s3.GetBucketLocationInput) (*request.Request, *s3.GetBucketLocationOutput)
	GetBucketLocationFn                         func(*s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error)
	GetBucketLoggingRequestFn                   func(*s3.GetBucketLoggingInput) (*request.Request, *s3.GetBucketLoggingOutput)
	GetBucketLoggingFn                          func(*s3.GetBucketLoggingInput) (*s3.GetBucketLoggingOutput, error)
	GetBucketNotificationRequestFn              func(*s3.GetBucketNotificationConfigurationRequest) (*request.Request, *s3.NotificationConfigurationDeprecated)
	GetBucketNotificationFn                     func(*s3.GetBucketNotificationConfigurationRequest) (*s3.NotificationConfigurationDeprecated, error)
	GetBucketNotificationConfigurationRequestFn func(*s3.GetBucketNotificationConfigurationRequest) (*request.Request, *s3.NotificationConfiguration)
	GetBucketNotificationConfigurationFn        func(*s3.GetBucketNotificationConfigurationRequest) (*s3.NotificationConfiguration, error)
	GetBucketPolicyRequestFn                    func(*s3.GetBucketPolicyInput) (*request.Request, *s3.GetBucketPolicyOutput)
	GetBucketPolicyFn                           func(*s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error)
	GetBucketReplicationRequestFn               func(*s3.GetBucketReplicationInput) (*request.Request, *s3.GetBucketReplicationOutput)
	GetBucketReplicationFn                      func(*s3.GetBucketReplicationInput) (*s3.GetBucketReplicationOutput, error)
	GetBucketRequestPaymentRequestFn            func(*s3.GetBucketRequestPaymentInput) (*request.Request, *s3.GetBucketRequestPaymentOutput)
	GetBucketRequestPaymentFn                   func(*s3.GetBucketRequestPaymentInput) (*s3.GetBucketRequestPaymentOutput, error)
	GetBucketTaggingRequestFn                   func(*s3.GetBucketTaggingInput) (*request.Request, *s3.GetBucketTaggingOutput)
	GetBucketTaggingFn                          func(*s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error)
	GetBucketVersioningRequestFn                func(*s3.GetBucketVersioningInput) (*request.Request, *s3.GetBucketVersioningOutput)
	GetBucketVersioningFn                       func(*s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error)
	GetBucketWebsiteRequestFn                   func(*s3.GetBucketWebsiteInput) (*request.Request, *s3.GetBucketWebsiteOutput)
	GetBucketWebsiteFn                          func(*s3.GetBucketWebsiteInput) (*s3.GetBucketWebsiteOutput, error)
	GetObjectRequestFn                          func(*s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput)
	GetObjectFn                                 func(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	GetObjectAclRequestFn                       func(*s3.GetObjectAclInput) (*request.Request, *s3.GetObjectAclOutput)
	GetObjectAclFn                              func(*s3.GetObjectAclInput) (*s3.GetObjectAclOutput, error)
	GetObjectTorrentRequestFn                   func(*s3.GetObjectTorrentInput) (*request.Request, *s3.GetObjectTorrentOutput)
	GetObjectTorrentFn                          func(*s3.GetObjectTorrentInput) (*s3.GetObjectTorrentOutput, error)
	HeadBucketRequestFn                         func(*s3.HeadBucketInput) (*request.Request, *s3.HeadBucketOutput)
	HeadBucketFn                                func(*s3.HeadBucketInput) (*s3.HeadBucketOutput, error)
	HeadObjectRequestFn                         func(*s3.HeadObjectInput) (*request.Request, *s3.HeadObjectOutput)
	HeadObjectFn                                func(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
	ListBucketsRequestFn                        func(*s3.ListBucketsInput) (*request.Request, *s3.ListBucketsOutput)
	ListBucketsFn                               func(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	ListMultipartUploadsRequestFn               func(*s3.ListMultipartUploadsInput) (*request.Request, *s3.ListMultipartUploadsOutput)
	ListMultipartUploadsFn                      func(*s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error)
	ListMultipartUploadsPagesFn                 func(*s3.ListMultipartUploadsInput, func(*s3.ListMultipartUploadsOutput, bool) bool) error
	ListObjectVersionsRequestFn                 func(*s3.ListObjectVersionsInput) (*request.Request, *s3.ListObjectVersionsOutput)
	ListObjectVersionsFn                        func(*s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error)
	ListObjectVersionsPagesFn                   func(*s3.ListObjectVersionsInput, func(*s3.ListObjectVersionsOutput, bool) bool) error
	ListObjectsRequestFn                        func(*s3.ListObjectsInput) (*request.Request, *s3.ListObjectsOutput)
	ListObjectsFn                               func(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	ListObjectsPagesFn                          func(*s3.ListObjectsInput, func(*s3.ListObjectsOutput, bool) bool) error
	//	ListObjectsV2RequestFn                      func(*s3.ListObjectsV2Input) (*request.Request, *s3.ListObjectsV2Output)
	//	ListObjectsV2Fn                             func(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	//	ListObjectsV2PagesFn                        func(*s3.ListObjectsV2Input, func(*s3.ListObjectsV2Output, bool) bool) error
	ListPartsRequestFn func(*s3.ListPartsInput) (*request.Request, *s3.ListPartsOutput)
	ListPartsFn        func(*s3.ListPartsInput) (*s3.ListPartsOutput, error)
	ListPartsPagesFn   func(*s3.ListPartsInput, func(*s3.ListPartsOutput, bool) bool) error
	//	PutBucketAccelerateConfigurationRequestFn   func(*s3.PutBucketAccelerateConfigurationInput) (*request.Request, *s3.PutBucketAccelerateConfigurationOutput)
	//	PutBucketAccelerateConfigurationFn          func(*s3.PutBucketAccelerateConfigurationInput) (*s3.PutBucketAccelerateConfigurationOutput, error)
	PutBucketAclRequestFn                       func(*s3.PutBucketAclInput) (*request.Request, *s3.PutBucketAclOutput)
	PutBucketAclFn                              func(*s3.PutBucketAclInput) (*s3.PutBucketAclOutput, error)
	PutBucketCorsRequestFn                      func(*s3.PutBucketCorsInput) (*request.Request, *s3.PutBucketCorsOutput)
	PutBucketCorsFn                             func(*s3.PutBucketCorsInput) (*s3.PutBucketCorsOutput, error)
	PutBucketLifecycleRequestFn                 func(*s3.PutBucketLifecycleInput) (*request.Request, *s3.PutBucketLifecycleOutput)
	PutBucketLifecycleFn                        func(*s3.PutBucketLifecycleInput) (*s3.PutBucketLifecycleOutput, error)
	PutBucketLifecycleConfigurationRequestFn    func(*s3.PutBucketLifecycleConfigurationInput) (*request.Request, *s3.PutBucketLifecycleConfigurationOutput)
	PutBucketLifecycleConfigurationFn           func(*s3.PutBucketLifecycleConfigurationInput) (*s3.PutBucketLifecycleConfigurationOutput, error)
	PutBucketLoggingRequestFn                   func(*s3.PutBucketLoggingInput) (*request.Request, *s3.PutBucketLoggingOutput)
	PutBucketLoggingFn                          func(*s3.PutBucketLoggingInput) (*s3.PutBucketLoggingOutput, error)
	PutBucketNotificationRequestFn              func(*s3.PutBucketNotificationInput) (*request.Request, *s3.PutBucketNotificationOutput)
	PutBucketNotificationFn                     func(*s3.PutBucketNotificationInput) (*s3.PutBucketNotificationOutput, error)
	PutBucketNotificationConfigurationRequestFn func(*s3.PutBucketNotificationConfigurationInput) (*request.Request, *s3.PutBucketNotificationConfigurationOutput)
	PutBucketNotificationConfigurationFn        func(*s3.PutBucketNotificationConfigurationInput) (*s3.PutBucketNotificationConfigurationOutput, error)
	PutBucketPolicyRequestFn                    func(*s3.PutBucketPolicyInput) (*request.Request, *s3.PutBucketPolicyOutput)
	PutBucketPolicyFn                           func(*s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error)
	PutBucketReplicationRequestFn               func(*s3.PutBucketReplicationInput) (*request.Request, *s3.PutBucketReplicationOutput)
	PutBucketReplicationFn                      func(*s3.PutBucketReplicationInput) (*s3.PutBucketReplicationOutput, error)
	PutBucketRequestPaymentRequestFn            func(*s3.PutBucketRequestPaymentInput) (*request.Request, *s3.PutBucketRequestPaymentOutput)
	PutBucketRequestPaymentFn                   func(*s3.PutBucketRequestPaymentInput) (*s3.PutBucketRequestPaymentOutput, error)
	PutBucketTaggingRequestFn                   func(*s3.PutBucketTaggingInput) (*request.Request, *s3.PutBucketTaggingOutput)
	PutBucketTaggingFn                          func(*s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error)
	PutBucketVersioningRequestFn                func(*s3.PutBucketVersioningInput) (*request.Request, *s3.PutBucketVersioningOutput)
	PutBucketVersioningFn                       func(*s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error)
	PutBucketWebsiteRequestFn                   func(*s3.PutBucketWebsiteInput) (*request.Request, *s3.PutBucketWebsiteOutput)
	PutBucketWebsiteFn                          func(*s3.PutBucketWebsiteInput) (*s3.PutBucketWebsiteOutput, error)
	PutObjectRequestFn                          func(*s3.PutObjectInput) (*request.Request, *s3.PutObjectOutput)
	PutObjectFn                                 func(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	PutObjectAclRequestFn                       func(*s3.PutObjectAclInput) (*request.Request, *s3.PutObjectAclOutput)
	PutObjectAclFn                              func(*s3.PutObjectAclInput) (*s3.PutObjectAclOutput, error)
	RestoreObjectRequestFn                      func(*s3.RestoreObjectInput) (*request.Request, *s3.RestoreObjectOutput)
	RestoreObjectFn                             func(*s3.RestoreObjectInput) (*s3.RestoreObjectOutput, error)
	UploadPartRequestFn                         func(*s3.UploadPartInput) (*request.Request, *s3.UploadPartOutput)
	UploadPartFn                                func(*s3.UploadPartInput) (*s3.UploadPartOutput, error)
	UploadPartCopyRequestFn                     func(*s3.UploadPartCopyInput) (*request.Request, *s3.UploadPartCopyOutput)
	UploadPartCopyFn                            func(*s3.UploadPartCopyInput) (*s3.UploadPartCopyOutput, error)
	WaitUntilBucketExistsFn                     func(*s3.HeadBucketInput) error
	WaitUntilBucketNotExistsFn                  func(*s3.HeadBucketInput) error
	WaitUntilObjectExistsFn                     func(*s3.HeadObjectInput) error
	WaitUntilObjectNotExistsFn                  func(*s3.HeadObjectInput) error
}

func (f *S3) AbortMultipartUploadRequest(input *s3.AbortMultipartUploadInput) (*request.Request, *s3.AbortMultipartUploadOutput) {
	if f.AbortMultipartUploadRequestFn == nil {
		return nil, nil
	}
	return f.AbortMultipartUploadRequestFn(input)
}
func (f *S3) AbortMultipartUpload(input *s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error) {
	if f.AbortMultipartUploadFn == nil {
		return nil, nil
	}
	return f.AbortMultipartUploadFn(input)
}
func (f *S3) CompleteMultipartUploadRequest(input *s3.CompleteMultipartUploadInput) (*request.Request, *s3.CompleteMultipartUploadOutput) {
	if f.CompleteMultipartUploadRequestFn == nil {
		return nil, nil
	}
	return f.CompleteMultipartUploadRequestFn(input)
}
func (f *S3) CompleteMultipartUpload(input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
	if f.CompleteMultipartUploadFn == nil {
		return nil, nil
	}
	return f.CompleteMultipartUploadFn(input)
}
func (f *S3) CopyObjectRequest(input *s3.CopyObjectInput) (*request.Request, *s3.CopyObjectOutput) {
	if f.CopyObjectRequestFn == nil {
		return nil, nil
	}
	return f.CopyObjectRequestFn(input)
}
func (f *S3) CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	if f.CopyObjectFn == nil {
		return nil, nil
	}
	return f.CopyObjectFn(input)
}
func (f *S3) CreateBucketRequest(input *s3.CreateBucketInput) (*request.Request, *s3.CreateBucketOutput) {
	if f.CreateBucketRequestFn == nil {
		return nil, nil
	}
	return f.CreateBucketRequestFn(input)
}
func (f *S3) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	if f.CreateBucketFn == nil {
		return nil, nil
	}
	return f.CreateBucketFn(input)
}
func (f *S3) CreateMultipartUploadRequest(input *s3.CreateMultipartUploadInput) (*request.Request, *s3.CreateMultipartUploadOutput) {
	if f.CreateMultipartUploadRequestFn == nil {
		return nil, nil
	}
	return f.CreateMultipartUploadRequestFn(input)
}
func (f *S3) CreateMultipartUpload(input *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
	if f.CreateMultipartUploadFn == nil {
		return nil, nil
	}
	return f.CreateMultipartUploadFn(input)
}
func (f *S3) DeleteBucketRequest(input *s3.DeleteBucketInput) (*request.Request, *s3.DeleteBucketOutput) {
	if f.DeleteBucketRequestFn == nil {
		return nil, nil
	}
	return f.DeleteBucketRequestFn(input)
}
func (f *S3) DeleteBucket(input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	if f.DeleteBucketFn == nil {
		return nil, nil
	}
	return f.DeleteBucketFn(input)
}
func (f *S3) DeleteBucketCorsRequest(input *s3.DeleteBucketCorsInput) (*request.Request, *s3.DeleteBucketCorsOutput) {
	if f.DeleteBucketCorsRequestFn == nil {
		return nil, nil
	}
	return f.DeleteBucketCorsRequestFn(input)
}
func (f *S3) DeleteBucketCors(input *s3.DeleteBucketCorsInput) (*s3.DeleteBucketCorsOutput, error) {
	if f.DeleteBucketCorsFn == nil {
		return nil, nil
	}
	return f.DeleteBucketCorsFn(input)
}
func (f *S3) DeleteBucketLifecycleRequest(input *s3.DeleteBucketLifecycleInput) (*request.Request, *s3.DeleteBucketLifecycleOutput) {
	if f.DeleteBucketLifecycleRequestFn == nil {
		return nil, nil
	}
	return f.DeleteBucketLifecycleRequestFn(input)
}
func (f *S3) DeleteBucketLifecycle(input *s3.DeleteBucketLifecycleInput) (*s3.DeleteBucketLifecycleOutput, error) {
	if f.DeleteBucketLifecycleFn == nil {
		return nil, nil
	}
	return f.DeleteBucketLifecycleFn(input)
}
func (f *S3) DeleteBucketPolicyRequest(input *s3.DeleteBucketPolicyInput) (*request.Request, *s3.DeleteBucketPolicyOutput) {
	if f.DeleteBucketPolicyRequestFn == nil {
		return nil, nil
	}
	return f.DeleteBucketPolicyRequestFn(input)
}
func (f *S3) DeleteBucketPolicy(input *s3.DeleteBucketPolicyInput) (*s3.DeleteBucketPolicyOutput, error) {
	if f.DeleteBucketPolicyFn == nil {
		return nil, nil
	}
	return f.DeleteBucketPolicyFn(input)
}
func (f *S3) DeleteBucketReplicationRequest(input *s3.DeleteBucketReplicationInput) (*request.Request, *s3.DeleteBucketReplicationOutput) {
	if f.DeleteBucketReplicationRequestFn == nil {
		return nil, nil
	}
	return f.DeleteBucketReplicationRequestFn(input)
}
func (f *S3) DeleteBucketReplication(input *s3.DeleteBucketReplicationInput) (*s3.DeleteBucketReplicationOutput, error) {
	if f.DeleteBucketReplicationFn == nil {
		return nil, nil
	}
	return f.DeleteBucketReplicationFn(input)
}
func (f *S3) DeleteBucketTaggingRequest(input *s3.DeleteBucketTaggingInput) (*request.Request, *s3.DeleteBucketTaggingOutput) {
	if f.DeleteBucketTaggingRequestFn == nil {
		return nil, nil
	}
	return f.DeleteBucketTaggingRequestFn(input)
}
func (f *S3) DeleteBucketTagging(input *s3.DeleteBucketTaggingInput) (*s3.DeleteBucketTaggingOutput, error) {
	if f.DeleteBucketTaggingFn == nil {
		return nil, nil
	}
	return f.DeleteBucketTaggingFn(input)
}
func (f *S3) DeleteBucketWebsiteRequest(input *s3.DeleteBucketWebsiteInput) (*request.Request, *s3.DeleteBucketWebsiteOutput) {
	if f.DeleteBucketWebsiteRequestFn == nil {
		return nil, nil
	}
	return f.DeleteBucketWebsiteRequestFn(input)
}
func (f *S3) DeleteBucketWebsite(input *s3.DeleteBucketWebsiteInput) (*s3.DeleteBucketWebsiteOutput, error) {
	if f.DeleteBucketWebsiteFn == nil {
		return nil, nil
	}
	return f.DeleteBucketWebsiteFn(input)
}
func (f *S3) DeleteObjectRequest(input *s3.DeleteObjectInput) (*request.Request, *s3.DeleteObjectOutput) {
	if f.DeleteObjectRequestFn == nil {
		return nil, nil
	}
	return f.DeleteObjectRequestFn(input)
}
func (f *S3) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	if f.DeleteObjectFn == nil {
		return nil, nil
	}
	return f.DeleteObjectFn(input)
}
func (f *S3) DeleteObjectsRequest(input *s3.DeleteObjectsInput) (*request.Request, *s3.DeleteObjectsOutput) {
	if f.DeleteObjectsRequestFn == nil {
		return nil, nil
	}
	return f.DeleteObjectsRequestFn(input)
}
func (f *S3) DeleteObjects(input *s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error) {
	if f.DeleteObjectsFn == nil {
		return nil, nil
	}
	return f.DeleteObjectsFn(input)
}

//func (f *S3) GetBucketAccelerateConfigurationRequest(input *s3.GetBucketAccelerateConfigurationInput) (*request.Request, *s3.GetBucketAccelerateConfigurationOutput) {
//	if f.GetBucketAccelerateConfigurationRequestFn == nil {
//		return nil, nil
//	}
//	return f.GetBucketAccelerateConfigurationRequestFn(input)
//}
//func (f *S3) GetBucketAccelerateConfiguration(input *s3.GetBucketAccelerateConfigurationInput) (*s3.GetBucketAccelerateConfigurationOutput, error) {
//	if f.GetBucketAccelerateConfigurationFn == nil {
//		return nil, nil
//	}
//	return f.GetBucketAccelerateConfigurationFn(input)
//}
func (f *S3) GetBucketAclRequest(input *s3.GetBucketAclInput) (*request.Request, *s3.GetBucketAclOutput) {
	if f.GetBucketAclRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketAclRequestFn(input)
}
func (f *S3) GetBucketAcl(input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error) {
	if f.GetBucketAclFn == nil {
		return nil, nil
	}
	return f.GetBucketAclFn(input)
}
func (f *S3) GetBucketCorsRequest(input *s3.GetBucketCorsInput) (*request.Request, *s3.GetBucketCorsOutput) {
	if f.GetBucketCorsRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketCorsRequestFn(input)
}
func (f *S3) GetBucketCors(input *s3.GetBucketCorsInput) (*s3.GetBucketCorsOutput, error) {
	if f.GetBucketCorsFn == nil {
		return nil, nil
	}
	return f.GetBucketCorsFn(input)
}
func (f *S3) GetBucketLifecycleRequest(input *s3.GetBucketLifecycleInput) (*request.Request, *s3.GetBucketLifecycleOutput) {
	if f.GetBucketLifecycleRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketLifecycleRequestFn(input)
}
func (f *S3) GetBucketLifecycle(input *s3.GetBucketLifecycleInput) (*s3.GetBucketLifecycleOutput, error) {
	if f.GetBucketLifecycleFn == nil {
		return nil, nil
	}
	return f.GetBucketLifecycleFn(input)
}
func (f *S3) GetBucketLifecycleConfigurationRequest(input *s3.GetBucketLifecycleConfigurationInput) (*request.Request, *s3.GetBucketLifecycleConfigurationOutput) {
	if f.GetBucketLifecycleConfigurationRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketLifecycleConfigurationRequestFn(input)
}
func (f *S3) GetBucketLifecycleConfiguration(input *s3.GetBucketLifecycleConfigurationInput) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	if f.GetBucketLifecycleConfigurationFn == nil {
		return nil, nil
	}
	return f.GetBucketLifecycleConfigurationFn(input)
}
func (f *S3) GetBucketLocationRequest(input *s3.GetBucketLocationInput) (*request.Request, *s3.GetBucketLocationOutput) {
	if f.GetBucketLocationRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketLocationRequestFn(input)
}
func (f *S3) GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	if f.GetBucketLocationFn == nil {
		return nil, nil
	}
	return f.GetBucketLocationFn(input)
}
func (f *S3) GetBucketLoggingRequest(input *s3.GetBucketLoggingInput) (*request.Request, *s3.GetBucketLoggingOutput) {
	if f.GetBucketLoggingRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketLoggingRequestFn(input)
}
func (f *S3) GetBucketLogging(input *s3.GetBucketLoggingInput) (*s3.GetBucketLoggingOutput, error) {
	if f.GetBucketLoggingFn == nil {
		return nil, nil
	}
	return f.GetBucketLoggingFn(input)
}
func (f *S3) GetBucketNotificationRequest(input *s3.GetBucketNotificationConfigurationRequest) (*request.Request, *s3.NotificationConfigurationDeprecated) {
	if f.GetBucketNotificationRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketNotificationRequestFn(input)
}
func (f *S3) GetBucketNotification(input *s3.GetBucketNotificationConfigurationRequest) (*s3.NotificationConfigurationDeprecated, error) {
	if f.GetBucketNotificationFn == nil {
		return nil, nil
	}
	return f.GetBucketNotificationFn(input)
}
func (f *S3) GetBucketNotificationConfigurationRequest(input *s3.GetBucketNotificationConfigurationRequest) (*request.Request, *s3.NotificationConfiguration) {
	if f.GetBucketNotificationConfigurationRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketNotificationConfigurationRequestFn(input)
}
func (f *S3) GetBucketNotificationConfiguration(input *s3.GetBucketNotificationConfigurationRequest) (*s3.NotificationConfiguration, error) {
	if f.GetBucketNotificationConfigurationFn == nil {
		return nil, nil
	}
	return f.GetBucketNotificationConfigurationFn(input)
}
func (f *S3) GetBucketPolicyRequest(input *s3.GetBucketPolicyInput) (*request.Request, *s3.GetBucketPolicyOutput) {
	if f.GetBucketPolicyRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketPolicyRequestFn(input)
}
func (f *S3) GetBucketPolicy(input *s3.GetBucketPolicyInput) (*s3.GetBucketPolicyOutput, error) {
	if f.GetBucketPolicyFn == nil {
		return nil, nil
	}
	return f.GetBucketPolicyFn(input)
}
func (f *S3) GetBucketReplicationRequest(input *s3.GetBucketReplicationInput) (*request.Request, *s3.GetBucketReplicationOutput) {
	if f.GetBucketReplicationRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketReplicationRequestFn(input)
}
func (f *S3) GetBucketReplication(input *s3.GetBucketReplicationInput) (*s3.GetBucketReplicationOutput, error) {
	if f.GetBucketReplicationFn == nil {
		return nil, nil
	}
	return f.GetBucketReplicationFn(input)
}
func (f *S3) GetBucketRequestPaymentRequest(input *s3.GetBucketRequestPaymentInput) (*request.Request, *s3.GetBucketRequestPaymentOutput) {
	if f.GetBucketRequestPaymentRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketRequestPaymentRequestFn(input)
}
func (f *S3) GetBucketRequestPayment(input *s3.GetBucketRequestPaymentInput) (*s3.GetBucketRequestPaymentOutput, error) {
	if f.GetBucketRequestPaymentFn == nil {
		return nil, nil
	}
	return f.GetBucketRequestPaymentFn(input)
}
func (f *S3) GetBucketTaggingRequest(input *s3.GetBucketTaggingInput) (*request.Request, *s3.GetBucketTaggingOutput) {
	if f.GetBucketTaggingRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketTaggingRequestFn(input)
}
func (f *S3) GetBucketTagging(input *s3.GetBucketTaggingInput) (*s3.GetBucketTaggingOutput, error) {
	if f.GetBucketTaggingFn == nil {
		return nil, nil
	}
	return f.GetBucketTaggingFn(input)
}
func (f *S3) GetBucketVersioningRequest(input *s3.GetBucketVersioningInput) (*request.Request, *s3.GetBucketVersioningOutput) {
	if f.GetBucketVersioningRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketVersioningRequestFn(input)
}
func (f *S3) GetBucketVersioning(input *s3.GetBucketVersioningInput) (*s3.GetBucketVersioningOutput, error) {
	if f.GetBucketVersioningFn == nil {
		return nil, nil
	}
	return f.GetBucketVersioningFn(input)
}
func (f *S3) GetBucketWebsiteRequest(input *s3.GetBucketWebsiteInput) (*request.Request, *s3.GetBucketWebsiteOutput) {
	if f.GetBucketWebsiteRequestFn == nil {
		return nil, nil
	}
	return f.GetBucketWebsiteRequestFn(input)
}
func (f *S3) GetBucketWebsite(input *s3.GetBucketWebsiteInput) (*s3.GetBucketWebsiteOutput, error) {
	if f.GetBucketWebsiteFn == nil {
		return nil, nil
	}
	return f.GetBucketWebsiteFn(input)
}
func (f *S3) GetObjectRequest(input *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	if f.GetObjectRequestFn == nil {
		return nil, nil
	}
	return f.GetObjectRequestFn(input)
}
func (f *S3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if f.GetObjectFn == nil {
		return nil, nil
	}
	return f.GetObjectFn(input)
}
func (f *S3) GetObjectAclRequest(input *s3.GetObjectAclInput) (*request.Request, *s3.GetObjectAclOutput) {
	if f.GetObjectAclRequestFn == nil {
		return nil, nil
	}
	return f.GetObjectAclRequestFn(input)
}
func (f *S3) GetObjectAcl(input *s3.GetObjectAclInput) (*s3.GetObjectAclOutput, error) {
	if f.GetObjectAclFn == nil {
		return nil, nil
	}
	return f.GetObjectAclFn(input)
}
func (f *S3) GetObjectTorrentRequest(input *s3.GetObjectTorrentInput) (*request.Request, *s3.GetObjectTorrentOutput) {
	if f.GetObjectTorrentRequestFn == nil {
		return nil, nil
	}
	return f.GetObjectTorrentRequestFn(input)
}
func (f *S3) GetObjectTorrent(input *s3.GetObjectTorrentInput) (*s3.GetObjectTorrentOutput, error) {
	if f.GetObjectTorrentFn == nil {
		return nil, nil
	}
	return f.GetObjectTorrentFn(input)
}
func (f *S3) HeadBucketRequest(input *s3.HeadBucketInput) (*request.Request, *s3.HeadBucketOutput) {
	if f.HeadBucketRequestFn == nil {
		return nil, nil
	}
	return f.HeadBucketRequestFn(input)
}
func (f *S3) HeadBucket(input *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	if f.HeadBucketFn == nil {
		return nil, nil
	}
	return f.HeadBucketFn(input)
}
func (f *S3) HeadObjectRequest(input *s3.HeadObjectInput) (*request.Request, *s3.HeadObjectOutput) {
	if f.HeadObjectRequestFn == nil {
		return nil, nil
	}
	return f.HeadObjectRequestFn(input)
}
func (f *S3) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	if f.HeadObjectFn == nil {
		return nil, nil
	}
	return f.HeadObjectFn(input)
}
func (f *S3) ListBucketsRequest(input *s3.ListBucketsInput) (*request.Request, *s3.ListBucketsOutput) {
	if f.ListBucketsRequestFn == nil {
		return nil, nil
	}
	return f.ListBucketsRequestFn(input)
}
func (f *S3) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	if f.ListBucketsFn == nil {
		return nil, nil
	}
	return f.ListBucketsFn(input)
}
func (f *S3) ListMultipartUploadsRequest(input *s3.ListMultipartUploadsInput) (*request.Request, *s3.ListMultipartUploadsOutput) {
	if f.ListMultipartUploadsRequestFn == nil {
		return nil, nil
	}
	return f.ListMultipartUploadsRequestFn(input)
}
func (f *S3) ListMultipartUploads(input *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
	if f.ListMultipartUploadsFn == nil {
		return nil, nil
	}
	return f.ListMultipartUploadsFn(input)
}

func (f *S3) ListMultipartUploadsPages(input *s3.ListMultipartUploadsInput, fn func(*s3.ListMultipartUploadsOutput, bool) bool) error {
	if f.ListMultipartUploadsPagesFn == nil {
		return nil
	}
	return f.ListMultipartUploadsPagesFn(input, fn)
}
func (f *S3) ListObjectVersionsRequest(input *s3.ListObjectVersionsInput) (*request.Request, *s3.ListObjectVersionsOutput) {
	if f.ListObjectVersionsRequestFn == nil {
		return nil, nil
	}
	return f.ListObjectVersionsRequestFn(input)
}
func (f *S3) ListObjectVersions(input *s3.ListObjectVersionsInput) (*s3.ListObjectVersionsOutput, error) {
	if f.ListObjectVersionsFn == nil {
		return nil, nil
	}
	return f.ListObjectVersionsFn(input)
}

func (f *S3) ListObjectVersionsPages(input *s3.ListObjectVersionsInput, fn func(*s3.ListObjectVersionsOutput, bool) bool) error {
	if f.ListObjectVersionsPagesFn == nil {
		return nil
	}
	return f.ListObjectVersionsPagesFn(input, fn)
}
func (f *S3) ListObjectsRequest(input *s3.ListObjectsInput) (*request.Request, *s3.ListObjectsOutput) {
	if f.ListObjectsRequestFn == nil {
		return nil, nil
	}
	return f.ListObjectsRequestFn(input)
}
func (f *S3) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	if f.ListObjectsFn == nil {
		return nil, nil
	}
	return f.ListObjectsFn(input)
}

func (f *S3) ListObjectsPages(input *s3.ListObjectsInput, fn func(*s3.ListObjectsOutput, bool) bool) error {
	if f.ListObjectsPagesFn == nil {
		return nil
	}
	return f.ListObjectsPagesFn(input, fn)
}

//func (f *S3) ListObjectsV2Request(*s3.ListObjectsV2Input) (*request.Request, *s3.ListObjectsV2Output) {
//	if f.AttachUserPolicyRequestFn == nil {
//		return nil, nil
//	}
//	return f.AttachUserPolicyRequestFn(input)
//}
//func (f *S3) ListObjectsV2(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
//	if f.AttachUserPolicyRequestFn == nil {
//		return nil, nil
//	}
//	return f.AttachUserPolicyRequestFn(input)
//}
//func (f *S3) ListObjectsV2Pages(*s3.ListObjectsV2Input) error {
//	if f.ListObjectsV2PagesFn == nil {
//		return nil
//	}
//	return f.AttachUserPolicyRequestFn(input)
//}
func (f *S3) ListPartsRequest(input *s3.ListPartsInput) (*request.Request, *s3.ListPartsOutput) {
	if f.ListPartsRequestFn == nil {
		return nil, nil
	}
	return f.ListPartsRequestFn(input)
}
func (f *S3) ListParts(input *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
	if f.ListPartsFn == nil {
		return nil, nil
	}
	return f.ListPartsFn(input)
}

func (f *S3) ListPartsPages(input *s3.ListPartsInput, fn func(*s3.ListPartsOutput, bool) bool) error {
	if f.ListPartsPagesFn == nil {
		return nil
	}
	return f.ListPartsPagesFn(input, fn)
}

//func (f *S3) PutBucketAccelerateConfigurationRequest(*s3.PutBucketAccelerateConfigurationInput) (*request.Request, *s3.PutBucketAccelerateConfigurationOutput) {
//	if f.AttachUserPolicyRequestFn == nil {
//		return nil, nil
//	}
//	return f.AttachUserPolicyRequestFn(input)
//}
//func (f *S3) PutBucketAccelerateConfiguration(*s3.PutBucketAccelerateConfigurationInput) (*s3.PutBucketAccelerateConfigurationOutput, error) {
//	if f.AttachUserPolicyRequestFn == nil {
//		return nil, nil
//	}
//	return f.AttachUserPolicyRequestFn(input)
//}
func (f *S3) PutBucketAclRequest(input *s3.PutBucketAclInput) (*request.Request, *s3.PutBucketAclOutput) {
	if f.PutBucketAclRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketAclRequestFn(input)
}
func (f *S3) PutBucketAcl(input *s3.PutBucketAclInput) (*s3.PutBucketAclOutput, error) {
	if f.PutBucketAclFn == nil {
		return nil, nil
	}
	return f.PutBucketAclFn(input)
}
func (f *S3) PutBucketCorsRequest(input *s3.PutBucketCorsInput) (*request.Request, *s3.PutBucketCorsOutput) {
	if f.PutBucketCorsRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketCorsRequestFn(input)
}
func (f *S3) PutBucketCors(input *s3.PutBucketCorsInput) (*s3.PutBucketCorsOutput, error) {
	if f.PutBucketCorsFn == nil {
		return nil, nil
	}
	return f.PutBucketCorsFn(input)
}
func (f *S3) PutBucketLifecycleRequest(input *s3.PutBucketLifecycleInput) (*request.Request, *s3.PutBucketLifecycleOutput) {
	if f.PutBucketLifecycleRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketLifecycleRequestFn(input)
}
func (f *S3) PutBucketLifecycle(input *s3.PutBucketLifecycleInput) (*s3.PutBucketLifecycleOutput, error) {
	if f.PutBucketLifecycleFn == nil {
		return nil, nil
	}
	return f.PutBucketLifecycleFn(input)
}
func (f *S3) PutBucketLifecycleConfigurationRequest(input *s3.PutBucketLifecycleConfigurationInput) (*request.Request, *s3.PutBucketLifecycleConfigurationOutput) {
	if f.PutBucketLifecycleConfigurationRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketLifecycleConfigurationRequestFn(input)
}
func (f *S3) PutBucketLifecycleConfiguration(input *s3.PutBucketLifecycleConfigurationInput) (*s3.PutBucketLifecycleConfigurationOutput, error) {
	if f.PutBucketLifecycleConfigurationFn == nil {
		return nil, nil
	}
	return f.PutBucketLifecycleConfigurationFn(input)
}
func (f *S3) PutBucketLoggingRequest(input *s3.PutBucketLoggingInput) (*request.Request, *s3.PutBucketLoggingOutput) {
	if f.PutBucketLoggingRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketLoggingRequestFn(input)
}
func (f *S3) PutBucketLogging(input *s3.PutBucketLoggingInput) (*s3.PutBucketLoggingOutput, error) {
	if f.PutBucketLoggingFn == nil {
		return nil, nil
	}
	return f.PutBucketLoggingFn(input)
}
func (f *S3) PutBucketNotificationRequest(input *s3.PutBucketNotificationInput) (*request.Request, *s3.PutBucketNotificationOutput) {
	if f.PutBucketNotificationRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketNotificationRequestFn(input)
}
func (f *S3) PutBucketNotification(input *s3.PutBucketNotificationInput) (*s3.PutBucketNotificationOutput, error) {
	if f.PutBucketNotificationFn == nil {
		return nil, nil
	}
	return f.PutBucketNotificationFn(input)
}
func (f *S3) PutBucketNotificationConfigurationRequest(input *s3.PutBucketNotificationConfigurationInput) (*request.Request, *s3.PutBucketNotificationConfigurationOutput) {
	if f.PutBucketNotificationConfigurationRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketNotificationConfigurationRequestFn(input)
}
func (f *S3) PutBucketNotificationConfiguration(input *s3.PutBucketNotificationConfigurationInput) (*s3.PutBucketNotificationConfigurationOutput, error) {
	if f.PutBucketNotificationConfigurationFn == nil {
		return nil, nil
	}
	return f.PutBucketNotificationConfigurationFn(input)
}
func (f *S3) PutBucketPolicyRequest(input *s3.PutBucketPolicyInput) (*request.Request, *s3.PutBucketPolicyOutput) {
	if f.PutBucketPolicyRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketPolicyRequestFn(input)
}
func (f *S3) PutBucketPolicy(input *s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error) {
	if f.PutBucketPolicyFn == nil {
		return nil, nil
	}
	return f.PutBucketPolicyFn(input)
}
func (f *S3) PutBucketReplicationRequest(input *s3.PutBucketReplicationInput) (*request.Request, *s3.PutBucketReplicationOutput) {
	if f.PutBucketReplicationRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketReplicationRequestFn(input)
}
func (f *S3) PutBucketReplication(input *s3.PutBucketReplicationInput) (*s3.PutBucketReplicationOutput, error) {
	if f.PutBucketReplicationFn == nil {
		return nil, nil
	}
	return f.PutBucketReplicationFn(input)
}
func (f *S3) PutBucketRequestPaymentRequest(input *s3.PutBucketRequestPaymentInput) (*request.Request, *s3.PutBucketRequestPaymentOutput) {
	if f.PutBucketRequestPaymentRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketRequestPaymentRequestFn(input)
}
func (f *S3) PutBucketRequestPayment(input *s3.PutBucketRequestPaymentInput) (*s3.PutBucketRequestPaymentOutput, error) {
	if f.PutBucketRequestPaymentFn == nil {
		return nil, nil
	}
	return f.PutBucketRequestPaymentFn(input)
}
func (f *S3) PutBucketTaggingRequest(input *s3.PutBucketTaggingInput) (*request.Request, *s3.PutBucketTaggingOutput) {
	if f.PutBucketTaggingRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketTaggingRequestFn(input)
}
func (f *S3) PutBucketTagging(input *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {
	if f.PutBucketTaggingFn == nil {
		return nil, nil
	}
	return f.PutBucketTaggingFn(input)
}
func (f *S3) PutBucketVersioningRequest(input *s3.PutBucketVersioningInput) (*request.Request, *s3.PutBucketVersioningOutput) {
	if f.PutBucketVersioningRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketVersioningRequestFn(input)
}
func (f *S3) PutBucketVersioning(input *s3.PutBucketVersioningInput) (*s3.PutBucketVersioningOutput, error) {
	if f.PutBucketVersioningFn == nil {
		return nil, nil
	}
	return f.PutBucketVersioningFn(input)
}
func (f *S3) PutBucketWebsiteRequest(input *s3.PutBucketWebsiteInput) (*request.Request, *s3.PutBucketWebsiteOutput) {
	if f.PutBucketWebsiteRequestFn == nil {
		return nil, nil
	}
	return f.PutBucketWebsiteRequestFn(input)
}
func (f *S3) PutBucketWebsite(input *s3.PutBucketWebsiteInput) (*s3.PutBucketWebsiteOutput, error) {
	if f.PutBucketWebsiteFn == nil {
		return nil, nil
	}
	return f.PutBucketWebsiteFn(input)
}
func (f *S3) PutObjectRequest(input *s3.PutObjectInput) (*request.Request, *s3.PutObjectOutput) {
	if f.PutObjectRequestFn == nil {
		return nil, nil
	}
	return f.PutObjectRequestFn(input)
}
func (f *S3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	if f.PutObjectFn == nil {
		return nil, nil
	}
	return f.PutObjectFn(input)
}
func (f *S3) PutObjectAclRequest(input *s3.PutObjectAclInput) (*request.Request, *s3.PutObjectAclOutput) {
	if f.PutObjectAclRequestFn == nil {
		return nil, nil
	}
	return f.PutObjectAclRequestFn(input)
}
func (f *S3) PutObjectAcl(input *s3.PutObjectAclInput) (*s3.PutObjectAclOutput, error) {
	if f.PutObjectAclFn == nil {
		return nil, nil
	}
	return f.PutObjectAclFn(input)
}
func (f *S3) RestoreObjectRequest(input *s3.RestoreObjectInput) (*request.Request, *s3.RestoreObjectOutput) {
	if f.RestoreObjectRequestFn == nil {
		return nil, nil
	}
	return f.RestoreObjectRequestFn(input)
}
func (f *S3) RestoreObject(input *s3.RestoreObjectInput) (*s3.RestoreObjectOutput, error) {
	if f.RestoreObjectFn == nil {
		return nil, nil
	}
	return f.RestoreObjectFn(input)
}
func (f *S3) UploadPartRequest(input *s3.UploadPartInput) (*request.Request, *s3.UploadPartOutput) {
	if f.UploadPartRequestFn == nil {
		return nil, nil
	}
	return f.UploadPartRequestFn(input)
}
func (f *S3) UploadPart(input *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
	if f.UploadPartFn == nil {
		return nil, nil
	}
	return f.UploadPartFn(input)
}
func (f *S3) UploadPartCopyRequest(input *s3.UploadPartCopyInput) (*request.Request, *s3.UploadPartCopyOutput) {
	if f.UploadPartCopyRequestFn == nil {
		return nil, nil
	}
	return f.UploadPartCopyRequestFn(input)
}
func (f *S3) UploadPartCopy(input *s3.UploadPartCopyInput) (*s3.UploadPartCopyOutput, error) {
	if f.UploadPartCopyFn == nil {
		return nil, nil
	}
	return f.UploadPartCopyFn(input)
}
func (f *S3) WaitUntilBucketExists(input *s3.HeadBucketInput) error {
	if f.WaitUntilBucketExistsFn == nil {
		return nil
	}
	return f.WaitUntilBucketExistsFn(input)
}
func (f *S3) WaitUntilBucketNotExists(input *s3.HeadBucketInput) error {
	if f.WaitUntilBucketNotExistsFn == nil {
		return nil
	}
	return f.WaitUntilBucketNotExistsFn(input)
}
func (f *S3) WaitUntilObjectExists(input *s3.HeadObjectInput) error {
	if f.WaitUntilObjectExistsFn == nil {
		return nil
	}
	return f.WaitUntilObjectExistsFn(input)
}
func (f *S3) WaitUntilObjectNotExists(input *s3.HeadObjectInput) error {
	if f.WaitUntilObjectNotExistsFn == nil {
		return nil
	}
	return f.WaitUntilObjectNotExistsFn(input)
}
