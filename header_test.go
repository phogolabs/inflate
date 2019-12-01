package reflectify_test

import (
	"net/http"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/reflectify"
)

var _ = Describe("Header", func() {
	var (
		provider *reflectify.HeaderProvider
		ctx      *reflectify.Context
	)

	BeforeEach(func() {
		ctx = &reflectify.Context{
			Field: "X-MyHeader",
		}

		provider = &reflectify.HeaderProvider{
			Header: http.Header{},
		}
	})

	Describe("NewHeaderDecoder", func() {
		It("creates a new header decoder", func() {
			decoder := reflectify.NewHeaderDecoder(http.Header{})
			Expect(decoder).NotTo(BeNil())
		})
	})

	Describe("Value", func() {
		BeforeEach(func() {
			provider.Header.Set("X-MyHeader", "5")
			ctx.Options = []string{"simple"}
		})

		Context("when the value is primitive type", func() {
			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(Equal("5"))
			})

			Context("when the header is not found", func() {
				BeforeEach(func() {
					ctx.Field = "name"
				})

				It("returns a nil value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(BeNil())
				})
			})

			Context("when the simple option is not provided", func() {
				BeforeEach(func() {
					ctx.Options = []string{}
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(Succeed())
					Expect(value).To(Equal("5"))
				})
			})
		})

		Context("when the value is array type", func() {
			BeforeEach(func() {
				provider.Header.Set("X-MyHeader", "3,4,5")
				ctx.FieldKind = reflect.Array
				ctx.Options = []string{"simple"}
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(HaveLen(3))
				Expect(value).To(ContainElement("3"))
				Expect(value).To(ContainElement("4"))
				Expect(value).To(ContainElement("5"))
			})

			Context("when the header is not found", func() {
				BeforeEach(func() {
					ctx.Field = "X-TheirHeader"
				})

				It("returns a nil value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(BeNil())
				})
			})

			Context("when the simple option is not provided", func() {
				BeforeEach(func() {
					ctx.Options = []string{}
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
				provider.Header.Set("X-MyHeader", "role,admin,firstName,Alex")
				ctx.FieldKind = reflect.Map
				ctx.Options = []string{"simple"}
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(HaveLen(2))
				Expect(value).To(HaveKeyWithValue("role", "admin"))
				Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
			})

			Context("when the header is not found", func() {
				BeforeEach(func() {
					ctx.Field = "X-TheirHeader"
				})

				It("returns a nil value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(0))
				})
			})

			Context("when the value is invalid", func() {
				BeforeEach(func() {
					provider.Header.Set("X-MyHeader", ",firstName,Alex")
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("header: object value: [ firstName Alex] invalid"))
					Expect(value).To(HaveLen(0))
				})
			})

			Context("when the simple option is not provided", func() {
				BeforeEach(func() {
					ctx.Options = []string{}
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})
			})

			Context("when the explode option is available", func() {
				BeforeEach(func() {
					provider.Header.Set("X-MyHeader", "role=admin,firstName=Alex")
					ctx.Options = []string{"simple", "explode"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})

				Context("when the header is not found", func() {
					BeforeEach(func() {
						ctx.Field = "X-TheirHeader"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(0))
					})
				})

				Context("when the value is invalid", func() {
					BeforeEach(func() {
						provider.Header.Set("X-MyHeader", ",firstName=Alex")
					})

					It("returns an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("header: object value: [ firstName=Alex] invalid"))
						Expect(value).To(HaveLen(0))
					})
				})

				Context("when the field name is not provided", func() {
					BeforeEach(func() {
						provider.Header.Set("X-MyHeader", "=admin,firstName=Alex")
					})

					It("returns an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("header: object value: [=admin firstName=Alex] invalid"))
						Expect(value).To(HaveLen(0))
					})
				})

				Context("when the simple option is not provided", func() {
					BeforeEach(func() {
						provider.Header.Set("X-MyHeader", "role=admin,firstName=Alex")
						ctx.Options = []string{"explode"}
					})

					It("returns a error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("header: field 'X-MyHeader' option: [simple] not provided"))
						Expect(value).To(BeNil())
					})
				})
			})
		})
	})
})
