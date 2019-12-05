package inflate

import (
	"net/http"
)

func ExampleDecoder_header() {
	type Tag struct {
		RequestID string `header:"X-Request-ID"`
	}

	header := http.Header{}
	header.Set("X-Request-ID", "123456")

	tag := &Tag{}

	if err := NewHeaderDecoder(header).Decode(tag); err != nil {
		panic(err)
	}

	// Output:
	// &{RequestID:123456}
}
