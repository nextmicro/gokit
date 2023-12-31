package valueonlycontext

import (
	"context"
	"time"
)

type valueOnlyContext struct {
	context.Context
}

func (c valueOnlyContext) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}

func (valueOnlyContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (valueOnlyContext) Done() <-chan struct{} {
	return nil
}

func (valueOnlyContext) Err() error {
	return nil
}

// ValueOnlyFrom takes all values from the given ctx, without deadline and error control.
func ValueOnlyFrom(ctx context.Context) context.Context {
	return valueOnlyContext{
		Context: ctx,
	}
}
