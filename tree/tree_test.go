package tree

import (
	"testing"

	"github.com/yasin-wu/utils/tool"
)

func TestTree(t *testing.T) {
	var nodes []Tree
	nodes = append(nodes, Tree{ID: "0", Name: "all", ParentID: ""})
	nodes = append(nodes, Tree{ID: "1", Name: "第一级", ParentID: "0", Index: 1})
	nodes = append(nodes, Tree{ID: "2", Name: "第二级", ParentID: "1", Index: 1})
	nodes = append(nodes, Tree{ID: "3", Name: "第一级2", ParentID: "0", Index: 2})

	var root Tree
	var all = make(map[string]Tree)
	root.ID = "0"
	root.Name = "all"
	root.Level = 0
	root.MakeTree(nodes)
	root.FindAll(&all)
	tool.PrintlnFmt(root)
	tool.PrintlnFmt(all)
}
