package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(serviceName string, options ...Option) (*zap.SugaredLogger, error) {
	if serviceName == "" {
		return nil, errors.New("service name must not be empty")
	}
	core := newCore(serviceName, options...)
	logger, err := zap.NewProduction(zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return core.newTee()
		}))
	if err != nil {
		return nil, err
	}
	return logger.With(zap.String("service", serviceName)).Sugar(), nil
}
