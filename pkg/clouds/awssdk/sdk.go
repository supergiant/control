package awssdk

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//SDK is a wrapper around aws client library handling concerns like authentication
//currently it only supports accessKey/secret auth method, in the future it may support MFA
type SDK struct {
	EC2 *ec2.EC2
}

//New creates SDK using keyID/secret and optionally token if temporary credentials auth method is used
func New(region string, keyID, secret, token string) (*SDK, error) {
	if keyID == "" || secret == "" {
		return nil, ErrInvalidCreds
	}

	sdk := &SDK{}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(keyID, secret, token),
		},
	})

	if err != nil {
		return nil, err
	}

	sdk.EC2 = ec2.New(sess)
	return sdk, nil
}
