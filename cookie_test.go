package inflate_test

import (
	"net/http"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/inflate"
)

var _ = Describe("Cookie", func() {
	var (
		provider *inflate.CookieProvider
		ctx      *inflate.Context
	)

	BeforeEach(func() {
		ctx = &inflate.Context{
			Field: "id",
			Type:  reflect.TypeOf(""),
			Tag: &inflate.Tag{
				Key:  "fake",
				Name: "id",
			},
		}

		provider = &inflate.CookieProvider{
			[]*http.Cookie{{Name: "id"}},
		}
	})

	Describe("NewCookieDecoder", func() {
		It("creates a new path decoder", func() {
			decoder := inflate.NewCookieDecoder([]*http.Cookie{})
			Expect(decoder).NotTo(BeNil())
		})
	})

	Context("when the value is primitive type", func() {
		Context("when the cookie is not found", func() {
			BeforeEach(func() {
				ctx.Tag.Name = "name"
			})

			It("returns a nil value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(BeNil())
			})
		})

		Context("when the form option is not provided", func() {
			BeforeEach(func() {
				ctx.Tag.Options = []string{}
				provider.Cookies[0].Value = "5"
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(Equal("5"))
			})
		})

		Context("when the option is unknown", func() {
			BeforeEach(func() {
				ctx.Tag.Options = []string{"unknown"}
			})

			It("returns the an error", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(MatchError("cookie: field: 'id' option: [form] not provided"))
				Expect(value).To(BeNil())
			})
		})
	})

	Describe("Value", func() {
		Context("when the value is array type", func() {
			BeforeEach(func() {
				provider.Cookies[0].Value = "3,4,5"
				ctx.Type = reflect.TypeOf([]interface{}{})
				ctx.Tag.Options = []string{"form"}
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(HaveLen(3))
				Expect(value).To(ContainElement("3"))
				Expect(value).To(ContainElement("4"))
				Expect(value).To(ContainElement("5"))
			})

			Context("when the cookie is not found", func() {
				BeforeEach(func() {
					ctx.Tag.Name = "name"
				})

				It("returns a nil value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(BeNil())
				})
			})

			Context("when the explode option is not provided", func() {
				BeforeEach(func() {
					ctx.Tag.Options = append(ctx.Tag.Options, "explode")
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("cookie: field: 'id' option: [explode] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the option is unknown", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"unknown"}
				})

				It("returns the an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("cookie: field: 'id' option: [form] not provided"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the form option is not provided", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{}
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(3))
					Expect(value).To(ContainElement("3"))
					Expect(value).To(ContainElement("4"))
					Expect(value).To(ContainElement("5"))
				})
			})
		})

		Context("when the value is map type", func() {
			BeforeEach(func() {
				provider.Cookies[0].Value = "role,admin,firstName,Alex"
				ctx.Type = reflect.TypeOf(make(map[string]interface{}))
				ctx.Tag.Options = []string{"form"}
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(HaveLen(2))
				Expect(value).To(HaveKeyWithValue("role", "admin"))
				Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
			})

			Context("when the explode option is  provided", func() {
				BeforeEach(func() {
					ctx.Tag.Options = append(ctx.Tag.Options, "explode")
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("cookie: field: 'id' option: [explode] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the cookie is not found", func() {
				BeforeEach(func() {
					ctx.Tag.Name = "name"
				})

				It("returns a nil value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(0))
				})
			})

			Context("when the option is unknown", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"unknown"}
				})

				It("returns the an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("cookie: field: 'id' option: [form] not provided"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the value is invalid", func() {
				BeforeEach(func() {
					provider.Cookies[0].Value = ",firstName,Alex"
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("cookie: object value: [ firstName Alex] invalid"))
					Expect(value).To(HaveLen(0))
				})
			})

			Context("when the form option is not provided", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{}
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})
			})

			Context("when the value is invalid", func() {
				BeforeEach(func() {
					provider.Cookies[0].Value = "role,admin,firstName"
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("cookie: object value: [role admin firstName] invalid"))
					Expect(value).To(BeNil())
				})
			})
		})
	})
})
