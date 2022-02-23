package platformclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlatformClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PlatformClient Suite")
}
