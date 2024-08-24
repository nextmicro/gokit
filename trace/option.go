package trace

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	KindNoop     = "noop"
	KindFile     = "file"
	KindStdout   = "stdout"
	KindZipkin   = "zipkin"
	KindOtlpGrpc = "otlpgrpc"
	KindOtlpHttp = "otlphttp"
	// TraceName represents the tracing name.
	TraceName = "opentelemetry"
)

type Option interface {
	apply(*Config)
}

// A Config is a opentelemetry config.
type Config struct {
	Endpoint string
	Sampler  float64
	Batcher  string
	// OtlpHeaders represents the headers for OTLP gRPC or HTTP transport.
	// For example:
	//  uptrace-dsn: 'http://project2_secret_token@localhost:14317/2'
	OtlpHeaders map[string]string
	// OtlpHttpPath represents the path for OTLP HTTP transport.
	// For example
	// /v1/traces
	OtlpHttpPath string
	Attributes   []attribute.KeyValue
}

type OptionFunc func(*Config)

func (fn OptionFunc) apply(cfg *Config) {
	fn(cfg)
}

// WithName sets the service name.
func WithName(name string) Option {
	return OptionFunc(func(o *Config) {
		o.Attributes = append(o.Attributes, attribute.Key("service.name").String(name))
	})
}

// WithEndpoint sets the endpoint.
func WithEndpoint(endpoint string) Option {
	return OptionFunc(func(o *Config) {
		o.Endpoint = endpoint
	})
}

// WithSampler sets the sampler.
func WithSampler(sampler float64) Option {
	return OptionFunc(func(o *Config) {
		o.Sampler = sampler
	})
}

// WithBatcher sets the batcher.
func WithBatcher(batcher string) Option {
	return OptionFunc(func(o *Config) {
		o.Batcher = batcher
	})
}

// WithOtlpHeaders sets the headers for OTLP gRPC or HTTP transport.
func WithOtlpHeaders(headers map[string]string) Option {
	return OptionFunc(func(o *Config) {
		o.OtlpHeaders = headers
	})
}

// WithOtlpHttpPath sets the path for OTLP HTTP transport.
func WithOtlpHttpPath(path string) Option {
	return OptionFunc(func(o *Config) {
		o.OtlpHttpPath = path
	})
}

// WithAttributes adds attributes to the configured Resource.
func WithAttributes(attributes ...attribute.KeyValue) Option {
	return OptionFunc(func(o *Config) {
		o.Attributes = attributes
	})
}
