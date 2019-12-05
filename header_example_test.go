package inflate_test

import (
	"fmt"
	"net/http"

	"github.com/phogolabs/inflate"
)

func ExampleDecoder_header() {
	type Tag struct {
		RequestID string `header:"X-Request-ID"`
	}

	header := http.Header{}
	header.Set("X-Request-ID", "123456")

	tag := &Tag{}

	if err := inflate.NewHeaderDecoder(header).Decode(tag); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", tag)

	// Output:
	// &{RequestID:123456}
}
