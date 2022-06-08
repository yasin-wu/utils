package logger

import (
	"io"
)

var defaultOutput = Output{
	path:        "./log",
	level:       "info",
	stdout:      true,
	jsonEncoder: true,
}

type Output struct {
	path        string
	level       string
	stdout      bool
	jsonEncoder bool
	writer      []io.Writer
}

func NewOutput(options ...OutputOption) Output {
	output := defaultOutput
	for _, f := range options {
		f(&output)
	}
	if len(output.writer) > 0 {
		output.jsonEncoder = true
	}
	return output
}
