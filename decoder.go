package inflate

import (
	"encoding/json"
	"reflect"
)

const (
	// OptionSimple is the simple opt
	OptionSimple = "simple"
	// OptionForm is the form opt
	OptionForm = "form"
	// OptionLabel is the simple opt
	OptionLabel = "label"
	// OptionMatrix is the simple opt
	OptionMatrix = "matrix"
	// OptionExplode is the simple opt
	OptionExplode = "explode"
	// OptionDeepObject is the simple opt
	OptionDeepObject = "deep-object"
	// OptionSpaceDelimited is the space-delimited opt
	OptionSpaceDelimited = "space-delimited"
	// OptionPipeDelimited is the pipe-delimited opt
	OptionPipeDelimited = "pipe-delimited"
)

// Context is the context
type Context struct {
	Field  string
	Type   reflect.Type
	IsZero bool
	Tag    *Tag
}

//go:generate counterfeiter -fake-name ValueProvider -o ./fake/value_provider.go . ValueProvider

// ValueProvider provides a value
type ValueProvider interface {
	Value(ctx *Context) (interface{}, error)
}

//go:generate counterfeiter -fake-name ValueConverter -o ./fake/value_converter.go . ValueConverter

// ValueConverter converts source to target
type ValueConverter interface {
	Convert(source, target interface{}) error
}

// Decoder decodes the values from given source
type Decoder struct {
	TagName   string
	Provider  ValueProvider
	Converter ValueConverter
}

// Decode decodes the values to given target
func (d *Decoder) Decode(value interface{}) error {
	target, err := check("target", value)
	if err != nil {
		return err
	}

	return d.decode(StructOf(d.TagName, target))
}

func (d *Decoder) decode(ch *Struct) error {
	for _, field := range ch.Fields() {
		target := refer(field.Value)

		if field.Tag.Name == "~" {
			if target.Kind() == reflect.Struct {
				if err := d.decode(StructOf(d.TagName, target)); err != nil {
					return err
				}
			}

			set(field.Value, target)
			continue
		}

		ctx := &Context{
			Field:  field.Name,
			Tag:    field.Tag,
			Type:   target.Type(),
			IsZero: field.Value.IsZero(),
		}

		value, err := d.Provider.Value(ctx)
		if err != nil {
			return err
		}

		source := elem(reflect.ValueOf(value))

		if err := d.Converter.Convert(source, target); err != nil {
			return err
		}

		set(field.Value, target)
	}

	return nil
}

// Set sets the value
func Set(source, target interface{}) error {
	converter := &Converter{
		TagName: "field",
	}

	return converter.Convert(source, target)
}

// SetDefault set the default values
func SetDefault(target interface{}) error {
	decoder := &Decoder{
		TagName: "default",
		Converter: &Converter{
			TagName: "default",
		},
		Provider: &DefaultProvider{},
	}

	return decoder.Decode(target)
}

// DefaultProvider returns the default for given field
type DefaultProvider struct{}

// Value returns the default value if specified
func (p *DefaultProvider) Value(ctx *Context) (interface{}, error) {
	if !ctx.IsZero {
		return nil, nil
	}

	value := ctx.Tag.Name

	if canUnmarshalText(ctx.Type) {
		return value, nil
	}

	switch ctx.Type.Kind() {
	case reflect.Map, reflect.Struct:
		kv := make(map[string]interface{})

		if err := json.Unmarshal([]byte(value), &kv); err != nil {
			return nil, err
		}

		return kv, nil
	case reflect.Array, reflect.Slice:
		arr := []interface{}{}

		if err := json.Unmarshal([]byte(value), &arr); err != nil {
			return nil, err
		}

		return arr, nil
	default:
		return value, nil
	}
}
