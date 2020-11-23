package retry_test

import (
	retries "accountapi-client/http/retry"
	"log"
	"net/http"
	"time"
)

func ExampleRetry_Execute() {
	config := retries.RetriesConfig{
		MaxRetries: 3,
		Delay:      time.Millisecond * 500,
		Factor:     1.3,
	}

	retry, err := retries.NewRetries(config)

	if err != nil {
		log.Fatal(err)
	}
	retry.Execute(func() error {
		response, err := http.Get("http://localhost")
		// We need retries only on 500s and higher
		if response.StatusCode >= 500 {
			return &retries.RetryableError{Err: err}
		}
		return err
	})
}

func ExampleNewRetries() {
	client, err := retries.NewRetries(retries.RetriesConfig{
		MaxRetries: 3,
		Delay:      time.Millisecond * 500,
		Factor:     1.3,
	})

	client, err = retries.NewRetries(retries.RetriesConfig{
		MaxRetries: 5,
		Delay:      time.Second,
		Factor:     2,
	})

	client, err = retries.NewRetries(retries.RetriesConfig{
		MaxRetries: 1,
		Delay:      time.Second,
		Factor:     1,
	})

	log.Print(client, err)
}
