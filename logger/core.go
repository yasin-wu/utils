package logger

import (
	"github.com/yasin-wu/utils/logger/internal"
	"github.com/yasin-wu/utils/logger/stdout"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func wrapCore(outputs ...internal.Corer) zap.Option {
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		var cores []zapcore.Core
		if len(outputs) == 0 {
			op := stdout.New("info")
			cores = append(cores, zapcore.NewCore(op.Encoder(), op.WriteSyncer(), op.AtomicLevel()))
		}
		for _, op := range outputs {
			cores = append(cores, zapcore.NewCore(op.Encoder(), op.WriteSyncer(), op.AtomicLevel()))
		}
		return zapcore.NewTee(cores...)
	})
}
