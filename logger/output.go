package logger

import (
	"io"
)

type Output struct {
	config *Config
	writer []io.Writer
}

type Config struct {
	Topics      []string
	Filename    string
	Level       string
	Stdout      bool
	JsonEncoder bool
}

var defaultConfig = &Config{
	Filename:    "./log/main.log",
	Level:       "info",
	Stdout:      true,
	JsonEncoder: true,
}

var defaultOutput = Output{
	config: defaultConfig,
}

func NewOutput(config *Config, writers ...io.Writer) Output {
	output := Output{
		config: config,
		writer: writers,
	}

	if len(output.writer) > 0 {
		output.config.JsonEncoder = true
	}

	return output
}

func NewDefaultOutput() Output {
	return defaultOutput
}
