package test

import (
	"testing"

	"github.com/yasin-wu/utils/tool"

	"github.com/yasin-wu/utils/tree"
)

func TestTree(t *testing.T) {
	var nodes []tree.Tree
	nodes = append(nodes, tree.Tree{Id: "0", Name: "all", ParentId: ""})
	nodes = append(nodes, tree.Tree{Id: "1", Name: "第一级", ParentId: "0", Index: 1})
	nodes = append(nodes, tree.Tree{Id: "2", Name: "第二级", ParentId: "1", Index: 1})
	nodes = append(nodes, tree.Tree{Id: "3", Name: "第一级2", ParentId: "0", Index: 2})

	var root tree.Tree
	var all = make(map[string]tree.Tree)
	root.Id = "0"
	root.Name = "all"
	root.Level = 0
	root.MakeTree(nodes)
	root.FindAll(&all)
	tool.Println(root)
	tool.Println(all)
}
