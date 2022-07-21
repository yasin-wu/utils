package logger

import (
	"errors"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

type WrapCore struct {
	ServiceName string
	Path        string
	Level       string
	Stdout      bool
}

func (w *WrapCore) New(options ...Option) (zap.Option, error) {
	if w == nil {
		return nil, errors.New("wrap core is nil")
	}
	output := NewOutput(WithLevel(w.Level), WithStdout(w.Stdout),
		WithPath(w.Path))
	errOutput := NewOutput(WithLevel("error"), WithStdout(false),
		WithPath(w.Path))
	core := newCore(w.ServiceName, append(options, WithOutputs(output, errOutput))...)
	return wrapCore(core), nil
}

func wrapCore(core Core) zap.Option {
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return core.newTee()
	})
}
