package main

import (
	"flag"
	"fmt"
	"github.com/sunfmin/goapigen/parser"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var pkg = flag.String("pkg", "", "Put a full go package path like 'github.com/theplant/qortexapi', make sure you did 'go get github.com/theplant/qortexapi'")
var lang = flag.String("lang", "javascript", "put language like 'javascript', 'objc', 'java'")
var outdir = flag.String("outdir", ".", "the dir to output the generated source code")
var impl = flag.String("impl", "", "implementation package like 'github.com/theplant/qortex/services'")

func main() {
	flag.Parse()

	buildpkg, err := build.Default.Import(*pkg, "", 0)

	if err != nil {
		die(err)
	}

	apis := parser.Parse(buildpkg.Dir)

	switch *lang {
	case "javascript":
		printjavascript(*outdir, apis)
	case "server":
		printserver(*outdir, apis, *pkg, *impl)
	case "objc":
		printobjc(*outdir, apis)
	}

}

func die(message interface{}) {
	fmt.Println(message)
	os.Exit(0)
}

func dieIf(message interface{}, condition bool) {
	if !condition {
		return
	}
	die(message)
}

func codeTemplate() (tpl *template.Template) {
	tpl = template.New("")
	tpl = tpl.Funcs(template.FuncMap{
		"title":       strings.Title,
		"downcase":    strings.ToLower,
		"dotlastname": dotLastName,
	})
	tpl = template.Must(tpl.Parse(Templates))
	return
}

func dotLastName(pkg string) (r string) {
	names := strings.Split(pkg, "/")
	r = names[len(names)-1]
	return
}

func printserver(dir string, apiset *parser.APISet, apipkg string, impl string) {
	if impl == "" {
		die("must use -impl=your.package/full/path to give implementation package")
	}

	apiset.ServerImports = []string{
		"time",
		"encoding/json",
		apipkg,
		impl,
		"net/http",
		"github.com/sunfmin/govalidations",
	}
	apiset.ImplPkg = impl

	tpl := codeTemplate()

	p := filepath.Join(dir, apiset.Name+"httpimpl", "gen.go")
	os.Mkdir(filepath.Dir(p), 0755)
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	err = tpl.ExecuteTemplate(f, "httpserver", apiset)
	if err != nil {
		panic(err)
	}
}

func printobjc(dir string, apiset *parser.APISet) {
	tpl := codeTemplate()
	hfile, err1 := os.Create(filepath.Join(dir, apiset.Name+".h"))
	dieIf(err1, err1 != nil)
	mfile, err2 := os.Create(filepath.Join(dir, apiset.Name+".m"))
	dieIf(err2, err2 != nil)
	err3 := tpl.ExecuteTemplate(hfile, "objc/h", apiset)
	dieIf(err3, err3 != nil)
	err4 := tpl.ExecuteTemplate(mfile, "objc/m", apiset)
	dieIf(err4, err4 != nil)
}

func printjavascript(dir string, apiset *parser.APISet) {
	tpl := codeTemplate()
	p := filepath.Join(dir, apiset.Name+".js")
	f, err := os.Create(p)
	dieIf(err, err != nil)
	err = tpl.ExecuteTemplate(f, "javascript/interfaces", apiset)
	dieIf(err, err != nil)
}
