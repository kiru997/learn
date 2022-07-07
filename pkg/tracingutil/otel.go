package tracingutil

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kiru997/go-ex/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitTelemetry(c *configs.BaseConfig, zapLogger *zap.Logger) *prometheus.Exporter {
	var (
		err error
		pe  *prometheus.Exporter
	)

	if c.StatsEnabled {
		config := prometheus.Config{
			DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
		}
		c := controller.New(
			processor.NewFactory(
				selector.NewWithHistogramDistribution(
					histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
				),
				aggregation.CumulativeTemporalitySelector(),
				processor.WithMemory(true),
			),
		)
		pe, err = prometheus.New(config, c)
		if err != nil {
			zapLogger.Error("failed to initialize prometheus exporte", zap.Error(err))
		}

		global.SetMeterProvider(pe.MeterProvider())
	}

	if c.RemoteTrace.Enabled {
		tp, err := tracerProvider(c.RemoteTrace.TraceCollector, c.Name, c.Env)
		if err != nil {
			zapLogger.Panic("err stdout.NewExporter", zap.Error(err))
		}

		// Register our TracerProvider as the global so any imported
		// instrumentation in the future will default to using it.
		otel.SetTracerProvider(tp)
	}

	return pe
}

const tracerKey = "otel-go-contrib-tracer"

type tracerKeyInt int

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	switch c := ctx.(type) {
	case *gin.Context:
		tracer, ok := c.Get(tracerKey)
		if ok {
			c, span := tracer.(trace.Tracer).Start(c.Request.Context(), spanName, opts...)
			c = context.WithValue(c, tracerKeyInt(0), tracer)
			return c, span
		}
	default:
		tracer, ok := ctx.Value(tracerKeyInt(0)).(trace.Tracer)
		if ok {
			return tracer.Start(ctx, spanName, opts...)
		}
	}

	return otel.Tracer("unknown tracer").Start(ctx, spanName, opts...)
}

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url, serviceName, env string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", env),
		)),
	)
	return tp, nil
}
