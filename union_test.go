package goavro_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/karrick/goavro"
)

func TestUnion(t *testing.T) {
	testCodecBidirectional(t, `["null"]`, goavro.Union("null", nil), []byte("\x00"))
	testCodecBidirectional(t, `["null","int"]`, goavro.Union("null", nil), []byte("\x00"))
	testCodecBidirectional(t, `["int","null"]`, goavro.Union("null", nil), []byte("\x02"))

	testCodecBidirectional(t, `["null","int"]`, goavro.Union("int", 3), []byte("\x02\x06"))
	testCodecBidirectional(t, `["null","long"]`, goavro.Union("long", 3), []byte("\x02\x06"))

	testCodecBidirectional(t, `["int","null"]`, goavro.Union("int", 3), []byte("\x00\x06"))
	testCodecEncoder(t, `["int","null"]`, goavro.Union("int", 3), []byte("\x00\x06")) // can encode a bare 3
}

func TestUnionRejectInvalidType(t *testing.T) {
	t.Skip("TODO: ensure these returns error")
	testCodecBidirectional(t, `["null","long"]`, goavro.Union("int", 3), []byte("\x02\x06"))
	testCodecBidirectional(t, `["null","int","long","float"]`, goavro.Union("double", float64(3.5)), []byte("\x06\x00\x00\x60\x40"))
}

func TestUnionWillCoerceTypeIfPossible(t *testing.T) {
	testCodecBidirectional(t, `["null","long","float","double"]`, goavro.Union("long", int32(3)), []byte("\x02\x06"))
	testCodecBidirectional(t, `["null","int","float","double"]`, goavro.Union("int", int64(3)), []byte("\x02\x06"))
	testCodecBidirectional(t, `["null","int","long","double"]`, goavro.Union("double", float32(3.5)), []byte("\x06\x00\x00\x00\x00\x00\x00\f@"))
	testCodecBidirectional(t, `["null","int","long","float"]`, goavro.Union("float", float64(3.5)), []byte("\x06\x00\x00\x60\x40"))
}

func TestUnionWithArray(t *testing.T) {
	testCodecBidirectional(t, `["null",{"type":"array","items":"int"}]`, goavro.Union("null", nil), []byte("\x00"))

	testCodecBidirectional(t, `["null",{"type":"array","items":"int"}]`, goavro.Union("array", []interface{}{}), []byte("\x02\x00"))
	testCodecBidirectional(t, `["null",{"type":"array","items":"int"}]`, goavro.Union("array", []interface{}{1}), []byte("\x02\x02\x02\x00"))
	testCodecBidirectional(t, `["null",{"type":"array","items":"int"}]`, goavro.Union("array", []interface{}{1, 2}), []byte("\x02\x04\x02\x04\x00"))
}

func TestUnionWithMap(t *testing.T) {
	testCodecBidirectional(t, `["null",{"type":"map","values":"string"}]`, goavro.Union("null", nil), []byte("\x00"))
	testCodecBidirectional(t, `["string",{"type":"map","values":"string"}]`, goavro.Union("map", map[string]interface{}{"He": "Helium"}), []byte("\x02\x02\x04He\x0cHelium\x00"))
	testCodecBidirectional(t, `["string",{"type":"array","items":"string"}]`, goavro.Union("string", "Helium"), []byte("\x00\x0cHelium"))
}

func TestUnionOfEnumsWithSameType(t *testing.T) {
	_, err := goavro.NewCodec(`[{"type":"enum","name":"com.example.foo","symbols":["alpha","bravo"]},"com.example.foo"]`)
	if err == nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, "non-nil")
	}
}

func TestUnionOfEnumsWithDifferentTypeButInvalidString(t *testing.T) {
	codec, err := goavro.NewCodec(`[{"type":"enum","name":"com.example.colors","symbols":["red","green","blue"]},{"type":"enum","name":"com.example.animals","symbols":["dog","cat"]}]`)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := codec.BinaryEncode(nil, "bravo")
	if err == nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, "non-nil")
	}
	if actual, expected := buf, []byte{}; !bytes.Equal(buf, expected) {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
}

func TestUnionOfEnumsWithSameNames(t *testing.T) {
	_, err := goavro.NewCodec(`[{"type":"enum","name":"com.example.one","symbols":["red","green","blue"]},{"type":"enum","name":"one","namespace":"com.example","symbols":["dog","cat"]}]`)
	if err == nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, "non-nil")
	}
}

func TestUnionOfEnumsWithDifferentTypeValidString(t *testing.T) {
	codec, err := goavro.NewCodec(`[{"type":"enum","name":"com.example.colors","symbols":["red","green","blue"]},{"type":"enum","name":"com.example.animals","symbols":["dog","cat"]}]`)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := codec.BinaryEncode(nil, goavro.Union("com.example.animals", "dog"))
	if err != nil {
		t.Fatal(err)
	}
	if actual, expected := buf, []byte{0x2, 0x0}; !bytes.Equal(buf, expected) {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}

	// round trip back to native string
	value, buf, err := codec.BinaryDecode(buf)
	if err != nil {
		t.Fatal(err)
	}
	if actual, expected := buf, []byte{}; !bytes.Equal(buf, expected) {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		t.Fatalf("Actual: %#v; Expected: %#v", ok, false)
	}
	if actual, expected := len(valueMap), 1; actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
	datum, ok := valueMap["com.example.animals"]
	if !ok {
		t.Fatalf("Actual: %#v; Expected: %#v", valueMap, "have `com.example.animals` key")
	}
	if actual, expected := datum.(string), "dog"; actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
}

func TestUnionMapRecordFitsInRecord(t *testing.T) {
	// when encoding union with child object, named types, such as records, enums, and fixed, are named

	// union value may be either map or a record
	codec, err := goavro.NewCodec(`["null",{"type":"map","values":"double"},{"type":"record","name":"com.example.record","fields":[{"name":"field1","type":"int"},{"name":"field2","type":"float"}]}]`)
	if err != nil {
		t.Fatal(err)
	}

	// the provided datum value could be encoded by either the map or the record schemas above
	datum := map[string]interface{}{
		"field1": 3,
		"field2": 3.5,
	}
	datumIn := goavro.Union("com.example.record", datum)

	buf, err := codec.BinaryEncode(nil, datumIn)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte{
		0x04,                   // prefer record (union item 2) over map (union item 1)
		0x06,                   // field1 == 3
		0x00, 0x00, 0x60, 0x40, // field2 == 3.5
	}) {
		t.Errorf("Actual: %#v; Expected: %#v", buf, []byte{byte(2)})
	}

	// round trip
	datumOut, buf, err := codec.BinaryDecode(buf)
	if err != nil {
		t.Fatal(err)
	}
	if actual, expected := len(buf), 0; actual != expected {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}

	datumOutMap, ok := datumOut.(map[string]interface{})
	if !ok {
		t.Fatalf("Actual: %#v; Expected: %#v", ok, false)
	}
	if actual, expected := len(datumOutMap), 1; actual != expected {
		t.Fatalf("Actual: %#v; Expected: %#v", actual, expected)
	}
	datumValue, ok := datumOutMap["com.example.record"]
	if !ok {
		t.Fatalf("Actual: %#v; Expected: %#v", datumOutMap, "have `com.example.record` key")
	}
	datumValueMap, ok := datumValue.(map[string]interface{})
	if !ok {
		t.Errorf("Actual: %#v; Expected: %#v", ok, true)
	}
	if actual, expected := len(datumValueMap), len(datum); actual != expected {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	for k, v := range datum {
		if actual, expected := fmt.Sprintf("%v", datumValueMap[k]), fmt.Sprintf("%v", v); actual != expected {
			t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
		}
	}
}
