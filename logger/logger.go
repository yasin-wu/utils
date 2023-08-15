package logger

import (
	"errors"
	"github.com/yasin-wu/utils/logger/internal"
	"go.uber.org/zap"
)

// New default:stdout info
func New(serviceName string, outputs ...internal.Corer) (*zap.SugaredLogger, error) {
	if serviceName == "" {
		return nil, errors.New("service name must not be empty")
	}
	opts := []zap.Option{zap.AddCaller(), wrapCore(outputs...)}
	logger, err := zap.NewProduction(opts...)
	if err != nil {
		return nil, err
	}
	return logger.With(zap.String("service", serviceName)).Sugar(), nil
}

// NewZapOption default:stdout info
func NewZapOption(outputs ...internal.Corer) (zap.Option, error) {
	return wrapCore(outputs...), nil
}
