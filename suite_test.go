package inflate_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEncoding(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Encoding Suite")
}

type TValue struct {
	Value uint `fake:"99"`
}

type User struct {
	Name string `fake:"value"`
}

type Account struct {
	User *User `fake:"user"`
}

type Nested struct {
	ID   string                 `fake:"id"`
	User *User                  `fake:"~"`
	Map  map[string]interface{} `fake:"~"`
}

type Text struct {
	Value string `fake:"value"`
	Error error  `fake:"error"`
}

func (t *Text) UnmarshalText(data []byte) error {
	if t.Error != nil {
		return t.Error
	}

	*t = Text{Value: string(data)}
	return nil
}

func (t Text) MarshalText() ([]byte, error) {
	if t.Error != nil {
		return nil, t.Error
	}

	return []byte(t.Value), nil
}
