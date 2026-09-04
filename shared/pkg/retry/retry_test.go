package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fleet/shared/pkg/retry"
)

var errTransient = errors.New("transient error")

type permanentError struct{ error }

func (permanentError) Permanent() bool { return true }

func TestWithRetry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		maxAttempts   int
		fn            func() func() error
		expectedError error
	}{
		{
			name:        "works correctly on first attempt",
			maxAttempts: 3,
			fn: func() func() error {
				return func() error { return nil }
			},
		},
		{
			name:        "works correctly after retry",
			maxAttempts: 3,
			fn: func() func() error {
				calls := 0
				return func() error {
					calls++
					if calls < 2 {
						return errTransient
					}
					return nil
				}
			},
		},
		{
			name:        "fails when permanent error skips retry",
			maxAttempts: 3,
			fn: func() func() error {
				calls := 0
				return func() error {
					calls++
					return permanentError{errTransient}
				}
			},
			expectedError: permanentError{errTransient},
		},
		{
			name:        "fails when max retries exhausted",
			maxAttempts: 2,
			fn: func() func() error {
				return func() error { return errTransient }
			},
			expectedError: errTransient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := retry.WithRetry(context.Background(), tt.maxAttempts, time.Millisecond, tt.fn())

			assert.Equal(t, tt.expectedError, err)
		})
	}
}
