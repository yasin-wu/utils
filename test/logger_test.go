package test

import (
	"strings"
	"testing"

	"github.com/yasin-wu/utils/logger/writer"

	"github.com/yasin-wu/utils/logger"
)

func TestLogger(t *testing.T) {
	esConf := &writer.ESConfig{
		Addresses: strings.Split("http://47.108.155.25:9200", ","),
		Username:  "elastic",
		Password:  "yasinwu",
	}
	esWriter, _ := writer.NewESWriter("yasin_logs", esConf)
	errOutput := logger.NewOutput(logger.WithPath("./log"),
		logger.WithLevel("error"), logger.WithWriter(esWriter))
	defaultOutput := logger.NewOutput()
	log1, _ := logger.New("yasin", logger.WithOutputs(defaultOutput, errOutput))
	log2, _ := logger.New("yasin", logger.WithOutputs(defaultOutput, errOutput))
	log1.Info("info test1")
	log1.Error("error test1")
	log2.Info("info test2")
	log2.Error("error test2")
}
