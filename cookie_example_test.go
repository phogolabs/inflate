package inflate

import (
	"net/http"
)

func ExampleDecoder_cookie() {
	type Session struct {
		Token string `cookie:"token"`
	}

	cookies := []*http.Cookie{
		{Name: "token", Value: "123456"},
	}

	session := &Session{}

	if err := NewCookieDecoder(cookies).Decode(session); err != nil {
		panic(err)
	}

	// Output:
	// &{Token:123456}
}
