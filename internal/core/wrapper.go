package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func WithLevel(level zapcore.Level) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &coreWrapper{core, level}
	})
}

type (
	coreWrapper struct {
		zapcore.Core
		level zapcore.Level
	}
)

func (c *coreWrapper) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

func (c *coreWrapper) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *coreWrapper) With(fields []zapcore.Field) zapcore.Core {
	return &coreWrapper{
		c.Core.With(fields),
		c.level,
	}
}
