package model

import "time"

type File struct {
	ID        uint   `gorm:"primaryKey"`
	FileName  string `gorm:"not null"`
	Type      string `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
	UserId    uint   `gorm:"not null"`
	Post      Post   `gorm:"foreignKey:PostId; constraint:OnDelete:CASCADE"`
	PostId    uint   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
