package ib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
	read(serverVersion int64, b *bufio.Reader) error
}

type writable interface {
	write(serverVersion int64, buf *bytes.Buffer) error
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

func readInt(b *bufio.Reader) (int64, error) {
	str, err := readString(b)
	if err != nil {
		return -1, err
	}
	if str == "" {
		return math.MaxInt64, nil
	}
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err
}

// readIntList reads an IB pipe-separated string of integers into a Go slice.
func readIntList(b *bufio.Reader) ([]int, error) {
	s, err := readString(b)
	if err != nil {
		return nil, err
	}
	split := strings.Split(s, "|")
	r := make([]int, len(split))
	for i, val := range split {
		r[i], err = strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func readFloat(b *bufio.Reader) (float64, error) {
	str, err := readString(b)
	if err != nil {
		return -1., err
	}
	if str == "" {
		return math.MaxFloat64, nil
	}
	f, err := strconv.ParseFloat(str, 64)
	return f, err
}

// readBool is equivalent of IB API EReader.readBoolFromInt.
func readBool(b *bufio.Reader) (bool, error) {
	i, err := readInt(b)
	if err != nil {
		return false, err
	}
	return (i > 0), nil
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

// read 4 bytes on wire in network order
func readUInt32(b io.Reader) (uint32, error) {

	data := make([]byte, 4)

	n, err := b.Read(data)
	if err != nil {
		return 0, err
	}

	if n != len(data) {
		// error?
	}

	buf := bytes.NewReader(data)

	var num uint32
	err = binary.Read(buf, binary.BigEndian, &num)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func writeString(b *bytes.Buffer, s string) error {
	_, err := b.WriteString(s + "\000")
	return err
}

func writeInt(b *bytes.Buffer, i int64) error {
	return writeString(b, strconv.FormatInt(i, 10))
}

func writeFloat(b *bytes.Buffer, f float64) error {
	return writeString(b, strconv.FormatFloat(f, 'g', 10, 64))
}

func writeTagValue(b *bytes.Buffer, options []TagValue) error {
	var optionsBuf bytes.Buffer

	optionsBuf.WriteString("")
	for _, opt := range options {
		optionsBuf.WriteString(opt.Tag)
		optionsBuf.WriteString("=")
		optionsBuf.WriteString(opt.Value)
		optionsBuf.WriteString(";")
	}
	return writeString(b, optionsBuf.String())
}

// TODO: this never errors. Is it expected?
func writeMaxFloat(b *bytes.Buffer, f float64) error {
	if f >= math.MaxFloat64 {
		b.WriteByte('\000')
	} else {
		writeFloat(b, f)
	}
	return nil
}

// TODO: this never errors. Is it expected?
func writeMaxInt(b *bytes.Buffer, i int64) error {
	if i == math.MaxInt64 {
		b.WriteByte('\000')
	} else {
		writeInt(b, i)
	}
	return nil
}

func writeBool(b *bytes.Buffer, bo bool) error {
	s := "0"
	if bo {
		s = "1"
	}
	return writeString(b, s)
}

func writeTime(b *bytes.Buffer, t time.Time, f timeFmt) error {
	switch f {
	case timeWriteUTC:
		return writeString(b, t.UTC().Format("20060102 15:04:05")+" UTC")
	case timeWriteLocalTime:
		return writeString(b, t.Format("20060102 15:04:05"))
	}
	return fmt.Errorf("goib: cannot write time format '%v'", f)
}

func detectTime(timeString string) (timeFmt, error) {
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

	if _, err := strconv.ParseInt(timeString, 10, 64); err == nil {
		return timeReadEpoch, nil
	}

	return -1, fmt.Errorf("ibgo: '%s' has unknown time format", timeString)
}

func decodeUnicodeEscapedString(str string) string {
	// TODO
	return str
	//String v = new String(str);
	//
	//try {
	//    for (;;) {
	//        int escapeIndex = v.indexOf("\\u");

	//        if (escapeIndex == -1
	//         || v.length() - escapeIndex < 6) {
	//            break;
	//        }

	//        String escapeString = v.substring(escapeIndex ,  escapeIndex + 6);
	//        int hexVal = Integer.parseInt(escapeString.replace("\\u", ""), 16);

	//        v = v.replace(escapeString, "" + (char)hexVal);
	//    }
	//} catch (NumberFormatException e) { }
	//
	//return v;
}

func parseLastTradeDate(contract *ContractDetails, isBond bool, lastTradeDateOrContractMonthStr string) {
	if lastTradeDateOrContractMonthStr != "" {
		//String[] splitted = lastTradeDateOrContractMonth.split("\\s+");
		//if (splitted.length > 0) {
		//    if (isBond) {
		//        contract.maturity(splitted[0]);
		//    } else {
		//        contract.contract().lastTradeDateOrContractMonth(splitted[0]);
		//    }
		//}
		//if (splitted.length > 1) {
		//    contract.lastTradeTime(splitted[1]);
		//}
		//if (isBond && splitted.length > 2) {
		//    contract.timeZoneId(splitted[2]);
		//}
	}
}
