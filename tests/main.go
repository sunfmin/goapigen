package main

import (
	"fmt"
	"github.com/sunfmin/goapigen"
	"go/build"
	"os"
)

func main() {
	pkg, err := build.Default.Import("github.com/theplant/qortexapi", "", 0)

	if err != nil {
		println(err)
		os.Exit(0)
	}

	apis := goapigen.Parse(pkg.Dir)

	for _, inf := range apis.Interfaces {
		for _, m := range inf.Methods {
			fmt.Printf("%+v\n\n", m)
		}
	}

}
