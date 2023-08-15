package logger

import (
	"github.com/yasin-wu/utils/logger/elasticsearch"
	"github.com/yasin-wu/utils/logger/stdout"
	"testing"
)

func TestLogger(t *testing.T) {
	defaultOutput := stdout.New("debug")
	esOutPut, err := elasticsearch.New("yasin", "error", &elasticsearch.ESConfig{
		Addresses: []string{"http://localhost:80"},
		Username:  "elastic",
		Password:  "yasinwu",
	}, elasticsearch.WithIndexType("day"))
	if err != nil {
		t.Fatal(err)
	}
	log, _ := New("yasin", defaultOutput, esOutPut)
	log.Debug("this is debug")
	log.Info("this is info")
	log.Error("this is error")
}
