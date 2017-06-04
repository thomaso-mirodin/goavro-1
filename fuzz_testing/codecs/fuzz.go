package codecs

import "github.com/karrick/goavro"

func Fuzz(data []byte) int {
	_, err := goavro.NewCodec(string(data))
	if err != nil {
		return 0
	}

	return 1
}
