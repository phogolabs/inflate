package inflate

import (
	"fmt"
	"net/http"
	"net/textproto"
	"reflect"
	"strings"
)

// NewHeaderDecoder creates a header decoder
func NewHeaderDecoder(header http.Header) *Decoder {
	return &Decoder{
		TagName: "path",
		Converter: &Converter{
			TagName: "path",
		},
		Provider: &HeaderProvider{
			Header: header,
		},
	}
}

var _ ValueProvider = &HeaderProvider{}

// HeaderProvider represents a parameter provider that fetches values from
// incoming request's header
type HeaderProvider struct {
	Header http.Header
}

// Value returns a primitive value
func (p *HeaderProvider) Value(ctx *Context) (interface{}, error) {
	if len(ctx.Tag.Options) == 0 {
		ctx.Tag.AddOption(OptionSimple)
	}

	switch ctx.Kind {
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

func (p *HeaderProvider) valueOf(ctx *Context) (interface{}, error) {
	header := p.header(ctx.Tag.Name)

	if header == nil {
		return nil, nil
	}

	if !ctx.Tag.HasOption(OptionSimple) {
		return nil, p.notProvided(ctx, OptionSimple)
	}

	return *header, nil
}

func (p *HeaderProvider) arrayOf(ctx *Context) ([]interface{}, error) {
	header := p.header(ctx.Tag.Name)

	if header == nil {
		return nil, nil
	}

	if !ctx.Tag.HasOption(OptionSimple) {
		return nil, p.notProvided(ctx, OptionSimple)
	}

	var (
		separator = ","
		parts     = strings.Split(*header, separator)
		result    = make([]interface{}, len(parts))
	)

	for index, part := range parts {
		result[index] = part
	}

	return result, nil
}

func (p *HeaderProvider) mapOf(ctx *Context) (m map[string]interface{}, err error) {
	header := p.header(ctx.Tag.Name)

	if header == nil {
		return nil, nil
	}

	if !ctx.Tag.HasOption(OptionSimple) {
		return nil, p.notProvided(ctx, OptionSimple)
	}

	var (
		separator = ","
		parts     = strings.Split(*header, separator)
	)

	if ctx.Tag.HasOption(OptionExplode) {
		m, err = explodeMap(parts)
	} else {
		m, err = convertMap(parts)
	}

	if err != nil {
		return nil, p.errorf(err.Error())
	}

	return m, err
}

func (p *HeaderProvider) header(name string) *string {
	key := textproto.CanonicalMIMEHeaderKey(name)

	if _, ok := p.Header[key]; ok {
		value := p.Header.Get(name)
		return &value
	}

	return nil
}

func (p *HeaderProvider) notProvided(ctx *Context, opts ...string) error {
	return p.errorf("field '%v' option: %v not provided", ctx.Tag.Name, opts)
}

func (p *HeaderProvider) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("header: %s", msg)
}
