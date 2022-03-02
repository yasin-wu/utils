package logger

import "io"

var defaultLogger = &Logger{
	filename:    "./log/main.log",
	level:       "info",
	maxSize:     128,
	maxBackups:  30,
	maxAge:      7,
	compress:    true,
	dev:         true,
	stdout:      true,
	jsonEncoder: true,
}

type Option func(logger *Logger)

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:05
 * @params: filename string
 * @return: Option
 * @description: 日志文件路径,default:./log/main.log
 */
func WithFilename(filename string) Option {
	return func(logger *Logger) {
		if filename != "" {
			logger.filename = filename
		}
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:05
 * @params: level string;debug,info,warn,error,dpanic,panic,fatal
 * @return: Option
 * @description: 日志输出级别,default:info
 */
func WithLevel(level string) Option {
	return func(logger *Logger) {
		if level != "" {
			logger.level = level
		}
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:06
 * @params: maxSize int
 * @return: Option
 * @description: 每个日志文件大小,单位:MB,default:128
 */
func WithMaxSize(maxSize int) Option {
	return func(logger *Logger) {
		logger.maxSize = maxSize
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:07
 * @params: maxBackups int
 * @return: Option
 * @description: 日志最大备份数,default:30
 */
func WithMaxBackups(maxBackups int) Option {
	return func(logger *Logger) {
		if maxBackups > 0 {
			logger.maxBackups = maxBackups
		}
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:08
 * @params: maxAge int
 * @return: Option
 * @description: 日志保存最大天数,default:7
 */
func WithMaxAge(maxAge int) Option {
	return func(logger *Logger) {
		if maxAge > 0 {
			logger.maxAge = maxAge
		}
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:09
 * @params: compress bool
 * @return: Option
 * @description: 日志是否压缩,default:true
 */
func WithCompress(compress bool) Option {
	return func(logger *Logger) {
		logger.compress = compress
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:09
 * @params: dev bool
 * @return: Option
 * @description: 日志是否dev,default:true
 */
func WithDev(dev bool) Option {
	return func(logger *Logger) {
		logger.dev = dev
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:10
 * @params: stdout bool
 * @return: Option
 * @description: 日志是否控制台输出,default:true
 */
func WithStdout(stdout bool) Option {
	return func(logger *Logger) {
		logger.stdout = stdout
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 16:36
 * @params: jsonEncoder bool
 * @return: Option
 * @description: 设置默认encoder,default:true
 */
func WithJsonEncoder(jsonEncoder bool) Option {
	return func(logger *Logger) {
		logger.jsonEncoder = jsonEncoder
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/21 10:02
 * @params: w ...io.Writer
 * @return: Option
 * @description: 设置io writer,default:file
 */
func WithWriter(w ...io.Writer) Option {
	return func(logger *Logger) {
		logger.writer = w
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/21 10:02
 * @params: errorFile string
 * @return: Option
 * @description: 设置error file,default:filename
 */
func WithErrorFile(errorFile string) Option {
	return func(logger *Logger) {
		logger.errorfile = errorFile
	}
}
