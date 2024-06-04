package awsmanager

import (
	"crypto/x509"
	"encoding/pem"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type IAwsS3Manager interface {
	GetDownloadLink(key string) (string, error)
	DeleteFile(key string) error
	UploadFile(objectKey string, file multipart.FileHeader) error
}

// TODO: 構造体の名前は要検討
type awsS3Manager struct {
	session    *session.Session
	bucketName string
	client     *s3.S3
}

func NewAwsS3Manager() IAwsS3Manager {
	session := NewS3Session()
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	client := s3.New(session)
	return &awsS3Manager{
		session:    session,
		bucketName: bucketName,
		client:     client,
	}
}

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

func (am *awsS3Manager) UploadFile(objectKey string, file multipart.FileHeader) error {
	u := s3manager.NewUploader(am.session)
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(am.bucketName),
		Key:    aws.String(objectKey),
		Body:   src,
	})
	if err != nil {
		return err
	}

	return nil
}

func (am *awsS3Manager) GetDownloadLink(key string) (string, error) {
	req, _ := am.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(am.bucketName),
		Key:    aws.String(key),
	})

	if presignedURL, err := req.Presign(10 * time.Second); err != nil {
		return "", nil
	} else {
		return presignedURL, nil
	}
}

func (am *awsS3Manager) DeleteFile(key string) error {
	_, err := am.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(am.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	err = am.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(am.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}
	return nil
}

func GenerateSignedURL(resourcePath string) (string, error) {
	// CloudFrontの設定
	keyPairID := os.Getenv("AWS_CLOUDFRONT_KEY_PAIR")

	// 秘密鍵の読み込み
	privateKeyBytes, err := os.ReadFile("private_key.pem")
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(privateKeyBytes))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	signer := sign.NewURLSigner(keyPairID, key)
	signedURL, err := signer.Sign(resourcePath, time.Now().Add(3*time.Second))
	if err != nil {
		return "", err
	}
	return signedURL, nil
}
