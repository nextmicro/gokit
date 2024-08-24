package trace

import (
	"context"
	"testing"
)

func TestMain(t *testing.M) {
	tracing, err := New(WithBatcher(KindNoop), WithName(TraceName))
	if err != nil {
		panic(err)
	}

	defer func() {
		tracing.Shutdown(context.Background())
	}()

	t.Run()
}
