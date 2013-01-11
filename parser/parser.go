package parser

import (
	// "fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func Parse(dir string) (r *APISet) {

	fset := token.NewFileSet()

	astPkgs, _ := parser.ParseDir(fset, dir, nil, 0)

	var foundPkg *ast.Package
	for _, astPkg := range astPkgs {
		foundPkg = astPkg
	}

	w := &Walker{fset: fset, APISet: &APISet{}}

	// ast.Print(fset, foundPkg)
	ast.Walk(w, foundPkg)
	updateConstructors(w.APISet)
	r = w.APISet

	return
}

type Walker struct {
	fset              *token.FileSet
	APISet            *APISet
	currentName       string
	currentInterface  *Interface
	currentDataObject *DataObject
	currentMethod     *Method
	currentFieldList  *[]Field
}

// Visit implements the ast.Visitor interface.
func (w *Walker) Visit(node ast.Node) ast.Visitor {

	switch n := node.(type) {
	case *ast.Package:
		w.APISet.Name = n.Name
	case *ast.TypeSpec:
		w.currentName = n.Name.Name
	case *ast.StructType:
		w.currentDataObject = &DataObject{Name: w.currentName}
		w.currentInterface = nil
		w.currentFieldList = nil
		w.APISet.DataObjects = append(w.APISet.DataObjects, w.currentDataObject)
	case *ast.InterfaceType:
		w.currentInterface = &Interface{Name: w.currentName}
		w.currentDataObject = nil
		w.currentFieldList = nil
		w.APISet.Interfaces = append(w.APISet.Interfaces, w.currentInterface)
	case *ast.FuncType:
		if w.currentInterface != nil {
			w.currentMethod = &Method{Name: w.currentName}
			for _, param := range n.Params.List {
				f := Field{}
				parseField(param, &f)
				w.currentMethod.Params = append(w.currentMethod.Params, f)
			}

			for _, result := range n.Results.List {
				f := Field{}
				parseField(result, &f)
				w.currentMethod.Results = append(w.currentMethod.Results, f)
			}
			w.currentInterface.Methods = append(w.currentInterface.Methods, w.currentMethod)
			w.currentDataObject = nil
			w.currentFieldList = nil
		}

	case *ast.Field:
		if w.currentInterface != nil && len(n.Names) > 0 {
			w.currentName = n.Names[0].Name
			// fmt.Println(w.currentName)
		}

		if w.currentDataObject != nil && len(n.Names) > 0 {
			f := Field{}
			parseField(n, &f)
			w.currentDataObject.Fields = append(w.currentDataObject.Fields, f)
		}
	}
	return w
}

func updateConstructors(apiset *APISet) {
	for _, inf := range apiset.Interfaces {
		for _, inftarget := range apiset.Interfaces {
			for _, m := range inftarget.Methods {
				for _, f := range m.Results {
					if f.Type == inf.Name {
						m.ConstructorForInterface = inf
					}
				}
			}
		}
	}
}

func parseField(n *ast.Field, f *Field) {
	f.Name = n.Names[0].Name
	switch nt := n.Type.(type) {
	case *ast.Ident:
		f.Type = nt.Name
	case *ast.SelectorExpr:
		f.Type = nt.X.(*ast.Ident).Name + "." + nt.Sel.Name
	case *ast.StarExpr:
		switch xt := nt.X.(type) {
		case *ast.Ident:
			f.Type = xt.Name
		case *ast.SelectorExpr:
			f.Type = xt.X.(*ast.Ident).Name + "." + xt.Sel.Name
		}

	case *ast.ArrayType:
		var tname *ast.Ident
		st, isstar := nt.Elt.(*ast.StarExpr)
		if isstar {
			tname = st.X.(*ast.Ident)
		} else {
			tname = nt.Elt.(*ast.Ident)
		}
		f.Type = tname.Name
		f.IsArray = true

	}
}
