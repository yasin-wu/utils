package writer

import (
	"bytes"
	"errors"
	"io"

	"github.com/rs/xid"

	"github.com/elastic/go-elasticsearch/v7"
)

type ESWriter struct {
	client *elasticsearch.Client
	index  string
	idFunc IdFunc
}

var _ io.Writer = (*ESWriter)(nil)

type ESConfig elasticsearch.Config

type IdFunc func() string

func NewESWriter(index string, config *ESConfig, idFunc ...IdFunc) (*ESWriter, error) {
	esWriter := &ESWriter{index: index}
	esWriter.idFunc = esWriter.IdFunc
	if config == nil {
		return nil, errors.New("elasticsearch config is nil")
	}
	if idFunc != nil {
		esWriter.idFunc = idFunc[0]
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config(*config))
	if err != nil {
		return nil, err
	}
	esWriter.client = client
	return esWriter, nil

}

func (w *ESWriter) Write(message []byte) (int, error) {
	resp, err := w.client.Create(w.index, w.idFunc(), bytes.NewReader(message))
	defer resp.Body.Close()
	if err != nil {
		return -1, err
	}
	if resp.IsError() {
		return -1, errors.New(resp.String())
	}
	return 1, nil
}

func (w *ESWriter) IdFunc() string {
	return xid.New().String()
}
