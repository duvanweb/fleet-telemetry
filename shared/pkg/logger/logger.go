package logger

import "context"

// Logger defines the structured logging interface used across all services.
type Logger interface {
	Infow(ctx context.Context, msg string, keysAndValues ...any)
	Warnw(ctx context.Context, msg string, keysAndValues ...any)
	Errorw(ctx context.Context, msg string, keysAndValues ...any)
}
