package test

import (
	"github.com/yasin-wu/utils/logger/file"
	"github.com/yasin-wu/utils/logger/stdout"
	"testing"

	"github.com/yasin-wu/utils/logger"
)

func TestLogger(t *testing.T) {
	defaultOutput := stdout.New("debug")
	fileOutput := file.New("info")
	fileErrOutput := file.New("error")
	log, _ := logger.New("yasin", defaultOutput, fileOutput, fileErrOutput)
	log.Debug("this is debug")
	log.Info("this is info")
	log.Error("this is error")
}
