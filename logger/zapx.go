package logger

import (
	path2 "path"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func WrapCore(path, name, level string, stdout bool) zap.Option {
	output := NewOutput(WithLevel(level), WithStdout(stdout),
		WithFilename(path2.Join(path, name+".log")))
	errOutput := NewOutput(WithLevel("error"), WithStdout(false),
		WithFilename(path2.Join(path, name+"-error.log")))
	core := newCore(WithOutputs(output, errOutput))
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return core.newTee()
	})
}
