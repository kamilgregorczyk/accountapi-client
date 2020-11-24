// HTTP client for Form3 Organisation/Account resource.
// https://api-docs.form3.tech/api.html#organisation-accounts
//
// Provides functionality of Create, Fetch, List and Delete operations via http call which are retryable.
//
// Retries can be configured in ClientConfig by providing retry.RetriesConfig
// 	account.ClientConfig{
//		Timeout: time.Second,
//		Logging: true,
//		Url: url.URL{
//			Scheme: "http",
//			Host:   "localhost:8080"},
//		RetriesConfig: retry.RetriesConfig{
//			MaxRetries: 3,
//			Delay:      time.Second,
//			Factor:     1.5,
//		},
//	}
package account

import (
	"accountapi-client/http"
	"accountapi-client/retry"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ClientConfig struct {
	Timeout       time.Duration
	Logging       bool
	Url           *url.URL
	RetriesConfig *retry.RetriesConfig
}

type Client struct {
	Url    *url.URL
	Client *http.Client
}

// Creates new instance of Client.
//
// If ClientConfig.Timeout is zero or bellow it returns TimeoutZeroError.
//
// If ClientConfig.RetriesConfig has any errors, those will be also returned to the caller,
// not providing those values is not possible as retries are required on all of the endpoints.
//
// If Logging is enabled, every outgoing request will be logged along with its execution time, including retries.
func NewClient(config ClientConfig) (*Client, error) {
	client, err := http.NewClient(http.ClientConfig{
		Timeout: config.Timeout,
		Logging: config.Logging,
		Retries: config.RetriesConfig,
		Headers: http.Headers{
			"Content-Type": "application/vnd.api+json",
			"Accept":       "application/vnd.api+json",
		}})
	if err != nil {
		return nil, err
	}
	return &Client{
		Url:    config.Url,
		Client: client,
	}, nil
}

// Creates Account https://api-docs.form3.tech/api.html#organisation-accounts-create
//
// In case of network, parsing or io error (non http related) it will return ClientError.
//
// In case of an http related error (>400 status code) it will return ClientHttpError along with returned status code.
func (c *Client) Create(ctx context.Context, request *CreateAccountRequest) (*CreateAccountResponse, error) {
	err := request.Validate()

	if err != nil {
		return nil, err
	}

	path, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/organisation/accounts", c.Url.String()))
	if err != nil {
		return nil, err
	}
	var createAccountResponse *CreateAccountResponse
	err = c.Client.Post(ctx, path, request, &createAccountResponse)
	return createAccountResponse, err
}

// Retrieves Account https://api-docs.form3.tech/api.html#organisation-accounts-fetch
//
// In case of network, parsing or io error (non http related) it will return ClientError.
//
// In case of an http related error (>400 status code) it will return ClientHttpError along with returned status code.
//
// In case of invalid FetchAccountRequest it will return ValidationError
func (c *Client) Fetch(ctx context.Context, request *FetchAccountRequest) (*FetchAccountResponse, error) {
	err := request.Validate()

	if err != nil {
		return nil, err
	}

	path, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/organisation/accounts/%s", c.Url.String(), request.Id))
	if err != nil {
		return nil, err
	}

	var fetchAccountResponse *FetchAccountResponse
	err = c.Client.Get(ctx, path, &fetchAccountResponse)
	return fetchAccountResponse, err
}

// Lists Account https://api-docs.form3.tech/api.html#organisation-accounts-list
//
// In case of network, parsing or io error (non http related) it will return ClientError.
//
// In case of an http related error (>400 status code) it will return ClientHttpError along with returned status code.
//
// In case of invalid ListAccountsRequest it will return ValidationError
func (c *Client) List(ctx context.Context, request *ListAccountsRequest) (*ListAccountResponse, error) {
	err := request.Validate()

	if err != nil {
		return nil, err
	}

	path, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/organisation/accounts", c.Url.String()))
	if err != nil {
		return nil, err
	}

	query := path.Query()
	query.Set("page[number]", strconv.Itoa(request.PageNumber))
	query.Set("page[size]", strconv.Itoa(request.PageSize))
	path.RawQuery = query.Encode()
	var listAccountsResponse *ListAccountResponse
	err = c.Client.Get(ctx, path, &listAccountsResponse)
	return listAccountsResponse, err
}

// Deletes Account https://api-docs.form3.tech/api.html#organisation-accounts-delete
//
// In case of network, parsing or io error (non http related) it will return ClientError.
//
// In case of an http related error (>400 status code) it will return ClientHttpError along with returned status code.
//
// In case of invalid DeleteAccountRequest it will return ValidationError
func (c *Client) Delete(ctx context.Context, request *DeleteAccountRequest) error {
	err := request.Validate()

	if err != nil {
		return err
	}

	path, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/organisation/accounts/%s", c.Url.String(), request.Id))
	if err != nil {
		return err
	}

	query := path.Query()
	query.Set("version", strconv.Itoa(request.Version))
	path.RawQuery = query.Encode()
	err = c.Client.Delete(ctx, path)
	return err
}
