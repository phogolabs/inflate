package reflectify_test

import (
	"net/url"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/reflectify"
)

var _ = Describe("Query", func() {
	var (
		provider *reflectify.QueryProvider
		ctx      *reflectify.Context
	)

	BeforeEach(func() {
		ctx = &reflectify.Context{
			Field: "id",
		}

		provider = &reflectify.QueryProvider{
			Query: url.Values{},
		}
	})

	Describe("NewQueryDecoder", func() {
		It("creates a new query decoder", func() {
			decoder := reflectify.NewQueryDecoder(url.Values{})
			Expect(decoder).NotTo(BeNil())
		})
	})

	Describe("Value", func() {
		Context("when the value is primitive type", func() {
			BeforeEach(func() {
				provider.Query.Set("id", "5")
				ctx.Options = []string{"form"}
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(Equal("5"))
			})

			Context("when the cookie is not found", func() {
				BeforeEach(func() {
					ctx.Field = "name"
				})

				It("returns a nil value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(BeNil())
				})
			})

			Context("when the space-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Options = []string{"space-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [space-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the pipe-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Options = []string{"pipe-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [pipe-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the deep-object option is on", func() {
				BeforeEach(func() {
					ctx.Options = []string{"deep-object"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [deep-object] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the form option is not provided", func() {
				BeforeEach(func() {
					ctx.Options = []string{}
				})

				It("returns a error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(Equal("5"))
				})
			})

			Context("when the simple option is on", func() {
				BeforeEach(func() {
					ctx.Options = []string{"simple"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [form space-delimited deep-object] not provided"))
					Expect(value).To(BeNil())
				})
			})
		})

		Context("when the value is array type", func() {
			BeforeEach(func() {
				ctx.FieldKind = reflect.Array
			})

			Context("when the form option is provided", func() {
				BeforeEach(func() {
					provider.Query.Add("id", "3,4,5")
					ctx.Options = []string{"form"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(3))
					Expect(value).To(ContainElement("3"))
					Expect(value).To(ContainElement("4"))
					Expect(value).To(ContainElement("5"))
				})

				Context("when the explode option is provided", func() {
					BeforeEach(func() {
						provider.Query = url.Values{}
						provider.Query.Add("id", "3")
						provider.Query.Add("id", "4")
						provider.Query.Add("id", "5")
						ctx.Options = []string{"form", "explode"}
					})

					It("returns the value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(3))
						Expect(value).To(ContainElement("3"))
						Expect(value).To(ContainElement("4"))
						Expect(value).To(ContainElement("5"))
					})
				})
			})

			Context("when the space-delimited option is provided", func() {
				BeforeEach(func() {
					provider.Query.Add("id", "3 4 5")
					ctx.Options = []string{"space-delimited"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(3))
					Expect(value).To(ContainElement("3"))
					Expect(value).To(ContainElement("4"))
					Expect(value).To(ContainElement("5"))
				})
			})

			Context("when the pipe-delimited option is provided", func() {
				BeforeEach(func() {
					provider.Query.Add("id", "3|4|5")
					ctx.Options = []string{"pipe-delimited"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(3))
					Expect(value).To(ContainElement("3"))
					Expect(value).To(ContainElement("4"))
					Expect(value).To(ContainElement("5"))
				})
			})

			Context("when the deep-object option is provided", func() {
				BeforeEach(func() {
					provider.Query.Add("id", "3|4|5")
					ctx.Options = []string{"deep-object"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [deep-object] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the wrong option is provided", func() {
				BeforeEach(func() {
					provider.Query.Add("id", "3|4|5")
					ctx.Options = []string{"wrong"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [form space-delimited pipe-delimited] not provided"))
					Expect(value).To(BeNil())
				})
			})
		})

		Context("when the value is map type", func() {
			BeforeEach(func() {
				provider.Query.Add("id", "role,admin,firstName,Alex")
				ctx.FieldKind = reflect.Map
				ctx.Options = []string{"form"}
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(HaveLen(2))
				Expect(value).To(HaveKeyWithValue("role", "admin"))
				Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
			})

			Context("when the explode option is provided", func() {
				BeforeEach(func() {
					provider.Query = url.Values{}
					provider.Query.Add("role", "admin")
					provider.Query.Add("firstName", "Alex")

					ctx.Options = append(ctx.Options, "explode")
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})

				Context("when the deep-object option is on", func() {
					BeforeEach(func() {
						ctx.Options = []string{"deep-object", "explode"}
					})

					It("returns an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("query: field: 'id' option: [explode] not supported"))
						Expect(value).To(BeNil())
					})
				})
			})

			Context("when the form option is not provided", func() {
				BeforeEach(func() {
					provider.Query = url.Values{}
					provider.Query.Add("role", "admin")
					provider.Query.Add("firstName", "Alex")
					ctx.Options = []string{}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})
			})

			Context("when the deep-object is on", func() {
				BeforeEach(func() {
					provider.Query = url.Values{}
					provider.Query.Add("id[role]", "admin")
					provider.Query.Add("id[firstName]", "Alex")
					ctx.Options = []string{"deep-object"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})
			})

			Context("when the space-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Options = []string{"space-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [space-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the pipe-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Options = []string{"pipe-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [pipe-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the simple option is on", func() {
				BeforeEach(func() {
					ctx.Options = []string{"simple"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [form space-delimited pipe-delimited] not provided"))
					Expect(value).To(BeNil())
				})
			})
		})
	})
})
