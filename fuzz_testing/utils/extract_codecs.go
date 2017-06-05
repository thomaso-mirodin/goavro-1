package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/karrick/goavro"
)

func init() {
	flag.Parse()
}

func main() {
	m, err := filepath.Glob(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	sort.Strings(m)

	for i, name := range m {
		fmt.Println(name)
		f, err := os.Open(name)
		if err != nil {
			panic(err)
		}
		ocfr, err := goavro.NewOCFReader(f)
		if err != nil {
			log.Println(err)
			continue
		}
		//if err := ioutil.WriteFile(strconv.Itoa(i)+".json", []byte(ocfr.Schema()), os.ModePerm); err != nil {
		//	log.Println(err)
		//}
		_ = i
		fmt.Println(ocfr.Schema())
	}
}
