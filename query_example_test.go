package inflate_test

import (
	"net/url"
)

func ExampleDecoder_query() {
	type User struct {
		ID   string `query:"id"`
		Name string `query:"name"`
	}

	query, err := url.ParseQuery("id=1&name=Jack")
	if err != nil {
		panic(err)
	}

	user := &User{}

	if err := NewQueryDecoder(query).Decode(user); err != nil {
		panic(err)
	}

	// Output:
	// &{ID:1 Name:Jack}
}
