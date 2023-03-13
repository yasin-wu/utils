package test

import (
	"github.com/yasin-wu/utils/logger/output"
	"testing"

	"github.com/yasin-wu/utils/logger"
)

func TestLogger(t *testing.T) {
	errOutput := output.New(
		output.WithPath("./log"),
		output.WithLevel("error"),
		output.WithJSONEncoder(false),
	)
	defaultOutput := output.New()
	log1, _ := logger.New("yasin", logger.WithOutputs(defaultOutput, errOutput))
	log2, _ := logger.New("yasin", logger.WithOutputs(defaultOutput, errOutput))
	log1.Info("info test1")
	log1.Error("error test1")
	log2.Info("info test2")
	log2.Error("error test2")
}
