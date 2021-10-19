package file_parser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/yasin-wu/go-tika/tika"

	"yasin-wu/utils/common"
)

type Parser struct {
	tika   string
	header http.Header
	client *http.Client
}

func New(tika string, header http.Header, client *http.Client) (*Parser, error) {
	if tika == "" {
		tika = defaultTika
	}
	if header == nil {
		header = defaultHeader
	}
	return &Parser{tika: tika, header: header, client: client}, nil
}

func (this *Parser) Parser(fileName string, isFormat bool) (*FileInfo, error) {
	if fileName == "" {
		return nil, errors.New("fileName is nil")
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileInfo := this.parseFileInfo(file)
	ok := this.checkFileType(fileInfo.FileType)
	if !ok {
		return nil, errors.New("unsupported file type")
	}
	client := tika.NewClient(this.client, this.tika)
	body, err := client.Parse(context.Background(), file, this.header)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("client parse err:%v", err.Error()))
	}
	if isFormat {
		body = this.handleBody(body)
	}
	fileInfo.Content = body
	return fileInfo, nil
}

func (this *Parser) parseFileInfo(file *os.File) *FileInfo {
	fileName := file.Name()
	f, err := os.Stat(fileName)
	var size int64 = 0
	if err == nil {
		size = f.Size()
	}
	fileType := strings.Replace(path.Ext(path.Base(fileName)), ".", "", -1)
	fileInfo := &FileInfo{
		Name:     path.Base(fileName),
		Path:     fileName,
		FileType: fileType,
		Size:     size,
	}
	return fileInfo
}

func (this *Parser) checkFileType(fileType string) bool {
	for _, o := range FileTypes {
		if o == fileType {
			return true
		}
	}
	return false
}

func (this *Parser) handleBody(body string) string {
	body = strings.Replace(body, "\n", "", -1)
	body = strings.Replace(body, "\t", "", -1)
	body = strings.Replace(body, "\r", "", -1)
	body = strings.Replace(body, " ", "", -1)
	body = strings.Replace(body, "\u00a0", "", -1)
	body = strings.Replace(body, "\u200b", "", -1)
	body = common.RemoveHtml(body)
	return body
}
