package shutdown

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShutdown_WaitForShutdown(t *testing.T) {
	s := New(WithTimeout(time.Millisecond * 100))

	var (
		fistFunctionCalled, secondFunctionCalled bool
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

	errChan := make(chan error)

	go func() {
		errChan <- s.WaitForShutdown()
	}()

	s.sigchan <- syscall.SIGINT

	assert.NoError(t, <-errChan)

	assert.True(t, fistFunctionCalled)
	assert.True(t, secondFunctionCalled)
}

func TestShutdown_WaitForShutdownTimeout(t *testing.T) {
	s := New(WithTimeout(time.Millisecond * 100))

	s.RegisterShutdownFunc(func() {
		time.Sleep(time.Millisecond * 200)
	})

	errChan := make(chan error)

	go func() {
		errChan <- s.WaitForShutdown()
	}()

	s.sigchan <- syscall.SIGINT

	assert.EqualError(t, <-errChan, ErrShutdownTimeout.Error())
}
