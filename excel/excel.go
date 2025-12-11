package excel

import (
	"errors"
	"fmt"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const defaultBatchSize = 10000

type Header struct {
	Key   string `json:"key"`   //表头名称对应数据key,英文
	Value string `json:"value"` //表头名称,中文
}

type Excel struct {
	mx        sync.Mutex
	fileName  string
	colWidth  float64
	startRow  int
	password  string
	batchSize int64
	xlsx      *excelize.File
}

type Option func(e *Excel)

func New(fileName string, opts ...Option) *Excel {
	e := &Excel{
		fileName:  fileName,
		colWidth:  20,
		batchSize: defaultBatchSize,
		xlsx:      excelize.NewFile(),
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

func WithPassword(password string) Option {
	return func(e *Excel) {
		e.password = password
	}
}

func WithBatchSize(batchSize int64) Option {
	return func(e *Excel) {
		e.batchSize = batchSize
	}
}

func (e *Excel) Write(sheetName string, headers []Header, data []byte) error {
	e.mx.Lock()
	defer e.mx.Unlock()
	index, err := e.xlsx.NewSheet(sheetName)
	if err != nil {
		return err
	}
	startCol, _ := excelize.ColumnNumberToName(1)
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	if err := e.xlsx.SetColWidth(sheetName, startCol, endCol, e.colWidth); err != nil {
		return err
	}
	if err := e.writeHeader(sheetName, headers); err != nil {
		return err
	}
	if err := e.write(sheetName, headers, data); err != nil {
		return err
	}
	e.xlsx.SetActiveSheet(index)
	return e.xlsx.SaveAs(e.fileName, e.opts()...)
}

func (e *Excel) Append(sheetName string, headers []Header, data []byte) error {
	e.mx.Lock()
	defer e.mx.Unlock()
	index, err := e.xlsx.NewSheet(sheetName)
	if err != nil {
		return err
	}
	startRow, err := e.GetRow(sheetName)
	if err != nil {
		return err
	}
	e.SetStartRow(startRow - 1)
	if err := e.write(sheetName, headers, data); err != nil {
		return err
	}
	e.xlsx.SetActiveSheet(index)
	return e.xlsx.SaveAs(e.fileName, e.opts()...)
}

func (e *Excel) Read(sheetName string) ([]byte, error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	excelFile, err := excelize.OpenFile(e.fileName, e.opts()...)
	if err != nil {
		return nil, err
	}
	defer func(excelFile *excelize.File) {
		_ = excelFile.Close()
	}(excelFile)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("not found rows")
	}
	var keys []string
	var data []map[string]interface{}
	keys = append(keys, rows[0]...)
	for i := 1; i < len(rows); i++ {
		j := make(map[string]interface{})
		for k, v := range rows[i] {
			j[keys[k]] = v
		}
		data = append(data, j)
	}
	return json.Marshal(data)
}

func (e *Excel) DeleteSheet(sheetName string) error {
	err := e.xlsx.DeleteSheet(sheetName)
	if err != nil {
		return err
	}
	return e.xlsx.SaveAs(e.fileName, e.opts()...)
}

func (e *Excel) Close() {
	_ = e.xlsx.Close()
}

func (e *Excel) SetColWidth(width float64) {
	if width != 0 {
		e.colWidth = width
	}
}

func (e *Excel) SetStartRow(startRow int) {
	e.startRow = startRow
}

func (e *Excel) GetRow(sheetName string) (int, error) {
	rows, err := e.xlsx.Rows(sheetName)
	if err != nil {
		return 0, err
	}
	defer func(rows *excelize.Rows) {
		_ = rows.Close()
	}(rows)
	var rowCount int
	for rows.Next() {
		rowCount++
	}
	if err := rows.Error(); err != nil {
		return 0, err
	}
	return rowCount, err
}

func (e *Excel) NewStyle(style *excelize.Style) (int, error) {
	id, err := e.xlsx.NewStyle(style)
	return id, err
}

func (e *Excel) SetStyle(sheet string, hCell string, vCell string, styleID int) error {
	return e.xlsx.SetCellStyle(sheet, hCell, vCell, styleID)
}

// StreamWrite 流式写入excel
// nolint:funlen
func (e *Excel) StreamWrite(sheetName string, skip, limit int64, searcher Searcher) error {
	index, err := e.xlsx.NewSheet(sheetName)
	if err != nil {
		return err
	}
	e.xlsx.SetActiveSheet(index)
	streamWriter, err := e.xlsx.NewStreamWriter(sheetName)
	if err != nil {
		return fmt.Errorf("create stream writer failed: %v", err)
	}
	var keys []string
	var headers []any
	for _, header := range searcher.Headers() {
		keys = append(keys, header.Key)
		headers = append(headers, header.Value)
	}
	if err = streamWriter.SetColWidth(1, len(headers), e.colWidth); err != nil {
		return fmt.Errorf("set col width failed: %v", err)
	}
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		return fmt.Errorf("create coordinates failed: %v", err)
	}
	if err := streamWriter.SetRow(cell, headers); err != nil {
		return fmt.Errorf("set row failed: %v", err)
	}
	remainder := limit % e.batchSize
	batch := limit/e.batchSize + 1
	rowID := 2
	if remainder == 0 {
		batch--
	}
	for i := 0; i < int(batch); i++ {
		reqSkip := int64(i)*e.batchSize + skip
		reqLimit := e.batchSize
		if i == int(batch)-1 && remainder != 0 {
			reqLimit = remainder
		}
		buf, err := searcher.Search(reqSkip, reqLimit)
		if err != nil {
			return fmt.Errorf("search data failed: %v", err)
		}
		if len(buf) == 0 {
			break
		}
		var data []map[string]any
		if err := json.Unmarshal(buf, &data); err != nil {
			return fmt.Errorf("unmarshal data failed: %v", err)
		}
		for _, v := range data {
			var row []any
			for _, key := range keys {
				row = append(row, v[key])
			}
			cell, err := excelize.CoordinatesToCellName(1, rowID)
			if err != nil {
				return fmt.Errorf("create coordinates failed: %v", err)
			}
			if err := streamWriter.SetRow(cell, row); err != nil {
				return fmt.Errorf("set row failed: %v", err)
			}
			rowID++
		}
	}
	if err := streamWriter.Flush(); err != nil {
		return fmt.Errorf("flush excel failed: %v", err)
	}
	return e.xlsx.SaveAs(e.fileName)
}

func (e *Excel) writeHeader(sheetName string, headers []Header) error {
	for k, v := range headers {
		col, err := excelize.ColumnNumberToName(k + 1)
		if err != nil {
			return err
		}
		if err = e.xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, 1+e.startRow), v.Value); err != nil {
			return err
		}
	}
	return nil
}

func (e *Excel) write(sheetName string, headers []Header, data []byte) error {
	var buffer []map[string]interface{}
	if err := json.Unmarshal(data, &buffer); err != nil {
		return err
	}
	for k, v := range buffer {
		for i, header := range headers {
			col, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				return err
			}
			if err = e.xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, k+2+e.startRow), v[header.Key]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Excel) opts() []excelize.Options {
	var opts []excelize.Options
	if e.password != "" {
		opts = append(opts, excelize.Options{Password: e.password})
	}
	return opts
}
