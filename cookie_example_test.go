package inflate_test

import (
	"fmt"
	"net/http"

	"github.com/phogolabs/inflate"
)

func ExampleDecoder_cookie() {
	type Session struct {
		Token string `cookie:"token"`
	}

	cookies := []*http.Cookie{
		{Name: "token", Value: "123456"},
	}

	session := &Session{}

	if err := inflate.NewCookieDecoder(cookies).Decode(session); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", session)

	// Output:
	// &{Token:123456}
}
