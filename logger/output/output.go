package output

import (
	"io"
	"strings"
)

var defaultOutput = Output{
	Path:        "./log",
	Level:       "info",
	Stdout:      true,
	JSONEncoder: true,
}

type Output struct {
	Path        string
	Level       string
	Stdout      bool
	JSONEncoder bool
	Writer      []io.Writer
}

type Option func(output *Output)

func New(options ...Option) Output {
	output := defaultOutput
	for _, f := range options {
		f(&output)
	}
	if len(output.Writer) > 0 {
		output.JSONEncoder = true
	}
	return output
}

func WithPath(path string) Option {
	return func(output *Output) {
		if path != "" {
			output.Path = path
		}
	}
}

func WithLevel(level string) Option {
	return func(output *Output) {
		if level != "" {
			output.Level = strings.ToLower(level)
		}
	}
}

func WithStdout(stdout bool) Option {
	return func(output *Output) {
		output.Stdout = stdout
	}
}

func WithJSONEncoder(jsonEncoder bool) Option {
	return func(output *Output) {
		output.JSONEncoder = jsonEncoder
	}
}

func WithWriter(writer ...io.Writer) Option {
	return func(output *Output) {
		if len(writer) > 0 {
			output.Writer = append(output.Writer, writer...)
		}
	}
}
