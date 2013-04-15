package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	NodeName() string
	Children() []Node
	AddChild(n Node)
}

type DataObject struct {
	Name       string
	Fields     []*Field
	ChildNodes []Node
}

type Constructor struct {
	FromInterface *Interface
	Method        *Method
}

type Interface struct {
	Name        string
	Methods     []*Method
	Constructor *Constructor
	ChildNodes  []Node
}

func (do *DataObject) NodeName() string {
	return do.Name
}

func (do *DataObject) Children() []Node {
	return do.ChildNodes
}

func (do *DataObject) AddChild(n Node) {
	do.ChildNodes = append(do.ChildNodes, n)
	return
}

func (inf *Interface) NodeName() string {
	return inf.Name
}

func (inf *Interface) Children() []Node {
	return inf.ChildNodes
}

func (inf *Interface) AddChild(n Node) {
	inf.ChildNodes = append(inf.ChildNodes, n)
	return
}

type Method struct {
	Name                    string
	Params                  []*Field
	Results                 []*Field
	ConstructorForInterface *Interface
}

func (m *Method) ResultsForJavascriptFunction(prefix string) (r string) {
	rs := []string{}
	for _, r := range m.Results {
		rs = append(rs, prefix+"."+strings.Title(r.Name))
	}
	r = strings.Join(rs, ", ")
	return
}

func (m *Method) ParamsForJavascriptFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, p.Name)
	}
	r = strings.Join(ps, ", ")
	return
}

func (m *Method) ParamsForObjcFunction() (r string) {
	if len(m.Params) == 0 {
		r = m.Name
		return
	}

	ps := []string{}
	for i, p := range m.Params {
		op := p.ToLanguageField("objc")
		name := op.Name
		if i == 0 {
			name = m.Name
		}
		ps = append(ps, name+":("+op.FullObjcTypeName()+")"+op.Name)
	}
	r = strings.Join(ps, " ")
	return
}

func (m *Method) ParamsForGoServerFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, "p.Params."+strings.Title(p.Name))
	}
	r = strings.Join(ps, ", ")
	return
}

func (m *Method) ParamsForGoServerConstructorFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, "p.This."+strings.Title(p.Name))
	}
	r = strings.Join(ps, ", ")
	return
}

func (m *Method) ResultsForGoServerFunction(prefix string) (r string) {
	rs := []string{}
	for _, r := range m.Results {
		rs = append(rs, prefix+"."+strings.Title(r.Name))
	}
	r = strings.Join(rs, ", ")
	return
}

func (m *Method) ParamsForJson() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, `"`+strings.Title(p.Name)+`": `+p.Name)
	}
	r = strings.Join(ps, ", ")
	r = "{ " + r + " }"
	return
}

type Field struct {
	IsArray                     bool
	Type                        string
	Name                        string
	Star                        bool
	ImportName                  string
	PropertyAnnotation          string
	SetPropertyConvertFormatter string
	GetPropertyConvertFormatter string
	Primitive                   bool
	ConstructorType             string
}

func (f Field) FullGoTypeName() (r string) {
	if f.IsArray {
		r = r + "[]"
	}
	if f.Star {
		r = r + "*"
	}
	if f.ImportName != "" {
		r = r + f.ImportName + "."
	}
	r = r + f.Type
	return
}

func (f Field) FullObjcTypeName() (r string) {
	if f.IsArray {
		return "NSArray *"
	}
	r = f.Type
	return
}

func (f Field) SetPropertyFromObjcDict(key string) (r string) {
	r = fmt.Sprintf(f.SetPropertyConvertFormatter, "[dict valueForKey:@\""+key+"\"]")
	return
}

func (f Field) GetPropertyToObjcDict(key string) (r string) {
	r = fmt.Sprintf(f.GetPropertyConvertFormatter, key)
	return
}

func (f Field) GetPropertyObjc() (r string) {
	r = "self." + strings.Title(f.Name)
	return
}

func findDefiniationNode(t string, apiset *APISet) (r Node) {
	for _, do := range apiset.DataObjects {
		if t == do.Name {
			return do
		}
	}
	for _, inf := range apiset.Interfaces {
		if t == inf.Name {
			return inf
		}
	}
	return
}

func (f *Field) Update(apiset *APISet, parentNode Node) {
	n := findDefiniationNode(f.Type, apiset)
	f.Primitive = true
	if n != nil {
		f.ImportName = apiset.Name
		f.Primitive = false
		parentNode.AddChild(n)
	}
}

func (f Field) ToLanguageField(language string) (r Field) {
	languageMap, ok := TypeMapping[language]
	if !ok {
		panic(language + " not supported.")
	}

	r.Name = f.Name
	r.IsArray = f.IsArray
	r.Star = f.Star
	r.ImportName = f.ImportName
	t := languageMap.TypeOf(f)
	r.Type = t.Type
	r.Primitive = f.Primitive
	r.PropertyAnnotation = t.PropertyAnnotation
	r.SetPropertyConvertFormatter = t.SetPropertyConvertFormatter
	r.GetPropertyConvertFormatter = t.GetPropertyConvertFormatter
	r.ConstructorType = t.ConstructorType
	return
}

type APISet struct {
	Name          string
	ImplPkg       string
	ServerImports []string
	Interfaces    []*Interface
	DataObjects   []*DataObject
}
