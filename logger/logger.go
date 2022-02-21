package logger

import (
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	conf   *config
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:11
 * @params: options ...Option
 * @return: *Logger
 * @description: New Logger
 */
func New(options ...Option) *Logger {
	conf := defaultConfig
	for _, f := range options {
		f(conf)
	}
	logger := &Logger{conf: conf}
	core := zapcore.NewCore(logger.encoder(), logger.writeSyncer(), logger.level())
	zapOptions := []zap.Option{zap.AddCaller()}
	if conf.dev {
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

func (l *Logger) encoder() zapcore.Encoder {
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
	if l.conf.jsonEncoder {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	return encoder
}

func (l *Logger) writeSyncer() zapcore.WriteSyncer {
	hook := &lumberjack.Logger{
		Filename:   l.conf.filename,
		MaxSize:    l.conf.maxSize,
		MaxBackups: l.conf.maxBackups,
		MaxAge:     l.conf.maxAge,
		Compress:   l.conf.compress,
	}
	var sync []zapcore.WriteSyncer
	sync = append(sync, zapcore.AddSync(hook))
	if l.conf.stdout {
		sync = append(sync, zapcore.AddSync(os.Stdout))
	}
	for _, w := range l.conf.writer {
		sync = append(sync, zapcore.AddSync(w))
	}
	return zapcore.NewMultiWriteSyncer(sync...)
}

func (l *Logger) level() zap.AtomicLevel {
	logLevel := zapcore.InfoLevel
	switch strings.ToLower(l.conf.level) {
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
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(logLevel)
	return atomicLevel
}
