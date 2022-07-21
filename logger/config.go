package logger

import (
	"io"
	"strings"
)

type Option func(core *Core)

type OutputOption func(output *Output)

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

func WithOutputs(outputs ...Output) Option {
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

func WithPath(path string) OutputOption {
	return func(output *Output) {
		if path != "" {
			output.path = path
		}
	}
}

func WithLevel(level string) OutputOption {
	return func(output *Output) {
		if level != "" {
			output.level = strings.ToLower(level)
		}
	}
}

func WithStdout(stdout bool) OutputOption {
	return func(output *Output) {
		output.stdout = stdout
	}
}

func WithJsonEncoder(jsonEncoder bool) OutputOption {
	return func(output *Output) {
		output.jsonEncoder = jsonEncoder
	}
}

func WithWriter(writer ...io.Writer) OutputOption {
	return func(output *Output) {
		if len(writer) > 0 {
			output.writer = append(output.writer, writer...)
		}
	}
}
