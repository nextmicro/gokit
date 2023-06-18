package context

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValueOnlyFrom(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	ctx = ValueOnlyFrom(ctx)
	_, ok := ctx.Deadline()
	assert.NotEqual(t, ok, true)

	assert.Nil(t, ctx.Done())
}

func TestTraceValueOnlyFrom(t *testing.T) {
	ctx := context.WithValue(context.Background(), "trace", "12345556")
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()
	fmt.Println(ctx)

	ctx = ValueOnlyFrom(ctx)
	fmt.Println(ctx)
	_, ok := ctx.Deadline()
	assert.Equal(t, ok, false)
}
