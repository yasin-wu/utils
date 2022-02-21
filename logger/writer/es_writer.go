package writer

import (
	"bytes"
	"errors"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
)

type ESWriter struct {
	client *elasticsearch.Client
	index  string
}

var _ io.Writer = (*ESWriter)(nil)

type ESConfig elasticsearch.Config

func NewESWriter(index string, config *ESConfig) (*ESWriter, error) {
	esWriter := &ESWriter{index: index}
	if config == nil {
		return nil, errors.New("elasticsearch config is nil")
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config(*config))
	if err != nil {
		return nil, err
	}
	esWriter.client = client
	return esWriter, nil

}

func (w *ESWriter) Write(message []byte) (int, error) {
	resp, err := w.client.Index(w.index, bytes.NewReader(message))
	defer resp.Body.Close()
	if err != nil {
		return -1, err
	}
	if resp.IsError() {
		return -1, errors.New(resp.String())
	}
	return 1, nil
}
