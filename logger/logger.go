package logger

import (
	"context"
	stdlog "log"

	"git.ipc/samatil3/logger/internal/core"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetLogger(logger *zap.Logger) {
	_logger = logger
}

func AttachLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, _key, logger)
}

func CloneWithLevel(ctx context.Context, level zapcore.Level) *zap.Logger {
	return from(ctx).WithOptions(
		core.WithLevel(level),
	)
}

func Error(ctx context.Context, message string, fields ...zap.Field) {
	var (
		logger = from(ctx)
		span   = getSpan(ctx)
	)
	if span != nil {
		spanLog(span, "error", message, castToOpentracing(fields...)...)
		collectSpanInfo(logger, span)
	}
	logger.Error(message, fields...)
}

func Warn(ctx context.Context, message string, fields ...zap.Field) {
	var (
		logger = from(ctx)
		span   = getSpan(ctx)
	)
	if span != nil {
		spanLog(span, "warn", message, castToOpentracing(fields...)...)
		collectSpanInfo(logger, span)
	}
	logger.Warn(message, fields...)
}

func Info(ctx context.Context, message string, fields ...zap.Field) {
	var (
		logger = from(ctx)
		span   = getSpan(ctx)
	)
	if span != nil {
		spanLog(span, "info", message, castToOpentracing(fields...)...)
		collectSpanInfo(logger, span)
	}
	logger.Info(message, fields...)
}

func Debug(ctx context.Context, message string, fields ...zap.Field) {
	var (
		logger = from(ctx)
		span   = getSpan(ctx)
	)
	if span != nil {
		spanLog(span, "debug", message, castToOpentracing(fields...)...)
		collectSpanInfo(logger, span)
	}
	logger.Debug(message, fields...)
}

func Fatal(ctx context.Context, message string, fields ...zap.Field) {
	var (
		logger = from(ctx)
		span   = getSpan(ctx)
	)
	if span != nil {
		spanLog(span, "fatal", message, castToOpentracing(fields...)...)
		collectSpanInfo(logger, span)
	}
	logger.Fatal(message, fields...)
}

func spanLog(span opentracing.Span, level string, message string, fields ...log.Field) {
	span.LogFields(append(fields, log.String("level", level), log.String("event", message))...)
}

func collectSpanInfo(logger *zap.Logger, span opentracing.Span) {
	if spanCtx, ok := span.Context().(jaeger.SpanContext); ok {
		logger.With(
			zap.String("span", spanCtx.String()), // trace-id, span-id, parent-id
		)
	}
}

func getSpan(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

type (
	ctxKey struct{}
)

var (
	_key    = ctxKey{}
	_logger *zap.Logger
)

func from(ctx context.Context) *zap.Logger {
	if loggerCtx, ok := ctx.Value(_key).(*zap.Logger); ok {
		return loggerCtx
	}
	return _logger
}

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		stdlog.Fatal(err)
	}
	_logger = logger
}
