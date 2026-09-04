package retry

import (
	"context"
	"math/rand"
	"time"
)

// IsPermanent marks an error as permanent so WithRetry will not retry it.
type IsPermanent interface {
	Permanent() bool
}

// WithRetry calls fn up to maxAttempts times, applying exponential backoff with
// random jitter between attempts. It stops immediately if fn returns a permanent
// error (one implementing IsPermanent with Permanent() == true) or if ctx is
// cancelled.
func WithRetry(ctx context.Context, maxAttempts int, base time.Duration, fn func() error) error {
	var err error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if pe, ok := err.(IsPermanent); ok && pe.Permanent() {
			return err
		}

		if attempt == maxAttempts-1 {
			break
		}

		delay := backoff(base, attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return err
}

// backoff computes the exponential backoff delay with jitter for the given attempt (0-indexed).
func backoff(base time.Duration, attempt int) time.Duration {
	exp := base * (1 << uint(attempt))
	jitter := time.Duration(rand.Int63n(int64(exp) / 2))
	return exp + jitter
}
