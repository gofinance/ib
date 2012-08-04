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

func encode(buf *bytes.Buffer, tag int, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var s string

	switch v.Type().Kind() {
	case reflect.Int, reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)
	case reflect.Float64:
		s = strconv.FormatFloat(v.Float(), 'g', 10, 64)
	case reflect.String:
		s = string(v.String())
	case reflect.Bool:
		if v.Int() > 0 {
			s = "1"
		} else {
			s = "0"
		}
	case reflect.Slice:
		// encode size
		if err := encode(buf, 0, reflect.ValueOf(v.Len())); err != nil {
			return err
		}
		// encode elements
		for i := 0; i < v.Len(); i++ {
			if err := encode(buf, 0, v.Index(i)); err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		switch v.Interface().(type) {
		// custom encoding for time
		case time.Time:
			var t time.Time = v.Interface().(time.Time)
			s = t.Format(ibTimeFormat)
		default:
			// encode fields
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				if err := encode(buf, 0, f); err != nil {
					return err
				}
			}
		}
		if s == "" {
			return nil
		}
	default:
		return failEncode(v.Type().Name(), v.Interface(), v.Type())
	}

	if _, err := buf.WriteString(s + "\000"); err != nil {
		return err
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

func decode(b *bufio.Reader, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	kind := v.Type().Kind()

	// special processing for struct and slice

	switch kind {
	case reflect.Slice:
		/*		
			case reflect.Slice:
				// encode size
				if err := encode(buf, 0, v.Len()); err != nil {
					return err
				}
				// encode elements
				for i := 0; i < v.Len(); i++ {
					if err := encode(buf, 0, v.Index(i).Interface()); err != nil {
						return err
					}
				}
				return nil
		*/	case reflect.Struct:
		switch v.Interface().(type) {
		case time.Time:
			// do nothing
		default:
			// decode fields
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				if err := decode(b, f); err != nil {
					return err
				}
			}

			return nil
		}
	}

	var s string

	if data, err := b.ReadString(0); err != nil {
		return err
	} else {
		s = string(data[:len(data)-1]) // trim \000
	}

	switch kind {
	case reflect.Int64:
		if x, err := strconv.ParseInt(s, 10, 64); err != nil {
			return failDecode(s, v.Type().Name(), v.Type())
		} else {
			v.SetInt(x)
		}
	case reflect.Float64:
		if x, err := strconv.ParseFloat(s, 64); err != nil {
			return failDecode(s, v.Type().Name(), v.Type())
		} else {
			v.SetFloat(x)
		}
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		if x, err := strconv.ParseInt(s, 10, 64); err != nil {
			return failDecode(s, v.Type().Name(), v.Type())
		} else {
			v.SetBool(x > 0)
		}
	case reflect.Struct:
		switch v.Interface().(type) {
		case time.Time:
			if x, err := time.Parse(ibTimeFormat, s); err != nil {
				return failDecode(s, v.Type().Name(), v.Type())
			} else {
				v.Set(reflect.ValueOf(x))
			}
		default:
			return failDecode(s, v.Type().Name(), v.Type())
		}
	default:
		return failDecode(s, v.Type().Name(), v.Type())
	}

	return nil
}
