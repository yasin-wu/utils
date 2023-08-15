package elasticsearch

import (
	"bytes"
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/yasin-wu/utils/logger/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"time"
)

type es struct {
	index      string
	indexType  string
	level      string
	timeLayout string
	depth      int
	stacktrace bool
	client     *elasticsearch.Client
}

var _ io.Writer = (*es)(nil)
var _ internal.Corer = (*es)(nil)

type ESConfig elasticsearch.Config
type Option func(e *es)

func New(index, level string, config *ESConfig, options ...Option) (internal.Corer, error) {
	if config == nil {
		return nil, errors.New("elasticsearch config is nil")
	}
	e := &es{
		index:      index + "_" + level,
		level:      level,
		timeLayout: "2006-01-02 15:04:05",
	}
	for _, f := range options {
		f(e)
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config(*config))
	if err != nil {
		return nil, err
	}
	e.client = client
	return e, nil
}

func (e *es) Encoder() zapcore.Encoder {
	encoderConfig := internal.DefaultEncoderConfig
	encoderConfig.EncodeCaller = internal.CallerEncoder(e.depth)
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(e.timeLayout)
	if e.stacktrace {
		encoderConfig.StacktraceKey = "stacktrace"
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func (e *es) WriteSyncer() zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(e))
}

func (e *es) AtomicLevel() zap.AtomicLevel {
	return internal.AtomicLevel(e.level)
}

func (e *es) Write(message []byte) (int, error) {
	resp, err := e.client.Index(e.handleIndex(), bytes.NewReader(message))
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	if err != nil {
		return -1, err
	}
	if resp.IsError() {
		return -1, errors.New(resp.String())
	}
	return 1, nil
}

func (e *es) handleIndex() string {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")
	var buf bytes.Buffer
	buf.WriteString(e.index)
	switch strings.ToLower(e.indexType) {
	case "y", "year":
		buf.WriteString("_")
		buf.WriteString(year)
	case "m", "month":
		buf.WriteString("_")
		buf.WriteString(year)
		buf.WriteString("_")
		buf.WriteString(month)
	case "d", "day":
		buf.WriteString("_")
		buf.WriteString(year)
		buf.WriteString("_")
		buf.WriteString(month)
		buf.WriteString("_")
		buf.WriteString(day)
	}
	return buf.String()
}

func WithTimeLayout(timeLayout string) Option {
	return func(e *es) {
		if timeLayout != "" {
			e.timeLayout = timeLayout
		}
	}
}

func WithDepth(depth int) Option {
	return func(e *es) {
		e.depth = depth
	}
}

func WithStacktrace(stacktrace bool) Option {
	return func(e *es) {
		e.stacktrace = stacktrace
	}
}

func WithIndexType(indexType string) Option {
	return func(e *es) {
		e.indexType = indexType
	}
}
