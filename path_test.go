package reflectify_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/go-chi/chi"
	"github.com/phogolabs/reflectify"
)

var _ = Describe("Path", func() {
	var (
		provider *reflectify.PathProvider
		ctx      *reflectify.Context
	)

	BeforeEach(func() {
		ctx = &reflectify.Context{
			Field: "id",
		}

		provider = &reflectify.PathProvider{
			Param: &chi.RouteParams{},
		}
	})

	Describe("NewPathDecoder", func() {
		It("creates a new path decoder", func() {
			decoder := reflectify.NewPathDecoder(&chi.RouteParams{})
			Expect(decoder).NotTo(BeNil())
		})
	})

	Describe("Value", func() {
		Context("when the value is primitive type", func() {
			Context("when simple option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "5")
					ctx.Options = []string{"simple"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(Equal("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Field = "name"
					})

					It("returns a nil value successfully", func() {
						value, err := provider.Value(ctx)
						Expect(err).To(BeNil())
						Expect(value).To(BeNil())
					})
				})
			})

			Context("when label option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", ".5")
					ctx.Options = []string{"label"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(Equal("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Field = "name"
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
					ctx.Options = []string{"matrix"}
				})

				It("returns the value successfully", func() {
					value, err := provider.Value(ctx)
					Expect(err).To(BeNil())
					Expect(value).To(Equal("5"))
				})

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Field = "name"
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
				ctx.FieldKind = reflect.Array
			})

			Context("when simple option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "3,4,5")
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

				Context("when the param is not found", func() {
					BeforeEach(func() {
						ctx.Field = "name"
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
					ctx.Options = []string{"label"}
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
						ctx.Field = "name"
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
						ctx.Options = append(ctx.Options, "explode")
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
							ctx.Field = "name"
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
					ctx.Options = []string{"matrix"}
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
						ctx.Field = "name"
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
						ctx.Options = append(ctx.Options, "explode")
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
							ctx.Field = "name"
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
				ctx.FieldKind = reflect.Map
			})

			Context("when simple option is on", func() {
				BeforeEach(func() {
					provider.Param = &chi.RouteParams{}
					provider.Param.Add("id", "role,admin,firstName,Alex")
					ctx.Options = []string{"simple"}
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
						ctx.Field = "name"
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
						ctx.Options = append(ctx.Options, "explode")
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
					ctx.Options = []string{"label"}
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
						ctx.Options = append(ctx.Options, "explode")
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
					ctx.Options = []string{"matrix"}
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
						ctx.Options = append(ctx.Options, "explode")
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
