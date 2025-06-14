package sorts

import (
	"fmt"
	"testing"

	"github.com/yasin-wu/utils/strings"

	"github.com/yasin-wu/utils/util"
)

type ByIndex []map[string]interface{}

func (b ByIndex) Len() int      { return len(b) }
func (b ByIndex) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByIndex) Key(i int) int64 {
	return b[i]["index"].(int64)
}
func (b ByIndex) Less(i, j int) bool {
	return b[i]["index"].(int64) < b[j]["index"].(int64)
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
	ByInt64(ByIndex(data))
	util.PrintlnFmt(data)
	//降序
	Flip(ByIndex(data))
	util.PrintlnFmt(data)
}

func TestStringIn(t *testing.T) {
	src := []string{"1", "2", "3"}
	fmt.Println(strings.TargetIn("1", src))
	fmt.Println(strings.TargetIn("4", src))
}
