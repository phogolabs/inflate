package inflate_test

import (
	"bytes"
	"fmt"
	"io"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/inflate"
	"github.com/phogolabs/schema"
)

var _ = Describe("Converter", func() {
	var converter *inflate.Converter

	BeforeEach(func() {
		converter = &inflate.Converter{
			TagName: "fake",
		}
	})

	Context("when the target value is string", func() {
		var target string

		BeforeEach(func() {
			target = ""
		})

		Context("when the source value is string", func() {
			var source string

			BeforeEach(func() {
				source = "phogo"
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(Equal("phogo"))
			})
		})

		Context("when the source value is bool", func() {
			var source bool

			BeforeEach(func() {
				source = true
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(Equal("1"))
			})

			Context("when the source value is false", func() {
				BeforeEach(func() {
					source = false
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal("0"))
				})
			})
		})

		Context("when the source value is int", func() {
			var source int

			BeforeEach(func() {
				source = 10
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(Equal("10"))
			})
		})

		Context("when the source value is uint", func() {
			var source uint

			BeforeEach(func() {
				source = 10
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(Equal("10"))
			})
		})

		Context("when the source value is float", func() {
			var source float32

			BeforeEach(func() {
				source = 10
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(Equal("10"))
			})
		})

		Context("when the source value is struct", func() {
			It("returns an error", func() {
				source := struct{}{}
				Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert struct value {} to type string"))
				Expect(target).To(BeEmpty())
			})

			Context("when the struct implements TextMarshaller", func() {
				var source Text

				BeforeEach(func() {
					source = Text{Value: "John"}
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal("John"))
				})

				Context("when the TextMarshaller fails", func() {
					BeforeEach(func() {
						source.Error = fmt.Errorf("oh no")
					})

					It("returns an error", func() {
						Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert struct value {Value:John Error:oh no} to type string: oh no"))
						Expect(target).To(BeEmpty())
					})
				})
			})
		})
	})

	Context("when the target value is bool", func() {
		var target bool

		BeforeEach(func() {
			target = false
		})

		Context("when the source value is string", func() {
			var source string

			BeforeEach(func() {
				source = "1"
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(BeTrue())
			})

			Context("when the value is empty", func() {
				BeforeEach(func() {
					source = ""
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(BeFalse())
				})
			})

			Context("when the value cannot be parsed", func() {
				BeforeEach(func() {
					source = "unknonw"
				})

				It("returns an error", func() {
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert string value unknonw to type bool: strconv.ParseBool: parsing \"unknonw\": invalid syntax"))
					Expect(target).To(BeFalse())
				})
			})
		})

		Context("when the source value is bool", func() {
			var source bool

			BeforeEach(func() {
				source = true
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(BeTrue())
			})

			Context("when the source value is false", func() {
				BeforeEach(func() {
					source = false
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(BeFalse())
				})
			})
		})

		Context("when the source value is int", func() {
			var source int

			BeforeEach(func() {
				source = 10
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(BeTrue())
			})
		})

		Context("when the source value is uint", func() {
			var source uint

			BeforeEach(func() {
				source = 10
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(BeTrue())
			})
		})

		Context("when the source value is float", func() {
			var source float32

			BeforeEach(func() {
				source = 10
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(BeTrue())
			})
		})

		Context("when the source value is struct", func() {
			It("returns an error", func() {
				source := struct{}{}
				Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert struct value {} to type bool"))
				Expect(target).To(BeFalse())
			})
		})

		Context("when the target value is int", func() {
			var target int

			BeforeEach(func() {
				target = 0
			})

			Context("when the source value is string", func() {
				var source string

				BeforeEach(func() {
					source = "10"
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(10))
				})

				Context("when the source value cannot be parsed", func() {
					BeforeEach(func() {
						source = "unknown"
					})

					It("returns an error", func() {
						Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert string value unknown to type int: strconv.ParseInt: parsing \"unknown\": invalid syntax"))
						Expect(target).To(Equal(0))
					})
				})
			})

			Context("when the source value is bool", func() {
				var source bool

				BeforeEach(func() {
					source = false
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(0))
				})

				Context("when the source value is true", func() {
					BeforeEach(func() {
						source = true
					})

					It("converts the value successfully", func() {
						Expect(converter.Convert(&source, &target)).To(Succeed())
						Expect(target).To(Equal(1))
					})
				})
			})

			Context("when the source value is int", func() {
				var source int

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(10))
				})
			})

			Context("when the source value is uint", func() {
				var source uint

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(10))
				})
			})

			Context("when the source value is float", func() {
				var source float32

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(10))
				})
			})

			Context("when the source value is struct", func() {
				It("returns an error", func() {
					source := struct{}{}
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert struct value {} to type int"))
					Expect(target).To(Equal(0))
				})
			})
		})

		Context("when the target value is uint", func() {
			var target uint

			BeforeEach(func() {
				target = 0
			})

			Context("when the source value is string", func() {
				var source string

				BeforeEach(func() {
					source = "10"
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(uint(10)))
				})

				Context("when the source value cannot be parsed", func() {
					var source string

					BeforeEach(func() {
						source = "unknown"
					})

					It("returns an error", func() {
						Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert string value unknown to type uint: strconv.ParseUint: parsing \"unknown\": invalid syntax"))
						Expect(target).To(Equal(uint(0)))
					})
				})
			})

			Context("when the source value is bool", func() {
				var source bool

				BeforeEach(func() {
					source = true
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(uint(1)))
				})

				Context("when the source value is false", func() {
					BeforeEach(func() {
						source = false
					})

					It("converts the value successfully", func() {
						Expect(converter.Convert(&source, &target)).To(Succeed())
						Expect(target).To(Equal(uint(0)))
					})
				})
			})

			Context("when the source value is int", func() {
				var source int

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(uint(10)))
				})
			})

			Context("when the source value is uint", func() {
				var source uint

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(uint(10)))
				})
			})

			Context("when the source value is float", func() {
				var source float32

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(uint(10)))
				})
			})

			Context("when the source value is struct", func() {
				It("returns an error", func() {
					source := struct{}{}
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert struct value {} to type uint"))
					Expect(target).To(Equal(uint(0)))
				})
			})

		})

		Context("when the target value is float", func() {
			var target float32

			BeforeEach(func() {
				target = 0
			})

			Context("when the source value is string", func() {
				var source string

				BeforeEach(func() {
					source = "10"
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(float32(10)))
				})

				Context("when the source value cannot be parsed", func() {
					BeforeEach(func() {
						source = "unknown"
					})

					It("returns an error", func() {
						Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert string value unknown to type float32: strconv.ParseFloat: parsing \"unknown\": invalid syntax"))
						Expect(target).To(Equal(float32(0)))
					})
				})
			})

			Context("when the source value is bool", func() {
				var source bool

				BeforeEach(func() {
					source = true
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(float32(1)))
				})

				Context("when the source value is false", func() {
					BeforeEach(func() {
						source = false
					})

					It("converts the value successfully", func() {
						Expect(converter.Convert(&source, &target)).To(Succeed())
						Expect(target).To(Equal(float32(0)))
					})
				})
			})

			Context("when the source value is int", func() {
				var source int

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(float32(10)))
				})
			})

			Context("when the source value is uint", func() {
				var source uint

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(float32(10)))
				})
			})

			Context("when the source value is float", func() {
				var source float32

				BeforeEach(func() {
					source = 10
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal(float32(10)))
				})
			})

			Context("when the source value is struct", func() {
				It("returns an error", func() {
					source := struct{}{}
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert struct value {} to type float32"))
					Expect(target).To(Equal(float32(0)))
				})
			})
		})
	})

	Context("when the target value is struct", func() {
		var target Text

		BeforeEach(func() {
			target = Text{}
		})

		Context("when the source value is string", func() {
			var source string

			BeforeEach(func() {
				source = "John"
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target.Value).To(Equal("John"))
			})

			Context("when the target is not TextUnmarshaller", func() {
				var target User

				BeforeEach(func() {
					target = User{}
				})

				It("returns an error", func() {
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert string value John to type struct"))
				})
			})
		})

		Context("when the source value is struct", func() {
			var source User

			BeforeEach(func() {
				source = User{
					Name: "John",
				}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target.Value).To(Equal("John"))
			})

			Context("when the source is int", func() {
				It("returns an error", func() {
					source := 10
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert int value 10 to type struct"))
				})
			})
		})

		Context("when the source value is equal to the target", func() {
			var source Text

			BeforeEach(func() {
				source = Text{Value: "John"}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target.Value).To(Equal("John"))
			})
		})

		Context("when the source value is map", func() {
			var source map[string]interface{}

			BeforeEach(func() {
				source = map[string]interface{}{
					"value": "John",
				}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target.Value).To(Equal("John"))
			})

			Context("when the map is empty", func() {
				BeforeEach(func() {
					source = map[string]interface{}{}
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target.Value).To(BeEmpty())
				})
			})
		})
	})

	Context("when the target value is a nested struct", func() {
		var target Account

		BeforeEach(func() {
			target = Account{}
		})

		Context("when the source value is map", func() {
			var source map[string]interface{}

			BeforeEach(func() {
				source = make(map[string]interface{})
				source["user"] = map[string]interface{}{
					"value": "John",
				}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target.User).NotTo(BeNil())
				Expect(target.User.Name).To(Equal("John"))
			})
		})
	})

	Context("when the target value is map", func() {
		var target map[float32]int

		BeforeEach(func() {
			target = make(map[float32]int)
		})

		Context("when the source value is map", func() {
			var source map[string]float32

			BeforeEach(func() {
				source = make(map[string]float32)
				source["10"] = 99
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(HaveKeyWithValue(float32(10), 99))
			})
		})

		Context("when the source value is the same type as target ", func() {
			var source map[float32]int

			BeforeEach(func() {
				source = make(map[float32]int)
				source[10] = 99
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(HaveKeyWithValue(float32(10), 99))
			})
		})

		Context("when the source value is struct", func() {
			var source TValue

			BeforeEach(func() {
				source = TValue{
					Value: 10,
				}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(HaveKeyWithValue(float32(99), 10))
			})

			Context("when the source is int", func() {
				It("returns an error", func() {
					source := 10
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert int value 10 to type map"))
				})
			})
		})
	})

	Context("when the target value is slice", func() {
		var target []string

		Context("when the source value is string", func() {
			var source string

			BeforeEach(func() {
				source = "John"
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("John"))
			})
		})

		Context("when the source value is equal to the target", func() {
			var source []string

			BeforeEach(func() {
				source = []string{"John"}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(HaveLen(1))
				Expect(target).To(ContainElement("John"))
			})
		})

		Context("when the source value is struct", func() {
			var source User

			BeforeEach(func() {
				source = User{
					Name: "John",
				}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("John"))
			})
		})

		Context("when the source value is map", func() {
			var source map[int]string

			BeforeEach(func() {
				source = map[int]string{
					10: "John",
					9:  "Jack",
				}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("John"))
				Expect(target).To(ContainElement("Jack"))
			})
		})

		Context("when the source value is array", func() {
			var source [2]int

			BeforeEach(func() {
				source[0] = 5
				source[1] = 4
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("5"))
				Expect(target).To(ContainElement("4"))
			})
		})

		Context("when the source value is slice", func() {
			var source []int

			BeforeEach(func() {
				source = []int{5, 4}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("5"))
				Expect(target).To(ContainElement("4"))
			})
		})
	})

	Context("when the target value is array", func() {
		var target [1]string

		Context("when the source value is array", func() {
			var source [2]int

			BeforeEach(func() {
				source[0] = 5
				source[1] = 4
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("5"))
			})
		})

		Context("when the source value is equal to the target", func() {
			var source [1]int

			BeforeEach(func() {
				source[0] = 5
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("5"))
			})
		})

		Context("when the source value is slice", func() {
			var source []int

			BeforeEach(func() {
				source = []int{5, 4}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(ContainElement("5"))
			})
		})

		Context("when the source value is map", func() {
			var source map[int]string

			BeforeEach(func() {
				source = map[int]string{
					10: "John",
					9:  "Jack",
				}
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target).To(HaveLen(1))
			})
		})
	})

	Context("when the target value is uuid", func() {
		var target schema.UUID

		Context("when the source value is string", func() {
			var source string

			BeforeEach(func() {
				source = "5f68ada0-5c63-4f82-9619-a51bc07d38ec"
			})

			It("converts the value successfully", func() {
				Expect(converter.Convert(&source, &target)).To(Succeed())
				Expect(target.String()).To(Equal("5f68ada0-5c63-4f82-9619-a51bc07d38ec"))
			})

			Context("when the value is wrong", func() {
				BeforeEach(func() {
					source = "wrong"
				})

				It("returns an error", func() {
					Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert string value wrong to type array: uuid: incorrect UUID length: wrong"))
				})
			})
		})
	})

	Context("when the target is nil", func() {
		var target string

		It("returns an error", func() {
			var source *string

			Expect(converter.Convert(source, &target)).To(MatchError("the source must be addressable (a pointer)"))
		})
	})

	Context("when the target is interface", func() {
		Context("when the target is io.Reader", func() {
			var target io.Reader

			BeforeEach(func() {
				target = nil
			})

			Context("when the source is *bytes.Buffer", func() {
				var source bytes.Buffer

				BeforeEach(func() {
					source = bytes.Buffer{}
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).NotTo(BeNil())
				})
			})

			Context("when the source is time", func() {
				var now time.Time

				BeforeEach(func() {
					now = time.Now()
				})

				It("returns an error", func() {
					Expect(converter.Convert(&now, &target)).To(HaveOccurred())
				})
			})
		})

		Context("when the target is generic", func() {
			var target interface{}

			BeforeEach(func() {
				target = nil
			})

			Context("when the source is string", func() {
				var source string

				BeforeEach(func() {
					source = "Jack"
				})

				It("converts the value successfully", func() {
					Expect(converter.Convert(&source, &target)).To(Succeed())
					Expect(target).To(Equal("Jack"))
				})

				Context("when the target is initialized", func() {
					BeforeEach(func() {
						target = 0
						source = "1"
					})

					It("converts the value successfully", func() {
						Expect(converter.Convert(&source, &target)).To(Succeed())
						Expect(target).To(Equal(1))
					})

					Context("when the source is not valid", func() {
						BeforeEach(func() {
							source = "unknown"
						})

						It("returns an error", func() {
							Expect(converter.Convert(&source, &target)).To(MatchError("cannot convert string value unknown to type int: strconv.ParseInt: parsing \"unknown\": invalid syntax"))
						})
					})
				})
			})
		})
	})

	Context("when the target is not a pointer", func() {
		It("returns an error", func() {
			var (
				target string
				source = "jack"
			)

			Expect(converter.Convert(&source, target)).To(MatchError("the target must be a pointer"))
		})
	})

	Context("when the target is not addressable", func() {
		It("returns an error", func() {
			var (
				target *string
				source = "jack"
			)

			Expect(converter.Convert(&source, target)).To(MatchError("the target must be addressable (a pointer)"))
		})
	})
})
