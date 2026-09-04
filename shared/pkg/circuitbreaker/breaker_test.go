package circuitbreaker_test

import (
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"

	"fleet/shared/pkg/circuitbreaker"
)

var errTest = errors.New("test error")

func newBreaker(failureThreshold uint32) *circuitbreaker.Breaker {
	return circuitbreaker.New(circuitbreaker.Settings{
		Name:             "test",
		MaxRequests:      1,
		Interval:         time.Second,
		Timeout:          100 * time.Millisecond,
		FailureThreshold: failureThreshold,
	})
}

func TestBreaker_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setup         func(b *circuitbreaker.Breaker)
		fn            func() error
		expectedError error
	}{
		{
			name:  "works correctly when closed passes through",
			setup: func(_ *circuitbreaker.Breaker) {},
			fn:    func() error { return nil },
		},
		{
			name: "fails when breaker opens on consecutive failures",
			setup: func(b *circuitbreaker.Breaker) {
				for i := 0; i < 3; i++ {
					_ = b.Execute(func() error { return errTest })
				}
			},
			fn:            func() error { return nil },
			expectedError: gobreaker.ErrOpenState,
		},
		{
			name: "works correctly when half-open recovers to closed",
			setup: func(b *circuitbreaker.Breaker) {
				for i := 0; i < 3; i++ {
					_ = b.Execute(func() error { return errTest })
				}
				// Wait for Timeout to transition to HALF-OPEN
				time.Sleep(150 * time.Millisecond)
			},
			fn: func() error { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := newBreaker(3)
			tt.setup(b)

			err := b.Execute(tt.fn)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}
