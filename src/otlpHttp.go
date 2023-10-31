package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	ctx := context.Background()

	// resources := resource.NewWithAttributes(
	// 	semconv.SchemaURL,
	// 	semconv.ServiceNameKey.String("service"),
	// 	semconv.ServiceVersionKey.String("v0.0.0"),
	// )

	// Instantiate the OTLP HTTP exporter
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// Instantiate the OTLP HTTP exporter
	// meterProvider_ := sdk.NewMeterProvider(
	// 	sdk.WithResource(resources),
	// 	sdk.WithReader(sdk.NewPeriodicReader(exporter)),
	// )

	meterProvider := sdk.NewMeterProvider(sdk.WithReader(sdk.NewPeriodicReader(exporter)))

	defer func() {
		err := meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()

	otel.SetMeterProvider(meterProvider)

	// Create an instance on a meter for the given instrumentation scope
	meter := otel.GetMeterProvider().Meter(("my-service-meter"))
	apiCounter, err := meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// for {

		// 	time.Sleep(time.Duration(1) * time.Millisecond)

		apiCounter.Add(ctx, 1)
		fmt.Println("Counter: ", apiCounter)
		// }
	})

	http.ListenAndServe(":8081", nil)
}
