package aws

import (
	"log"
	"shyft/internal/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewS3Session(conf *models.S3Config) *session.Session {
	// open s3 session
	creds := credentials.NewStaticCredentials(
		conf.AccessKey,
		conf.SecretKey,
		"",
	)
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(conf.Endpoint),
		Region:           aws.String(conf.Region),
		Credentials:      creds,
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Fatalln("error_create_s3_session: ", err)
	}
	return sess
}
