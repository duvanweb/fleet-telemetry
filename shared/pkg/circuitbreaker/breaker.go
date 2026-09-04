package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

// Breaker wraps a gobreaker.CircuitBreaker with a simpler Execute interface.
type Breaker struct {
	cb *gobreaker.CircuitBreaker
}

// Settings configures the circuit breaker thresholds.
type Settings struct {
	Name             string
	MaxRequests      uint32
	Interval         time.Duration
	Timeout          time.Duration
	FailureThreshold uint32
}

// New creates a Breaker with the given settings.
func New(s Settings) *Breaker {
	st := gobreaker.Settings{
		Name:        s.Name,
		MaxRequests: s.MaxRequests,
		Interval:    s.Interval,
		Timeout:     s.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= s.FailureThreshold
		},
	}
	return &Breaker{cb: gobreaker.NewCircuitBreaker(st)}
}

// Execute runs fn through the circuit breaker. Returns gobreaker.ErrOpenState
// when the breaker is open.
func (b *Breaker) Execute(fn func() error) error {
	_, err := b.cb.Execute(func() (any, error) {
		return nil, fn()
	})
	return err
}

// State returns the current state of the circuit breaker.
func (b *Breaker) State() gobreaker.State {
	return b.cb.State()
}
