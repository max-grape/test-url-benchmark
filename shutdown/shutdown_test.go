package shutdown_test

import (
	"testing"
	"time"

	"github.com/max-grape/test-revo/shutdown"
	"github.com/stretchr/testify/assert"
)

func TestShutdown_Shutdown(t *testing.T) {
	s := shutdown.New(shutdown.WithTimeout(time.Millisecond * 100))

	var (
		fistFunctionCalled   bool
		secondFunctionCalled bool
	)

	s.RegisterShutdownFunc(func() {
		if !secondFunctionCalled {
			assert.Fail(t, "invalid functions order")
		}

		fistFunctionCalled = true
	})

	s.RegisterShutdownFunc(func() {
		if fistFunctionCalled {
			assert.Fail(t, "invalid functions order")
		}

		secondFunctionCalled = true
	})

	assert.NoError(t, s.Shutdown())

	assert.True(t, fistFunctionCalled)
	assert.True(t, secondFunctionCalled)
}

func TestShutdown_ShutdownTimeout(t *testing.T) {
	s := shutdown.New(shutdown.WithTimeout(time.Millisecond * 100))

	s.RegisterShutdownFunc(func() {
		time.Sleep(time.Millisecond * 200)
	})

	assert.EqualError(t, s.Shutdown(), shutdown.ErrShutdownTimeout.Error())
}
