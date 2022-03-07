package logger

import (
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	maxSize    int  //default:128,MB
	maxBackups int  //default:30
	maxAge     int  //default:7,day
	compress   bool //default:true
	dev        bool //default:true
	outputs    []Output
	logger     *zap.Logger
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:11
 * @params: options ...Option
 * @return: *Logger
 * @description: New Logger
 */
func New(options ...Option) *Logger {
	logger := defaultLogger
	for _, f := range options {
		f(logger)
	}
	if len(logger.outputs) == 0 {
		logger.outputs = append(logger.outputs, defaultOutput)
	}
	var cores []zapcore.Core
	for _, output := range logger.outputs {
		cores = append(cores, logger.newCore(output))
	}
	core := zapcore.NewTee(cores...)
	zapOptions := []zap.Option{zap.AddCaller()}
	if logger.dev {
		zapOptions = append(zapOptions, zap.Development())
	}
	logger.logger = zap.New(core, zapOptions...)
	return logger
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:11
 * @params: service string
 * @return: *zap.SugaredLogger
 * @description: New SugaredLogger
 */
func (l *Logger) SugaredLogger(service string) *zap.SugaredLogger {
	return l.logger.With(zap.String("service", service)).Sugar()
}

func (l *Logger) newCore(output Output) zapcore.Core {
	return zapcore.NewCore(l.encoder(output), l.writeSyncer(output), l.atomicLevel(output))
}

func (l *Logger) encoder(output Output) zapcore.Encoder {
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

func (l *Logger) writeSyncer(output Output) zapcore.WriteSyncer {
	hook := &lumberjack.Logger{
		Filename:   output.filename,
		MaxSize:    l.maxSize,
		MaxBackups: l.maxBackups,
		MaxAge:     l.maxAge,
		Compress:   l.compress,
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

func (l *Logger) atomicLevel(output Output) zap.AtomicLevel {
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
