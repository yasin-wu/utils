package test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/yasin-wu/utils/file_parser"
)

func TestFileParser(t *testing.T) {
	parser, err := file_parser.New("http://47.108.155.25:9998", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}
	fileInfo, err := parser.Parser("../conf/startTika.sh", true)
	if err != nil {
		t.Error(err)
		return
	}
	spew.Dump(fileInfo)
}
