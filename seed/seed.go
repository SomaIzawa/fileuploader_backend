package main

import (
	"file-uploader-api/db"
	"file-uploader-api/model"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main(){
	dbConn := db.NewDB()
	if err := user(dbConn); err != nil {
		log.Fatalf("failed to seed: %s", err)
	}
	fmt.Println("seed success")
}

func user(db *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := model.User{
		ID: 1,
		Name: "管理ユーザー",
		Email: "admin@example.com",
		AllowedAt: time.Now(),
		HashedPassword: string(hash),
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}