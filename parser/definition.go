package parser

import (
	"strings"
)

type DataObject struct {
	Name   string
	Fields []*Field
}

type Constructor struct {
	FromInterface *Interface
	Method        *Method
}

type Interface struct {
	Name        string
	Methods     []*Method
	Constructor *Constructor
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
	IsArray    bool
	Type       string
	Name       string
	Star       bool
	ImportName string
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

func (f Field) ToLanguageField(language string) (r Field) {
	languageMap, ok := TypeMapping[language]
	if !ok {
		panic(language + " not supported.")
	}

	r.Name = f.Name
	r.Type = languageMap.TypeOf(f)
	return
}

type APISet struct {
	Name          string
	ImplPkg       string
	ServerImports []string
	Interfaces    []*Interface
	DataObjects   []*DataObject
}
