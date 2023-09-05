package strings

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

const (
	gruleName     = "TagEngine"
	gruleVersion  = "1.0.0"
	gruleMaxCycle = 1
)

type Tag struct {
	Result bool
	value  []string
}

func Match(selector string, src []string) (bool, error) {
	t := &Tag{
		Result: false,
		value:  src,
	}
	dataContext := ast.NewDataContext()
	err := dataContext.Add("Tag", t)
	if err != nil {
		return false, fmt.Errorf("dataContext.Add err: %v", err.Error())
	}
	rule := t.rule(selector)
	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	ruleResource := pkg.NewBytesResource([]byte(rule))
	err = ruleBuilder.BuildRuleFromResource(gruleName, gruleVersion, ruleResource)
	if err != nil {
		return false, fmt.Errorf("ruleBuilder.BuildRuleFromResource err: %v", err.Error())
	}
	kb, err := lib.NewKnowledgeBaseInstance(gruleName, gruleVersion)
	if err != nil {
		return false, fmt.Errorf("new knowledge base instance err: %v", err.Error())
	}
	eng := &engine.GruleEngine{MaxCycle: gruleMaxCycle}
	err = eng.Execute(dataContext, kb)
	if err != nil {
		return false, fmt.Errorf("eng.Execute err: %v", err.Error())
	}
	return t.Result, nil
}

func (t *Tag) Do(value string) bool {
	temp := make([]string, len(t.value))
	copy(temp, t.value)
	sort.Strings(temp)
	index := sort.SearchStrings(temp, value)
	if index < len(temp) && temp[index] == value {
		return true
	}
	return false
}

func (t *Tag) rule(selector string) string {
	var buf bytes.Buffer
	buf.WriteString(`rule TagCheck "TagCheck" { when `)
	for _, v := range strings.Split(selector, " ") {
		switch {
		case v == "||" || v == "&&":
			buf.WriteString(v + " ")
		case strings.HasPrefix(v, "(") && strings.HasSuffix(v, ")"):
			buf.WriteString("(")
			buf.WriteString(fmt.Sprintf(`Tag.Do("%s")`, v[1:strings.LastIndex(v, ")")]))
			buf.WriteString(") ")
		case strings.HasSuffix(v, ")"):
			buf.WriteString(fmt.Sprintf(`Tag.Do("%s")`, v[:strings.LastIndex(v, ")")]))
			buf.WriteString(") ")
		case strings.HasPrefix(v, "("):
			buf.WriteString("(")
			buf.WriteString(fmt.Sprintf(`Tag.Do("%s") `, v[1:]))
		default:
			buf.WriteString(fmt.Sprintf(`Tag.Do("%s") `, v))
		}
	}
	buf.WriteString(`then Tag.Result=true; Retract("TagCheck"); }`)
	return buf.String()
}
