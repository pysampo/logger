package logger

import (
	"os"
	"path"
	"strings"

	"git.ipc/samatil3/logger/internal/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(serviceName string, config *Config) (func(), error) {

	// cores
	var (
		cores   = make([]zapcore.Core, 0, 2)
		encoder = zapcore.NewJSONEncoder(*getEncoderConfig(serviceName))
	)

	if config.ConsoleEnabled {
		consoleLevelEnabler, err := getLevelEnabler(config.ConsoleLevel)
		if err != nil {
			return nil, err
		}
		cores = append(cores, getCore(
			encoder, zapcore.Lock(os.Stderr), consoleLevelEnabler),
		)
	}
	if config.FileEnabled {
		fileLevelEnabler, err := getLevelEnabler(config.FileLevel)
		if err != nil {
			return nil, err
		}
		writeSyncer, err := newRollingFile(config)
		if err != nil {
			return nil, err
		}
		cores = append(cores, getCore(
			encoder, writeSyncer, fileLevelEnabler),
		)
	}

	var (
		core    = zapcore.NewTee(cores...)
		options []zap.Option
	)

	// options
	{
		hostname, _ := os.Hostname()
		options = append(options, zap.Fields(
			zap.String("ip", util.GetLocalIP()),
			zap.String("instance_id", hostname),
			zap.String("service_name", serviceName),
		))
		if config.Caller {
			options = append(options, zap.AddCaller(), zap.AddCallerSkip(1))
		}
	}

	// logger
	logger := zap.New(core, options...)

	return func() {
		logger.Sync()
	}, nil
}

func getLevelEnabler(levelStr string) (zapcore.LevelEnabler, error) {
	var (
		level zapcore.Level
	)
	if err := level.Set(strings.ToLower(levelStr)); err != nil {
		return nil, err
	}
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	}), nil
}

func getCore(encoder zapcore.Encoder, writeSyncer zapcore.WriteSyncer, levelEnabler zapcore.LevelEnabler) zapcore.Core {
	return zapcore.NewCore(encoder, writeSyncer, levelEnabler)
}

func getEncoderConfig(serviceName string) *zapcore.EncoderConfig {
	const (
		timeLayout    = "2006-01-02T15:04:05.000"
		keyTime       = "time"
		keyLevel      = "lvl"
		keyCaller     = "call"
		keyMessage    = "msg"
		keyStackTrace = "stacktrace"
	)
	return &zapcore.EncoderConfig{
		TimeKey:        keyTime,
		LevelKey:       keyLevel,
		NameKey:        serviceName,
		CallerKey:      keyCaller,
		MessageKey:     keyMessage,
		StacktraceKey:  keyStackTrace,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeLayout),
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newRollingFile(config *Config) (zapcore.WriteSyncer, error) {
	if err := os.MkdirAll(config.FileDirectory, 0744); err != nil {
		return nil, err
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.FileDirectory, config.Filename),
		MaxSize:    config.FileMaxSize,    // megabytes
		MaxAge:     config.FileMaxAge,     // days
		MaxBackups: config.FileMaxBackups, // files
		Compress:   config.FileCompress,
		LocalTime:  true,
	}), nil
}
