package api

import (
	"cheunn-panaa/golang-microservice/pkg/version"
	"context"

	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *Server) initTracer(ctx context.Context) {
	// Ensure default SDK resources and the required service name are set.
	client := otlptracegrpc.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		s.logger.Error("creating OTLP trace exporter", zap.Error(err))
	}

	s.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(s.config.App.Name),
			semconv.ServiceVersionKey.String(version.VERSION),
		)),
	)

	otel.SetTracerProvider(s.tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	s.tracer = s.tracerProvider.Tracer(
		"tracer-template",
		trace.WithInstrumentationVersion(version.VERSION),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
}

func (s *Server) newOpenTelemetryMiddleware() {
	s.router.Use(otelchi.Middleware(s.config.App.Name, otelchi.WithChiRoutes(s.router)))
}
