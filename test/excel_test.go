package test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

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
	var data []map[string]interface{}
	for i := 0; i < 10; i++ {
		j := make(map[string]interface{})
		j["name"] = fmt.Sprintf("书名%d", i)
		j["author"] = fmt.Sprintf("作者%d", i)
		j["time"] = time.Now().String()
		data = append(data, j)
	}
	buffer, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	if err = execl.Write("Sheet1", headers, buffer); err != nil {
		log.Fatal(err)
	}
	if err = execl.Write("Sheet2", headers, buffer); err != nil {
		log.Fatal(err)
	}
}

func TestRead(t *testing.T) {
	execl := excel.New("./log/book.xlsx")
	defer execl.Close()
	buffer, err := execl.Read("Sheet1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buffer))
	buffer, err = execl.Read("Sheet2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buffer))
}
