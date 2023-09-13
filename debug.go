package gohue

import (
	"context"
	"log/slog"
)

type debugKey struct{}

func WithDebugging(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, debugKey{}, log)
}

func getDebugValue(ctx context.Context) *slog.Logger {
	val := ctx.Value(debugKey{})
	if val == nil {
		return nil
	}

	return val.(*slog.Logger)
}
