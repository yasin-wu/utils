package writer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

type ESWriter struct {
	client      *elasticsearch.Client
	indexPrefix string
}

var _ io.Writer = (*ESWriter)(nil)

type ESConfig elasticsearch.Config

func NewESWriter(indexPrefix string, config *ESConfig) (*ESWriter, error) {
	esWriter := &ESWriter{indexPrefix: indexPrefix}
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
	resp, err := w.client.Index(w.handleIndex(), bytes.NewReader(message))
	defer resp.Body.Close()
	if err != nil {
		return -1, err
	}
	if resp.IsError() {
		return -1, errors.New(resp.String())
	}
	return 1, nil
}

func (w *ESWriter) handleIndex() string {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")
	return fmt.Sprintf("%s_%s_%s_%s", w.indexPrefix, year, month, day)
}
