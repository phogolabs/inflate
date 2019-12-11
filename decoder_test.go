package inflate_test

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/inflate"
	"github.com/phogolabs/inflate/fake"
)

var _ = Describe("Decoder", func() {
	var (
		decoder   *inflate.Decoder
		provider  *fake.ValueProvider
		converter *fake.ValueConverter
	)

	type User struct {
		Name string `fake:"name"`
	}

	BeforeEach(func() {
		provider = &fake.ValueProvider{}
		converter = &fake.ValueConverter{}

		converter.ConvertStub = func(source, target interface{}) error {
			target.(reflect.Value).Set(source.(reflect.Value))
			return nil
		}

		decoder = &inflate.Decoder{
			TagName:   "fake",
			Converter: converter,
			Provider:  provider,
		}
	})

	It("decodes the target successfully", func() {
		provider.ValueReturns("Jack", nil)

		user := &User{}
		Expect(decoder.Decode(user)).To(Succeed())
		Expect(user.Name).To(Equal("Jack"))
	})

	Context("when the provider fails", func() {
		BeforeEach(func() {
			provider.ValueReturns(nil, fmt.Errorf("oh no"))
		})

		It("returns an error", func() {
			user := &User{}
			Expect(decoder.Decode(user)).To(MatchError("oh no"))
			Expect(user.Name).To(BeEmpty())
		})
	})

	Context("when the converter fails", func() {
		BeforeEach(func() {
			converter.ConvertReturns(fmt.Errorf("oh no"))
		})

		It("returns an error", func() {
			user := &User{}
			Expect(decoder.Decode(user)).To(MatchError("oh no"))
			Expect(user.Name).To(BeEmpty())
		})
	})

	Context("when there is a squashed type", func() {
		type Account struct {
			User *User `fake:"~"`
		}

		It("decodes the target successfully", func() {
			provider.ValueReturns("Jack", nil)

			account := &Account{}

			Expect(decoder.Decode(account)).To(Succeed())
			Expect(account.User).NotTo(BeNil())
			Expect(account.User.Name).To(Equal("Jack"))
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				account := &Account{}
				Expect(decoder.Decode(account)).To(MatchError("oh no"))
				Expect(account.User).To(BeNil())
			})
		})

		Context("when the converter fails", func() {
			BeforeEach(func() {
				converter.ConvertReturns(fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				account := &Account{}
				Expect(decoder.Decode(account)).To(MatchError("oh no"))
				Expect(account.User).To(BeNil())
			})
		})
	})
})

var _ = Describe("SetDefault", func() {
	type Account struct {
		Category string `default:"unknown"`
		User     *User  `default:"{\"name\":\"Peter\"}"`
	}

	It("sets the defaults successfully", func() {
		account := &Account{}
		Expect(inflate.SetDefault(account)).To(Succeed())
		Expect(account.Category).To(Equal("unknown"))
		Expect(account.User).NotTo(BeNil())
		Expect(account.User.Name).To(Equal("Peter"))
	})

	Context("when the value is set", func() {
		It("does not set the value", func() {
			account := &Account{Category: "Jack"}
			Expect(inflate.SetDefault(account)).To(Succeed())
			Expect(account.Category).To(Equal("Jack"))
		})

		Context("whne the subproperty has a non zero field", func() {
			It("does not set the value", func() {
				account := &Account{User: &User{Name: "Peter"}}
				Expect(inflate.SetDefault(account)).To(Succeed())
				Expect(account.Category).To(Equal("unknown"))
				Expect(account.User).NotTo(BeNil())
				Expect(account.User.Name).To(Equal("Peter"))
			})
		})
	})
})

var _ = Describe("Set", func() {
	type Order struct {
		ID string `field:"order_id"`
	}

	type OrderItem struct {
		OrderID string `field:"order_id"`
	}

	It("sets the values successfully", func() {
		source := &Order{ID: "0000123"}
		target := &OrderItem{}

		Expect(inflate.Set(target, source)).To(Succeed())
		Expect(target.OrderID).To(Equal(source.ID))
	})
})
