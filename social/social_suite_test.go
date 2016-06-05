package social_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSocial(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Social Suite")
}
