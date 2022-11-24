package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/yasin-wu/utils/tool"

	"github.com/yasin-wu/utils/excel"
)

func TestFunctions(t *testing.T) {
	fmt.Println(excelize.ColumnNumberToName(1))
	fmt.Println(excelize.ColumnNumberToName(3))
}

func TestWrite(t *testing.T) {
	execl := excel.New("./log/book.xlsx")
	defer execl.Close()
	headers := []excel.Header{
		{"name", "书名"},
		{"author", "作者"},
		{"time", "时间"},
	}
	var data []excel.Data
	for i := 0; i < 10; i++ {
		j := make(excel.Data)
		j["name"] = fmt.Sprintf("书名%d", i)
		j["author"] = fmt.Sprintf("作者%d", i)
		j["time"] = time.Now().String()
		data = append(data, j)
	}
	err := execl.Write("Sheet1", headers, data)
	if err != nil {
		log.Fatal(err)
	}
	err = execl.Write("Sheet2", headers, data)
	if err != nil {
		log.Fatal(err)
	}
}

func TestRead(t *testing.T) {
	execl := excel.New("./log/book.xlsx")
	defer execl.Close()
	data, err := execl.Read("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	tool.Println(data)
	data, err = execl.Read("Sheet2")
	if err != nil {
		log.Fatal(err)
	}
	tool.Println(data)
}
