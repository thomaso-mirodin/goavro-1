package goavro

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

////////////////////////////////////////
// Binary Decode
////////////////////////////////////////

func bytesBinaryDecoder(buf []byte) (interface{}, []byte, error) {
	if len(buf) < 1 {
		return nil, nil, io.ErrShortBuffer
	}
	var decoded interface{}
	var err error
	if decoded, buf, err = longBinaryDecoder(buf); err != nil {
		return nil, buf, fmt.Errorf("bytes: %s", err)
	}
	size := decoded.(int64) // longDecoder always returns int64
	if size < 0 {
		return nil, buf, fmt.Errorf("bytes: negative length: %d", size)
	}
	if size > int64(len(buf)) {
		return nil, buf, io.ErrShortBuffer
	}
	return buf[:size], buf[size:], nil
}

func stringBinaryDecoder(buf []byte) (interface{}, []byte, error) {
	d, b, err := bytesBinaryDecoder(buf)
	if err != nil {
		return nil, buf, err
	}
	return string(d.([]byte)), b, nil
}

////////////////////////////////////////
// Binary Encode
////////////////////////////////////////

func bytesBinaryEncoder(buf []byte, datum interface{}) ([]byte, error) {
	var value []byte
	switch v := datum.(type) {
	case []byte:
		value = v
	case string:
		value = []byte(v)
	default:
		return buf, fmt.Errorf("bytes: expected: Go string or []byte; received: %T", v)
	}
	// longEncoder only fails when given non int, so elide error checking
	buf, _ = longBinaryEncoder(buf, len(value))
	// append datum bytes
	return append(buf, value...), nil
}

func stringBinaryEncoder(buf []byte, datum interface{}) ([]byte, error) {
	return bytesBinaryEncoder(buf, datum)
}

////////////////////////////////////////
// Text Decode
////////////////////////////////////////

func bytesTextDecoder(buf []byte) (interface{}, []byte, error) {
	return nil, nil, errors.New("TODO")
}

func stringTextDecoder(buf []byte) (interface{}, []byte, error) {
	// FIXME: this assumes remainder of buf is the string
	// TODO: process unicode bytes
	buflen := len(buf)
	if buflen < 2 || buf[0] != '"' || buf[buflen-1] != '"' {
		return nil, buf, io.ErrShortBuffer
	}
	newBytes := make([]byte, 0, buflen-2)
	var escaped bool
	for i, l := 1, buflen-1; i < l; i++ { // first and last byte are quotes: do not process
		b := buf[i]
		if escaped {
			switch b {
			case '"':
				newBytes = append(newBytes, '"')
			case '\\':
				newBytes = append(newBytes, '\\')
			case '/':
				newBytes = append(newBytes, '/')
			case 'b':
				newBytes = append(newBytes, '\b')
			case 'f':
				newBytes = append(newBytes, '\f')
			case 'n':
				newBytes = append(newBytes, '\n')
			case 'r':
				newBytes = append(newBytes, '\r')
			case 't':
				newBytes = append(newBytes, '\t')
			case 'u':
				if i > buflen-6 { // FIXME 6 --> 5 ???
					return nil, buf, io.ErrShortBuffer
				}
				blob := buf[i+1 : i+5]
				v, err := strconv.ParseUint(string(blob), 16, 64)
				if err != nil {
					return nil, buf, err
				}
				r := rune(v)
				_ = r

			default:
				newBytes = append(newBytes, b)
			}
			escaped = false
			continue
		}
		if b == '\\' {
			escaped = true
			continue
		}
		newBytes = append(newBytes, buf[i])
	}
	return string(newBytes), buf[buflen:], nil
}

////////////////////////////////////////
// Text Encode
////////////////////////////////////////

func bytesTextEncoder(buf []byte, datum interface{}) ([]byte, error) {
	var someBytes []byte
	switch v := datum.(type) {
	case []byte:
		someBytes = v
	case string:
		someBytes = []byte(v)
	default:
		return buf, fmt.Errorf("bytes: expected: Go string or []byte; received: %T", v)
	}
	buf = append(buf, '"')
	for i := 0; i < len(someBytes); i++ {
		switch b := someBytes[i]; b {
		case '"', '\\', '/':
			buf = append(buf, []byte{'\\', b}...)
		case '\b':
			buf = append(buf, []byte("\\b")...)
		case '\f':
			buf = append(buf, []byte("\\f")...)
		case '\n':
			buf = append(buf, []byte("\\n")...)
		case '\r':
			buf = append(buf, []byte("\\r")...)
		case '\t':
			buf = append(buf, []byte("\\t")...)
		default:
			buf = append(buf, b)
		}
	}
	return append(buf, '"'), nil
}

func stringTextEncoder(buf []byte, datum interface{}) ([]byte, error) {
	return bytesTextEncoder(buf, datum)
}
