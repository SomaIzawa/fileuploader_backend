package main

import (
	"context"
	"file-uploader-api/observability"
	"file-uploader-api/router"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
)

func main()  {
	e := router.NewRouter()
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}
	ctx := context.Background()
	conn, err := observability.InitConn()
	if err != nil {
		log.Fatalln(err)
	}
	tracer, err := observability.TraceSetup(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}
	otel.SetTracerProvider(tracer)
	if err != nil {
		log.Fatalln(err)
	}
	//ここに「observability.TraceSetup()」を追記

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}