package usecase

import (
	awsutil "file-uploader-api/aws"
	"file-uploader-api/model"
	"file-uploader-api/repository"
	"file-uploader-api/util"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"gorm.io/gorm"
)

type IFileUsecase interface {
	Download(id string) (file model.File, presignedURL string, err error)
	Delete(id string) error
}

type fileUsecase struct {
	fr repository.IFileRepository
	db *gorm.DB
}

func NewFileUsecase(fr repository.IFileRepository, db *gorm.DB) IFileUsecase {
	return &fileUsecase{
		fr: fr, 
		db: db,
	}
}

func (fu *fileUsecase) Download(id string) (model.File, string, error) {
	file := model.File{}
	uintId, err := util.AtoUint(id)
	if err != nil {
		return model.File{}, "", err
	}
	if err := fu.fr.GetFile(&file, uintId); err != nil {
		return model.File{}, "", err
	}

	session := awsutil.NewS3Session()
	s3client := s3.New(session)

	key := fmt.Sprintf("%d.%s", file.ID, file.Type)
	awsBucketName := os.Getenv("AWS_BUCKET_NAME")

	req, _ := s3client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String(key),
	})

	presignedURL, err := req.Presign(10 * time.Second) // 有効期限を15分に設定
	if err != nil {
		return model.File{}, "", err
	}

	return file, presignedURL, nil
}

func (fu *fileUsecase) Delete(id string) error {
	uintId, err := util.AtoUint(id)
	if err != nil {
		return err
	}

	tx := fu.db.Begin()

	file := model.File{}
	fu.fr.DeleteFile(&file, uintId)

	session := awsutil.NewS3Session()
	s3client := s3.New(session)
	awsBucketName := os.Getenv("AWS_BUCKET_NAME")
	key := fmt.Sprintf("%d.%s", file.ID, file.Type)

	_, err = s3client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	err = s3client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}
