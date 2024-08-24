package trace

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type Tracing struct {
	op       *Config
	provider *sdk.TracerProvider
}

func New(opts ...Option) (*Tracing, error) {
	op := &Config{
		Sampler: 1.0,
		Batcher: KindStdout,
	}
	for _, opt := range opts {
		opt.apply(op)
	}

	o := &Tracing{op: op}

	r, err := resource.New(context.Background(),
		resource.WithOS(),
		resource.WithHost(),
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
	case KindZipkin:
		return zipkin.New(t.op.Endpoint)
	case KindOtlpGrpc:
		// Always treat trace exporter as optional component, so we use nonblock here,
		// otherwise this would slow down app start up even set a dial timeout here when
		// endpoint can not reach.
		// If the connection not dial success, the global otel ErrorHandler will catch error
		// when reporting data like other exporters.
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(t.op.Endpoint),
		}
		if len(t.op.OtlpHeaders) > 0 {
			opts = append(opts, otlptracegrpc.WithHeaders(t.op.OtlpHeaders))
		}
		return otlptracegrpc.New(context.Background(), opts...)
	case KindOtlpHttp:
		// Not support flexible configuration now.
		opts := []otlptracehttp.Option{
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(t.op.Endpoint),
		}
		if len(t.op.OtlpHeaders) > 0 {
			opts = append(opts, otlptracehttp.WithHeaders(t.op.OtlpHeaders))
		}
		if len(t.op.OtlpHttpPath) > 0 {
			opts = append(opts, otlptracehttp.WithURLPath(t.op.OtlpHttpPath))
		}
		return otlptracehttp.New(context.Background(), opts...)
	case KindStdout:
		return stdouttrace.New()
	case KindFile:
		f, err := os.OpenFile(t.op.Endpoint, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("file exporter endpoint error: %s", err.Error())
		}
		return stdouttrace.New(stdouttrace.WithWriter(f))
	case KindNoop:
		return tracetest.NewNoopExporter(), nil
	default:
		return nil, fmt.Errorf("unknown exporter: %s", t.op.Batcher)
	}
}

// Shutdown shuts down the span processors in the order they were registered.
func (t *Tracing) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
