package inflate

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// CookieProvider represents a parameter provider that fetches values from
// incoming request's cookies
type CookieProvider struct {
	Cookies []*http.Cookie
}

// NewCookieDecoder creates a cookie decoder
func NewCookieDecoder(cookies []*http.Cookie) *Decoder {
	return &Decoder{
		TagName: "path",
		Converter: &Converter{
			TagName: "path",
		},
		Provider: &CookieProvider{
			Cookies: cookies,
		},
	}
}

// Value returns a primitive value
func (p *CookieProvider) Value(ctx *Context) (interface{}, error) {
	if len(ctx.Tag.Options) == 0 {
		ctx.Tag.AddOption(OptionForm)

		switch ctx.Kind {
		case reflect.Map, reflect.Struct:
		case reflect.Array, reflect.Slice:
		default:
			ctx.Tag.AddOption(OptionExplode)
		}
	}

	switch ctx.Kind {
	case reflect.Map, reflect.Struct:
		return p.mapOf(ctx)
	case reflect.Array, reflect.Slice:
		return p.arrayOf(ctx)
	default:
		return p.valueOf(ctx)
	}
}

func (p *CookieProvider) valueOf(ctx *Context) (interface{}, error) {
	cookie := p.cookie(ctx.Tag.Name)

	if cookie == nil {
		return nil, nil
	}

	if !ctx.Tag.HasOption(OptionForm) {
		return nil, p.notProvided(ctx, OptionForm)
	}

	return cookie.Value, nil
}

func (p *CookieProvider) arrayOf(ctx *Context) ([]interface{}, error) {
	cookie := p.cookie(ctx.Tag.Name)

	if cookie == nil {
		return nil, nil
	}

	if !ctx.Tag.HasOption(OptionForm) {
		return nil, p.notProvided(ctx, OptionForm)
	}

	if ctx.Tag.HasOption(OptionExplode) {
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
	cookie := p.cookie(ctx.Tag.Name)

	if cookie == nil {
		return nil, nil
	}

	if !ctx.Tag.HasOption(OptionForm) {
		return nil, p.notProvided(ctx, OptionForm)
	}

	if ctx.Tag.HasOption(OptionExplode) {
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

func (p *CookieProvider) notProvided(ctx *Context, opts ...string) error {
	return p.errorf("field: '%v' option: %v not provided", ctx.Tag.Name, opts)
}

func (p *CookieProvider) notSupported(ctx *Context, opt string) error {
	return p.errorf("field: '%v' option: [%v] not supported", ctx.Tag.Name, opt)
}

func (p *CookieProvider) errorf(msg string, values ...interface{}) error {
	msg = fmt.Sprintf(msg, values...)
	return fmt.Errorf("cookie: %s", msg)
}
