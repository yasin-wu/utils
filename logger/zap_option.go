package logger

import (
	"github.com/yasin-wu/utils/logger/core"
	"go.uber.org/zap"
)

// New default:stdout info
func NewZapOption(outputs ...core.Corer) (zap.Option, error) {
	return wrapCore(outputs...), nil
}
