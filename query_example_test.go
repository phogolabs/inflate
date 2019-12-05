package inflate_test

import (
	"fmt"
	"net/url"

	"github.com/phogolabs/inflate"
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

	if err := inflate.NewQueryDecoder(query).Decode(user); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", user)

	// Output:
	// &{ID:1 Name:Jack}
}
