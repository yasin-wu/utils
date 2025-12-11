package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Errorf(template string, args ...any)
	Infof(template string, args ...any)
}

var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct{}

func NewDefaultLogger() Logger {
	return &defaultLogger{}
}
func (d *defaultLogger) Errorf(template string, args ...any) {
	fmt.Printf("[\033[31mERROR\033[0m] - %s - [\033[31m%s\033[0m]\n",
		time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(template, args...))
}

func (d *defaultLogger) Infof(template string, args ...any) {
	fmt.Printf("[INFO] - %s - %s\n",
		time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(template, args...))
}
