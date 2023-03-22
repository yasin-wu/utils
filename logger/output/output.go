package output

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
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

type Option func(output *Output)

func New(options ...Option) Output {
	output := defaultOutput
	for _, f := range options {
		f(&output)
	}
	return output
}

func (op Output) Filename(serviceName string) string {
	return path.Join(op.path, fmt.Sprintf("%s-%s.log", serviceName, op.level))
}

func (op Output) WriteSyncer() []zapcore.WriteSyncer {
	var sync []zapcore.WriteSyncer
	if op.stdout {
		sync = append(sync, zapcore.AddSync(os.Stdout))
	}
	for _, w := range op.writer {
		sync = append(sync, zapcore.AddSync(w))
	}
	return sync
}

func (op Output) Encoder(stacktrace bool, depth int) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   callerEncoder(depth),
		EncodeName:     zapcore.FullNameEncoder,
	}
	if stacktrace {
		encoderConfig.StacktraceKey = "stacktrace"
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if op.jsonEncoder {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	return encoder
}

func (op Output) AtomicLevel() zap.AtomicLevel {
	logLevel := zapcore.InfoLevel
	switch strings.ToLower(op.level) {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "dpanic":
		logLevel = zapcore.DPanicLevel
	case "panic":
		logLevel = zapcore.PanicLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	}
	return zap.NewAtomicLevelAt(logLevel)
}

func WithPath(path string) Option {
	return func(output *Output) {
		if path != "" {
			output.path = path
		}
	}
}

func WithLevel(level string) Option {
	return func(output *Output) {
		if level != "" {
			output.level = strings.ToLower(level)
		}
	}
}

func WithStdout(stdout bool) Option {
	return func(output *Output) {
		output.stdout = stdout
	}
}

func WithJSONEncoder(jsonEncoder bool) Option {
	return func(output *Output) {
		output.jsonEncoder = jsonEncoder
	}
}

func WithWriter(writer ...io.Writer) Option {
	return func(output *Output) {
		if len(writer) > 0 {
			output.writer = append(output.writer, writer...)
		}
	}
}

func callerEncoder(depth int) zapcore.CallerEncoder {
	if depth == 0 {
		return zapcore.FullCallerEncoder
	}
	if depth == -1 {
		return zapcore.ShortCallerEncoder
	}
	return func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		var temp []string
		files := strings.Split(caller.File, "/")
		if depth > len(files) {
			depth = len(files)
		}
		for i := depth; i > 0; i-- {
			temp = append(temp, files[len(files)-i])
		}
		line := strings.Join(temp, "/") + ":" + strconv.Itoa(caller.Line)
		enc.AppendString(line)
	}
}
