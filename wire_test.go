package ib

import (
	"bufio"
	"bytes"
	"strconv"
	"testing"
	"time"
)

func makebuf() *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, 0, 1024))
}

func TestWriteString(t *testing.T) {
	b := makebuf()
	writeString(b, "foobar")
	expected := "foobar\000"
	if b.String() != expected {
		t.Fatalf("writeString('foobar') = %s, want %s", b.String(), expected)
	}
}

func TestWriteInt(t *testing.T) {
	b := makebuf()
	writeInt(b, int64(57))
	expected := "57\000"
	if b.String() != expected {
		t.Fatalf("writeInt(57) = %s, want %s", b.String(), expected)
	}
}

func TestWriteTime(t *testing.T) {
	b := makebuf()
	ts := time.Now()
	writeTime(b, ts)
	expected := ts.Format(ibTimeFormat) + "\000"
	if b.String() != expected {
		t.Fatalf("writeTime(%s) = %s, want %s", ts, b.String(), expected)
	}
}

func TestWriteFloat(t *testing.T) {
	f := 0.535
	b := makebuf()
	writeFloat(b, f)
	expected := strconv.FormatFloat(f, 'g', 10, 64) + "\000"
	if b.String() != expected {
		t.Fatalf("writeFloat(%g) = %s, want %s", f, b.String(), expected)
	}
}

func TestReadString(t *testing.T) {
	x := "foobar"
	b := makebuf()

	writeString(b, x)
	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readString(r)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x != y {
		t.Fatalf("expected %d but got %d", x, y)
	}
}

func TestReadInt(t *testing.T) {
	x := int64(57)
	b := makebuf()

	writeInt(b, x)
	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readInt(r)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x != y {
		t.Fatalf("expected %d but got %d", x, y)
	}
}

func TestReadTime(t *testing.T) {
	x := time.Now()
	x = x.Add(time.Duration(-1 * x.Nanosecond()))
	b := makebuf()

	writeTime(b, x)
	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readTime(r)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x != y {
		t.Fatalf("expected %v but got %v", x, y)
	}
}

func TestReadFloat(t *testing.T) {
	x := 0.545
	b := makebuf()

	writeFloat(b, x)
	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readFloat(r)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x != y {
		t.Fatalf("expected %v but got %v", x, y)
	}
}
