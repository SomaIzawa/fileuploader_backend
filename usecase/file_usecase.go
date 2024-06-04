package usecase

import (
	"file-uploader-api/awsmanager"
	"file-uploader-api/model"
	"file-uploader-api/repository"
	"file-uploader-api/util"
	"fmt"

	"gorm.io/gorm"
)

type IFileUsecase interface {
	Download(id string) (file model.File, presignedURL string, err error)
	Delete(id string) error
}

type fileUsecase struct {
	fr repository.IFileRepository
	am awsmanager.IAwsS3Manager
	db *gorm.DB
}

func NewFileUsecase(fr repository.IFileRepository, am awsmanager.IAwsS3Manager, db *gorm.DB) IFileUsecase {
	return &fileUsecase{
		fr: fr,
		am: am,
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

	key := fmt.Sprintf("%d.%s", file.ID, file.Type)

	presignedURL, err := fu.am.GetDownloadLink(key)
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

	key := fmt.Sprintf("%d.%s", file.ID, file.Type)

	err = fu.am.DeleteFile(key)

	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}
