package wire

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

const ibTimeFormat = "20060102 15:04:05.000000 MST"

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

type EncodeError struct {
	Name  string
	Value interface{}
	Type  reflect.Type
}

func (e *EncodeError) Error() string {
	return fmt.Sprintf("ibtws: cannot encode field %s of type %v with value %v",
		e.Name, e.Type, e.Value)
}

func failEncode(n string, v interface{}, t reflect.Type) error {
	return &EncodeError{
		Name:  n,
		Value: v,
		Type:  t,
	}
}

func encode(buf *bytes.Buffer, tag int, v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	buf.Reset()

	if val.Kind() != reflect.Struct {
		panic("ibtws: value given to encode is not a structure")
	}

	var s string

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		switch f.Type().Kind() {
		case reflect.Int64:
			s = strconv.FormatInt(f.Int(), 10)
		case reflect.Float64:
			s = strconv.FormatFloat(f.Float(), 'g', 10, 64)
		case reflect.String:
			s = string(f.String())
		case reflect.Bool:
			if f.Int() > 0 {
				s = "1"
			} else {
				s = "0"
			}
		case reflect.Struct:
			switch f.Interface().(type) {
			case time.Time:
				var t time.Time = f.Interface().(time.Time)
				s = t.Format(ibTimeFormat)
			default:
				if err := encode(buf, 0, f.Interface()); err != nil {
					return err
				}
				continue
			}
		default:
			return failEncode(f.Type().Field(i).Name, f.Interface(), f.Type())
		}

		if _, err := buf.WriteString(s + "\000"); err != nil {
			return err
		}
	}

	return nil
}

type DecodeError struct {
	Data string
	Name string
	Type reflect.Type
}

func (e *DecodeError) Error() string {
        return fmt.Sprintf("ibtws: cannot decode '%v' into field %s of type %v",
        e.Data, e.Name, e.Type)
}

func failDecode(d string, n string, t reflect.Type) error {
	return &DecodeError{
		Data: d,
		Name: n,
		Type: t,
	}
}

func decode(b *bufio.Reader, v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("ibtws: value given to decode is not a structure")
	}

	// parse and set field

	var s string

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)

		if data, err := b.ReadString(0); err != nil {
			return err
		} else {
			// don't want the trailing \000
			s = string(data[:len(data)-1])
		}
		switch f.Type().Kind() {
		case reflect.Int64:
			if x, err := strconv.ParseInt(s, 10, 64); err != nil {
				return failDecode(s, f.Type().Field(i).Name, f.Type())
			} else {
				f.SetInt(x)
			}
		case reflect.Float64:
			if x, err := strconv.ParseFloat(s, 64); err != nil {
				return failDecode(s, f.Type().Field(i).Name, f.Type())
			} else {
				f.SetFloat(x)
			}
		case reflect.String:
			f.SetString(s)
		case reflect.Bool:
			if x, err := strconv.ParseInt(s, 10, 64); err != nil {
				return failDecode(s, f.Type().Field(i).Name, f.Type())
			} else {
				if x > 0 {
					f.SetBool(true)
				} else {
					f.SetBool(false)
				}
			}
		case reflect.Struct:
			switch f.Interface().(type) {
			case time.Time:
				if x, err := time.Parse(ibTimeFormat, s); err != nil {
					return failDecode(s, f.Type().Field(i).Name, f.Type())
				} else {
					f.Set(reflect.ValueOf(x))
				}
			default:
				if err := decode(b, f.Interface()); err != nil {
					return err
				}
			}
		default:
			failDecode(s, f.Type().Field(i).Name, f.Type())
		}
	}

	return nil
}
