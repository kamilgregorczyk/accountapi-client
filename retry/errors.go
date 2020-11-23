package retry

import "errors"

// Errors returned during creation of the Retry by NewRetries
var (
	MaxRetriesZeroError = errors.New("maxRetries has to be larger than 0")
	DelayZeroError      = errors.New("delay has to be larger than 0")
	FactorZeroError     = errors.New("factor has to be larger than 0")
)

// Returned by the caller within Retry.Execute whenever there's a need to do a retry.
// Returning an error of a different type means there should be no retry
type RetryableError struct {
	Err error
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

func (e *RetryableError) Error() string {
	if e.Err == nil {
		return "retryable: <nil>"
	}
	return "retryable: " + e.Err.Error()
}
