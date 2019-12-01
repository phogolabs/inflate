package reflectify

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

// Option represents a tag option
type Option string

const (
	// OptionSimple is the simple opt
	OptionSimple Option = "simple"
	// OptionForm is the form opt
	OptionForm Option = "form"
	// OptionLabel is the simple opt
	OptionLabel Option = "label"
	// OptionMatrix is the simple opt
	OptionMatrix Option = "matrix"
	// OptionExplode is the simple opt
	OptionExplode Option = "explode"
	// OptionDeepObject is the simple opt
	OptionDeepObject Option = "deep-object"
	// OptionSpaceDelimited is the space-delimited opt
	OptionSpaceDelimited Option = "space-delimited"
	// OptionPipeDelimited is the pipe-delimited opt
	OptionPipeDelimited Option = "pipe-delimited"
)

// Options represents a list of options
type Options []string

// IsEmpty returns true if the options are empty
func (opts Options) IsEmpty() bool {
	return len(opts) == 0
}

// Has returns true if the option is available
func (opts Options) Has(opt Option) bool {
	for _, key := range opts {
		if strings.EqualFold(string(opt), key) {
			return true
		}
	}

	return false
}

// Field represents a field
type Field struct {
	Name    string
	Options Options
	Value   reflect.Value
}

// String returns the field as string
func (f *Field) String() string {
	return f.Name
}

// Kind returns the field kind
func (f *Field) Kind() reflect.Kind {
	return kind(f.Value)
}

// NewField creates a new field
func NewField(value string) *Field {
	parts := strings.Split(value, ",")

	return &Field{
		Name:    parts[0],
		Options: Options(parts[1:]),
	}
}

func convertMap(parts []string) (map[string]interface{}, error) {
	var (
		count  = len(parts)
		result = make(map[string]interface{})
	)

	if count%2 != 0 {
		return nil, fmt.Errorf("object value: %s invalid", parts)
	}

	for index := 1; index < count; index = index + 2 {
		prev := index - 1

		var (
			key   = parts[prev]
			value = parts[index]
		)

		if key == "" {
			return nil, fmt.Errorf("object value: %s invalid", parts)
		}

		result[key] = value
	}

	return result, nil
}

func explodeMap(parts []string) (map[string]interface{}, error) {
	var (
		count  = len(parts)
		result = make(map[string]interface{})
	)

	for index := 0; index < count; index++ {
		kv := strings.SplitN(parts[index], "=", 2)

		var (
			key   string
			value string
		)

		switch {
		case len(kv) > 1:
			key = kv[0]
			value = kv[1]
		case len(kv) > 0:
			key = kv[0]
		default:
			return nil, fmt.Errorf("object value: %s invalid", parts[index])
		}

		if key == "" {
			return nil, fmt.Errorf("object value: %s invalid", parts)
		}

		result[key] = value
	}

	return result, nil
}

func convertValue(values []interface{}) interface{} {
	if len(values) == 1 {
		return values[0]
	}

	return values
}

func convertArray(array []string) []interface{} {
	result := make([]interface{}, len(array))

	for index, item := range array {
		result[index] = item
	}

	return result
}

func kind(v reflect.Value) reflect.Kind {
	kind := v.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

func tryTextUnmarshaller(v reflect.Value) encoding.TextUnmarshaler {
	var (
		unmarshaller     encoding.TextUnmarshaler
		unmarshallerFrom = func(value reflect.Value) encoding.TextUnmarshaler {
			if unmarshaller, ok := value.Interface().(encoding.TextUnmarshaler); ok {
				return unmarshaller
			}

			return nil
		}
	)

	unmarshaller = unmarshallerFrom(v)

	if v.CanAddr() {
		unmarshaller = unmarshallerFrom(v.Addr())
	}

	return unmarshaller
}

func tryTextMarshaller(v reflect.Value) encoding.TextMarshaler {
	var (
		marshaller     encoding.TextMarshaler
		marshallerFrom = func(value reflect.Value) encoding.TextMarshaler {
			if unmarshaller, ok := value.Interface().(encoding.TextMarshaler); ok {
				return unmarshaller
			}

			return nil
		}
	)

	marshaller = marshallerFrom(v)

	if v.CanAddr() {
		marshaller = marshallerFrom(v.Addr())
	}

	return marshaller
}
