package reflectify

import (
	"fmt"
	"reflect"
	"strings"
)

// Encoding represents an reflectify option
type Encoding string

const (
	// EncodingText is the text reflectify
	EncodingText Encoding = "text"
)

// Encodings represents a list of options
type Encodings []string

// IsEmpty returns true if the options are empty
func (opts Encodings) IsEmpty() bool {
	return len(opts) == 0
}

// Has returns true if the option is available
func (opts Encodings) Has(opt Encoding) bool {
	for _, key := range opts {
		if strings.EqualFold(string(opt), key) {
			return true
		}
	}

	return false
}

// Context represents a provider context
type Context struct {
	Field      string
	FieldTag   string
	FieldKind  reflect.Kind
	FieldEmpty bool
	Options    Options
	Encoding   Encodings
}

//go:generate counterfeiter -fake-name Provider -o ./fake/provider.go . Provider

// Provider represents the decoder data provider
type Provider interface {
	// New returns a new provider
	New(value reflect.Value) Provider
	// Value returns a primitive value
	Value(ctx *Context) (interface{}, error)
}

// NewValueDecoder creates a default decoder
func NewValueDecoder(source interface{}) *Decoder {
	return &Decoder{
		Tag: "field",
		Provider: &ValueProvider{
			Var: reflect.ValueOf(source).Elem(),
		},
	}
}

// Set sets the values for given type
func Set(target, source interface{}) error {
	decoder := NewValueDecoder(source)
	return decoder.Decode(target)
}

var _ Provider = &ValueProvider{}

// ValueProvider provides a value from a field
type ValueProvider struct {
	Var reflect.Value
}

// New returns a new provider
func (p *ValueProvider) New(value reflect.Value) Provider {
	return &ValueProvider{
		Var: reflect.Indirect(value),
	}
}

// Value returns a primitive value
func (p *ValueProvider) Value(ctx *Context) (interface{}, error) {
	var (
		target     = p.Var
		targetType = target.Type()
		targetKind = target.Kind()
	)

	if targetKind == reflect.Struct {
		for index := 0; index < targetType.NumField(); index++ {
			item := targetType.Field(index)

			if item.PkgPath != "" {
				continue
			}

			field := NewField(item.Tag.Get(ctx.FieldTag))

			if field.Name == "" {
				continue
			}

			if field.Name == "-" {
				continue
			}

			if field.Name != ctx.Field {
				continue
			}

			value := target.FieldByIndex([]int{index})
			return value.Interface(), nil
		}
	}

	return p.Var.Interface(), nil
}

func (p *ValueProvider) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("field: %s", msg)
}

// NewDefaultDecoder creates a default decoder
func NewDefaultDecoder() *Decoder {
	return &Decoder{
		Tag:      "default",
		Provider: &DefaultProvider{},
	}
}

// SetDefaults sets the default values for given type
func SetDefaults(target interface{}) error {
	decoder := NewDefaultDecoder()
	return decoder.Decode(target)
}

// DefaultProvider provides a default value
type DefaultProvider struct{}

// New returns a new provider
func (p *DefaultProvider) New(value reflect.Value) Provider {
	return p
}

// Value returns a primitive value
func (p *DefaultProvider) Value(ctx *Context) (interface{}, error) {
	if ctx.FieldEmpty {
		return ctx.Field, nil
	}

	return nil, nil
}
