package reflectify

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
)

// NewPathDecoder creates a path decoder
func NewPathDecoder(r *chi.RouteParams) *Decoder {
	return &Decoder{
		Tag: "path",
		Provider: &PathProvider{
			Param: r,
		},
	}
}

var _ Provider = &PathProvider{}

// PathProvider represents a parameter provider that fetches values from
// incoming request's header
type PathProvider struct {
	Param *chi.RouteParams
}

// New returns a new provider
func (p *PathProvider) New(value reflect.Value) Provider {
	return &ValueProvider{
		Var: value,
	}
}

// Value returns a primitive value
func (p *PathProvider) Value(ctx *Context) (interface{}, error) {
	if ctx.Encoding.Has(EncodingText) {
		return p.valueOf(ctx)
	}

	switch ctx.FieldKind {
	case reflect.Map, reflect.Struct:
		return p.mapOf(ctx)
	case reflect.Array, reflect.Slice:
		values, err := p.arrayOf(ctx)
		if err != nil {
			return nil, err
		}

		return convertValue(values), nil
	default:
		return p.valueOf(ctx)
	}
}

func (p *PathProvider) valueOf(ctx *Context) (interface{}, error) {
	param := p.param(ctx.Field)

	if param == nil {
		return nil, nil
	}

	switch {
	case ctx.Options.IsEmpty(), ctx.Options.Has(OptionSimple):
		return *param, nil
	case ctx.Options.Has(OptionLabel):
		prefix := "."
		return strings.TrimPrefix(*param, prefix), nil
	case ctx.Options.Has(OptionMatrix):
		prefix := fmt.Sprintf(";%s=", ctx.Field)
		return strings.TrimPrefix(*param, prefix), nil
	default:
		return nil, p.notProvided(ctx,
			OptionSimple,
			OptionLabel,
			OptionMatrix,
		)
	}
}

func (p *PathProvider) arrayOf(ctx *Context) ([]interface{}, error) {
	param := p.param(ctx.Field)

	if param == nil {
		return nil, nil
	}

	var (
		prefix    = ""
		separator = ""
	)

	switch {
	case ctx.Options.IsEmpty(), ctx.Options.Has(OptionSimple):
		separator = ","
	case ctx.Options.Has(OptionLabel):
		prefix = "."
		separator = ","

		if ctx.Options.Has(OptionExplode) {
			separator = prefix
		}
	case ctx.Options.Has(OptionMatrix):
		separator = ","
		prefix = fmt.Sprintf(";%s=", ctx.Field)

		if ctx.Options.Has(OptionExplode) {
			separator = prefix
		}
	default:
		return nil, p.notProvided(ctx,
			OptionSimple,
			OptionLabel,
			OptionMatrix,
		)
	}

	var (
		value  = strings.TrimPrefix(*param, prefix)
		parts  = strings.Split(value, separator)
		result = make([]interface{}, len(parts))
	)

	for index, part := range parts {
		result[index] = part
	}

	return result, nil
}

func (p *PathProvider) mapOf(ctx *Context) (m map[string]interface{}, err error) {
	param := p.param(ctx.Field)

	if param == nil {
		return nil, nil
	}

	var (
		separator string
		prefix    string
	)

	switch {
	case ctx.Options.IsEmpty(), ctx.Options.Has(OptionSimple):
		separator = ","
	case ctx.Options.Has(OptionLabel):
		prefix = "."
		separator = ","

		if ctx.Options.Has(OptionExplode) {
			separator = prefix
		}
	case ctx.Options.Has(OptionMatrix):
		separator = ","
		prefix = fmt.Sprintf(";%s=", ctx.Field)

		if ctx.Options.Has(OptionExplode) {
			separator = ";"
			prefix = ";"
		}
	default:
		return nil, p.notProvided(ctx,
			OptionSimple,
			OptionLabel,
			OptionMatrix,
		)
	}

	var (
		value = strings.TrimPrefix(*param, prefix)
		parts = strings.Split(value, separator)
	)

	if ctx.Options.Has(OptionExplode) {
		m, err = explodeMap(parts)
	} else {
		m, err = convertMap(parts)
	}

	if err != nil {
		return nil, p.errorf(err.Error())
	}

	return m, nil
}

func (p *PathProvider) param(name string) *string {
	for index, k := range p.Param.Keys {
		if strings.EqualFold(k, name) {
			return &p.Param.Values[index]
		}
	}

	return nil
}

func (p *PathProvider) notProvided(ctx *Context, opts ...Option) error {
	return p.errorf("field: %v option: %v not provided", ctx.Field, opts)
}

func (p *PathProvider) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("path: %s", msg)
}
