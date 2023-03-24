package logger

import (
	"github.com/yasin-wu/utils/logger/output"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Core struct {
	serviceName string
	maxSize     int  //default:128,MB
	maxBackups  int  //default:30
	maxAge      int  //default:7,day
	depth       int  //default:0
	compress    bool //default:true
	stacktrace  bool //default:false
	stdout      bool //default:true
	outputs     []output.Output
}

var defaultCore = Core{
	maxSize:    128,
	maxBackups: 30,
	depth:      0,
	maxAge:     7,
	compress:   true,
	stdout:     true,
	stacktrace: false,
}

func newCore(serviceName string, options ...Option) Core {
	core := defaultCore
	core.serviceName = serviceName
	for _, f := range options {
		f(&core)
	}
	return core
}

func (c Core) newTee() zapcore.Core {
	var cores []zapcore.Core
	minLevel := c.minLevel()
	if len(c.outputs) == 0 || c.stdout {
		op := output.New()
		core := zapcore.NewCore(op.ConsoleEncoder(c.stacktrace, c.depth), c.writeStdout(), minLevel)
		cores = append(cores, core)
	}
	for _, op := range c.outputs {
		core := zapcore.NewCore(op.JSONEncoder(c.stacktrace, c.depth), c.writeSyncer(op), op.AtomicLevel())
		cores = append(cores, core)
	}
	return zapcore.NewTee(cores...)
}

// default:output to file,not output to stdout
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

// output to stdout
func (c Core) writeStdout() zapcore.WriteSyncer {
	var sync []zapcore.WriteSyncer
	sync = append(sync, zapcore.AddSync(os.Stdout))
	return zapcore.NewMultiWriteSyncer(sync...)
}

func (c Core) minLevel() zap.AtomicLevel {
	var minLevel = zapcore.InfoLevel
	for _, v := range c.outputs {
		if v.AtomicLevel().Level() < minLevel {
			minLevel = v.AtomicLevel().Level()
		}
	}
	return zap.NewAtomicLevelAt(minLevel)
}
