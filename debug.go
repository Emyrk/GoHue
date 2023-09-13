package gohue

import "context"

type debugKey struct{}

func WithDebugging(ctx context.Context) context.Context {
	return context.WithValue(ctx, debugKey{}, true)
}

func getDebugValue(ctx context.Context) bool {
	val := ctx.Value(debugKey{})
	if val == nil {
		return false
	}

	return val.(bool)
}
