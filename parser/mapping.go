package parser

var TypeMapping = map[string]TypeMap{
	"ObjectiveC": TypeMap{
		KnownMapping: map[string]string{
			"string":            "NSString *",
			"int64":             "NSInteger *",
			"int32":             "NSInteger *",
			"int":               "NSInteger *",
			"float64":           "NSNumber *",
			"float32":           "NSNumber *",
			"float":             "NSNumber *",
			"bool":              "bool",
			"error":             "NSError *",
			"template.HTML":     "NSString *",
			"template.HTMLAttr": "NSString *",
			"time.Time":         "NSDate *",
		},
		UnknownFunc: func(f Field) (r string) {
			r = f.Type + " *"
			return
		},
	},
}

type TypeMap struct {
	KnownMapping map[string]string
	UnknownFunc  func(f Field) (r string)
}

func (tm TypeMap) TypeOf(f Field) (r string) {
	t, ok := tm.KnownMapping[f.Type]
	if ok {
		r = t
		return
	}
	r = tm.UnknownFunc(f)
	return
}
