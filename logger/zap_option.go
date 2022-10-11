package logger

import (
	"errors"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func NewZapOption(serviceName string, options ...Option) (zap.Option, error) {
	if serviceName == "" {
		return nil, errors.New("serviceName is empty")
	}
	core := newCore(serviceName, options...)
	return wrapCore(core), nil
}

func wrapCore(core Core) zap.Option {
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return core.newTee()
	})
}
