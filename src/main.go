package main

import (
	"fmt"
	"os"
)

func configureExporter(exporterType string) {
	// Switch statement
	switch exporterType {
	case "prom":
		fmt.Println("PrometheusExporter")
		configureProm()

	case "otlp":
		fmt.Println("OtlpHttpExporter")
		configureOtlpHttp()

	default:
		fmt.Println("Invalid")
	}
}

func main() {
	// ctx := context.Background()

	exporterType := os.Getenv("EXPORTER_TYPE")
	configureExporter(exporterType)

}
