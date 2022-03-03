package logger

import (
	"io"
)

type Output struct {
	filename    string
	level       string
	stdout      bool
	jsonEncoder bool
	writer      []io.Writer
}

var defaultOutput = Output{
	filename:    "./log/main.log",
	level:       "info",
	stdout:      true,
	jsonEncoder: true,
}

func NewOutput(filename, level string, stdout, jsonEncoder bool, writers ...io.Writer) Output {
	output := Output{
		filename:    filename,
		level:       level,
		stdout:      stdout,
		jsonEncoder: jsonEncoder,
		writer:      writers,
	}

	if len(output.writer) > 0 {
		output.jsonEncoder = true
	}

	return output
}
