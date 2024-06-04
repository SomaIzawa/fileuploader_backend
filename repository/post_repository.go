package repository

import (
	"file-uploader-api/model"

	"gorm.io/gorm"
)

type IPostRepository interface {
	GetPosts(posts *[]model.Post) error
	GetPost(post *model.Post, id uint) error
	Create(post *model.Post) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) IPostRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) GetPosts(posts *[]model.Post) error {
	if err := pr.db.Preload("User").Preload("Category").Preload("Files").Find(&posts).Error; err != nil {
		return err
	}
	return nil
}

func (pr *postRepository) GetPost(post *model.Post, id uint) error {
	if err := pr.db.Where("id=?", id).Preload("User").Preload("Category").Preload("Files").First(&post).Error; err != nil {
		return err
	}
	return nil
}

func (pr *postRepository) Create(post *model.Post) error {
	if err := pr.db.Create(&post).Error; err != nil {
		return err
	}
	return nil
}
