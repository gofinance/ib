package wire

import (
	"bufio"
	"bytes"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type long int64

type intRec struct {
	I int64
}

type longRec struct {
	L long
}

type stringRec struct {
	S string
}

type timeRec struct {
	T time.Time
}

type floatRec struct {
	F float64
}

type item struct {
	S string
}

type sliceRec struct {
	Items []item
}

func makebuf() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 1024))
}

func testEncode(t *testing.T, v interface{}, s string) {
	b := makebuf()

	if err := encode(b, 0, v); err != nil {
		t.Fatal(err)
	}

	if b.String() != s+"\000" {
		t.Fatalf("encode(%v) = %s, want %s", v, b.String(), s)
	}
}

func TestEncodeInt(t *testing.T) {
	v := &intRec{I: 57}
	testEncode(t, v, "57")
}

func TestEncodeLong(t *testing.T) {
	v := &longRec{L: 57}
	testEncode(t, v, "57")
}

func TestEncodeString(t *testing.T) {
	v := &stringRec{S: "foobar"}
	testEncode(t, v, "foobar")
}

func TestEncodeTime(t *testing.T) {
	ts := time.Now()
	v := &timeRec{T: ts}
	testEncode(t, v, ts.Format(ibTimeFormat))
}

func TestEncodeFloat(t *testing.T) {
	f := 0.535
	v := &floatRec{F: f}
	testEncode(t, v, strconv.FormatFloat(f, 'g', 10, 64))
}

func TestEncodeEmptySlice(t *testing.T) {
	v := sliceRec{}
	testEncode(t, v, "0")
}

func TestEncodeSlice(t *testing.T) {
	v := sliceRec{Items: []item{{"foo"}, {"bar"}}}
	testEncode(t, v, "2\000foo\000bar")
}

func testDecode(t *testing.T, src interface{}, dst interface{}) {
	b := makebuf()

	if err := encode(b, 0, src); err != nil {
		t.Fatal(err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))

	if err := decode(r, reflect.ValueOf(dst)); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeInt(t *testing.T) {
	v1 := &intRec{I: 57}
	v2 := &intRec{}

	testDecode(t, v1, v2)

	if *v1 != *v2 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}

func TestDecodeString(t *testing.T) {
	v1 := &stringRec{S: "foobar"}
	v2 := &stringRec{}

	testDecode(t, v1, v2)

	if *v1 != *v2 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}

func TestDecodeTime(t *testing.T) {
	v1 := &timeRec{T: time.Now()}
	v2 := &timeRec{}

	testDecode(t, v1, v2)

	if *v1 != *v2 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}

func TestDecodeFloat(t *testing.T) {
	v1 := &floatRec{F: 0.545}
	v2 := &floatRec{}

	testDecode(t, v1, v2)

	if *v1 != *v2 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}
