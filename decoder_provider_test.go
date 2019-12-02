package reflectify_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/reflectify"
)

var _ = Describe("Set", func() {
	type Target struct {
		FirstName string `field:"first_name"`
		LastName  string `field:"last_name"`
	}

	type Source struct {
		Name   string `field:"first_name"`
		Family string `field:"last_name"`
	}

	It("sets the values successfully", func() {
		s := &Source{
			Name:   "John",
			Family: "Peterson",
		}
		t := &Target{}

		Expect(reflectify.Set(t, s)).To(Succeed())
		Expect(t.FirstName).To(Equal("John"))
		Expect(t.LastName).To(Equal("Peterson"))
	})

	Describe("New", func() {
		It("returns the provider", func() {
			provider := &reflectify.ValueProvider{}
			Expect(provider.New(reflect.Value{})).NotTo(BeNil())
		})
	})
})

var _ = Describe("Defaults", func() {
	type Target struct {
		FirstName string `default:"John"`
		LastName  string `default:"Doe"`
		Age       int    `default:"22"`
	}

	It("sets the default values successfully", func() {
		t := &Target{
			LastName: "Peterson",
		}

		Expect(reflectify.SetDefaults(t)).To(Succeed())
		Expect(t.FirstName).To(Equal("John"))
		Expect(t.LastName).To(Equal("Peterson"))
		Expect(t.Age).To(Equal(22))
	})

	Describe("New", func() {
		It("returns the provider", func() {
			provider := &reflectify.DefaultProvider{}
			Expect(provider).To(Equal(provider.New(reflect.Value{})))
		})
	})
})
