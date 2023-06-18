package trace

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// ExtractSpanId returns the current Span's SpanID.
func ExtractSpanId(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

// ExtractTraceId returns the current Span's TraceID.
func ExtractTraceId(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}

// MetadataFromContext Extracting contextual meta-information.
func MetadataFromContext(ctx context.Context) (md metadata.MD) {
	propagators := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
	md = metadata.MD{}
	propagators.Inject(ctx, &metadataSupplier{metadata: &md})
	return md
}

// StartSpanFromMetadata creates a new context with incoming md attached.
// takes all values from the given ctx, without deadline and error control.
func StartSpanFromMetadata(ctx context.Context, spanName string, md metadata.MD, opts ...trace.SpanStartOption) context.Context {
	tr := NewTracer(trace.SpanKindInternal)
	ctx = tr.Extract(ctx, &metadataSupplier{&md})

	var span trace.Span
	ctx, span = tr.Start(ctx, spanName, opts...)
	defer span.End()
	return ctx
}
