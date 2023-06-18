package trace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
	opt    *tracerOptions
}

// TracerOption is tracing option.
type TracerOption func(*tracerOptions)

type tracerOptions struct {
	propagator propagation.TextMapPropagator
}

// WithPropagator with tracer propagator.
func WithPropagator(propagator propagation.TextMapPropagator) TracerOption {
	return func(opts *tracerOptions) {
		opts.propagator = propagator
	}
}

// NewTracer create tracer instance
func NewTracer(kind trace.SpanKind, opts ...TracerOption) *Tracer {
	op := tracerOptions{
		propagator: propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}),
	}

	for _, o := range opts {
		o(&op)
	}

	tr := &Tracer{tracer: otel.Tracer(TraceName), kind: kind, opt: &op}
	switch kind {
	case trace.SpanKindInternal:
		return tr
	case trace.SpanKindClient, trace.SpanKindProducer:
		return tr
	case trace.SpanKindServer, trace.SpanKindConsumer:
		return tr
	default:
		panic(fmt.Sprintf("[otel] span kind: %v", kind))
	}
}

// Start tracing span.
func (t *Tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	var span trace.Span
	ctx, span = t.tracer.Start(ctx,
		spanName,
		append(opts, trace.WithSpanKind(t.kind))...,
	)

	return ctx, span
}

// Inject set cross-cutting concerns from the Context into the carrier.
func (t *Tracer) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	t.opt.propagator.Inject(ctx, carrier)
}

// Extract reads cross-cutting concerns from the carrier into a Context.
func (t *Tracer) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return t.opt.propagator.Extract(ctx, carrier)
}
