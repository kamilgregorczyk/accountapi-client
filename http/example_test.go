package http_test

import (
	"accountapi-client/http"
	"accountapi-client/http/retry"
	"context"
	"log"
	"time"
)

func ExampleClient_Get() {
	// Basic, valid config
	config := http.ClientConfig{
		Retries: retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
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

	// Actual call
	err = client.Get(context.Background(), "http://localhost:8000", &dummyResponse)

	if err != nil {
		log.Fatal(err)
	}

}

func ExampleClient_Delete() {
	// Basic, valid config
	config := http.ClientConfig{
		Retries: retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
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

	// Actual call
	err = client.Delete(context.Background(), "http://localhost:8000", &dummyResponse)

	if err != nil {
		log.Fatal(err)
	}

}

func ExampleClient_Post() {
	// Basic, valid config
	config := http.ClientConfig{
		Retries: retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
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

	// Actual call
	err = client.Post(context.Background(), "http://localhost:8000", &dummyRequest, &dummyResponse)

	if err != nil {
		log.Fatal(err)
	}

}
