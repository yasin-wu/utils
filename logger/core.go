package logger

import (
	"fmt"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Core struct {
	serviceName string
	maxSize     int  //default:128,MB
	maxBackups  int  //default:30
	maxAge      int  //default:7,day
	compress    bool //default:true
	outputs     []Output
}

var defaultCore = Core{
	maxSize:    128,
	maxBackups: 30,
	maxAge:     7,
	compress:   true,
}

func newCore(serviceName string, options ...Option) Core {
	core := defaultCore
	core.serviceName = serviceName
	for _, f := range options {
		f(&core)
	}
	if len(core.outputs) == 0 {
		core.outputs = append(core.outputs, defaultOutput)
	}
	return core
}

func (c Core) newTee() zapcore.Core {
	var cores []zapcore.Core
	for _, output := range c.outputs {
		cores = append(cores, zapcore.NewCore(c.encoder(output), c.writeSyncer(output), c.atomicLevel(output)))
	}
	return zapcore.NewTee(cores...)
}

func (c Core) encoder(output Output) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if output.jsonEncoder {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	return encoder
}

func (c Core) writeSyncer(output Output) zapcore.WriteSyncer {
	hook := &lumberjack.Logger{
		Filename:   path.Join(output.path, fmt.Sprintf("%s-%s.log", c.serviceName, output.level)),
		MaxSize:    c.maxSize,
		MaxBackups: c.maxBackups,
		MaxAge:     c.maxAge,
		Compress:   c.compress,
	}
	var sync []zapcore.WriteSyncer
	sync = append(sync, zapcore.AddSync(hook))
	if output.stdout {
		sync = append(sync, zapcore.AddSync(os.Stdout))
	}
	for _, w := range output.writer {
		sync = append(sync, zapcore.AddSync(w))
	}
	return zapcore.NewMultiWriteSyncer(sync...)
}

func (c Core) atomicLevel(output Output) zap.AtomicLevel {
	logLevel := zapcore.InfoLevel
	switch strings.ToLower(output.level) {
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
