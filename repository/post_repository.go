package repository

import (
	"context"
	"file-uploader-api/model"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

const name = "postRepository"

var (
	tracer = otel.Tracer(name)
)

type IPostRepository interface {
	GetPosts(posts *[]model.Post) error
	GetPost(post *model.Post, id uint) error
	Create(ctx context.Context, post *model.Post) error
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

func (pr *postRepository) Create(ctx context.Context, post *model.Post) error {
	// トレース
	ctx, span := tracer.Start(ctx, "create")
	defer span.End()
	
	if err := pr.db.Create(&post).Error; err != nil {
		return err
	}
	return nil
}
