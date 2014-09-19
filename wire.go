package ib

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type timeFmt int

const (
	timeWriteUTC          timeFmt = iota // write datetime as UTC converted time string with explicit UTC TZ designator
	timeWriteLocalTime                   // write datetime as local time without an explicit TZ designator
	timeReadAutoDetect                   // read datetime (or date) by auto-detecting likely format (first checks epoch)
	timeReadEpoch                        // read datetime as epoch in seconds since 1/1/70 UTC
	timeReadLocalDateTime                // read datetime which is in local time (any TZ designator is ignored)
	timeReadLocalDate                    // read date which is in local time and does not have any timezone designator
	timeReadLocalTime                    // read time which is in local time and does not have any timezone designator
)

type readable interface {
	read(b *bufio.Reader) error
}

type writable interface {
	write(buf *bytes.Buffer) error
}

func readString(b *bufio.Reader) (s string, err error) {
	if s, err = b.ReadString(0); err != nil {
		return
	}
	return string(s[:len(s)-1]), nil
}

// readStringList reads an IB delimited separated string of strings into a Go slice.
func readStringList(b *bufio.Reader, sep string) (r []string, err error) {
	s, err := readString(b)
	if err != nil {
		return
	}
	return strings.Split(s, sep), nil
}

func readInt(b *bufio.Reader) (i int64, err error) {
	var str string
	if str, err = readString(b); err != nil {
		return
	}
	if str == "" {
		return math.MaxInt64, nil
	}
	i, err = strconv.ParseInt(str, 10, 64)
	return
}

// readIntList reads an IB pipe-separated string of integers into a Go slice.
func readIntList(b *bufio.Reader) (r []int, err error) {
	s, err := readString(b)
	if err != nil {
		return
	}
	split := strings.Split(s, "|")
	r = make([]int, len(split))
	for i, val := range split {
		r[i], err = strconv.Atoi(val)
		if err != nil {
			return
		}
	}
	return
}

func readFloat(b *bufio.Reader) (f float64, err error) {
	var str string
	if str, err = readString(b); err != nil {
		return
	}
	if str == "" {
		return math.MaxFloat64, nil
	}
	f, err = strconv.ParseFloat(str, 64)
	return
}

// readBool is equivalent of IB API EReader.readBoolFromInt.
func readBool(b *bufio.Reader) (bo bool, err error) {
	var i int64
	if i, err = readInt(b); err != nil {
		return
	}
	return (i > 0), err
}

// readTime reads a string and then parses it according to the given time format.
// Returned times are always in the local timezone (convert with time.UTC()).
func readTime(b *bufio.Reader, f timeFmt) (t time.Time, err error) {
	var timeString string
	if timeString, err = readString(b); err != nil {
		return
	}

	if f == timeReadAutoDetect {
		f, err = detectTime(timeString)
		if err != nil {
			return
		}
	}

	if f == timeReadEpoch {
		var epochSecs int64
		if epochSecs, err = strconv.ParseInt(timeString, 10, 64); err != nil {
			return
		}
		return time.Unix(epochSecs, 0), nil
	}

	if f == timeReadLocalDateTime {
		format := "20060102 15:04:05"
		if len(timeString) < len(format) {
			return time.Now(), fmt.Errorf("ibgo: '%s' too short to be datetime format '%s'", timeString, format)
		}

		// Truncate any portion this parse does not require (ie timezones)
		fields := strings.Fields(timeString)
		if len(fields) < 2 {
			return time.Now(), fmt.Errorf("ibgo: '%s' does not contain expected whitespace for datetime format '%s'", timeString, format)
		}
		timeString = fields[0] + " " + fields[1]

		return time.ParseInLocation(format, timeString, time.Local)
	}

	if f == timeReadLocalDate {
		format := "20060102"
		if len(timeString) != len(format) {
			return time.Now(), fmt.Errorf("ibgo: '%s' wrong length to be datetime format '%s'", timeString, format)
		}
		return time.ParseInLocation(format, timeString, time.Local)
	}

	if f == timeReadLocalTime {
		formatShort := "15:04"
		formatLong := "15:04:05"
		switch len(timeString) {
		case len(formatShort):
			return time.ParseInLocation(formatShort, timeString, time.Local)
		case len(formatLong):
			return time.ParseInLocation(formatLong, timeString, time.Local)
		default:
			return time.Now(), fmt.Errorf("ibgo: '%s' wrong length to be time format '%s' or '%s'", timeString, formatShort, formatLong)
		}
	}

	return time.Now(), fmt.Errorf("ibgo: unsupported read time format '%v'", f)

}

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

func writeMaxFloat(b *bytes.Buffer, f float64) (err error) {
	if f >= math.MaxFloat64 {
		b.WriteByte('\000')
	} else {
		writeFloat(b, f)
	}
	return
}

func writeMaxInt(b *bytes.Buffer, i int64) (err error) {
	if i == math.MaxInt64 {
		b.WriteByte('\000')
	} else {
		writeInt(b, i)
	}
	return
}

func writeBool(b *bytes.Buffer, bo bool) (err error) {
	s := "0"
	if bo {
		s = "1"
	}
	return writeString(b, s)
}

func writeTime(b *bytes.Buffer, t time.Time, f timeFmt) (err error) {
	switch f {
	case timeWriteUTC:
		return writeString(b, t.UTC().Format("20060102 15:04:05")+" UTC")
	case timeWriteLocalTime:
		return writeString(b, t.Format("20060102 15:04:05"))
	}
	return fmt.Errorf("goib: cannot write time format '%v'", f)
}

func detectTime(timeString string) (f timeFmt, err error) {
	if len(timeString) == len("14:52") && strings.Contains(timeString, ":") {
		return timeReadLocalTime, nil
	}

	if len(timeString) == len("14:52:02") && strings.Contains(timeString, ":") {
		return timeReadLocalTime, nil
	}

	if len(timeString) == len("20060102 15:04:05") && strings.Contains(timeString, ":") {
		return timeReadLocalDateTime, nil
	}

	if len(timeString) >= len("20060102 15:04:05 ") && strings.Contains(timeString, ":") {
		return timeReadLocalDateTime, nil
	}

	// 8 character time strings are ambiguous as they can be a yyyymmdd date
	// or an epoch. So try yyyymmdd first, as 8 char epochs are less likely
	if len(timeString) == len("20060102") {
		return timeReadLocalDate, nil
	}

	if _, err = strconv.ParseInt(timeString, 10, 64); err == nil {
		return timeReadEpoch, nil
	}

	return f, fmt.Errorf("ibgo: '%s' has unknown time format", timeString)
}
