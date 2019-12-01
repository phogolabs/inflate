# reflectify
A Golang reflection package on steroids

[![Documentation][godoc-img]][godoc-url]
[![License][license-img]][license-url]
[![Build Status][action-img]][action-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

## Motivation

The package exists to support serialization of parameters based on [OpenAPI
spec](https://swagger.io/docs/specification/serialization/).

## Usage

Let's assume that we have an incoming http request `r`. In order to decode the
following structure:

```golang
type Input struct {
   RequestID string                 `header:"X-Header-ID"`
   Filter    map[string]interface{} `query:"filter,explode"`
   UserID    string                 `path:"user_id"`
   Secret    string                 `cookie:"secret"`
}
```

You can decode a object from a `http.Header` by using the following code:

```golang
if err := reflectify.NewHeaderDecoder(r.Header).Decode(obj) {
  panic(err)
}
```

You can decode a object from a `http.Cookie` by using the following code:

```golang
if err := reflectify.NewHeaderDecoder(r.Cookie()).Decode(obj) {
  panic(err)
}
```

You can decode a object from a `http.Values` by using the following code:

```golang
if err := reflectify.NewQueryDecoder(r.URL.Query()).Decode(obj) {
  panic(err)
}
```

You can decode a object from a `chi.RouteParams` by using the following code:

```golang
if ctx, ok := r.Context().Value(chi.RouteCtxKey).(*chi.Context); ok {
	if err = reflectify.NewPathDecoder(&ctx.URLParams).Decode(obj); err != nil {
    panic(err)
  }
}
```

## Installation

```console
$ go get -u github.com/phogolabs/reflectify
```

## Contributing

We are open for any contributions. Just fork the
[project](https://github.com/phogolabs/reflectify).

[report-img]: https://goreportcard.com/badge/github.com/phogolabs/reflectify
[report-url]: https://goreportcard.com/report/github.com/phogolabs/reflectify
[logo-author-url]: https://www.freepik.com/free-vector/abstract-cross-logo-template_1185919.htm
[logo-license]: http://creativecommons.org/licenses/by/3.0/
[codecov-url]: https://codecov.io/gh/phogolabs/reflectify
[codecov-img]: https://codecov.io/gh/phogolabs/reflectify/branch/master/graph/badge.svg
[action-img]: https://github.com/phogolabs/reflectify/workflows/pipeline/badge.svg
[action-url]: https://github.com/phogolabs/reflectify/actions
[godoc-url]: https://godoc.org/github.com/phogolabs/reflectify
[godoc-img]: https://godoc.org/github.com/phogolabs/reflectify?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: LICENSE
