package inflate

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"text/scanner"
)

// NewQueryDecoder creates a path decoder
func NewQueryDecoder(query url.Values) *Decoder {
	return &Decoder{
		TagName: "query",
		Converter: &Converter{
			TagName: "query",
		},
		Provider: &QueryProvider{
			Query: query,
		},
	}
}

// NewFormDecoder creates a path decoder
func NewFormDecoder(query url.Values) *Decoder {
	return &Decoder{
		TagName: "form",
		Converter: &Converter{
			TagName: "form",
		},
		Provider: &QueryProvider{
			Query: query,
		},
	}
}

var _ ValueProvider = &QueryProvider{}

// QueryProvider represents a parameter provider that fetches values from
// incoming request's cookies
type QueryProvider struct {
	Query url.Values
}

// Value returns a primitive value
func (p *QueryProvider) Value(ctx *Context) (interface{}, error) {
	if len(ctx.Tag.Options) == 0 {
		ctx.Tag.AddOption(OptionForm)
		ctx.Tag.AddOption(OptionExplode)
	}

	if canUnmarshalText(ctx.Type) {
		return p.valueOf(ctx)
	}

	switch ctx.Type.Kind() {
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
	values := p.queryArray(ctx.Tag.Name)
	if values == nil || len(values) == 0 {
		return nil, nil
	}

	switch {
	case ctx.Tag.HasOption(OptionForm):
		return values[0], nil
	case ctx.Tag.HasOption(OptionSpaceDelimited):
		return nil, p.notSupported(ctx, OptionSpaceDelimited)
	case ctx.Tag.HasOption(OptionPipeDelimited):
		return nil, p.notSupported(ctx, OptionPipeDelimited)
	case ctx.Tag.HasOption(OptionDeepObject):
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
	values := p.queryArray(ctx.Tag.Name)

	if values == nil || len(values) == 0 {
		return nil, nil
	}

	separator := ""

	switch {
	case ctx.Tag.HasOption(OptionForm):
		separator = ","
	case ctx.Tag.HasOption(OptionSpaceDelimited):
		separator = " "
	case ctx.Tag.HasOption(OptionPipeDelimited):
		separator = "|"
	case ctx.Tag.HasOption(OptionDeepObject):
		return nil, p.notSupported(ctx, OptionDeepObject)
	default:
		return nil, p.notProvided(ctx,
			OptionForm,
			OptionSpaceDelimited,
			OptionPipeDelimited,
		)
	}

	if !ctx.Tag.HasOption(OptionExplode) {
		values = strings.Split(values[0], separator)
	}

	return convertArray(values), nil
}

func (p *QueryProvider) mapOf(ctx *Context) (map[string]interface{}, error) {
	switch {
	case ctx.Tag.HasOption(OptionForm):
		if ctx.Tag.HasOption(OptionExplode) {
			return p.queryMap(), nil
		}

		values := p.queryArray(ctx.Tag.Name)

		if values == nil || len(values) == 0 {
			return nil, nil
		}

		values = strings.Split(values[0], ",")
		return convertMap(values)
	case ctx.Tag.HasOption(OptionSpaceDelimited):
		return nil, p.notSupported(ctx, OptionSpaceDelimited)
	case ctx.Tag.HasOption(OptionPipeDelimited):
		return nil, p.notSupported(ctx, OptionPipeDelimited)
	case ctx.Tag.HasOption(OptionDeepObject):
		if ctx.Tag.HasOption(OptionExplode) {
			return nil, p.notSupported(ctx, OptionExplode)
		}

		return p.deepObject(ctx)
	default:
		return nil, p.notProvided(ctx,
			OptionForm,
			OptionSpaceDelimited,
			OptionPipeDelimited,
		)
	}
}

func (p *QueryProvider) deepObject(ctx *Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for k, v := range p.queryMap() {
		k = strings.TrimPrefix(k, ctx.Tag.Name)

		var (
			values    = result
			keys, err = p.path(k)
			count     = len(keys)
		)

		if err != nil {
			return nil, p.notParsed(ctx, err)
		}

		for index, key := range keys {
			if index == count-1 {
				values[key] = v
				continue
			}

			next, ok := values[key].(map[string]interface{})

			if !ok {
				next = make(map[string]interface{})
				values[key] = next
			}

			values = next
		}
	}

	return result, nil
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
		}
	}

	return m
}

func (p *QueryProvider) path(k string) ([]string, error) {
	iter := &scanner.Scanner{}
	iter.Init(bytes.NewBufferString(k))

	var (
		result = []string{}
		err    = p.errorf("cannot parse key: %s", k)
	)

	const (
		left  = '['
		right = ']'
	)

	started := false

	for {
		token := iter.Scan()

		switch token {
		case left:
			if started {
				return nil, err
			}

			started = true
			continue
		case right:
			if !started {
				return nil, err
			}

			started = false
			continue
		case scanner.EOF:
			if started {
				return nil, err
			}

			return result, nil
		default:
			if !started {
				return nil, err
			}
		}

		result = append(result, iter.TokenText())
	}
}

func (p *QueryProvider) notProvided(ctx *Context, opts ...string) error {
	return p.errorf("field: '%v' option: %v not provided", ctx.Tag.Name, opts)
}

func (p *QueryProvider) notSupported(ctx *Context, opt string) error {
	return p.errorf("field: '%v' option: [%v] not supported", ctx.Tag.Name, opt)
}

func (p *QueryProvider) notParsed(ctx *Context, err error) error {
	return p.errorf("field: '%v' not parsed: %v", ctx.Tag.Name, err)
}

func (p *QueryProvider) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("query: %s", msg)
}
