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

func readHistDataTime(b *bufio.Reader) (t time.Time, err error) {
	// The date can either be in YYYYMMDD format (if daily bars were requested)
	// or in seconds since 1/1/1970 UTC (since we passed 2 as formatDate to the request)
	var timeString string
	if timeString, err = readString(b); err != nil {
		return
	}
	if len(timeString) == 8 {
		// handle YYYYMMDD format (received daily bars)
		if t, err = time.Parse("20060102", timeString); err != nil {
			return
		}
	} else {
		// handle bars that are less than daily (seconds since epoch)
		var epochSecs int64
		if epochSecs, err = strconv.ParseInt(timeString, 10, 64); err != nil {
			return
		}
		t = time.Unix(epochSecs, 0)
	}
	return
}

// Encode

func writeString(b *bytes.Buffer, s string) (err error) {
	_, err = b.WriteString(s + "\000")
	return
}

func writeInt(b *bytes.Buffer, i int64) (err error) {
	return writeString(b, strconv.FormatInt(i, 10))
}

func writeFloat(b *bytes.Buffer, f float64) (err error) {
	return writeString(b, strconv.FormatFloat(f, 'g', 10, 64))
}

func writeBool(b *bytes.Buffer, bo bool) (err error) {
	s := "0"
	if bo {
		s = "1"
	}
	return writeString(b, s)
}

func writeTime(b *bytes.Buffer, t time.Time) (err error) {
	return writeString(b, t.Format(ibTimeFormat))
}
