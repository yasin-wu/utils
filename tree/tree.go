package tree

type Tree struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	ParentId   string   `json:"parent_id"`
	ParentName string   `json:"parent_name"`
	Parents    string   `json:"parents"`
	ParentList []string `json:"parent_list"`
	Level      int      `json:"level"`
	Index      int      `json:"index"`
	ChildCount int      `json:"child_count"`
	Child      []Tree   `json:"child"`
}

func (t *Tree) MakeTree(nodes []Tree) {
	for _, v := range nodes {
		if v.ParentId == t.Id {
			v.ParentList = append(v.ParentList, t.ParentList...)
			v.ParentList = append(v.ParentList, t.Id)
			v.Level = t.Level + 1
			v.ParentName = t.Name
			makeTree(&v, nodes)
			t.Child = append(t.Child, v)
			t.ChildCount = len(t.Child)
		}
	}
}

func (t *Tree) FindAll(nodes *map[string]Tree) {
	(*nodes)[t.Id] = *t
	for _, v := range t.Child {
		v.FindAll(nodes)
	}
}

func makeTree(node *Tree, groups []Tree) {
	child := findChild(*node, groups)
	for _, v := range child {
		if has(v, groups) {
			makeTree(&v, groups)
		}
		node.Child = append(node.Child, v)
		node.ChildCount = len(node.Child)
	}
}

func findChild(node Tree, groups []Tree) []Tree {
	var result []Tree
	for _, v := range groups {
		if node.Id == v.ParentId {
			v.ParentList = append(v.ParentList, node.ParentList...)
			v.ParentList = append(v.ParentList, node.Id)
			v.Level = node.Level + 1
			v.ParentName = node.Name
			result = append(result, v)
		}
	}
	return result
}

func has(node Tree, groups []Tree) bool {
	has := false
	for _, v := range groups {
		if node.Id == v.ParentId {
			has = true
			break
		}
	}
	return has
}
