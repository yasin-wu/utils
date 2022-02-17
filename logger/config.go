package logger

type config struct {
	filename    string
	level       string
	maxSize     int
	maxBackups  int
	maxAge      int
	compress    bool
	dev         bool
	stdout      bool
	jsonEncoder bool
}

var defaultConfig = &config{
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

type Option func(config *config)

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:05
 * @params: filename string
 * @return: Option
 * @description: 日志文件路径,default:./log/main.log
 */
func WithFilename(filename string) Option {
	return func(config *config) {
		if filename != "" {
			config.filename = filename
		}
	}
}

/**
 * @author: yasinWu
 * @date: 2022/2/17 14:05
 * @params: level string
 * @return: Option
 * @description: 日志输出级别,default:info
 */
func WithLevel(level string) Option {
	return func(config *config) {
		if level != "" {
			config.level = level
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
	return func(config *config) {
		config.maxSize = maxSize
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
	return func(config *config) {
		if maxBackups > 0 {
			config.maxBackups = maxBackups
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
	return func(config *config) {
		if maxAge > 0 {
			config.maxAge = maxAge
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
	return func(config *config) {
		config.compress = compress
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
	return func(config *config) {
		config.dev = dev
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
	return func(config *config) {
		config.stdout = stdout
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
	return func(config *config) {
		config.jsonEncoder = jsonEncoder
	}
}
