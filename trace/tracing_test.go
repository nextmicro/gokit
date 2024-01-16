package trace

import (
	"context"
	"testing"

	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func TestMain(t *testing.M) {
	tracing, err := New(WithBatcher(KindOtlpHttp), WithEndpoint("127.0.0.1:4317"), WithAttributes(semconv.ServiceNameKey.String(TraceName)))
	if err != nil {
		panic(err)
	}

	defer func() {
		tracing.Shutdown(context.Background())
	}()

	t.Run()
}
