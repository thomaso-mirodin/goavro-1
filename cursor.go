package goavro

type cursor []byte

func (f *cursor) DecodeBoolean() (v bool, e error) {
	var i interface{}
	if i, *f, e = booleanBinaryDecoder(*f); e == nil {
		v = i.(bool)
	}
	return
}

func (f *cursor) EncodeBoolean(v bool) error {
	*f, _ = booleanBinaryEncoder(*f, v) // only fails for bad type
	return nil
}
