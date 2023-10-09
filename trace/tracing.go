package trace

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Tracing struct {
	op       *Config
	provider *sdk.TracerProvider
}

func New(opts ...Option) (*Tracing, error) {
	op := &Config{
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler:  1.0,
		Batcher:  "jaeger",
	}
	for _, opt := range opts {
		opt.apply(op)
	}

	o := &Tracing{op: op}

	r, err := resource.New(context.Background(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithHostID(),
		resource.WithFromEnv(), // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		resource.WithProcess(), // This option configures a set of Detectors that discover process information
		resource.WithTelemetrySDK(),
		resource.WithAttributes(op.Attributes...),
	)
	if err != nil {
		return nil, err
	}

	options := []sdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		sdk.WithSampler(sdk.ParentBased(sdk.TraceIDRatioBased(op.Sampler))),
		// Record information about this application in an Resource.
		sdk.WithResource(r),
	}

	var exp sdk.SpanExporter
	exp, err = o.createExporter()
	if err != nil {
		return nil, err
	}
	// Always be sure to batch in production.
	options = append(options, sdk.WithBatcher(exp))
	o.provider = sdk.NewTracerProvider(options...)
	otel.SetTracerProvider(o.provider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Printf("[otel] error: %v", err)
	}))

	return o, nil
}

func (t *Tracing) createExporter() (sdk.SpanExporter, error) {
	// Just support jaeger and zipkin now, more for later
	switch t.op.Batcher {
	case kindJaeger:
		return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(t.op.Endpoint)))
	case kindZipkin:
		return zipkin.New(t.op.Endpoint)
	default:
		return nil, fmt.Errorf("trace: unknown exporter: %s", t.op.Batcher)
	}
}

// Shutdown shuts down the span processors in the order they were registered.
func (t *Tracing) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
