package bcode

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBcode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bcode Suite")
}
