package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log *zap.SugaredLogger
}

// NewLogger creates a production zap logger wrapped behind Logger.
func NewLogger() (Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	base, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{log: base.Sugar()}, nil
}

// Infow logs an info message with structured key/value pairs.
func (z *zapLogger) Infow(_ context.Context, msg string, keysAndValues ...any) {
	z.log.Infow(msg, keysAndValues...)
}

// Warnw logs a warning message with structured key/value pairs.
func (z *zapLogger) Warnw(_ context.Context, msg string, keysAndValues ...any) {
	z.log.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message with structured key/value pairs.
func (z *zapLogger) Errorw(_ context.Context, msg string, keysAndValues ...any) {
	z.log.Errorw(msg, keysAndValues...)
}
