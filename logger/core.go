package logger

import (
	"fmt"
	"github.com/yasin-wu/utils/logger/output"
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
	stacktrace  bool //default:false
	outputs     []output.Output
}

var defaultCore = Core{
	maxSize:    128,
	maxBackups: 30,
	maxAge:     7,
	compress:   true,
	stacktrace: false,
}

func newCore(serviceName string, options ...Option) Core {
	core := defaultCore
	core.serviceName = serviceName
	for _, f := range options {
		f(&core)
	}
	if len(core.outputs) == 0 {
		core.outputs = append(core.outputs, output.New())
	}
	return core
}

func (c Core) newTee() zapcore.Core {
	var cores []zapcore.Core
	for _, op := range c.outputs {
		cores = append(cores, zapcore.NewCore(c.encoder(op), c.writeSyncer(op), c.atomicLevel(op)))
	}
	return zapcore.NewTee(cores...)
}

func (c Core) encoder(op output.Output) zapcore.Encoder {
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
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	if c.stacktrace {
		encoderConfig.StacktraceKey = "stacktrace"
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if op.JSONEncoder {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	return encoder
}

func (c Core) writeSyncer(op output.Output) zapcore.WriteSyncer {
	hook := &lumberjack.Logger{
		Filename:   path.Join(op.Path, fmt.Sprintf("%s-%s.log", c.serviceName, op.Level)),
		MaxSize:    c.maxSize,
		MaxBackups: c.maxBackups,
		MaxAge:     c.maxAge,
		Compress:   c.compress,
		LocalTime:  true,
	}
	var sync []zapcore.WriteSyncer //nolint:prealloc
	sync = append(sync, zapcore.AddSync(hook))
	if op.Stdout {
		sync = append(sync, zapcore.AddSync(os.Stdout))
	}
	for _, w := range op.Writer {
		sync = append(sync, zapcore.AddSync(w))
	}
	return zapcore.NewMultiWriteSyncer(sync...)
}

func (c Core) atomicLevel(op output.Output) zap.AtomicLevel {
	logLevel := zapcore.InfoLevel
	switch strings.ToLower(op.Level) {
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
