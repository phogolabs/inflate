package reflectify

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// NewCookieDecoder creates a cookie decoder
func NewCookieDecoder(cookies []*http.Cookie) *Decoder {
	return &Decoder{
		Tag: "cookie",
		Provider: &CookieProvider{
			Cookies: cookies,
		},
	}
}

var _ Provider = &CookieProvider{}

// CookieProvider represents a parameter provider that fetches values from
// incoming request's cookies
type CookieProvider struct {
	Cookies []*http.Cookie
}

// New returns a new provider
func (p *CookieProvider) New(value reflect.Value) Provider {
	return &ValueProvider{
		Var: value,
	}
}

// Value returns a primitive value
func (p *CookieProvider) Value(ctx *Context) (interface{}, error) {
	if ctx.Options.IsEmpty() {
		ctx.Options = append(ctx.Options, OptionForm.String())

		switch ctx.FieldKind {
		case reflect.Map, reflect.Struct:
		case reflect.Array, reflect.Slice:
		default:
			ctx.Options = append(ctx.Options, OptionExplode.String())
		}
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

func (p *CookieProvider) valueOf(ctx *Context) (interface{}, error) {
	cookie := p.cookie(ctx.Field)

	if cookie == nil {
		return nil, nil
	}

	if !ctx.Options.Has(OptionForm) {
		return nil, p.notProvided(ctx, OptionForm)
	}

	return cookie.Value, nil
}

func (p *CookieProvider) arrayOf(ctx *Context) ([]interface{}, error) {
	cookie := p.cookie(ctx.Field)

	if cookie == nil {
		return nil, nil
	}

	if !ctx.Options.Has(OptionForm) {
		return nil, p.notProvided(ctx, OptionForm)
	}

	if ctx.Options.Has(OptionExplode) {
		return nil, p.notSupported(ctx, OptionExplode)
	}

	var (
		separator = ","
		parts     = strings.Split(cookie.Value, separator)
		result    = make([]interface{}, len(parts))
	)

	for index, part := range parts {
		result[index] = part
	}

	return result, nil
}

func (p *CookieProvider) mapOf(ctx *Context) (map[string]interface{}, error) {
	cookie := p.cookie(ctx.Field)

	if cookie == nil {
		return nil, nil
	}

	if !ctx.Options.Has(OptionForm) {
		return nil, p.notProvided(ctx, OptionForm)
	}

	if ctx.Options.Has(OptionExplode) {
		return nil, p.notSupported(ctx, OptionExplode)
	}

	var (
		separator = ","
		parts     = strings.Split(cookie.Value, separator)
	)

	m, err := convertMap(parts)
	if err != nil {
		return nil, p.errorf(err.Error())
	}

	return m, nil
}

func (p *CookieProvider) cookie(name string) *http.Cookie {
	for _, cookie := range p.Cookies {
		if strings.EqualFold(cookie.Name, name) {
			return cookie
		}
	}

	return nil
}

func (p *CookieProvider) notProvided(ctx *Context, opts ...Option) error {
	return p.errorf("field: '%v' option: %v not provided", ctx.Field, opts)
}

func (p *CookieProvider) notSupported(ctx *Context, opt Option) error {
	return p.errorf("field: '%v' option: [%v] not supported", ctx.Field, opt)
}

func (p *CookieProvider) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("cookie: %s", msg)
}
