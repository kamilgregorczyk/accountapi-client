package http_test

import (
	"accountapi-client/http"
	"accountapi-client/retry"
	"context"
	"log"
	"net/url"
	"time"
)

func ExampleClient_Get() {
	// Basic, valid config
	config := http.ClientConfig{
		Retries: &retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
		Timeout: time.Second,
		Headers: http.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	// New client
	client, err := http.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	// Expected response structure
	dummyResponse := struct {
		Id    int
		Title string
	}{}
	toCall, _ := url.Parse("http://localhost:8000/test")

	// Actual call
	err = client.Get(context.Background(), toCall, &dummyResponse)

	if err != nil {
		log.Fatal(err)
	}

}

func ExampleClient_Delete() {
	// Basic, valid config
	config := http.ClientConfig{
		Retries: &retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
		Timeout: time.Second,
		Headers: http.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	// New client
	client, err := http.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}
	toCall, _ := url.Parse("http://localhost:8000/test")

	// Actual call
	err = client.Delete(context.Background(), toCall)

	if err != nil {
		log.Fatal(err)
	}

}

func ExampleClient_Post() {
	// Basic, valid config
	config := http.ClientConfig{
		Retries: &retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
		Timeout: time.Second,
		Headers: http.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	// New client
	client, err := http.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	dummyRequest := struct {
		Title string
	}{Title: "John"}

	// Expected response structure
	dummyResponse := struct {
		Id    int
		Title string
	}{}
	toCall, _ := url.Parse("http://localhost:8000/test")

	// Actual call
	err = client.Post(context.Background(), toCall, &dummyRequest, &dummyResponse)

	if err != nil {
		log.Fatal(err)
	}

}
