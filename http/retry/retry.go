// Package that provides functionality of retries to http clients with delays
// calculated based on exponential backoff strategy.
//
// The NewRetries function creates new Retry with based on RetryConfig with validation.
//
// The Execute function should be called whenever caller needs to run an http query with eventual retries.
// Caller has to provide the logic that will be retried with RetryFunc.
//
// In order for Retry to execute again provided RetryFunc, the caller has to return RetryableError,
// otherwise the execution will be treated as successfully no matter it error of other type is returned or not
// (as some errors are not worth to retry)
//
// The delay between retries is calculated based on a simple exponential-backoff equation: delay * factor^currentTry
// Providing delay of 1 second, factor 2.0  and maximum number of retires will retry in 1s, 3s and 7s of delay between runs
//
package retry

import (
	"errors"
	"math"
	"net/http"
	"time"
)

type RetriesConfig struct {
	MaxRetries int
	Delay      time.Duration
	Factor     float64
}

// Constructs new Retry from RetriesConfig
// If RetriesConfig.MaxRetries is zero or below, it returns MaxRetriesZeroError
// If RetriesConfig.Delay is zero or below, it returns DelayZeroError
// If RetriesConfig.Factor is zero or below, it returns FactorZeroError
func NewRetries(config RetriesConfig) (*Retry, error) {
	if config.MaxRetries <= 0 {
		return nil, MaxRetriesZeroError
	}
	if config.Delay.Milliseconds() <= 0 {
		return nil, DelayZeroError
	}
	if config.Factor <= 0 {
		return nil, FactorZeroError
	}

	return &Retry{config: config}, nil
}

// Constructed with NewRetry, contains Execute function for running HTTP requests with retries
type Retry struct {
	config RetriesConfig
}

type RetryFunc func() (*http.Response, error)

// Runs HTTP requests with retries
// In order for Execute to run again provided RetryFunc, the caller has to return RetryableError,
// otherwise the execution will be treated as successfully no matter it error of other type is returned or not
// (as some errors are not worth to retry)
//
// The delay between retries is calculated based on a simple exponential-backoff equation: delay * factor^currentTry
// Providing delay of 1 second, factor 2.0  and maximum number of retires will retry in 1s, 3s and 7s of delay between runs
func (r *Retry) Execute(runnable RetryFunc) (*http.Response, error) {
	var tryCount int
	for {
		response, err := runnable()
		if err == nil || tryCount >= r.config.MaxRetries {
			return response, err
		}

		var retryError *RetryableError
		if !errors.As(err, &retryError) {
			return response, err
		}

		tryCount++
		select {
		case <-time.After(r.next(tryCount)):
		}

	}

}

func (r *Retry) next(currentTry int) time.Duration {
	delay := math.Abs(float64(r.config.Delay.Nanoseconds()) * (math.Pow(r.config.Factor, float64(currentTry)) - 1.0))
	return time.Duration(delay)

}
