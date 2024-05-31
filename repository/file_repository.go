package repository

import (
	"file-uploader-api/model"

	"gorm.io/gorm"
)

type IFileRepository interface {
	GetFile(file *model.File, id uint) error
	DeleteFile(file *model.File, id uint) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) IFileRepository {
	return &fileRepository{db: db}
}

func (fr *fileRepository) GetFile(file *model.File, id uint) error {
	if err := fr.db.Where("id=?", id).First(&file).Error; err != nil {
		return err
	}
	return nil
}

func (fr *fileRepository) DeleteFile(file *model.File, id uint) error {
	if err := fr.db.Where("id=?", id).First(&file).Error; err != nil {
		return err
	}
	if err := fr.db.Delete(&file).Error; err != nil {
		return err
	}
	return nil
}
