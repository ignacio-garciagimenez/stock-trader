package main

import (
	"context"
	"net/http"
	"stock-trader/portfolio-service/infrastructure"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	e := echo.New()

	ctx := context.Background()

	res, _ := resource.New(ctx)

	otlptracehttp.WithInsecure()
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("portfolio-service-otelcol:4318"),
		otlptracehttp.WithInsecure(),
	)
	exporter, _ := otlptrace.New(ctx, client)

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer tracerProvider.Shutdown(ctx)
	otel.SetTracerProvider(tracerProvider)

	tracer := otel.Tracer("portfolio-tracer")

	e.Validator = infrastructure.NewRequestValidator()
	db, err := infrastructure.ConnectDB()
	if err != nil {
		panic("Could not connect to the database")
	}

	e.GET("/", func(ctx echo.Context) error {
		_, span := tracer.Start(ctx.Request().Context(), "hello-portfolio")
		defer span.End()
		return ctx.String(http.StatusOK, "Hello from portfolio-service!")
	})
	e.POST("/portfolios", BuildOpenPortfolioFeature(db))

	e.Logger.Fatal(e.Start(":8080"))
}
