package telemetry

import (
	"context"
	"server/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var tp *sdktrace.TracerProvider

func InitTracer() (func(context.Context) error, error) {
	cfg := config.GlobalConfig.Telemetry

	exporter, err := zipkin.New(cfg.ZipkinEndpoint)
	if err != nil {
		return nil, err
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
