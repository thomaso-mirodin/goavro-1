package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/karrick/goavro/fuzz_testing/e2e"
)

func init() {
	flag.Parse()
}

func main() {
	files, err := filepath.Glob(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	for _, name := range files {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}

		fmt.Println(name, e2e.Fuzz(b))
	}
}
