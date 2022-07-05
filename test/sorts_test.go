package test

import (
	"testing"

	"github.com/yasin-wu/utils/tool"

	"github.com/yasin-wu/utils/sorts"

	js "github.com/bitly/go-simplejson"
)

type ByIndex []*js.Json

func (this ByIndex) Len() int      { return len(this) }
func (this ByIndex) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this ByIndex) Key(i int) int64 {
	return this[i].Get("index").MustInt64()
}
func (this ByIndex) Less(i, j int) bool {
	return this[i].Get("index").MustInt64() < this[j].Get("index").MustInt64()
}

func TestSort(t *testing.T) {
	var data []*js.Json
	j1 := js.New()
	j1.Set("index", 1)
	data = append(data, j1)
	j2 := js.New()
	j2.Set("index", 3)
	data = append(data, j2)
	j3 := js.New()
	j3.Set("index", 2)
	data = append(data, j3)
	//升序
	sorts.ByInt64(ByIndex(data))
	tool.Println(data)
	//降序
	sorts.Flip(ByIndex(data))
	tool.Println(data)
}
