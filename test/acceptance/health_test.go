package acceptance

import (
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Запрос /health", func() {

	Context("", func() {
		var response *http.Response

		It("Запрос успешно выполнен", func() {
			u := url.URL{
				Scheme: "http",
				Host:   internalServer,
				Path:   "/health",
			}

			query := u.Query()
			query.Set("search", "foo")

			u.RawQuery = query.Encode()

			request, err := http.NewRequest(http.MethodGet, u.String(), nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Код 200 получен в ответ", func() {
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		It("Тело ответа содержит один из допустимых статусов", func() {
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(body)).To(BeElementOf("healthy", "unhealthy"))
			Expect(response.Body.Close()).NotTo(HaveOccurred())
		})
	})
})
