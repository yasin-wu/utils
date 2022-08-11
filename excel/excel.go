package excel

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/xuri/excelize/v2"

	js "github.com/bitly/go-simplejson"
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
		fileName: fileName,
		colWidth: 20,
		xlsx:     excelize.NewFile(),
	}
}

func (e *Excel) Write(sheetName string, headers []Header, data []*js.Json) error {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.sheetName = sheetName
	index := e.xlsx.NewSheet(sheetName)
	startCol, _ := excelize.ColumnNumberToName(1)
	endCol, _ := excelize.ColumnNumberToName(len(headers))
	if err := e.xlsx.SetColWidth(sheetName, startCol, endCol, e.colWidth); err != nil {
		return err
	}
	e.writeHeader(headers)
	e.write(headers, data)
	e.xlsx.SetActiveSheet(index)
	return e.xlsx.SaveAs(e.fileName)
}

func (e *Excel) Read(sheetName string) ([]*js.Json, error) {
	e.mx.Lock()
	defer e.mx.Unlock()
	excelFile, err := excelize.OpenFile(e.fileName)
	defer func(excelFile *excelize.File) {
		_ = excelFile.Close()
	}(excelFile)
	if err != nil {
		return nil, err
	}
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("not found rows")
	}
	var keys []string //nolint:prealloc
	var data []*js.Json
	keys = append(keys, rows[0]...)
	for i := 1; i < len(rows); i++ {
		j := js.New()
		for k, v := range rows[i] {
			j.Set(keys[k], v)
		}
		data = append(data, j)
	}
	return data, nil
}

func (e *Excel) Close() {
	_ = e.xlsx.Close()
}

func (e *Excel) SetColWidth(width float64) {
	if width != 0 {
		e.colWidth = width
	}
}

func (e *Excel) writeHeader(headers []Header) {
	for k, v := range headers {
		col, err := excelize.ColumnNumberToName(k + 1)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if err = e.xlsx.SetCellValue(e.sheetName, fmt.Sprintf("%s1", col), v.Value); err != nil {
			log.Println(err.Error())
			continue
		}
	}
}

func (e *Excel) write(headers []Header, data []*js.Json) {
	for k, v := range data {
		for i, header := range headers {
			col, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if err = e.xlsx.SetCellValue(e.sheetName, fmt.Sprintf("%s%d", col, k+2), v.Get(header.Key).Interface()); err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}
}
