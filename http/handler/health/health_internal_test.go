package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_RegisterHealthCheckFunc(t *testing.T) {
	h := New()

	h.RegisterHealthCheckFunc(func() bool {
		return false
	})

	h.RegisterHealthCheckFunc(func() bool {
		return true
	})

	assert.Len(t, h.funcs, 2)
	assert.False(t, h.funcs[0]())
	assert.True(t, h.funcs[1]())
}
