package test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/yasin-wu/utils/file_parser"
)

func TestFileParser(t *testing.T) {
	url := "http://47.108.155.25:9998"
	parser := file_parser.New(url)
	fileInfo, err := parser.Parse("../../dsi_engine/sample/test.docx", true)
	if err != nil {
		log.Fatal(err)
	}
	fileInfos, _ := json.MarshalIndent(fileInfo, "", "\t")
	fmt.Println(string(fileInfos))
}
