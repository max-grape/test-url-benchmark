package acceptance

import (
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	externalServer = "backend:8080"
	internalServer = "backend:8090"
	httpClient     = &http.Client{
		Timeout: time.Second * 30,
	}
)

var _ = BeforeSuite(func() {
	//
})

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "acceptance suite")
}

var _ = AfterSuite(func() {
	//
})
