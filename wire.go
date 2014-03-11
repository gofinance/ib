package trade

import (
	"bufio"
	"bytes"
	"runtime"
	"strconv"
	"time"
)

//const ibTimeFormat = "20060102 15:04:05.000000 MST"
const ibTimeFormat = "20060102 15:04:05"

type readable interface {
	read(b *bufio.Reader)
}

type writable interface {
	write(buf *bytes.Buffer)
}

func unpanic() (err error) {
	if r := recover(); r != nil {
		if _, ok := r.(runtime.Error); ok {
			panic(r)
		}
		err = r.(error)
	}

	return nil
}

func read(b *bufio.Reader, v readable) error {
	defer unpanic()
	v.read(b)
	return nil
}

func write(buf *bytes.Buffer, v writable) error {
	defer unpanic()
	v.write(buf)
	return nil
}

// Decode

func readString(b *bufio.Reader) string {
	data, err := b.ReadString(0)

	if err != nil {
		panic(err)
	}

	return string(data[:len(data)-1])
}

func readInt(b *bufio.Reader) int64 {
	i, err := strconv.ParseInt(readString(b), 10, 64)

	if err != nil {
		panic(err)
	}

	return i
}

func readFloat(b *bufio.Reader) float64 {
	f, err := strconv.ParseFloat(readString(b), 64)

	if err != nil {
		panic(err)
	}

	return f
}

func readBool(b *bufio.Reader) bool {
	return (readInt(b) > 0)
}

func readTime(b *bufio.Reader) time.Time {
	timeString := readString(b)
	length := len(ibTimeFormat)
	timeString = timeString[0:length]
	t, err := time.ParseInLocation(ibTimeFormat, timeString, time.Local)

	if err != nil {
		panic(err)
	}

	return t
}

// Encode

func writeString(buf *bytes.Buffer, v string) {
	if _, err := buf.WriteString(v + "\000"); err != nil {
		panic(err)
	}
}

func writeInt(buf *bytes.Buffer, v int64) {
	writeString(buf, strconv.FormatInt(v, 10))
}

func writeFloat(buf *bytes.Buffer, v float64) {
	writeString(buf, strconv.FormatFloat(v, 'g', 10, 64))
}

func writeBool(buf *bytes.Buffer, v bool) {
	var s string

	if v {
		s = "1"
	} else {
		s = "0"
	}

	writeString(buf, s)
}

func writeTime(buf *bytes.Buffer, v time.Time) {
	writeString(buf, v.Format(ibTimeFormat))
}
