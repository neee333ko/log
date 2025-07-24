package log

import "context"

type key int

const (
	logContextKey key = iota
)

func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

func (logger *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, logger)
}

func FromContext(ctx context.Context) *Logger {
	if ctx != nil {
		if l := ctx.Value(logContextKey); l != nil {
			return l.(*Logger)
		}
	}

	return nil
}
