package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/max-grape/test-revo/omap"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Запрос /sites", func() {

	Context("Параметр search отсутствует", func() {
		var response *http.Response

		It("Запрос успешно выполнен", func() {
			u := url.URL{
				Scheme: "http",
				Host:   externalServer,
				Path:   "/sites",
			}

			request, err := http.NewRequest(http.MethodGet, u.String(), nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Код 400 получен в ответ", func() {
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Параметр search присутствует, но пуст", func() {
		var response *http.Response

		It("Запрос успешно выполнен", func() {
			u := url.URL{
				Scheme: "http",
				Host:   externalServer,
				Path:   "/sites",
			}

			query := u.Query()
			query.Set("search", "")

			u.RawQuery = query.Encode()

			request, err := http.NewRequest(http.MethodGet, u.String(), nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = httpClient.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Код 400 получен в ответ", func() {
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})

	Context("Параметр search присутствует", func() {
		var response *http.Response

		It("Запрос успешно выполнен", func() {
			u := url.URL{
				Scheme: "http",
				Host:   externalServer,
				Path:   "/sites",
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

		It("Тело ответа содержит хосты", func() {
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).NotTo(HaveOccurred())

			var actual map[string]interface{}

			Expect(json.Unmarshal(body, &actual)).NotTo(HaveOccurred())

			expected := omap.OrderedMap{
				omap.KeyVal{Key: "stub-external-resource-1:8080", Val: 1000},
				omap.KeyVal{Key: "stub-external-resource-2:8080", Val: 1000},
			}

			for _, item := range expected {
				Expect(actual).To(HaveKey(item.Key))
			}
		})
	})
})
