package logger

var defaultLogger = Logger{
	maxSize:    128,
	maxBackups: 30,
	maxAge:     7,
	compress:   true,
	dev:        true,
}

type Option func(logger *Logger)

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
 * @date: 2022/2/17 14:09
 * @params: outputs ...*Output
 * @return: Option
 * @description: 设置多个output,default:defaultOutput
 */
func WithOutputs(outputs ...Output) Option {
	return func(logger *Logger) {
		logger.outputs = append(logger.outputs, outputs...)
	}
}
