package reflectify

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// NewQueryDecoder creates a path decoder
func NewQueryDecoder(query url.Values) *Decoder {
	return &Decoder{
		Tag: "query",
		Provider: &QueryProvider{
			Query: query,
		},
	}
}

var _ Provider = &QueryProvider{}

// QueryProvider represents a parameter provider that fetches values from
// incoming request's cookies
type QueryProvider struct {
	Query url.Values
}

// New returns a new provider
func (p *QueryProvider) New(value reflect.Value) Provider {
	return &ValueProvider{
		Var: value,
	}
}

// Value returns a primitive value
func (p *QueryProvider) Value(ctx *Context) (interface{}, error) {
	if ctx.Options.IsEmpty() {
		ctx.Options = append(ctx.Options, OptionForm.String())
		ctx.Options = append(ctx.Options, OptionExplode.String())
	}

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

func (p *QueryProvider) valueOf(ctx *Context) (interface{}, error) {
	values := p.queryArray(ctx.Field)
	if values == nil || len(values) == 0 {
		return nil, nil
	}

	switch {
	case ctx.Options.Has(OptionForm):
		return values[0], nil
	case ctx.Options.Has(OptionSpaceDelimited):
		return nil, p.notSupported(ctx, OptionSpaceDelimited)
	case ctx.Options.Has(OptionPipeDelimited):
		return nil, p.notSupported(ctx, OptionPipeDelimited)
	case ctx.Options.Has(OptionDeepObject):
		return nil, p.notSupported(ctx, OptionDeepObject)
	default:
		return nil, p.notProvided(ctx,
			OptionForm,
			OptionSpaceDelimited,
			OptionDeepObject,
		)
	}
}

func (p *QueryProvider) arrayOf(ctx *Context) ([]interface{}, error) {
	values := p.queryArray(ctx.Field)

	if values == nil || len(values) == 0 {
		return nil, nil
	}

	separator := ""

	switch {
	case ctx.Options.Has(OptionForm):
		separator = ","
	case ctx.Options.Has(OptionSpaceDelimited):
		separator = " "
	case ctx.Options.Has(OptionPipeDelimited):
		separator = "|"
	case ctx.Options.Has(OptionDeepObject):
		return nil, p.notSupported(ctx, OptionDeepObject)
	default:
		return nil, p.notProvided(ctx,
			OptionForm,
			OptionSpaceDelimited,
			OptionPipeDelimited,
		)
	}

	if !ctx.Options.Has(OptionExplode) {
		values = strings.Split(values[0], separator)
	}

	return convertArray(values), nil
}

func (p *QueryProvider) mapOf(ctx *Context) (map[string]interface{}, error) {
	switch {
	case ctx.Options.Has(OptionForm):
		if ctx.Options.Has(OptionExplode) {
			return p.queryMap(), nil
		}

		values := p.queryArray(ctx.Field)

		if values == nil || len(values) == 0 {
			return nil, nil
		}

		values = strings.Split(values[0], ",")
		return convertMap(values)
	case ctx.Options.Has(OptionSpaceDelimited):
		return nil, p.notSupported(ctx, OptionSpaceDelimited)
	case ctx.Options.Has(OptionPipeDelimited):
		return nil, p.notSupported(ctx, OptionPipeDelimited)
	case ctx.Options.Has(OptionDeepObject):
		if ctx.Options.Has(OptionExplode) {
			return nil, p.notSupported(ctx, OptionExplode)
		}

		var (
			kv     = p.queryMap()
			result = make(map[string]interface{})
			prefix = ctx.Field + "["
			suffix = "]"
		)

		for k, v := range kv {
			key := k

			if !strings.HasPrefix(k, prefix) {
				continue
			}

			key = strings.TrimPrefix(k, prefix)

			if !strings.HasSuffix(k, suffix) {
				continue
			}

			key = strings.TrimSuffix(key, suffix)

			if key = strings.TrimSpace(key); key != "" {
				result[key] = v
			}
		}

		return result, nil
	default:
		return nil, p.notProvided(ctx,
			OptionForm,
			OptionSpaceDelimited,
			OptionPipeDelimited,
		)
	}
}

func (p *QueryProvider) queryArray(key string) []string {
	if values, ok := p.Query[key]; ok {
		return values
	}

	return nil
}

func (p *QueryProvider) queryMap() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range p.Query {
		if len(v) > 0 {
			m[k] = v[0]
		} else {
			m[k] = nil
		}
	}

	return m
}

func (p *QueryProvider) notProvided(ctx *Context, opts ...Option) error {
	return p.errorf("field: '%v' option: %v not provided", ctx.Field, opts)
}

func (p *QueryProvider) notSupported(ctx *Context, opt Option) error {
	return p.errorf("field: '%v' option: [%v] not supported", ctx.Field, opt)
}

func (p *QueryProvider) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("query: %s", msg)
}
