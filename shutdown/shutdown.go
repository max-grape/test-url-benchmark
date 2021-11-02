package shutdown

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultTimeout = time.Second * 10
)

var (
	ErrShutdownTimeout = errors.New("shutdown timeout")
)

type Func func()

type Shutdown struct {
	funcs   []Func
	timeout time.Duration
	sigchan chan os.Signal
}

func New(options ...Option) *Shutdown {
	s := &Shutdown{
		timeout: defaultTimeout,
		sigchan: make(chan os.Signal, 1),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *Shutdown) RegisterShutdownFunc(h Func) {
	s.funcs = append(s.funcs, h)
}

func (s *Shutdown) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	doneChan := make(chan struct{})

	go func() {
		for i := len(s.funcs) - 1; i >= 0; i-- {
			s.funcs[i]()
		}

		close(doneChan)
	}()

	select {
	case <-ctx.Done():
		return ErrShutdownTimeout
	case <-doneChan:
		return nil
	}
}

func (s *Shutdown) WaitForShutdown() error {
	signal.Notify(s.sigchan, syscall.SIGTERM, syscall.SIGINT)

	<-s.sigchan

	return s.Shutdown()
}
