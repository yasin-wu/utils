package logger

import (
	"errors"
	"github.com/yasin-wu/utils/logger/core"
	"github.com/yasin-wu/utils/logger/file"
	"go.uber.org/zap"
)

// New default:stdout info
func New(serviceName string, outputs ...core.Corer) (*zap.SugaredLogger, error) {
	if serviceName == "" {
		return nil, errors.New("service name must not be empty")
	}
	file.SetServiceName(serviceName)
	opts := []zap.Option{zap.AddCaller(), wrapCore(outputs...)}
	logger, err := zap.NewProduction(opts...)
	if err != nil {
		return nil, err
	}
	return logger.With(zap.String("service", serviceName)).Sugar(), nil
}
