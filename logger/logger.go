package logger

import (
	"errors"

	"go.uber.org/zap"
)

func New(serviceName string, options ...Option) (*zap.SugaredLogger, error) {
	if serviceName == "" {
		return nil, errors.New("service name must not be empty")
	}
	core := newCore(serviceName, options...)
	opts := []zap.Option{zap.AddCaller(), wrapCore(core)}
	if core.stacktrace {
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	}
	logger, err := zap.NewProduction(opts...)
	if err != nil {
		return nil, err
	}
	return logger.With(zap.String("service", serviceName)).Sugar(), nil
}
