package wire

import (
	"bytes"
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

func makebuf() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 1024))
}

func TestEncodeInt(t *testing.T) {
	v := &intRec{I: 57}
	b := makebuf()
	encode(b, 0, v)
	if b.String() != "57\000" {
		t.Errorf("encode(%v) = %s, want %s", v, b.String(), "57")
	}
}

func TestEncodeString(t *testing.T) {
	v := &stringRec{S: "foobar"}
	b := makebuf()
	encode(b, 0, v)
	if b.String() != "foobar\000" {
		t.Errorf("encode(%v) = %s, want %s", v, b.String(), "foobar")
	}
}

func TestEncodeTime(t *testing.T) {
	ts := time.Now()
	s := ts.Format(ibTimeFormat) + "\000"
	v := &timeRec{T: ts}
	b := makebuf()
	encode(b, 0, v)
	if b.String() != s {
		t.Errorf("encode(%v) = %s, want %s", v, b.String(), s)
	}
}
