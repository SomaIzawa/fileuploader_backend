package model

import "time"

type User struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Email string `gorm:"unique;not null"`
	HashedPassword string `gorm:"not null"`
	AllowedAt time.Time `gorm:"default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}