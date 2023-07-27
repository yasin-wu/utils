package logger

import (
	"errors"
	"github.com/yasin-wu/utils/logger/core"
	"github.com/yasin-wu/utils/logger/file"
	"github.com/yasin-wu/utils/logger/stdout"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/zero-contrib/logx/zapx"
	"go.uber.org/zap"
)

const (
	consoleStdout = "console"
	errorLevel    = "error"
)

type Logger struct {
	ServiceName string `mapstructure:"service_name"` //服务名称
	Mode        string `mapstructure:"mode"`         //日志打印方式,console,file
	Encoding    string `mapstructure:"encoding"`     //日志格式,json,plain
	TimeFormat  string `mapstructure:"time_format"`  //日期格式
	Path        string `mapstructure:"path"`         //日志目录
	Level       string `mapstructure:"level"`        //日志等级
	Stacktrace  bool   `mapstructure:"stacktrace"`   //日志的stacktrace
	Compress    bool   `mapstructure:"compress"`     //是否压缩
	Stat        bool   `mapstructure:"stat"`         //是否统计
	KeepDays    int    `mapstructure:"keep_days"`    //保留天数
	MaxBackups  int    `mapstructure:"max_backups"`  //日志最大备份数
	MaxSize     int    `mapstructure:"max_size"`     //日志切割限制
	Rotation    string `mapstructure:"rotation"`     //切割方式
}

type (
	Option func(options *logxOptions)

	logxOptions struct {
		depth   int
		outputs []core.Corer
	}
)

// NewWriter default:output to error file and stdout,
func (l *Logger) NewWriter(options ...Option) (logx.Writer, error) {
	if l.ServiceName == "" {
		return nil, errors.New("service name is required")
	}
	opt, err := l.logxOptions(l.Stacktrace, options...)
	if err != nil {
		return nil, err
	}
	cores, err := NewZapOption(l.ServiceName, opt.outputs...)
	if err != nil {
		return nil, err
	}
	return zapx.NewZapWriter(cores, zap.Fields(zap.String("service", l.ServiceName)))
}

// NewLogger default:output to error file and stdout,
func (l *Logger) NewLogger(options ...Option) (*zap.SugaredLogger, error) {
	if l.ServiceName == "" {
		return nil, errors.New("service name is required")
	}
	opt, err := l.logxOptions(l.Stacktrace, options...)
	if err != nil {
		return nil, err
	}
	return New(l.ServiceName, opt.outputs...)
}

func (l *Logger) logxOptions(stacktrace bool, options ...Option) (*logxOptions, error) {
	if l == nil {
		return nil, errors.New("config is nil")
	}
	opt := &logxOptions{
		depth: 0,
	}
	for _, f := range options {
		f(opt)
	}
	fileOptions := []file.Option{
		file.WithStacktrace(stacktrace),
		file.WithDepth(opt.depth),
		file.WithMaxSize(l.MaxSize),
		file.WithPath(l.Path),
		file.WithMaxBackups(l.MaxBackups),
		file.WithMaxAge(l.KeepDays),
		file.WithCompress(l.Compress),
	}
	var outputs []core.Corer
	errFileOP := file.New(errorLevel, fileOptions...)
	if l.Mode == consoleStdout {
		stdoutOP := stdout.New(l.Level, stdout.WithStacktrace(stacktrace), stdout.WithDepth(opt.depth))
		outputs = append(outputs, stdoutOP)
	}
	if l.Level != errorLevel {
		op := file.New(l.Level, fileOptions...)
		outputs = append(outputs, op)
	}
	outputs = append(outputs, errFileOP)
	opt.outputs = outputs
	return opt, nil
}

func WithDepth(depth int) Option {
	return func(options *logxOptions) {
		options.depth = depth
	}
}
