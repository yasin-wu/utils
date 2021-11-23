package test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"

	"github.com/yasin-wu/utils/file_parser"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func TestFileParser(t *testing.T) {
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConf, nil
	})
	fmt.Println("初始化Apollo配置成功")
	cache := client.GetConfigCache(apolloConf.NamespaceName)
	url, _ := cache.Get("tika.url")
	parser, err := file_parser.New(url.(string), nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := parser.Parser("../../dsi_engine/sample/test.docx", true)
	if err != nil {
		log.Fatal(err)
	}
	fileInfos, _ := json.MarshalIndent(fileInfo, "", "\t")
	fmt.Println(string(fileInfos))
}
