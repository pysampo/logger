package logger

import (
	"fmt"
	"math"
	"time"

	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func castToOpentracing(fields ...zap.Field) []log.Field {
	ret := make([]log.Field, 0, len(fields))
	for _, field := range fields {
		ret = append(ret, castFieldToOpentracing(field))
	}
	return ret
}

func castFieldToOpentracing(field zap.Field) log.Field {
	switch field.Type {
	case zapcore.BoolType:
		val := false
		if field.Integer >= 1 {
			val = true
		}
		return log.Bool(field.Key, val)
	case zapcore.Float32Type:
		return log.Float32(field.Key, math.Float32frombits(uint32(field.Integer)))
	case zapcore.Float64Type:
		return log.Float64(field.Key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.Int64Type:
		return log.Int64(field.Key, int64(field.Integer))
	case zapcore.Int32Type:
		return log.Int32(field.Key, int32(field.Integer))
	case zapcore.StringType:
		return log.String(field.Key, field.String)
	case zapcore.StringerType:
		return log.String(field.Key, field.Interface.(fmt.Stringer).String())
	case zapcore.Uint64Type:
		return log.Uint64(field.Key, uint64(field.Integer))
	case zapcore.Uint32Type:
		return log.Uint32(field.Key, uint32(field.Integer))
	case zapcore.DurationType:
		return log.String(field.Key, time.Duration(field.Integer).String())
	case zapcore.ErrorType:
		return log.Error(field.Interface.(error))
	default:
		return log.Object(field.Key, field.Interface)
	}
}
