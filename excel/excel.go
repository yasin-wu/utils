package excel

import (
	"errors"
	"fmt"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Excel struct {
	mx sync.Mutex

	fileName string
	colWidth float64
	startRow int
	password string
	xlsx     *excelize.File
}

type Option func(e *Excel)

func New(fileName string, opts ...Option) *Excel {
	e := &Excel{
		fileName: fileName,
		colWidth: 20,
		xlsx:     excelize.NewFile(),
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
	data, err := e.xlsx.GetRows(sheetName)
	return len(data), err
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
