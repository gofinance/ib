package trade

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

//const ibTimeFormat = "20060102 15:04:05.000000 MST"
const ibTimeFormat = "20060102 15:04:05 MST"

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
	return fmt.Sprintf("trade: cannot encode field %s of type %v with value %v",
		e.Name, e.Type, e.Value)
}

func failEncode(n string, v interface{}, t reflect.Type) error {
	return &EncodeError{
		Name:  n,
		Value: v,
		Type:  t,
	}
}

func skipField(v reflect.Value, f reflect.StructField) bool {
	// no tag on the structure field
	name := f.Tag.Get("when")
	if name == "" {
		return false
	}
	// invalid field in tag
	f1 := v.FieldByName(name)
	if !f1.IsValid() {
		return false
	}
	// target field is not a string
	if f1.Type().Kind() != reflect.String {
		return false
	}

	target := f1.String()
	value := f.Tag.Get("value")

	switch f.Tag.Get("cond") {
	case "is":
		// skip when target equals value
		return (target == value)
	case "not":
		return (target != value)
	}

	return false
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var s string

	if !v.IsValid() {
		v = reflect.ValueOf(0)
	}

	switch v.Type().Kind() {
	case reflect.Int, reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)
	case reflect.Float64:
		s = strconv.FormatFloat(v.Float(), 'g', 10, 64)
	case reflect.String:
		s = string(v.String())
	case reflect.Bool:
		if v.Bool() {
			s = "1"
		} else {
			s = "0"
		}
	case reflect.Slice:
		// encode size
		type size struct {
			N int64
		}
		sz := size{int64(v.Len())}
		if err := encode(buf, reflect.ValueOf(&sz)); err != nil {
			return err
		}
		// encode elements
		for i := 0; i < v.Len(); i++ {
			if err := encode(buf, v.Index(i)); err != nil {
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
				if skipField(v, v.Type().Field(i)) {
					continue
				}
				if err := encode(buf, f); err != nil {
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
	return fmt.Sprintf("trade: cannot decode '%v' into field %s of type %v",
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
		type size struct {
			N int64
		}

		sz := size{}

		if err := decode(b, reflect.ValueOf(&sz)); err != nil {
			return err
		}

		n := int(sz.N)
		v.Set(reflect.MakeSlice(v.Type(), n, n))

		// decode elements
		for i := 0; i < n; i++ {
			if err := decode(b, v.Index(i)); err != nil {
				return err
			}
		}

		return nil
	case reflect.Struct:
		switch v.Interface().(type) {
		case time.Time:
			// do nothing
		default:
			// decode fields
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				if skipField(v, v.Type().Field(i)) {
					// string field we depend on is empty
					continue
				}
				if err := decode(b, f); err != nil {
					fmt.Printf("Failed to decode structure %v, error %s\n", v, err)
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
