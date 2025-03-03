package internal

import (
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Corer interface {
	Encoder() zapcore.Encoder
	WriteSyncer() zapcore.WriteSyncer
	AtomicLevel() zap.AtomicLevel
}

var DefaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "line",
	MessageKey:     "message",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}

// CallerEncoder Delete depth elements
func CallerEncoder(depth int) zapcore.CallerEncoder {
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
			depth = 0
		}
		for i := depth; i < len(files); i++ {
			temp = append(temp, files[i])
		}
		line := strings.Join(temp, "/") + ":" + strconv.Itoa(caller.Line)
		enc.AppendString(line)
	}
}

func AtomicLevel(level string) zap.AtomicLevel {
	logLevel := zapcore.InfoLevel
	switch strings.ToLower(level) {
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
