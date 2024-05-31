package aws

import (
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewS3Session() *session.Session {
	// 環境変数からAWSのクレデンシャル情報を取得
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: creds,
	}))
	return session
}

func CreateUploader(s *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(s)
}

func S3Upload(u *s3manager.Uploader, bucketName string, objectKey string, file multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   src,
	})
	if err != nil {
		return err
	}

	return nil
}
