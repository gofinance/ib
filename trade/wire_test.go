package trade

import (
	"bufio"
	"bytes"
	"strconv"
	"testing"
	"time"
)

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

type tagRec struct {
	S string
	I long `when:"S" cond:"is" value:""`
}

func makebuf() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 1024))
}

func testEncode(t *testing.T, v interface{}, s string) {
	b := makebuf()

	if err := Encode(b, v); err != nil {
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
	v := &sliceRec{}
	testEncode(t, v, "0")
}

func TestEncodeSlice(t *testing.T) {
	v := &sliceRec{Items: []item{{"foo"}, {"bar"}}}
	testEncode(t, v, "2\000foo\000bar")
}

func TestNotEncodeTag(t *testing.T) {
	v := &tagRec{I: 10}
	testEncode(t, v, "")
}

func TestEncodeTag(t *testing.T) {
	v := &tagRec{S: "yes!", I: 10}
	testEncode(t, v, "yes!\00010")
}

func testDecode(t *testing.T, src interface{}, dst interface{}) {
	b := makebuf()

	if err := Encode(b, src); err != nil {
		t.Fatal(err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))

	if err := Decode(r, dst); err != nil {
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
	now := time.Now()
	now = now.Add(time.Duration(-1 * now.Nanosecond()))
	v1 := &timeRec{T: now}
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

func TestDecodeEmptySlice(t *testing.T) {
	v1 := &sliceRec{}
	v2 := &sliceRec{}

	testDecode(t, v1, v2)

	if len(v2.Items) != 0 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}

func TestDecodeSlice(t *testing.T) {
	v1 := &sliceRec{Items: []item{{"foo"}, {"bar"}}}
	v2 := &sliceRec{}

	testDecode(t, v1, v2)

	if len(v2.Items) != 2 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}

	if v2.Items[0] != v1.Items[0] || v2.Items[1] != v1.Items[1] {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}

func TestNotDecodeTag(t *testing.T) {
	v1 := &tagRec{I: 10}
	v2 := &tagRec{}

	testDecode(t, v1, v2)

	if v2.S != "" || v2.I != 0 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}

func TestDecodeTag(t *testing.T) {
	v1 := &tagRec{S: "yes!", I: 10}
	v2 := &tagRec{}

	testDecode(t, v1, v2)

	if *v1 != *v2 {
		t.Fatalf("decode got %v, want %v", v2, v1)
	}
}
