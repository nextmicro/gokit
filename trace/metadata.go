package trace

import (
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// MapCarrier is a TextMapCarrier that uses a map held in memory as a storage
// medium for propagated key-value pairs.
type metadataSupplier struct {
	metadata *metadata.MD
}

// assert that metadataSupplier implements the TextMapCarrier interface
var _ propagation.TextMapCarrier = (*metadataSupplier)(nil)

// Get returns the value associated with the passed key.
func (m *metadataSupplier) Get(key string) string {
	values := m.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

// Set stores the key-value pair.
func (m *metadataSupplier) Set(key, value string) {
	m.metadata.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (m *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*m.metadata))
	for key := range *m.metadata {
		out = append(out, key)
	}

	return out
}
