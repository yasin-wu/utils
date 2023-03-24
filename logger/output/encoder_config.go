package output

import (
	"go.uber.org/zap/zapcore"
)

var defaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "line",
	MessageKey:     "message",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}
