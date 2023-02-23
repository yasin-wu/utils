package test

import (
	"testing"

	"github.com/yasin-wu/utils/tool"

	"github.com/yasin-wu/utils/tree"
)

func TestTree(t *testing.T) {
	var nodes []tree.Tree
	nodes = append(nodes, tree.Tree{ID: "0", Name: "all", ParentID: ""})
	nodes = append(nodes, tree.Tree{ID: "1", Name: "第一级", ParentID: "0", Index: 1})
	nodes = append(nodes, tree.Tree{ID: "2", Name: "第二级", ParentID: "1", Index: 1})
	nodes = append(nodes, tree.Tree{ID: "3", Name: "第一级2", ParentID: "0", Index: 2})

	var root tree.Tree
	var all = make(map[string]tree.Tree)
	root.ID = "0"
	root.Name = "all"
	root.Level = 0
	root.MakeTree(nodes)
	root.FindAll(&all)
	tool.PrintlnFmt(root)
	tool.PrintlnFmt(all)
}
