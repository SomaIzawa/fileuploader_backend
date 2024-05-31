package repository

import (
	"file-uploader-api/model"

	"gorm.io/gorm"
)

type ICategoryRepository interface {
	GetCategories(categories *[]model.Category) error
	GetCategory(category *model.Category, id uint) error
	CreateCategory(category *model.Category) (id uint, err error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &categoryRepository{db: db}
}

func (cr *categoryRepository) GetCategories(categories *[]model.Category) error {
	if err := cr.db.Find(&categories).Error; err != nil {
		return err
	}
	return nil
}

func (cr *categoryRepository) GetCategory(category *model.Category, id uint) error {
	if err := cr.db.Where("id=?",id).First(&category).Error; err != nil {
		return err
	}
	return nil
}

func (cr *categoryRepository) CreateCategory(category *model.Category) (id uint, err error) {
	if err := cr.db.Create(category).Error; err != nil {
		return 0, err
	}
	return category.ID, nil
}

