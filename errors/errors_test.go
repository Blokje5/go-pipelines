package errors

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Optional", func() {
	var valOpt Optional
	var errOpt Optional

	BeforeEach(func() {
		valOpt = Optional{
			Val: 1,
			Error: nil,
		}

		errOpt = Optional{
			Val: nil,
			Error: errors.New("error"),
		}
	})

	Describe("MustGetVal", func() {
		It("Should return a value if the Optional contains a value", func() {
			val := valOpt.MustGetVal()
			Expect(val).To(Equal(1))
		})

		It("Should panic if there is an Error in the Optional", func() {
			Expect(func() {errOpt.MustGetVal()}).To(PanicWith("error"))
		})
	})
})