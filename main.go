package main

import (
	"flag"
	"fmt"
	"github.com/sunfmin/goapigen/parser"
	"go/build"
	"os"
	"text/template"
)

var pkg = flag.String("pkg", "", "Put a full go package path like 'github.com/theplant/qortexapi', make sure you did 'go get github.com/theplant/qortexapi'")
var lang = flag.String("lang", "javascript", "put language like 'javascript', 'objc', 'java'")
var outdir = flag.String("outdir", ".", "the dir to output the generated source code")

func main() {
	flag.Parse()

	pkg, err := build.Default.Import(*pkg, "", 0)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	apis := parser.Parse(pkg.Dir)

	switch *lang {
	case "javascript":
		printjavascript(*outdir, apis)
	}

}

func printjavascript(dir string, apiset *parser.APISet) {
	tpl := template.Must(template.New("").Parse(Templates))
	f, err := os.Create(apiset.Name + ".js")
	if err != nil {
		panic(err)
	}
	err = tpl.ExecuteTemplate(f, "javascript/interfaces", apiset)
	if err != nil {
		panic(err)
	}
}
