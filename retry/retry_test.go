package retry

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRetriesWithValidConfig(t *testing.T) {
	config := RetriesConfig{
		MaxRetries: 1,
		Delay:      time.Second,
		Factor:     1.0,
	}
	t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", config.MaxRetries, config.Delay, config.Factor)

	t.Logf("When creating Retry")
	retry, err := NewRetries(&config)

	t.Logf("Should not return any errors")

	assert.NoError(t, err)
	assert.NotNil(t, retry)
}

func TestNewRetriesWithInValidConfig(t *testing.T) {
	testCases := []struct {
		MaxRetries    int
		Delay         time.Duration
		Factor        float64
		ExpectedError error
	}{
		{MaxRetries: 0, Delay: time.Second, Factor: 1.0, ExpectedError: MaxRetriesZeroError},
		{MaxRetries: -1, Delay: time.Second, Factor: 1.0, ExpectedError: MaxRetriesZeroError},
		{MaxRetries: 1, Delay: 0 * time.Second, Factor: 1.0, ExpectedError: DelayZeroError},
		{MaxRetries: 1, Delay: -1 * time.Second, Factor: 1.0, ExpectedError: DelayZeroError},
		{MaxRetries: 1, Delay: time.Second, Factor: 0, ExpectedError: FactorZeroError},
		{MaxRetries: 1, Delay: time.Second, Factor: -1.0, ExpectedError: FactorZeroError},
	}
	for _, testCase := range testCases {
		t.Logf("Given invalid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", testCase.MaxRetries, testCase.Delay, testCase.Factor)
		config := RetriesConfig{
			MaxRetries: testCase.MaxRetries,
			Delay:      testCase.Delay,
			Factor:     testCase.Factor,
		}

		t.Logf("When creating Retry")
		retry, err := NewRetries(&config)

		t.Logf("Should return '%s' error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Nil(t, retry)
	}
}

func TestRetryWithSuccessAtFirstTry(t *testing.T) {
	maxRetries := 3
	delay := time.Millisecond
	factor := 1.0
	t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", maxRetries, delay, factor)
	config := RetriesConfig{
		MaxRetries: maxRetries,
		Delay:      delay,
		Factor:     factor,
	}
	t.Logf("And given Retry")
	retry, _ := NewRetries(&config)

	t.Logf("And given a func to run")
	var callCount int
	funcToRetry := func() error {
		callCount++
		return nil
	}

	t.Logf("When executing a func")
	err := retry.Execute(funcToRetry)

	t.Logf("Should call only once and not return any errors")
	assert.Equal(t, callCount, 1)
	assert.NoError(t, err)
}

func TestRetryWithInitialFailuresAndThenSuccess(t *testing.T) {
	numberOfRetries := []int{1, 2, 3}
	for _, retryCount := range numberOfRetries {
		maxRetries := 3
		delay := time.Millisecond
		factor := 1.0
		t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", maxRetries, delay, factor)
		config := RetriesConfig{
			MaxRetries: maxRetries,
			Delay:      delay,
			Factor:     factor,
		}
		t.Logf("And given Retry")
		retry, _ := NewRetries(&config)

		t.Logf("And given a func to run")
		var callCount int
		expectedCallCount := retryCount + 1
		funcToRetry := func() error {
			if retryCount > callCount {
				callCount++
				return &RetryableError{}
			}
			callCount++
			return nil
		}

		t.Logf("When executing a func")
		err := retry.Execute(funcToRetry)

		t.Logf("Should call function %d times and not return any errors", expectedCallCount)
		assert.Equal(t, callCount, expectedCallCount)
		assert.NoError(t, err)
	}

}

func TestRetryWithConstantFailures(t *testing.T) {
	maxRetries := 3
	delay := time.Millisecond
	factor := 1.0
	t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", maxRetries, delay, factor)
	config := RetriesConfig{
		MaxRetries: maxRetries,
		Delay:      delay,
		Factor:     factor,
	}
	t.Logf("And given Retry")
	retry, _ := NewRetries(&config)

	t.Logf("And given a func to run")
	var callCount int
	expectedError := errors.New("something is wrong")
	funcToRetry := func() error {
		callCount++
		return &RetryableError{Err: expectedError}
	}

	t.Logf("When executing a func")
	err := retry.Execute(funcToRetry)

	t.Logf("Should call function %d times and return unwrapped error", 4)
	assert.Equal(t, callCount, 4)
	assert.EqualError(t, err, expectedError.Error())
}

func TestExponentialBackoff(t *testing.T) {
	testCases := []struct {
		MaxRetries    int
		Delay         time.Duration
		Factor        float64
		RetryCount    int
		ExpectedDelay time.Duration
	}{

		{MaxRetries: 3, Delay: time.Second, Factor: 2.0, RetryCount: 1, ExpectedDelay: time.Second * 1},
		{MaxRetries: 3, Delay: time.Second, Factor: 2.0, RetryCount: 2, ExpectedDelay: time.Second * 3},
		{MaxRetries: 3, Delay: time.Second, Factor: 2.0, RetryCount: 3, ExpectedDelay: time.Second * 7},

		{MaxRetries: 3, Delay: time.Second, Factor: 1.5, RetryCount: 1, ExpectedDelay: time.Millisecond * 500},
		{MaxRetries: 3, Delay: time.Second, Factor: 1.5, RetryCount: 2, ExpectedDelay: time.Second*1 + time.Millisecond*250},
		{MaxRetries: 3, Delay: time.Second, Factor: 1.5, RetryCount: 3, ExpectedDelay: time.Second*2 + time.Millisecond*375},

		{MaxRetries: 3, Delay: time.Second, Factor: 1.6, RetryCount: 1, ExpectedDelay: time.Millisecond * 600},
		{MaxRetries: 3, Delay: time.Second, Factor: 1.5, RetryCount: 1, ExpectedDelay: time.Millisecond * 500},
		{MaxRetries: 3, Delay: time.Second, Factor: 0.5, RetryCount: 1, ExpectedDelay: time.Millisecond * 500},
	}

	for _, testCase := range testCases {
		t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", testCase.MaxRetries, testCase.Delay, testCase.Factor)
		config := RetriesConfig{
			MaxRetries: testCase.MaxRetries,
			Delay:      testCase.Delay,
			Factor:     testCase.Factor,
		}
		t.Logf("And given Retry")
		retry, _ := NewRetries(&config)

		t.Logf("When calculating backoff")
		delay := retry.next(testCase.RetryCount)

		t.Logf("Delay should be %s", testCase.ExpectedDelay.String())
		assert.Equal(t, testCase.ExpectedDelay, delay)
	}
}
