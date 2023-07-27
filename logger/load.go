package logger

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
)

type Config struct {
	Logger Logger
}

var defaultReplacer = strings.NewReplacer(".", "_")

var v *viper.Viper

func MustLoad(serverName, configFile string) *Logger {
	var conf Config
	v = viper.New()
	v.SetConfigFile(configFile)
	v.SetEnvKeyReplacer(defaultReplacer)
	v.AutomaticEnv()
	defaultConfig(serverName)
	if err := v.ReadInConfig(); err != nil {
		return nil
	}
	unmarshal(&conf)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		unmarshal(&conf)
	})
	return &conf.Logger
}

func unmarshal(conf *Config) {
	if err := v.Unmarshal(&conf); err != nil {
		logx.Errorf("unmarshal logger config failed: %v", err)
		os.Exit(1)
	}
}

func defaultConfig(serverName string) {
	for k, val := range loggerDefault(serverName) {
		v.SetDefault(k, val)
	}
}

type DefaultConfig map[string]any

func loggerDefault(serverName string) DefaultConfig {
	prefix := "logger."
	logger := DefaultConfig{
		prefix + "service_name":          serverName,
		prefix + "mode":                  "console",
		prefix + "encoding":              "json",
		prefix + "path":                  "logs",
		prefix + "level":                 "info",
		prefix + "stacktrace":            false,
		prefix + "compress":              true,
		prefix + "keep_days":             7,
		prefix + "stack_cooldown_millis": 100,
		prefix + "max_backups":           30,
		prefix + "max_size":              16,
		prefix + "rotation":              "size",
	}
	return logger
}
