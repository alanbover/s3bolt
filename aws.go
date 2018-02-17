package s3bolt

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

// SessionParameters TODO
type SessionParameters struct {
	AccessKey  string
	SecretKey  string
	Region     string
	IamRole    string
	IamSession string
}

// NewS3Client TODO
func NewS3Client(parameters *SessionParameters) (*s3.S3, error) {
	session, err := newAwsSession(parameters)
	if err != nil {
		return nil, err
	}
	return s3.New(session), nil
}

func newAwsSession(parameters *SessionParameters) (*session.Session, error) {
	if parameters.Region == "" {
		return nil, errors.New("missing aws Region (required)")
	}

	sess := session.New(&aws.Config{Region: aws.String(parameters.Region)})

	if parameters.AccessKey != "" && parameters.SecretKey != "" {
		sess = session.New(&aws.Config{
			Region:      aws.String(parameters.Region),
			Credentials: credentials.NewStaticCredentials(parameters.AccessKey, parameters.SecretKey, ""),
		})
	}

	if parameters.IamRole != "" {
		creds := assumeRoleCredentials(sess, parameters.IamRole, parameters.IamSession)
		sess.Config.Credentials = creds
	}

	return sess, nil
}

func assumeRoleCredentials(sess *session.Session, IamRole, IamSession string) *credentials.Credentials {

	if IamSession == "" {
		IamSession = "default"
	}

	creds := stscreds.NewCredentials(sess, IamRole, func(o *stscreds.AssumeRoleProvider) {
		o.Duration = time.Hour
		o.ExpiryWindow = 5 * time.Minute
		o.RoleSessionName = IamSession
	})
	return creds
}
