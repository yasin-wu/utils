package excel

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/xuri/excelize/v2"
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Excel struct {
	mx sync.Mutex

	fileName  string
	sheetName string
	colWidth  float64
	xlsx      *excelize.File
}

func New(fileName string) *Excel {
	return &Excel{
		mx:       sync.Mutex{},
		fileName: fileName,
		colWidth: 20,
		xlsx:     excelize.NewFile(),
	}
}

func (e *Excel) Write(sheetName string, headers []Header, data []byte) error {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.sheetName = sheetName
	index := e.xlsx.NewSheet(sheetName)
	startCol, _ := excelize.ColumnNumberToName(1)
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	if err := e.xlsx.SetColWidth(sheetName, startCol, endCol, e.colWidth); err != nil {
		return err
	}
	if err := e.writeHeader(headers); err != nil {
		return err
	}
	if err := e.write(headers, data); err != nil {
		return err
	}
	e.xlsx.SetActiveSheet(index)
	return e.xlsx.SaveAs(e.fileName)
}

func (e *Excel) Read(sheetName string) ([]byte, error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	excelFile, err := excelize.OpenFile(e.fileName)
	if err != nil {
		return nil, err
	}
	defer func(excelFile *excelize.File) { _ = excelFile.Close() }(excelFile)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("not found rows")
	}
	var keys []string
	var data []map[string]any
	keys = append(keys, rows[0]...)
	for i := 1; i < len(rows); i++ {
		j := make(map[string]any)
		for k, v := range rows[i] {
			j[keys[k]] = v
		}
		data = append(data, j)
	}
	return json.Marshal(data)
}

func (e *Excel) Close() {
	_ = e.xlsx.Close()
}

func (e *Excel) SetColWidth(width float64) {
	if width != 0 {
		e.colWidth = width
	}
}

func (e *Excel) writeHeader(headers []Header) error {
	for k, v := range headers {
		col, err := excelize.ColumnNumberToName(k + 1)
		if err != nil {
			return err
		}
		if err = e.xlsx.SetCellValue(e.sheetName, fmt.Sprintf("%s1", col), v.Value); err != nil {
			return err
		}
	}
	return nil
}

func (e *Excel) write(headers []Header, data []byte) error {
	var buffer []map[string]any
	if err := json.Unmarshal(data, &buffer); err != nil {
		return err
	}
	for k, v := range buffer {
		for i, header := range headers {
			col, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				return err
			}
			if err = e.xlsx.SetCellValue(e.sheetName, fmt.Sprintf("%s%d", col, k+2), v[header.Key]); err != nil {
				return err
			}
		}
	}
	return nil
}
