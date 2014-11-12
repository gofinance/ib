package ib

import (
	"bufio"
	"bytes"
	"math"
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
		t.Fatalf("expected %s but got %s", x, y)
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

func TestReadEmptyStringReturnsMaxInt(t *testing.T) {
	b := makebuf()
	writeString(b, "")
	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readInt(r)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if y != math.MaxInt64 {
		t.Fatalf("expected maximum int64 but got %d", y)
	}
}

func TestReadIntList(t *testing.T) {
	b := makebuf()
	writeString(b, "1|2|3")
	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readIntList(r)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if len(y) != 3 {
		t.Fatalf("expected 3 but got %d", len(y))
	}
}

func TestWriteTimeUTC(t *testing.T) {
	b := makebuf()
	ts := time.Now()
	writeTime(b, ts, timeWriteUTC)
	expected := ts.UTC().Format("20060102 15:04:05") + " UTC\000"
	if b.String() != expected {
		t.Fatalf("writeTime(%s) = %s, want %s", ts, b.String(), expected)
	}
}

func TestWriteTimeLocal(t *testing.T) {
	b := makebuf()
	ts := time.Now()
	writeTime(b, ts, timeWriteLocalTime)
	expected := ts.Format("20060102 15:04:05") + "\000"
	if b.String() != expected {
		t.Fatalf("writeTime(%s) = %s, want %s", ts, b.String(), expected)
	}
}

func TestReadTimeAutoDetect(t *testing.T) {
	if f, _ := detectTime("349834583"); f != timeReadEpoch {
		t.Fatalf("failed to detect epoch")
	}
	if f, _ := detectTime("20140322 23:59:31"); f != timeReadLocalDateTime {
		t.Fatalf("failed to detect local time")
	}
	if f, _ := detectTime("20140322 23:59:31 US"); f != timeReadLocalDateTime {
		t.Fatalf("failed to detect local time when timezone present")
	}
	if f, _ := detectTime("20140322"); f != timeReadLocalDate {
		t.Fatalf("failed to detect local date")
	}
	if f, _ := detectTime("15:33"); f != timeReadLocalTime {
		t.Fatalf("failed to detect local hh:mm time")
	}
	if f, _ := detectTime("15:33:42"); f != timeReadLocalTime {
		t.Fatalf("failed to detect local hh:mm:ss time")
	}
	if _, err := detectTime("2014/03/22"); err == nil {
		t.Fatalf("failed to return error for unknown time format")
	}
}

func TestReadTimeEpoch(t *testing.T) {
	x := time.Now()
	x = x.Add(time.Duration(-1 * x.Nanosecond()))
	b := makebuf()
	err := writeString(b, strconv.Itoa(int(x.Unix())))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readTime(r, timeReadEpoch)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x != y {
		t.Fatalf("expected %v but got %v", x, y)
	}
}

func TestReadTimeLocalDateTime(t *testing.T) {
	x := time.Now()
	x = x.Add(time.Duration(-1 * x.Nanosecond()))
	b := makebuf()

	err := writeTime(b, x, timeWriteLocalTime)
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readTime(r, timeReadLocalDateTime)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x != y {
		t.Fatalf("expected %v but got %v", x, y)
	}
}

func TestReadTimeLocalDate(t *testing.T) {
	x := time.Date(2014, 3, 22, 0, 0, 0, 0, time.Local)
	b := makebuf()
	err := writeString(b, x.Format("20060102"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readTime(r, timeReadLocalDate)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x != y {
		t.Fatalf("expected %v but got %v", x, y)
	}
}

func TestReadTimeLocalTimeShort(t *testing.T) {
	x := time.Date(2014, 3, 22, 14, 4, 0, 0, time.Local)
	b := makebuf()
	err := writeString(b, x.Format("15:04"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readTime(r, timeReadLocalTime)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x.Hour() != y.Hour() || x.Minute() != y.Minute() {
		t.Fatalf("expected %v but got %v", x, y)
	}
}

func TestReadTimeLocalTimeLong(t *testing.T) {
	x := time.Date(2014, 3, 22, 14, 4, 32, 0, time.Local)
	b := makebuf()
	err := writeString(b, x.Format("15:04:05"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	r := bufio.NewReader(bytes.NewReader(b.Bytes()))
	y, err := readTime(r, timeReadLocalTime)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if x.Hour() != y.Hour() || x.Minute() != y.Minute() || x.Second() != y.Second() {
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
