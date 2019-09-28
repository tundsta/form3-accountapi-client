package accountapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAccountapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Accountapi Suite")
}
