package bcode

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test bcode package", func() {
	It("Test New bcode funtion", func() {
		bcode := NewBcode(400, 4000, "test")
		Expect(bcode).ShouldNot(BeNil())
		Expect(bcode.Message).ShouldNot(BeNil())
		Expect(bcode.Error()).ShouldNot(BeNil())
	})
})
