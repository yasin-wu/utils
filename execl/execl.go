package execl

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/xuri/excelize/v2"

	js "github.com/bitly/go-simplejson"
)

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

/**
 * @author: yasinWu
 * @date: 2022/3/17 14:06
 * @params: sheetName string, headers [][]string, data []*js.Json
 * headers[0]为需要显示的列名，headers[1]对应列名在data中JSON的key
 * @return: error
 * @description: write excel
 */
func (e *Excel) Write(sheetName string, headers [][]string, data []*js.Json) error {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.sheetName = sheetName
	index := e.xlsx.NewSheet(sheetName)
	startCol, _ := excelize.ColumnNumberToName(1)
	endCol, _ := excelize.ColumnNumberToName(len(headers[0]))
	err := e.xlsx.SetColWidth(sheetName, startCol, endCol, e.colWidth)
	if err != nil {
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
	defer excelFile.Close()
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
	e.xlsx.Close()
}

func (e *Excel) SetColWidth(width float64) {
	if width != 0 {
		e.colWidth = width
	}
}

func (e *Excel) writeHeader(headers [][]string) {
	headerDesc := headers[0]
	for i := 0; i < len(headerDesc); i++ {
		col, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = e.xlsx.SetCellValue(e.sheetName, fmt.Sprintf("%s1", col), headerDesc[i])
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
}

func (e *Excel) write(headers [][]string, data []*js.Json) {
	headerKey := headers[1]
	for k, v := range data {
		for i := 0; i < len(headerKey); i++ {
			col, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			err = e.xlsx.SetCellValue(e.sheetName, fmt.Sprintf("%s%d", col, k+2), v.Get(headerKey[i]).Interface())
			if err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}
}
