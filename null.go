package goavro

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func nullBinaryDecoder(buf []byte) (interface{}, []byte, error) { return nil, buf, nil }

func nullBinaryEncoder(buf []byte, datum interface{}) ([]byte, error) {
	if datum != nil {
		return buf, fmt.Errorf("null: expected: Go nil; received: %T", datum)
	}
	return buf, nil
}

func nullTextDecoder(buf []byte) (interface{}, []byte, error) {
	if len(buf) < 4 {
		return nil, buf, io.ErrShortBuffer
	}
	if bytes.Equal(buf[:4], []byte("null")) {
		return nil, buf[4:], nil
	}
	return nil, buf, errors.New("expected: null")
}

func nullTextEncoder(buf []byte, datum interface{}) ([]byte, error) {
	if datum != nil {
		return buf, fmt.Errorf("null: expected: Go nil; received: %T", datum)
	}
	return append(buf, []byte("null")...), nil
}
