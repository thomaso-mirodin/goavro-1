package e2e

import (
	"bytes"
	"fmt"

	"github.com/karrick/goavro"
)

func init() {
	goavro.MaxBlockSize = 10 * 1024 * 1024
	goavro.MaxBlockCount = 1024
}

func Fuzz(data []byte) int {
	ocfr, err := goavro.NewOCFReader(bytes.NewReader(data))
	if err != nil {
		return 0
	}

	var datums []interface{}
	for ocfr.Scan() {
		datum, err := ocfr.Read()
		if err != nil {
			return 0
		}
		datums = append(datums, datum)
	}

	b := new(bytes.Buffer)
	ocfw, err := goavro.NewOCFWriter(
		goavro.OCFWriterConfig{
			W:           b,
			Compression: goavro.CompressionNull,
			Schema:      ocfr.Schema(),
		})
	if err != nil {
		fmt.Println("failed to create ocf writer")
		panic(err)
	}

	if err := ocfw.Append(datums); err != nil {
		panic(err)
	}

	return 1
}
