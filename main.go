package main

import (
	"context"
	"file-uploader-api/observability"
	"file-uploader-api/router"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func main()  {
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
	meter, err := observability.MeterSetup(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}
	otel.SetMeterProvider(meter)
	
	// test -----
	
	sampleTracer := otel.Tracer("test-tracer")
	sampleMeter := otel.Meter("test-meter")
	
	// Attributes represent additional key-value descriptors that can be bound
	// to a metric observer or recorder.
	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}
	
	runCount, err := sampleMeter.Int64Counter("run", metric.WithDescription("The number of times the iteration ran"))
	if err != nil {
		log.Fatal(err)
	}
	
	// Work begins
	ctx, span := sampleTracer.Start(
		ctx,
		"CollectorExporter-Example",
		trace.WithAttributes(commonAttrs...))
		defer span.End()
		for i := 0; i < 10; i++ {
			_, iSpan := sampleTracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
			//runCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))
			runCount.Add(ctx, 1)
			log.Printf("Doing really hard work (%d / 10)\n", i+1)
			
			<-time.After(time.Second)
			iSpan.End()
		}
		
		// test -----
		
	e := router.NewRouter()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}