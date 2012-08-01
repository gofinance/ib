package wire

import (
	"bytes"
	"reflect"
	"runtime"
	"strconv"
	"time"
	//"bufio"
)

const ibTime = "20060102 15:04:05 MST"
const maxInt = int(^uint(0) >> 1)

type TickType int

const (
	bid  TickType = 1
	Ask           = 2
	Last          = 4
)

type serverVersion struct {
	Version int
}

type clientVersion struct {
	Version int
}

type serverTime struct {
	Time time.Time
}

type Tick struct {
	TickerId       int
	Type           int
	Price          float64
	Size           int
	CanAutoExecute bool
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

func encode(buf *bytes.Buffer, tag int, v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	buf.Reset()

	if val.Kind() != reflect.Struct {
		panic("Value given to decode is not a structure")
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		switch f.Type().Kind() {
		case reflect.Int:
			buf.WriteString(strconv.FormatInt(f.Int(), 10))
		case reflect.Float64:
			buf.WriteString(strconv.FormatFloat(f.Float(), 'g', 10, 64))
		case reflect.String:
			var s string = string(f.String())
			buf.WriteString(s)
		case reflect.Bool:
			if f.Int() > 0 {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		}
		buf.WriteByte(0)
	}

	return nil
}
