package e2e_repro

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/karrick/goavro"
)

func init() {
	goavro.MaxBlockSize = 10 * 1024 * 1024
	goavro.MaxBlockCount = 1024
}

func readDatums(b []byte) ([]interface{}, *goavro.OCFReader, error) {
	ocfr, err := goavro.NewOCFReader(bytes.NewReader(b))
	if err != nil {
		return nil, nil, err
	}

	var datums []interface{}
	for ocfr.Scan() {
		datum, err := ocfr.Read()
		if err != nil {
			return nil, nil, err
		}
		datums = append(datums, datum)
	}

	return datums, ocfr, nil
}

func writeDatums(ocfr *goavro.OCFReader, datums []interface{}) []byte {
	b := new(bytes.Buffer)
	ocfw, err := goavro.NewOCFWriter(
		goavro.OCFWriterConfig{
			W:           b,
			Compression: ocfr.CompressionID(),
			Schema:      ocfr.Schema(),
		})
	if err != nil {
		fmt.Println("failed to create ocf writer")
		panic(err)
	}

	if err := ocfw.Append(datums); err != nil {
		panic(err)
	}

	return b.Bytes()
}

func Fuzz(data []byte) int {
	rawDatums, ocfr, err := readDatums(data)
	if err != nil {
		return 0
	}

	w1 := writeDatums(ocfr, rawDatums)
	for i := 0; i < 10; i++ {
		w2 := writeDatums(ocfr, rawDatums)
		if !bytes.Equal(w1, w2) {
			panic("Failed to re-create the same bytes from the same avro data")
		}
	}

	datums, ocfr, err := readDatums(w1)
	if err != nil {
		panic("Unable to re-read written data")
	}

	if !reflect.DeepEqual(rawDatums, datums) {
		panic("Datums didn't end up equal")
	}

	return 1
}
