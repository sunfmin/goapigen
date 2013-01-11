package parser

import (
	"strings"
)

type DataObject struct {
	Name   string
	Fields []Field
}

type Constructor struct {
	FromInterface *Interface
	Method        *Method
}

type Interface struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name                    string
	Params                  []Field
	Results                 []Field
	ConstructorForInterface *Interface
}

func (m *Method) ParamsForJavascriptFunction() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, p.Name)
	}
	r = strings.Join(ps, ", ")
	return
}

func (m *Method) ParamsForJson() (r string) {
	ps := []string{}
	for _, p := range m.Params {
		ps = append(ps, `"`+p.Name+`": `+p.Name)
	}
	r = strings.Join(ps, ", ")
	r = "{ " + r + " }"
	return
}

type Field struct {
	IsArray bool
	Type    string
	Name    string
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
	Name        string
	Interfaces  []*Interface
	DataObjects []*DataObject
}
