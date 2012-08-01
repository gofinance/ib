package wire

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

type stringRec struct {
	S string
}

type timeRec struct {
	T time.Time
}

type floatRec struct {
	F float64
}

func makebuf() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 1024))
}

func TestEncodeInt(t *testing.T) {
	v := &intRec{I: 57}
	b := makebuf()

	if err := encode(b, 0, v); err != nil {
		t.Fatal(err)
	}

	if b.String() != "57\000" {
		t.Fatalf("encode(%v) = %s, want %s", v, b.String(), "57")
	}
}

func TestEncodeString(t *testing.T) {
	v := &stringRec{S: "foobar"}
	b := makebuf()

	if err := encode(b, 0, v); err != nil {
		t.Fatal(err)
	}

	if b.String() != "foobar\000" {
		t.Fatalf("encode(%v) = %s, want %s", v, b.String(), "foobar")
	}
}

func TestEncodeTime(t *testing.T) {
	ts := time.Now()
	s := ts.Format(ibTimeFormat) + "\000"
	v := &timeRec{T: ts}
	b := makebuf()

	if err := encode(b, 0, v); err != nil {
		t.Fatal(err)
	}

	if b.String() != s {
		t.Fatalf("encode(%v) = %s, want %s", v, b.String(), s)
	}
}

func TestEncodeFloat(t *testing.T) {
	f := 0.535
	v := &floatRec{F: f}
	s := strconv.FormatFloat(f, 'g', 10, 64) + "\000"
	b := makebuf()

	if err := encode(b, 0, v); err != nil {
		t.Fatal(err)
	}

	if b.String() != s {
		t.Fatalf("encode(%v) = %s, want %s", v, b.String(), s)
	}
}

func testDecode(t *testing.T, src interface{}, dst interface{}) {
	b := makebuf()

	if err := encode(b, 0, src); err != nil {
        t.Fatal(err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))

	if err := decode(r, dst); err != nil {
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
