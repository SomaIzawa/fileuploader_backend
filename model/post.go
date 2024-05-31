package model

import "time"

type Post struct {
	ID uint `gorm:"primaryKey"`
	Title string `gorm:"not null"`
	Comment string
	ThumbnailType string
	User User `gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
	UserId uint `gorm:"not null"`
	Category Category `gorm:"foreignKey:CategoryId; constraint:OnDelete:CASCADE"`
	CategoryId uint `gorm:"not null"`
	Files []File    `gorm:"foreignKey:PostId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}