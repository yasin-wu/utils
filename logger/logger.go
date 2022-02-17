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
	core := newCore(conf)
	zapOptions := []zap.Option{zap.AddCaller()}
	if conf.dev {
		zapOptions = append(zapOptions, zap.Development())
	}

	return &Logger{logger: zap.New(core, zapOptions...)}
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

func newCore(conf *config) zapcore.Core {
	//encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if conf.jsonEncoder {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	//log output
	hook := &lumberjack.Logger{
		Filename:   conf.filename,
		MaxSize:    conf.maxSize,
		MaxBackups: conf.maxBackups,
		MaxAge:     conf.maxAge,
		Compress:   conf.compress,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook))
	if conf.stdout {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook))
	}

	//log level
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level(conf.level))
	return zapcore.NewCore(encoder, writeSyncer, atomicLevel)
}

func level(level string) zapcore.Level {
	logLevel := zapcore.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "panic":
		logLevel = zapcore.PanicLevel
	case "fatal":
		logLevel = zapcore.FatalLevel
	}
	return logLevel
}
