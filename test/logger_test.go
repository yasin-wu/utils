package test

import (
	"testing"

	"github.com/yasin-wu/utils/logger"
)

func TestLogger(t *testing.T) {
	log := logger.New(logger.WithJsonEncoder(false),
		logger.WithLevel("error"))
	log1 := log.SugaredLogger("test1")
	log2 := log.SugaredLogger("test2")
	log1.Info("info test1")
	log1.Error("error test1")
	log2.Info("info test2")
	log2.Error("error test2")
}
