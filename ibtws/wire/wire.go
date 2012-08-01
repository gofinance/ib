package wire

import (
	"bytes"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

const ibTimeFormat = "20060102 15:04:05 MST"
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

func decodeServerTime(data string) (time.Time, error) {
	return time.Parse(ibTimeFormat, data)
}

type EncodeError struct {
    Type  reflect.Type 
}
  
func (e *EncodeError) Error() string {
    return "ibtws: cannot encode type " + e.Type.String()
}

func encode(buf *bytes.Buffer, tag int, v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	buf.Reset()

	if val.Kind() != reflect.Struct {
		panic("ibtws: value given to decode is not a structure")
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		switch f.Type().Kind() {
		case reflect.Int64:
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
        case reflect.Struct:
            switch f.Interface().(type) {
            case time.Time:
                var t time.Time = f.Interface().(time.Time)
                buf.WriteString(t.Format(ibTimeFormat))
            default:
                return &EncodeError{Type: f.Type()}
            }
        default:
            return &EncodeError{Type: f.Type()}
		}
		buf.WriteByte(0)
	}

	return nil
}
