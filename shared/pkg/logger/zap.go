package logger

import (
	"context"

	"go.uber.org/zap"
)

type zapLogger struct {
	sugar *zap.SugaredLogger
}

// NewZapLogger creates and returns a Logger backed by a Zap production logger.
func NewZapLogger() (Logger, error) {
	base, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &zapLogger{sugar: base.Sugar()}, nil
}

// Errorw logs an error message with structured key-value pairs.
func (l *zapLogger) Errorw(_ context.Context, msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Infow logs an informational message with structured key-value pairs.
func (l *zapLogger) Infow(_ context.Context, msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warnw logs a warning message with structured key-value pairs.
func (l *zapLogger) Warnw(_ context.Context, msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}
