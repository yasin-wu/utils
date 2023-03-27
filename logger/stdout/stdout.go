package stdout

import (
	"github.com/yasin-wu/utils/logger/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type stdout struct {
	level      string
	timeLayout string
	depth      int
	stacktrace bool
}

var _ core.Corer = (*stdout)(nil)

type Option func(s *stdout)

func New(level string, options ...Option) core.Corer {
	corer := &stdout{
		level:      level,
		stacktrace: false,
		depth:      0,
		timeLayout: "2006-01-02 15:04:05",
	}
	for _, f := range options {
		f(corer)
	}
	return corer
}

func (s stdout) Encoder() zapcore.Encoder {
	encoderConfig := core.DefaultEncoderConfig
	encoderConfig.EncodeCaller = core.CallerEncoder(s.depth)
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(s.timeLayout)
	if s.stacktrace {
		encoderConfig.StacktraceKey = "stacktrace"
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (s stdout) WriteSyncer() zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
}

func (s stdout) AtomicLevel() zap.AtomicLevel {
	return core.AtomicLevel(s.level)
}

func WithTimelayout(timeLayout string) Option {
	return func(s *stdout) {
		if timeLayout == "" {
			s.timeLayout = timeLayout
		}
	}
}

func WithStacktrace(stacktrace bool) Option {
	return func(s *stdout) {
		s.stacktrace = stacktrace
	}
}

func WithDepth(depth int) Option {
	return func(s *stdout) {
		s.depth = depth
	}
}
