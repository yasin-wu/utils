package logger

import (
	"github.com/yasin-wu/utils/logger/output"
)

type Option func(core *Core)

func WithMaxSize(maxSize int) Option {
	return func(core *Core) {
		if maxSize > 0 {
			core.maxSize = maxSize
		}
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(core *Core) {
		if maxBackups > 0 {
			core.maxBackups = maxBackups
		}
	}
}

func WithMaxAge(maxAge int) Option {
	return func(core *Core) {
		if maxAge > 0 {
			core.maxAge = maxAge
		}
	}
}

func WithCompress(compress bool) Option {
	return func(core *Core) {
		core.compress = compress
	}
}

func WithOutputs(outputs ...output.Output) Option {
	return func(core *Core) {
		if len(outputs) > 0 {
			core.outputs = append(core.outputs, outputs...)
		}
	}
}

func WithStacktrace(stacktrace bool) Option {
	return func(core *Core) {
		core.stacktrace = stacktrace
	}
}
