package inflate_test

import (
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/phogolabs/inflate"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	var (
		provider *inflate.PathProvider
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

		provider = &inflate.PathProvider{
			Param: &chi.RouteParams{},
		}
	})

	Describe("NewPathDecoder", func() {
		It("creates a new path decoder", func() {
			decoder := inflate.NewPathDecoder(&chi.RouteParams{})
			Expect(decoder).NotTo(BeNil())
		})
	})

	Describe("Value", func() {
		Context("when the value is primitive type", func() {
			BeforeEach(func() {
				ctx.Type = reflect.TypeOf("")
			})

			Context("when simple option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "5")
					ctx.Tag.Options = []string{"simple"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(Equal("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Tag.Name = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(BeNil())
					})
				})

				Context("when the option is not provided", func() {
					BeforeEach(func() {
						ctx.Tag.Options = []string{}
					})

					It("returns the value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(Equal("5"))
					})
				})
			})

			Context("when the option is unknown", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "5")
					ctx.Tag.Options = []string{"unknown"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("path: field: id option: [simple label matrix] not provided"))
					Expect(value).To(BeNil())
				})
			})

			Context("when label option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", ".5")
					ctx.Tag.Options = []string{"label"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(Equal("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Tag.Name = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(BeNil())
					})
				})
			})

			Context("when matrix option is on", func() {
				BeforeEach(func() {
					provider.Param.Add("id", ";id=5")
					ctx.Tag.Options = []string{"matrix"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(Equal("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Tag.Name = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(BeNil())
					})
				})
			})
		})

		Context("when the value is array type", func() {
			BeforeEach(func() {
				ctx.Type = reflect.TypeOf([]interface{}{})
			})

			Context("when the option is unknown", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "3,4,5")
					ctx.Tag.Options = []string{"unknown"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("path: field: id option: [simple label matrix] not provided"))
					Expect(value).To(BeNil())
				})
			})

			Context("when simple option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "3,4,5")
					ctx.Tag.Options = []string{"simple"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(3))
					Expect(value).To(ContainElement("3"))
					Expect(value).To(ContainElement("4"))
					Expect(value).To(ContainElement("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Tag.Name = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(0))
					})
				})
			})

			Context("when label option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", ".3,4,5")
					ctx.Tag.Options = []string{"label"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(3))
					Expect(value).To(ContainElement("3"))
					Expect(value).To(ContainElement("4"))
					Expect(value).To(ContainElement("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Tag.Name = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(0))
					})
				})

				Context("when explode option is on", func() {
					BeforeEach(func() {
						provider.Param = &chi.RouteParams{}
						provider.Param.Add("id", ".3.4.5")
						ctx.Tag.Options = append(ctx.Tag.Options, "explode")
					})

					It("returns the value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(3))
						Expect(value).To(ContainElement("3"))
						Expect(value).To(ContainElement("4"))
						Expect(value).To(ContainElement("5"))
					})

					Context("when the param is not found", func() {
						BeforeEach(func() {
							ctx.Tag.Name = "name"
						})

						It("returns a nil value successfully", func() {
							value, err := provider.Value(ctx)
							Expect(err).To(BeNil())
							Expect(value).To(HaveLen(0))
						})
					})
				})
			})

			Context("when matrix option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", ";id=3,4,5")
					ctx.Tag.Options = []string{"matrix"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(3))
					Expect(value).To(ContainElement("3"))
					Expect(value).To(ContainElement("4"))
					Expect(value).To(ContainElement("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Tag.Name = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(0))
					})
				})

				Context("when explode option is on", func() {
					BeforeEach(func() {
						provider.Param = &chi.RouteParams{}
						provider.Param.Add("id", ";id=3;id=4;id=5")
						ctx.Tag.Options = append(ctx.Tag.Options, "explode")
					})

					It("returns the value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(3))
						Expect(value).To(ContainElement("3"))
						Expect(value).To(ContainElement("4"))
						Expect(value).To(ContainElement("5"))
					})

					Context("when the param is not found", func() {
						BeforeEach(func() {
							ctx.Tag.Name = "name"
						})

						It("returns a nil value successfully", func() {
							value, err := provider.Value(ctx)
							Expect(err).To(BeNil())
							Expect(value).To(HaveLen(0))
						})
					})
				})
			})
		})

		Context("when the value is map type", func() {
			BeforeEach(func() {
				ctx.Type = reflect.TypeOf(make(map[string]interface{}))
			})

			Context("when the option is unknown", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "role,admin,firstName,Alex")
					ctx.Tag.Options = []string{"unknown"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(MatchError("path: field: id option: [simple label matrix] not provided"))
					Expect(value).To(BeNil())
				})
			})

			Context("when simple option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "role,admin,firstName,Alex")
					ctx.Tag.Options = []string{"simple"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})

				Context("when the value is invalid", func() {
					BeforeEach(func() {
						provider.Param = &chi.RouteParams{}
						provider.Param.Add("id", "role,admin,firstName")
					})

					It("returns a error", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(MatchError("path: object value: [role admin firstName] invalid"))
						Expect(value).To(BeNil())
					})
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Tag.Name = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(0))
					})
				})

				Context("when explode option is on", func() {
					BeforeEach(func() {
						provider.Param = &chi.RouteParams{}
						provider.Param.Add("id", "role=admin,firstName=Alex")
						ctx.Tag.Options = append(ctx.Tag.Options, "explode")
					})

					It("returns the value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(2))
						Expect(value).To(HaveKeyWithValue("role", "admin"))
						Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
					})
				})
			})

			Context("when label option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", ".role,admin,firstName,Alex")
					ctx.Tag.Options = []string{"label"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})

				Context("when explode option is on", func() {
					BeforeEach(func() {
						provider.Param = &chi.RouteParams{}
						provider.Param.Add("id", ".role=admin.firstName=Alex")
						ctx.Tag.Options = append(ctx.Tag.Options, "explode")
					})

					It("returns the value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(2))
						Expect(value).To(HaveKeyWithValue("role", "admin"))
						Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
					})
				})
			})

			Context("when matrix option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", ";id=role,admin,firstName,Alex")
					ctx.Tag.Options = []string{"matrix"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(HaveLen(2))
					Expect(value).To(HaveKeyWithValue("role", "admin"))
					Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
				})

				Context("when explode option is on", func() {
					BeforeEach(func() {
						provider.Param = &chi.RouteParams{}
						provider.Param.Add("id", ";role=admin;firstName=Alex")
						ctx.Tag.Options = append(ctx.Tag.Options, "explode")
					})

					It("returns the value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(HaveLen(2))
						Expect(value).To(HaveKeyWithValue("role", "admin"))
						Expect(value).To(HaveKeyWithValue("firstName", "Alex"))
					})
				})
			})
		})
	})
})
