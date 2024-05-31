package repository

import (
	"file-uploader-api/model"
	"time"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUser(user *model.User, id uint) error
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
	UpdateAllowAt(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) GetUser(user *model.User, id uint) error {
	if err := ur.db.Where("id=?",id).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateAllowAt(id string) error {
	if err := ur.db.Model(&model.User{}).Where("id=?", id).Update("allowed_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

