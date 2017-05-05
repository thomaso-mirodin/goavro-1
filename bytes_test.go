package goavro_test

import "testing"

func TestSchemaPrimitiveCodecBytes(t *testing.T) {
	testSchemaPrimativeCodec(t, "bytes")
}

func TestPrimitiveBytesBinary(t *testing.T) {
	testBinaryEncodeFailBadDatumType(t, "bytes", 13)
	testBinaryDecodeFailShortBuffer(t, "bytes", nil)
	testBinaryDecodeFailShortBuffer(t, "bytes", []byte{2})
	testBinaryCodecPass(t, "bytes", []byte(""), []byte("\x00"))
	testBinaryCodecPass(t, "bytes", []byte("some bytes"), []byte("\x14some bytes"))
}

func TestPrimitiveBytesText(t *testing.T) {
	testTextEncodeFailBadDatumType(t, "bytes", 42)
	testTextDecodeFailShortBuffer(t, "bytes", []byte(``))
	testTextDecodeFailShortBuffer(t, "bytes", []byte(`"`))
	testTextDecodeFailShortBuffer(t, "bytes", []byte(`..`))
	testTextDecodeFailShortBuffer(t, "bytes", []byte(`".`))

	testTextCodecPass(t, "bytes", []byte(""), []byte("\"\""))
	testTextCodecPass(t, "bytes", []byte("a"), []byte("\"a\""))
	testTextCodecPass(t, "bytes", []byte("ab"), []byte("\"ab\""))
	testTextCodecPass(t, "bytes", []byte("a\"b"), []byte("\"a\\\"b\""))
	testTextCodecPass(t, "bytes", []byte("a\\b"), []byte("\"a\\\\b\""))
	testTextCodecPass(t, "bytes", []byte("a/b"), []byte("\"a\\/b\""))

	testTextCodecPass(t, "bytes", []byte("a\bb"), []byte(`"a\bb"`))
	testTextCodecPass(t, "bytes", []byte("a\fb"), []byte(`"a\fb"`))
	testTextCodecPass(t, "bytes", []byte("a\nb"), []byte(`"a\nb"`))
	testTextCodecPass(t, "bytes", []byte("a\rb"), []byte(`"a\rb"`))
	testTextCodecPass(t, "bytes", []byte("a\tb"), []byte(`"a\tb"`))
	testTextCodecPass(t, "bytes", []byte("a	b"), []byte(`"a\tb"`)) // tab byte between a and b

	testTextDecodeFail(t, "bytes", []byte("\"a\\u\""), "short buffer")
	testTextDecodeFail(t, "bytes", []byte("\"a\\u0\""), "short buffer")
	testTextDecodeFail(t, "bytes", []byte("\"a\\u00\""), "short buffer")
	testTextDecodeFail(t, "bytes", []byte("\"a\\u004\""), "short buffer")
	// testTextDecodePass(t, "bytes", []byte("\"AA\""), []byte("\"a\\u0041\""))
}

func TestSchemaPrimitiveStringCodec(t *testing.T) {
	testSchemaPrimativeCodec(t, "string")
}

func TestPrimitiveStringBinary(t *testing.T) {
	testBinaryEncodeFailBadDatumType(t, "string", 42)
	testBinaryDecodeFailShortBuffer(t, "string", nil)
	testBinaryDecodeFailShortBuffer(t, "string", []byte{2})
	testBinaryCodecPass(t, "string", "", []byte("\x00"))
	testBinaryCodecPass(t, "string", "some string", []byte("\x16some string"))
}

func TestPrimitiveStringText(t *testing.T) {
	testTextEncodeFailBadDatumType(t, "string", 42)
	testTextDecodeFailShortBuffer(t, "string", []byte(``))
	testTextDecodeFailShortBuffer(t, "string", []byte(`"`))
	testTextDecodeFailShortBuffer(t, "string", []byte(`..`))
	testTextDecodeFailShortBuffer(t, "string", []byte(`".`))

	testTextCodecPass(t, "string", "", []byte("\"\""))
	testTextCodecPass(t, "string", "a", []byte("\"a\""))
	testTextCodecPass(t, "string", "ab", []byte("\"ab\""))
	testTextCodecPass(t, "string", "a\"b", []byte("\"a\\\"b\""))
	testTextCodecPass(t, "string", "a\\b", []byte("\"a\\\\b\""))
	testTextCodecPass(t, "string", "a/b", []byte("\"a\\/b\""))

	testTextCodecPass(t, "string", "a\bb", []byte(`"a\bb"`))
	testTextCodecPass(t, "string", "a\fb", []byte(`"a\fb"`))
	testTextCodecPass(t, "string", "a\nb", []byte(`"a\nb"`))
	testTextCodecPass(t, "string", "a\rb", []byte(`"a\rb"`))
	testTextCodecPass(t, "string", "a\tb", []byte(`"a\tb"`))
	testTextCodecPass(t, "string", "a	b", []byte(`"a\tb"`)) // tab byte between a and b

	testTextDecodeFail(t, "string", []byte("\"a\\u\""), "short buffer")
	testTextDecodeFail(t, "string", []byte("\"a\\u0\""), "short buffer")
	testTextDecodeFail(t, "string", []byte("\"a\\u00\""), "short buffer")
	testTextDecodeFail(t, "string", []byte("\"a\\u004\""), "short buffer")
	// testTextDecodePass(t, "string", []byte("\"AA\""), []byte("\"a\\u0041\""))
}
