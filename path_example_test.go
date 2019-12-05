package inflate_test

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/phogolabs/inflate"
)

func ExampleDecoder_path() {
	type Member struct {
		ID string `path:"id"`
	}

	param := &chi.RouteParams{}
	param.Keys = append(param.Keys, "id")
	param.Values = append(param.Values, "123456")

	member := &Member{}

	if err := inflate.NewPathDecoder(param).Decode(member); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", member)

	// Output:
	// &{ID:123456}
}
