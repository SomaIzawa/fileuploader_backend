package main

import (
	"file-uploader-api/db"
	"file-uploader-api/model"
	"fmt"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Post{},
		&model.File{},
	)
}