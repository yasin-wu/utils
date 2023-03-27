package file

import (
	"fmt"
	"github.com/yasin-wu/utils/logger/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
)

type file struct {
	level      string
	path       string
	timeLayout string
	maxSize    int
	maxBackups int
	maxAge     int
	depth      int
	compress   bool
	stacktrace bool
}

var _ core.Corer = (*file)(nil)

var serviceName string

type Option func(f *file)

func New(level string, options ...Option) core.Corer {
	corer := &file{
		level:      level,
		path:       "./log",
		stacktrace: false,
		depth:      0,
		maxSize:    16,
		maxBackups: 30,
		maxAge:     7,
		compress:   true,
		timeLayout: "2006-01-02 15:04:05",
	}
	for _, f := range options {
		f(corer)
	}
	return corer
}

func (f file) Encoder() zapcore.Encoder {
	encoderConfig := core.DefaultEncoderConfig
	encoderConfig.EncodeCaller = core.CallerEncoder(f.depth)
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(f.timeLayout)
	if f.stacktrace {
		encoderConfig.StacktraceKey = "stacktrace"
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func (f file) WriteSyncer() zapcore.WriteSyncer {
	filename := path.Join(f.path, fmt.Sprintf("%s-%s.log", serviceName, f.level))
	hook := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    f.maxSize,
		MaxBackups: f.maxBackups,
		MaxAge:     f.maxAge,
		Compress:   f.compress,
		LocalTime:  true,
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook))
}

func (f file) AtomicLevel() zap.AtomicLevel {
	return core.AtomicLevel(f.level)
}

func WithPath(path string) Option {
	return func(f *file) {
		if path == "" {
			f.path = path
		}
	}
}

func WithTimelayout(timeLayout string) Option {
	return func(f *file) {
		if timeLayout == "" {
			f.timeLayout = timeLayout
		}
	}
}

func WithMaxSize(maxSize int) Option {
	return func(f *file) {
		if maxSize > 0 {
			f.maxSize = maxSize
		}
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(f *file) {
		if maxBackups > 0 {
			f.maxBackups = maxBackups
		}
	}
}

func WithMaxAge(maxAge int) Option {
	return func(f *file) {
		if maxAge > 0 {
			f.maxAge = maxAge
		}
	}
}

func WithDepth(depth int) Option {
	return func(f *file) {
		f.depth = depth
	}
}

func WithStacktrace(stacktrace bool) Option {
	return func(f *file) {
		f.stacktrace = stacktrace
	}
}

func WithCompress(compress bool) Option {
	return func(f *file) {
		f.compress = compress
	}
}

func SetServiceName(name string) {
	serviceName = name
}
