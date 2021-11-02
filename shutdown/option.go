package shutdown

import (
	"time"
)

type Option func(s *Shutdown)

func WithTimeout(timeout time.Duration) Option {
	return func(s *Shutdown) {
		s.timeout = timeout
	}
}
