package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const meterName = "github.com/open-telemetry/opentelemetry-go/example/prometheus"

// type exporterType string

// const (
// 	prometheusExporter exporterType = "prometheus"
// 	otlpGRPCExporter   exporterType = "otlp-grpc"
// 	otlpHTTPExporter   exporterType = "otlp-http"
// )

func newOTLPExporter(t string) (*metric.MeterProvider, func(), error) {
	log.Printf("Exporter Type selected: %v", t)

	var exporter metric.Exporter
	var err error

	switch t {
	case "otlp-grpc":
		exporter, err = otlpmetricgrpc.New(context.Background())

	case "otlp-http":
		exporter, err = otlpmetricgrpc.New(context.Background())
	default:
		panic("invalid exporter type")
	}
	if err != nil {
		return nil, nil, err
	}

	shutdown := func() {
		exporter.Shutdown(context.Background())
	}

	return metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(exporter, metric.WithInterval(3*time.Second)),
		),
	), shutdown, nil
}

// The exporter embeds a default OpenTelemetry Reader and
// implements prometheus.Collector, allowing it to be used as
// both a Reader and Collector.
func newPromExporter() (*metric.MeterProvider, func(), error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}

	return metric.NewMeterProvider(metric.WithReader(exporter)), func() {}, nil
}

func main() {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()

	exporterType := os.Getenv("EXPORTER_TYPE")

	var provider *metric.MeterProvider
	var shutdown func()
	var err error

	switch exporterType {
	case "prometheus":
		// Start the prometheus HTTP server and pass the exporter Collector to it
		go serveHTTP()
		provider, shutdown, err = newPromExporter()
	case "otlp-grpc", "otlp-http":
		provider, shutdown, err = newOTLPExporter(exporterType)

	}
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown()

	meter := provider.Meter(meterName)

	opt := api.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)

	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, opt)

	gauge, err := meter.Float64ObservableGauge("bar", api.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n, opt)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveHTTP() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
