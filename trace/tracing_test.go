package trace

import (
	"context"
	"testing"

	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func TestMain(t *testing.M) {
	tracing, err := New(WithBatcher(KindStdout), WithAttributes(semconv.ServiceNameKey.String(TraceName)))
	if err != nil {
		panic(err)
	}

	defer func() {
		tracing.Shutdown(context.Background())
	}()

	t.Run()
}
