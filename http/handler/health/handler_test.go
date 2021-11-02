package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/max-grape/test-revo/http/handler/health"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP_Unhealthy(t *testing.T) {
	h := health.New()

	h.RegisterHealthCheckFunc(func() bool {
		return false
	})

	h.RegisterHealthCheckFunc(func() bool {
		return true
	})

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.ServeHTTP(rr, r)

	assert.Equal(t, health.StatusUnhealthy, rr.Body.String())
}

func TestHandler_ServeHTTP_Healthy(t *testing.T) {
	h := health.New()

	h.RegisterHealthCheckFunc(func() bool {
		return true
	})

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.ServeHTTP(rr, r)

	assert.Equal(t, health.StatusHealthy, rr.Body.String())
}
