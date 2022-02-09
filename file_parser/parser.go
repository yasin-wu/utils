package file_parser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/yasin-wu/utils/file_parser/tika"

	"github.com/yasin-wu/utils/common"
)

/**
 * @author: yasin
 * @date: 2022/1/13 14:41
 * @description: 文件解析器配置项
 */
type Option func(parser *Parser)

/**
 * @author: yasin
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
 * @author: yasin
 * @date: 2022/1/13 14:41
 * @params: tika string, options ...Option
 * @return: *Parser
 * @description: 新建Parser Client
 */
func New(tika string, options ...Option) *Parser {
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
 * @author: yasin
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
 * @author: yasin
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
 * @author: yasin
 * @date: 2022/1/13 14:42
 * @params: filePath string, isFormat bool
 * @return: *FileInfo, error
 * @description: 解析文件
 */
func (p *Parser) Parse(filePath string, isFormat bool) (*FileInfo, error) {
	if filePath == "" {
		return nil, errors.New("filePath is nil")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileInfo := p.parseFileInfo(file)
	ok := p.checkFileType(fileInfo.FileType)
	if !ok {
		return nil, errors.New("unsupported file type")
	}
	client := tika.NewClient(p.client, p.tika)
	body, err := client.Parse(p.ctx, file, p.header)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("client parse err:%v", err.Error()))
	}
	if isFormat {
		body = p.handleBody(body)
	}
	fileInfo.Content = body
	return fileInfo, nil
}

func (p *Parser) parseFileInfo(file *os.File) *FileInfo {
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

func (p *Parser) checkFileType(fileType string) bool {
	for _, o := range FileTypes {
		if o == fileType {
			return true
		}
	}
	return false
}

func (p *Parser) handleBody(body string) string {
	body = strings.Replace(body, "\n", "", -1)
	body = strings.Replace(body, "\t", "", -1)
	body = strings.Replace(body, "\r", "", -1)
	body = strings.Replace(body, " ", "", -1)
	body = strings.Replace(body, "\u00a0", "", -1)
	body = strings.Replace(body, "\u200b", "", -1)
	body = common.RemoveHtml(body)
	return body
}
