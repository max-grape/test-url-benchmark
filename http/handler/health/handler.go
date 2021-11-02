package health

import (
	"net/http"

	"github.com/improbable-eng/go-httpwares/logging/logrus/ctxlogrus"
)

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

type CheckFunc func() bool

type Handler struct {
	funcs []CheckFunc
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterHealthCheckFunc(f CheckFunc) {
	if h == nil {
		return
	}

	h.funcs = append(h.funcs, f)
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctxlogrus.Extract(ctx)

	for _, f := range h.funcs {
		if !f() {
			if _, err := rw.Write([]byte(StatusUnhealthy)); err != nil {
				logger.Errorf("health server failed to write status unhealthy: %+v", err)
			}

			return
		}
	}

	if _, err := rw.Write([]byte(StatusHealthy)); err != nil {
		logger.Errorf("health server failed to write status healthy: %+v", err)
	}
}
