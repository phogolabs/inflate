package inflate_test

import (
	"net/url"
	"reflect"

	"github.com/phogolabs/inflate"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query", func() {
	var (
		provider *inflate.QueryProvider
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

		provider = &inflate.QueryProvider{
			Query: url.Values{},
		}
	})

	Describe("NewQueryDecoder", func() {
		It("creates a new query decoder", func() {
			decoder := inflate.NewQueryDecoder(url.Values{})
			Expect(decoder).NotTo(BeNil())
		})
	})

	Describe("NewFormDecoder", func() {
		It("creates a new form decoder", func() {
			decoder := inflate.NewFormDecoder(url.Values{})
			Expect(decoder).NotTo(BeNil())
		})
	})

	Describe("Value", func() {
		Context("when the value is primitive type", func() {
			BeforeEach(func() {
				provider.Query.Set("id", "5")
				ctx.Tag.Options = []string{"form"}
			})

			It("returns the value successfully", func() {
				value, err := provider.Value(ctx)
				Expect(err).To(BeNil())
				Expect(value).To(Equal("5"))
			})

			Context("when the option is unknown", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"unknown"}
				})

				It("returns the an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [form space-delimited deep-object] not provided"))
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
					Expect(value).To(BeNil())
				})
			})

			Context("when the space-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"space-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [space-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the pipe-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"pipe-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [pipe-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the deep-object option is on", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"deep-object"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [deep-object] not supported"))
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
					Expect(value).To(Equal("5"))
				})
			})

			Context("when the simple option is on", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"simple"}
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
				ctx.Type = reflect.TypeOf([]interface{}{})
			})

			Context("when the value is empty string", func() {
				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(BeNil())
				})
			})

			Context("when the form option is provided", func() {
				BeforeEach(func() {
					provider.Query.Add("id", "3,4,5")
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

				Context("when the option is unknown", func() {
					BeforeEach(func() {
						ctx.Tag.Options = []string{"unknown"}
					})

					It("returns the an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("query: field: 'id' option: [form space-delimited pipe-delimited] not provided"))
						Expect(value).To(BeNil())
					})
				})

				Context("when the explode option is provided", func() {
					BeforeEach(func() {
						provider.Query = url.Values{}
						provider.Query.Add("id", "3")
						provider.Query.Add("id", "4")
						provider.Query.Add("id", "5")
						ctx.Tag.Options = []string{"form", "explode"}
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
					ctx.Tag.Options = []string{"space-delimited"}
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
					ctx.Tag.Options = []string{"pipe-delimited"}
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
					ctx.Tag.Options = []string{"deep-object"}
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
					ctx.Tag.Options = []string{"wrong"}
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

			Context("when the value is empty string", func() {
				BeforeEach(func() {
					provider.Query = url.Values{}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(BeNil())
				})
			})

			Context("when the option is unknown", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"unknown"}
				})

				It("returns the an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [form space-delimited pipe-delimited] not provided"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the explode option is provided", func() {
				BeforeEach(func() {
					provider.Query = url.Values{}
					provider.Query.Add("role", "admin")
					provider.Query.Add("firstName", "Alex")

					ctx.Tag.Options = append(ctx.Tag.Options, "explode")
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
						ctx.Tag.Options = []string{"deep-object", "explode"}
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
					ctx.Tag.Options = []string{}
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
					provider.Query.Add("id[role][user]", "admin")
					provider.Query.Add("id[name][first]", "Alex")
					ctx.Tag.Options = []string{"deep-object"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))

					kv := value.(map[string]interface{})
					Expect(kv).To(HaveKey("role"))
					Expect(kv["role"]).To(HaveKeyWithValue("user", "admin"))
					Expect(kv).To(HaveKey("name"))
					Expect(kv["name"]).To(HaveKeyWithValue("first", "Alex"))
				})
			})

			Context("when the query is not valid", func() {
				Context("when the end ] is missing", func() {
					BeforeEach(func() {
						provider.Query = url.Values{}
						provider.Query.Add("id[role][user", "admin")
						ctx.Tag.Options = []string{"deep-object"}
					})

					It("returns an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("query: field: 'id' not parsed: query: cannot parse key: [role][user"))
						Expect(value).To(BeNil())
					})
				})

				Context("when the start [ is missing", func() {
					BeforeEach(func() {
						provider.Query = url.Values{}
						provider.Query.Add("id[role]user]", "admin")
						ctx.Tag.Options = []string{"deep-object"}
					})

					It("returns an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("query: field: 'id' not parsed: query: cannot parse key: [role]user]"))
						Expect(value).To(BeNil())
					})
				})

				Context("when the start [[ is missing", func() {
					BeforeEach(func() {
						provider.Query = url.Values{}
						provider.Query.Add("id[[role][user]", "admin")
						ctx.Tag.Options = []string{"deep-object"}
					})

					It("returns an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("query: field: 'id' not parsed: query: cannot parse key: [[role][user]"))
						Expect(value).To(BeNil())
					})
				})

				Context("when the start [[ is missing", func() {
					BeforeEach(func() {
						provider.Query = url.Values{}
						provider.Query.Add("id]]role][user]", "admin")
						ctx.Tag.Options = []string{"deep-object"}
					})

					It("returns an error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("query: field: 'id' not parsed: query: cannot parse key: ]]role][user]"))
						Expect(value).To(BeNil())
					})
				})
			})

			Context("when the space-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"space-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [space-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the pipe-delimited option is on", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"pipe-delimited"}
				})

				It("returns an error", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("query: field: 'id' option: [pipe-delimited] not supported"))
					Expect(value).To(BeNil())
				})
			})

			Context("when the simple option is on", func() {
				BeforeEach(func() {
					ctx.Tag.Options = []string{"simple"}
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
