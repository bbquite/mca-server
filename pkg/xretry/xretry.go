package xretry

import (
	"time"
)

// RetryPolicy is the retry policy
type RetryPolicy struct {
	immediateRetries   int
	retriesWithBackoff int
	delay              time.Duration
	backoffFactor      float64
}

// RetryPolicyOption is the option for the retry policy
type RetryPolicyOption func(*RetryPolicy)

// WithImmediateRetries sets the immediate retries
func WithImmediateRetries(retries int) RetryPolicyOption {
	return func(p *RetryPolicy) {
		p.immediateRetries = retries
	}
}

// WithRetriesWithBackoff sets the retries with backoff
func WithRetriesWithBackoff(retries int, delay time.Duration, backoffFactor float64) RetryPolicyOption {
	return func(p *RetryPolicy) {
		p.retriesWithBackoff = retries
		p.delay = delay
		p.backoffFactor = backoffFactor
	}
}

// NewRetryPolicy creates a new RetryPolicy
func NewRetryPolicy(opts ...RetryPolicyOption) RetryPolicy {
	p := RetryPolicy{
		immediateRetries:   0,
		retriesWithBackoff: 0,
		delay:              0,
		backoffFactor:      0,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// Retrier is the interface that wraps the Retry method
type Retrier struct {
	p RetryPolicy
}

// NewRetrier creates a new Retrier
func NewRetrier(p RetryPolicy) *Retrier {
	return &Retrier{
		p: p,
	}
}

// Retry will retry the given function
func (r *Retrier) Retry(f func() error) error {
	// err := immediatelyRetry(f, r.p.immediateRetries)
	// if err != nil {
	// 	err = retryWithBackoff(f, r.p.retriesWithBackoff, r.p.delay, r.p.backoffFactor)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	err := retryWithBackoff(f, r.p.retriesWithBackoff, r.p.delay, r.p.backoffFactor)
	if err != nil {
		return err
	}

	return nil
}

func retryWithBackoff(f func() error, retriesLeft int, delay time.Duration, backoff float64) error {
	err := f()
	if err == nil {
		return nil
	}

	if retriesLeft == 0 {
		return err
	}

	time.Sleep(delay)
	return retryWithBackoff(f, retriesLeft-1, time.Duration(float64(delay)*backoff), backoff)
}

func immediatelyRetry(f func() error, retriesLeft int) error {
	err := f()
	if err == nil {
		return nil
	}

	if retriesLeft == 0 {
		return err
	}

	return immediatelyRetry(f, retriesLeft-1)
}

// EXAMPLE
// Define a retry policy
// policy := xretry.NewRetryPolicy(
// 	xretry.WithImmediateRetries(3),
// 	xretry.WithRetriesWithBackoff(1, 1*time.Second, 1.5),
// )

// // Create a new retrier with the policy
// retrier := xretry.NewRetrier(policy)

// // Define a function that may fail
// f := unstableOperation

// // Use the retrier to retry the function if it fails
// err := retrier.Retry(f)
// if err != nil {
// 	fmt.Println("Operation failed after retries:", err)
// }
