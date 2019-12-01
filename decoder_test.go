package reflectify_test

import (
	"fmt"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/reflectify"
	"github.com/phogolabs/reflectify/fake"
)

var _ = Describe("Decoder", func() {
	var (
		provider *fake.Provider
		decoder  *reflectify.Decoder
	)

	BeforeEach(func() {
		provider = &fake.Provider{}
		provider.NewStub = func(value reflect.Value) reflectify.Provider {
			return &reflectify.ValueProvider{Var: value}
		}

		decoder = &reflectify.Decoder{
			Tag:      "fake",
			Provider: provider,
		}
	})

	Context("when the target value is string", func() {
		type Target struct {
			Name string `fake:"name"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("John", nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Name).To(Equal("John"))
			})
		})

		Context("when the source value is bool", func() {
			BeforeEach(func() {
				provider.ValueReturns(true, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Name).To(Equal("1"))
			})
		})

		Context("when the source value is int", func() {
			BeforeEach(func() {
				provider.ValueReturns(int(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Name).To(Equal("10"))
			})
		})

		Context("when the source value is uint", func() {
			BeforeEach(func() {
				provider.ValueReturns(uint(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Name).To(Equal("10"))
			})
		})

		Context("when the source value is float", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Name).To(Equal("10"))
			})
		})

		Context("when the source value is nil", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Name).To(BeEmpty())
			})
		})

		Context("when the source value is reflectify.TextMarshaller", func() {
			var name string

			BeforeEach(func() {
				now := time.Now()
				provider.ValueReturns(now, nil)

				data, err := now.MarshalText()
				Expect(err).To(BeNil())
				name = string(data)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Name).To(Equal(name))
			})
		})

		Context("when the source value is not supported", func() {
			BeforeEach(func() {
				m := make(map[string]interface{})
				provider.ValueReturns(m, nil)
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("field 'name' does not support value 'map[]'"))
			})
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("oh no"))
			})
		})
	})

	Context("when the target value is bool", func() {
		type Target struct {
			OK bool `fake:"ok"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("true", nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.OK).To(BeTrue())
			})

			Context("when the source value is empty", func() {
				BeforeEach(func() {
					provider.ValueReturns("", nil)
				})

				It("decodes the value successfully", func() {
					t := &Target{}
					Expect(decoder.Decode(t)).To(Succeed())
					Expect(t.OK).To(BeFalse())
				})
			})

			Context("when the source value cannot be decoded", func() {
				BeforeEach(func() {
					provider.ValueReturns("yes", nil)
				})

				It("returns an error", func() {
					t := &Target{}
					Expect(decoder.Decode(t)).To(MatchError("field 'ok' does not support value 'yes': strconv.ParseBool: parsing \"yes\": invalid syntax"))
				})
			})
		})

		Context("when the source value is bool", func() {
			BeforeEach(func() {
				provider.ValueReturns(true, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.OK).To(BeTrue())
			})
		})

		Context("when the source value is int", func() {
			BeforeEach(func() {
				provider.ValueReturns(int(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.OK).To(BeTrue())
			})
		})

		Context("when the source value is uint", func() {
			BeforeEach(func() {
				provider.ValueReturns(uint(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.OK).To(BeTrue())
			})
		})

		Context("when the source value is float", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.OK).To(BeTrue())
			})
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("oh no"))
			})
		})
	})

	Context("when the target value is int", func() {
		type Target struct {
			Result int `fake:"result"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("10", nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(10))
			})

			Context("when the source value cannot be decoded", func() {
				BeforeEach(func() {
					provider.ValueReturns("yes", nil)
				})

				It("returns an error", func() {
					t := &Target{}
					Expect(decoder.Decode(t)).To(MatchError("field 'result' does not support value 'yes': strconv.ParseInt: parsing \"yes\": invalid syntax"))
				})
			})
		})

		Context("when the source value is bool", func() {
			BeforeEach(func() {
				provider.ValueReturns(true, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(1))
			})
		})

		Context("when the source value is int", func() {
			BeforeEach(func() {
				provider.ValueReturns(int(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(10))
			})
		})

		Context("when the source value is uint", func() {
			BeforeEach(func() {
				provider.ValueReturns(uint(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(10))
			})
		})

		Context("when the source value is float", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(10))
			})
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("oh no"))
			})
		})
	})

	Context("when the target value is uint", func() {
		type Target struct {
			Result uint `fake:"result"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("10", nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(uint(10)))
			})

			Context("when the source value cannot be decoded", func() {
				BeforeEach(func() {
					provider.ValueReturns("yes", nil)
				})

				It("returns an error", func() {
					t := &Target{}
					Expect(decoder.Decode(t)).To(MatchError("field 'result' does not support value 'yes': strconv.ParseUint: parsing \"yes\": invalid syntax"))
				})
			})
		})

		Context("when the source value is bool", func() {
			BeforeEach(func() {
				provider.ValueReturns(true, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(uint(1)))
			})
		})

		Context("when the source value is int", func() {
			BeforeEach(func() {
				provider.ValueReturns(int(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(uint(10)))
			})
		})

		Context("when the source value is uint", func() {
			BeforeEach(func() {
				provider.ValueReturns(uint(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(uint(10)))
			})
		})

		Context("when the source value is float", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(uint(10)))
			})
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("oh no"))
			})
		})
	})

	Context("when the target value is float", func() {
		type Target struct {
			Result float32 `fake:"result"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("10", nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(float32(10)))
			})

			Context("when the source value cannot be decoded", func() {
				BeforeEach(func() {
					provider.ValueReturns("yes", nil)
				})

				It("returns an error", func() {
					t := &Target{}
					Expect(decoder.Decode(t)).To(MatchError("field 'result' does not support value 'yes': strconv.ParseFloat: parsing \"yes\": invalid syntax"))
				})
			})
		})

		Context("when the source value is bool", func() {
			BeforeEach(func() {
				provider.ValueReturns(true, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(float32(1)))
			})
		})

		Context("when the source value is int", func() {
			BeforeEach(func() {
				provider.ValueReturns(int(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(float32(10)))
			})
		})

		Context("when the source value is uint", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(float32(10)))
			})
		})

		Context("when the source value is float", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Result).To(Equal(float32(10)))
			})
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("oh no"))
			})
		})
	})

	Context("when the target value is struct", func() {
		type TargetUser struct {
			Name string `fake:"name"`
		}

		type Target struct {
			User *TargetUser `fake:"user"`
		}

		Context("when the source value is map", func() {
			BeforeEach(func() {
				m := make(map[string]interface{})
				m["name"] = "John"

				provider.ValueReturns(m, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.User).NotTo(BeNil())
				Expect(t.User.Name).To(Equal("John"))
			})
		})

		Context("when the source value is the same type", func() {
			BeforeEach(func() {
				m := &TargetUser{
					Name: "John",
				}

				provider.ValueReturns(m, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.User).NotTo(BeNil())
				Expect(t.User.Name).To(Equal("John"))
			})
		})

		Context("when the source value cannot be decoded", func() {
			BeforeEach(func() {
				provider.ValueReturns("yes", nil)
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("field 'user' does not support value 'string'"))
			})
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("oh no"))
			})
		})
	})

	Context("when the target value is map", func() {
		type Target struct {
			KV map[string]interface{} `fake:"kv"`
		}

		Context("when the source value is struct", func() {
			BeforeEach(func() {
				type T struct {
					Name string `fake:"name"`
				}

				t := &T{
					Name: "John",
				}

				provider.ValueReturns(t, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.KV).To(HaveKeyWithValue("name", "John"))
			})
		})

		Context("when the source value is map", func() {
			BeforeEach(func() {
				m := make(map[string]string)
				m["name"] = "John"

				provider.ValueReturns(m, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.KV).To(HaveKeyWithValue("name", "John"))
			})
		})

		Context("when the source value cannot be decoded", func() {
			BeforeEach(func() {
				provider.ValueReturns("yes", nil)
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("field 'kv' does not support value 'yes'"))
			})
		})
	})

	Context("when the target value is array", func() {
		type Target struct {
			List [5]string `fake:"categories"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("john", nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("john"))
			})
		})

		Context("when the source value is bool", func() {
			BeforeEach(func() {
				provider.ValueReturns(true, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("1"))
			})
		})

		Context("when the source value is int", func() {
			BeforeEach(func() {
				provider.ValueReturns(10, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("10"))
			})
		})

		Context("when the source value is uint", func() {
			BeforeEach(func() {
				provider.ValueReturns(uint(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("10"))
			})
		})

		Context("when the source value is float", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("10"))
			})
		})

		Context("when the source value cannot be decoded", func() {
			BeforeEach(func() {
				provider.ValueReturns(struct{}{}, nil)
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("field 'categories' does not support value '{}': field '~' does not support value '{}'"))
			})
		})

		Context("when the source value is array", func() {
			BeforeEach(func() {
				var n [1]string

				n[0] = "John"

				provider.ValueReturns(n, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("John"))
			})
		})

		Context("when the source value is slice", func() {
			BeforeEach(func() {
				provider.ValueReturns([]string{"John"}, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("John"))
			})
		})
	})

	Context("when the target value is slice", func() {
		type Target struct {
			List []string `fake:"categories"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("john", nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("john"))
			})
		})

		Context("when the source value is bool", func() {
			BeforeEach(func() {
				provider.ValueReturns(true, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("1"))
			})
		})

		Context("when the source value is int", func() {
			BeforeEach(func() {
				provider.ValueReturns(10, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("10"))
			})
		})

		Context("when the source value is uint", func() {
			BeforeEach(func() {
				provider.ValueReturns(uint(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("10"))
			})
		})

		Context("when the source value is float", func() {
			BeforeEach(func() {
				provider.ValueReturns(float32(10), nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("10"))
			})
		})

		Context("when the source value cannot be decoded", func() {
			BeforeEach(func() {
				provider.ValueReturns(struct{}{}, nil)
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("field 'categories' does not support value '{}': field '~' does not support value '{}'"))
			})
		})

		Context("when the source value is array", func() {
			BeforeEach(func() {
				var n [1]string

				n[0] = "John"

				provider.ValueReturns(n, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("John"))
			})
		})

		Context("when the source value is slice", func() {
			BeforeEach(func() {
				provider.ValueReturns([]string{"John"}, nil)
			})

			It("decodes the value successfully", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.List).To(ContainElement("John"))
			})
		})
	})

	Context("when the target value is interface", func() {
		type Target struct {
			Data interface{} `fake:"data"`
		}

		Context("when the source value is string", func() {
			BeforeEach(func() {
				provider.ValueReturns("John", nil)
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(Succeed())
				Expect(t.Data).To(Equal("John"))
			})
		})

		Context("when the provider fails", func() {
			BeforeEach(func() {
				provider.ValueReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func() {
				t := &Target{}
				Expect(decoder.Decode(t)).To(MatchError("oh no"))
			})
		})
	})
})
