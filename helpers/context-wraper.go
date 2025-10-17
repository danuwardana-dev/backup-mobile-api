package helpers

import "context"

type WrappedContext struct {
	context.Context
	parent context.Context
}

func WrapContext(ctx context.Context) context.Context {
	return &WrappedContext{
		Context: context.Background(),
		parent:  ctx,
	}
}

func (w *WrappedContext) Value(key interface{}) interface{} {
	return w.parent.Value(key)
}
