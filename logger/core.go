package logger

import (
	"github.com/yasin-wu/utils/logger/output"
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
	depth       int  //default:full
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
		cores = append(cores, zapcore.NewCore(op.Encoder(c.stacktrace, c.depth), c.writeSyncer(op), op.AtomicLevel()))
	}
	return zapcore.NewTee(cores...)
}

func (c Core) writeSyncer(op output.Output) zapcore.WriteSyncer {
	hook := &lumberjack.Logger{
		Filename:   op.Filename(c.serviceName),
		MaxSize:    c.maxSize,
		MaxBackups: c.maxBackups,
		MaxAge:     c.maxAge,
		Compress:   c.compress,
		LocalTime:  true,
	}
	var sync []zapcore.WriteSyncer
	sync = append(sync, zapcore.AddSync(hook))
	if len(op.WriteSyncer()) > 0 {
		sync = append(sync, op.WriteSyncer()...)
	}
	return zapcore.NewMultiWriteSyncer(sync...)
}
