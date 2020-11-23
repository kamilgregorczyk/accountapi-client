package account_test

import (
	"accountapi-client/account"
	"accountapi-client/retry"
	"context"
	"log"
	"net/url"
	"time"
)

func ExampleNewClient() {
	// Create config
	config := account.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080"},
		RetriesConfig: &retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.5,
		},
	}

	// Create new client
	account.NewClient(config)
}

func ExampleClient_Fetch() {
	// Create config
	config := account.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080"},
		RetriesConfig: &retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.5,
		},
	}

	// Create new client
	accountClient, err := account.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	// Fetch account
	fetchRequest := account.FetchAccountRequest{Id: "fb1ff76f-f360-403f-a324-4bfe2f215895"}
	fetchResponse, err := accountClient.Fetch(context.Background(), &fetchRequest)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(fetchResponse)
}

func ExampleClient_Create() {
	// Create config
	config := account.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080"},
		RetriesConfig: &retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.5,
		},
	}

	// Create new client
	accountClient, err := account.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	// Create account
	createRequest := account.CreateAccountRequest{Account: &account.Account{Id: "fb1ff76f-f360-403f-a324-4bfe2f215895"}}
	createResponse, err := accountClient.Create(context.Background(), &createRequest)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(createResponse)
}

func ExampleClient_List() {
	// Create config
	config := account.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080"},
		RetriesConfig: &retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.5,
		},
	}

	// Create new client
	accountClient, err := account.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	// List accounts
	listRequest := account.ListAccountsRequest{PageSize: 100, PageNumber: 3}
	listResponse, err := accountClient.List(context.Background(), &listRequest)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(listResponse)
}

func ExampleClient_Delete() {
	// Create config
	config := account.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: &url.URL{
			Scheme: "http",
			Host:   "localhost:8080"},
		RetriesConfig: &retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.5,
		},
	}

	// Create new client
	accountClient, err := account.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	// Deletes account
	deleteRequest := account.DeleteAccountRequest{Id: "fb1ff76f-f360-403f-a324-4bfe2f215895", Version: 3}
	err = accountClient.Delete(context.Background(), &deleteRequest)

	if err != nil {
		log.Fatal(err)
	}

}
