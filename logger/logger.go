package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func New(dev bool, options ...Option) Logger {
	core := NewCore(options...)
	zapOptions := []zap.Option{zap.AddCaller()}
	if dev {
		zapOptions = append(zapOptions, zap.Development())
	}
	return Logger{logger: zap.New(core.newTee(), zapOptions...)}
}

func (l Logger) SugaredLogger(service string) *zap.SugaredLogger {
	return l.logger.With(zap.String("service", service)).Sugar()
}
