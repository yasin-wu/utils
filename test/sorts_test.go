package test

import (
	"fmt"
	"github.com/yasin-wu/utils/strings"
	"testing"

	"github.com/yasin-wu/utils/tool"

	"github.com/yasin-wu/utils/sorts"
)

type ByIndex []map[string]interface{}

func (this ByIndex) Len() int      { return len(this) }
func (this ByIndex) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this ByIndex) Key(i int) int64 {
	return this[i]["index"].(int64)
}
func (this ByIndex) Less(i, j int) bool {
	return this[i]["index"].(int64) < this[j]["index"].(int64)
}

func TestSort(t *testing.T) {
	var data []map[string]interface{}
	j1 := make(map[string]interface{})
	j1["index"] = int64(1)
	data = append(data, j1)
	j2 := make(map[string]interface{})
	j2["index"] = int64(2)
	data = append(data, j2)
	j3 := make(map[string]interface{})
	j3["index"] = int64(3)
	data = append(data, j3)
	//升序
	sorts.ByInt64(ByIndex(data))
	tool.PrintlnFmt(data)
	//降序
	sorts.Flip(ByIndex(data))
	tool.PrintlnFmt(data)
}

func TestStringIn(t *testing.T) {
	src := []string{"1", "2", "3"}
	fmt.Println(strings.TargetIn("1", src))
	fmt.Println(strings.TargetIn("4", src))
}
