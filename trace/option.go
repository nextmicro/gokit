package trace

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"

	// TraceName represents the tracing name.
	TraceName = "opentelemetry"
)

type Option interface {
	apply(*Config)
}

// A Config is a opentelemetry config.
type Config struct {
	Endpoint   string
	Sampler    float64
	Batcher    string
	Attributes []attribute.KeyValue
}

type OptionFunc func(*Config)

func (fn OptionFunc) apply(cfg *Config) {
	fn(cfg)
}

func WithName(name string) Option {
	return OptionFunc(func(o *Config) {
		o.Attributes = append(o.Attributes, semconv.ServiceNameKey.String(name))
	})
}

func WithEndpoint(endpoint string) Option {
	return OptionFunc(func(o *Config) {
		o.Endpoint = endpoint
	})
}

func WithSampler(sampler float64) Option {
	return OptionFunc(func(o *Config) {
		o.Sampler = sampler
	})
}

func WithBatcher(batcher string) Option {
	return OptionFunc(func(o *Config) {
		o.Batcher = batcher
	})
}

// WithAttributes adds attributes to the configured Resource.
func WithAttributes(attributes ...attribute.KeyValue) Option {
	return OptionFunc(func(o *Config) {
		o.Attributes = attributes
	})
}
