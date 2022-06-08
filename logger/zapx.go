package logger

import (
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func WrapCore(serviceName, path, level string, stdout bool) zap.Option {
	output := NewOutput(WithLevel(level), WithStdout(stdout),
		WithPath(path))
	errOutput := NewOutput(WithLevel("error"), WithStdout(false),
		WithPath(path))
	core := newCore(serviceName, WithOutputs(output, errOutput))
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return core.newTee()
	})
}
