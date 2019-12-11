package inflate_test

import (
	"fmt"

	"github.com/phogolabs/inflate"
)

func ExampleDecoder_default() {
	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type Profile struct {
		Name    string  `default:"John"`
		Address Address `default:"{\"city\":\"London\",\"country\":\"UK\"}"`
	}

	profile := &Profile{}

	if err := inflate.SetDefault(profile); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", profile)

	// Output:
	// &{Name:John Address:{City:London Country:UK}}
}
