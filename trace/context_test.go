package trace

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestMetadataFromContext(t *testing.T) {
	ctx, span := otel.Tracer(TraceName).Start(context.TODO(), "HTTP Client Get /api/get")
	defer span.End()

	spanContext := trace.SpanContextFromContext(ctx)
	t.Logf("trace_id: %s", spanContext.TraceID().String())
	t.Logf("span_id: %s", spanContext.SpanID().String())

	md := MetadataFromContext(ctx)
	t.Log(md)
}

func TestStartSpanFromMetadata(t *testing.T) {
	ctx, span := otel.Tracer(TraceName).Start(context.TODO(), "HTTP Client Get /api/get")
	defer span.End()

	spanContext := trace.SpanContextFromContext(ctx)
	t.Logf("trace_id: %s", spanContext.TraceID().String())
	t.Logf("span_id: %s", spanContext.SpanID().String())

	md := MetadataFromContext(ctx)

	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second)
	defer cancelFunc()

	ctx = StartSpanFromMetadata(ctx, "TestStartSpanFromMetadata", md)

	spanContext = trace.SpanContextFromContext(ctx)
	t.Logf("trace_id: %s", spanContext.TraceID().String())
	t.Logf("span_id: %s", spanContext.SpanID().String())
}

func TestExtractTraceId(t *testing.T) {
	ctx, span := otel.Tracer(TraceName).Start(context.TODO(), "HTTP Client Get /api/get")
	defer span.End()

	traceId := ExtractTraceId(ctx)
	t.Logf("trace_id: %s", traceId)
}

func TestExtractSpanId(t *testing.T) {
	ctx, span := otel.Tracer(TraceName).Start(context.TODO(), "HTTP Client Get /api/get")
	defer span.End()

	spanId := ExtractSpanId(ctx)
	t.Logf("span_id: %s", spanId)
}
