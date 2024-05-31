package main

import (
	"file-uploader-api/router"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main()  {
	e := router.NewRouter()
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}