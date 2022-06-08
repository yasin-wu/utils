package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func New(options ...Option) *Logger {
	core := newCore(options...)
	log, err := zap.NewProduction(zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return core.newTee()
		}))
	if err != nil {
		return nil
	}
	return &Logger{logger: log}
}

func (l Logger) SugaredLogger(service string) *zap.SugaredLogger {
	return l.logger.With(zap.String("service", service)).Sugar()
}
