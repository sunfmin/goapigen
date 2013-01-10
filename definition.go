package goapigen

type DataObject struct {
	Name   string
	Fields []Field
}

type Interface struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name    string
	Params  []Field
	Results []Field
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
	Interfaces  []*Interface
	DataObjects []*DataObject
}
