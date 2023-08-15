package file

import (
	"context"
	"errors"
	"fmt"
	"github.com/yasin-wu/utils/file/internal"
	"net/http"
	"os"
	"path"
	"strings"

	strings2 "github.com/yasin-wu/utils/strings"
)

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:41
 * @description: 文件解析器配置项
 */
type Option func(parser *Parser)

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:41
 * @description: Parser Client
 */
type Parser struct {
	tika   string
	header http.Header
	client *http.Client
	ctx    context.Context
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:41
 * @params: tika string, options ...Option
 * @return: *Parser
 * @description: 新建Parser Client
 */
func NewParser(tika string, options ...Option) *Parser {
	if tika == "" {
		tika = defaultTika
	}
	parser := &Parser{tika: tika, ctx: context.Background()}
	for _, f := range options {
		f(parser)
	}
	if parser.header == nil {
		parser.header = defaultHeader
	}
	return parser
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:42
 * @params: header http.Header
 * @return: Option
 * @description: 配置http header
 */
func WithHeader(header http.Header) Option {
	return func(parser *Parser) {
		parser.header = header
	}
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:42
 * @params: client *http.Client
 * @return: Option
 * @description: 配置http client
 */
func WithClient(client *http.Client) Option {
	return func(parser *Parser) {
		parser.client = client
	}
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:42
 * @params: filePath string, isFormat bool
 * @return: *FileInfo, error
 * @description: 解析文件
 */
func (p *Parser) Parse(filePath string, isFormat bool) (*File, error) {
	if filePath == "" {
		return nil, errors.New("filePath is nil")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	fileInfo := p.parseFileInfo(file)
	ok := p.checkFileType(fileInfo.fileType)
	if !ok {
		return nil, errors.New("unsupported file type")
	}
	client := internal.NewClient(p.client, p.tika)
	body, err := client.Parse(p.ctx, file, p.header)
	if err != nil {
		return nil, fmt.Errorf("client parse err:%w", err)
	}
	if isFormat {
		body = p.handleBody(body)
	}
	fileInfo.content = body
	return fileInfo, nil
}

func (p *Parser) parseFileInfo(file *os.File) *File {
	fileName := file.Name()
	f, err := os.Stat(fileName)
	var size int64
	if err == nil {
		size = f.Size()
	}
	fileType := strings.ReplaceAll(path.Ext(path.Base(fileName)), ".", "")
	info := &File{
		name:     path.Base(fileName),
		path:     fileName,
		fileType: fileType,
		size:     size,
	}
	return info
}

func (p *Parser) checkFileType(fileType string) bool {
	for _, o := range FileTypes {
		if o == fileType {
			return true
		}
	}
	return false
}

func (p *Parser) handleBody(body string) string {
	body = strings.ReplaceAll(body, "\n", "")
	body = strings.ReplaceAll(body, "\t", "")
	body = strings.ReplaceAll(body, "\r", "")
	body = strings.ReplaceAll(body, " ", "")
	body = strings.ReplaceAll(body, "\u00a0", "")
	body = strings.ReplaceAll(body, "\u200b", "")
	body = strings2.DeleteHTML(body)
	return body
}
