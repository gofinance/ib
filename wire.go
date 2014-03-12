package trade

import (
	"bufio"
	"bytes"
	"strconv"
	"time"
)

//const ibTimeFormat = "20060102 15:04:05.000000 MST"
const ibTimeFormat = "20060102 15:04:05"

type readable interface {
	read(b *bufio.Reader) error
}

type writable interface {
	write(buf *bytes.Buffer) error
}

// Decode

func readString(b *bufio.Reader) (s string, err error) {
	if s, err = b.ReadString(0); err != nil {
		return
	}
	return string(s[:len(s)-1]), nil
}

func readInt(b *bufio.Reader) (i int64, err error) {
	var str string
	if str, err = readString(b); err != nil {
		return
	}
	i, err = strconv.ParseInt(str, 10, 64)
	return
}

func readFloat(b *bufio.Reader) (f float64, err error) {
	var str string
	if str, err = readString(b); err != nil {
		return
	}
	f, err = strconv.ParseFloat(str, 64)
	return
}

func readBool(b *bufio.Reader) (bo bool, err error) {
	var i int64
	if i, err = readInt(b); err != nil {
		return
	}
	return (i > 0), err
}

func readTime(b *bufio.Reader) (t time.Time, err error) {
	var timeString string
	if timeString, err = readString(b); err != nil {
		return
	}
	length := len(ibTimeFormat)
	timeString = timeString[0:length]

	t, err = time.ParseInLocation(ibTimeFormat, timeString, time.Local)
	return
}

// Encode

func writeString(buf *bytes.Buffer, v string) (err error) {
	_, err = buf.WriteString(v + "\000")
	return
}

func writeInt(buf *bytes.Buffer, v int64) (err error) {
	return writeString(buf, strconv.FormatInt(v, 10))
}

func writeFloat(buf *bytes.Buffer, v float64) (err error) {
	return writeString(buf, strconv.FormatFloat(v, 'g', 10, 64))
}

func writeBool(buf *bytes.Buffer, v bool) (err error) {
	s := "0"
	if v {
		s = "1"
	}
	return writeString(buf, s)
}

func writeTime(buf *bytes.Buffer, v time.Time) (err error) {
	return writeString(buf, v.Format(ibTimeFormat))
}
