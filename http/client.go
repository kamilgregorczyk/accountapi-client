// Provides an HTTP client, an abstraction over go's built in http client.
// It simplifies making http requests by doing boilerplate logic, serialising & deserializing JSON responses
// and by encapsulating all the errors that can happen along the way under once common format.
//
// The NewClient function creates a new instance of the Client by providing ClientConfig.
// It is required to pass all the fields from that config.
//
// For the time being it only provides GET, POST and DELETE operations.
//
// All calls support retries which can be defined in ClientConfig via retry.RetriesConfig, if you wish to disable them then set MaxRetries to 0
//
package http

import (
	"accountapi-client/retry"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	corehttp "net/http"
	"net/url"
	"time"
)

type ClientConfig struct {
	Timeout time.Duration
	Retries *retry.RetriesConfig
	Headers Headers
	Logging bool
}

type Client struct {
	client  *corehttp.Client
	retry   *retry.Retry
	headers Headers
	logging bool
}

type Headers map[string]string

// Creates a new instance of the Client.
//
// If ClientConfig.Timeout is zero or bellow it returns TimeoutZeroError.
//
// If ClientConfig.RetriesConfig has any errors, those will be also returned to the caller,
// not providing those values is not possible as retries are required on all of the endpoints.
//
// If Headers won't be empty, all the headers will be set on every outgoing http request.
//
// If Logging is enabled, every outgoing request will be logged along with its execution time, including retries.
func NewClient(config ClientConfig) (*Client, error) {
	if config.Timeout.Milliseconds() <= 0 {
		return nil, TimeoutZeroError
	}

	retry, err := retry.NewRetries(config.Retries)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:  &corehttp.Client{Timeout: config.Timeout},
		retry:   retry,
		headers: config.Headers,
		logging: config.Logging,
	}, nil
}

// Runs GET HTTP query for provided url, responseBody (pointer) will be written by json.Unmarshal.
//
// In case of network, parsing or io error (non http related) it will return ClientError.
//
// In case of an http related error (>400 status code) it will return ClientHttpError along with returned status code.
func (c *Client) Get(ctx context.Context, url *url.URL, responseBody interface{}) error {
	method := "GET"
	request, err := c.createRequest(ctx, method, url, nil)
	if err != nil {
		return err
	}

	err = c.executeWithRetry(request, responseBody)
	if err != nil {
		return err
	}

	return nil
}

// Runs DELETE HTTP query for provided url.
//
// In case of network, parsing or io error (non http related) it will return ClientError.
//
// In case of an http related error (>400 status code) it will return ClientHttpError along with returned status code & body.
func (c *Client) Delete(ctx context.Context, url *url.URL) error {
	method := "DELETE"
	request, err := c.createRequest(ctx, method, url, nil)
	if err != nil {
		return err
	}

	err = c.executeWithRetry(request, nil)
	if err != nil {
		return err
	}
	return nil
}

// Runs POST HTTP query for provided url, responseBody (pointer) will be written by json.Unmarshal.
//
// In case of network, parsing or io error (non http related) it will return ClientError.
//
// In case of an http related error (>400 status code) it will return ClientHttpError along with returned status code.
func (c *Client) Post(ctx context.Context, url *url.URL, requestBody interface{}, responseBody interface{}) error {
	method := "POST"
	request, err := c.createRequest(ctx, method, url, requestBody)
	if err != nil {
		return err
	}

	err = c.executeWithRetry(request, responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) createRequest(context context.Context, method string, url *url.URL, requestBody interface{}) (resp *corehttp.Request, err error) {
	marshaledBody, err := json.Marshal(requestBody)

	if err != nil {
		return nil, &ClientError{Message: "body parse error", Url: url, Err: err}
	}

	req, err := corehttp.NewRequestWithContext(context, method, url.String(), bytes.NewBuffer(marshaledBody))
	if err != nil {
		return nil, &ClientError{Message: "network error", Url: url, Err: err}
	}
	c.setHeaders(req)
	return req, nil
}

func (c *Client) setHeaders(req *corehttp.Request) {
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
}

func (c *Client) executeWithRetry(request *corehttp.Request, responseBody interface{}) error {
	return c.retry.Execute(func() error {
		c.logNewRequest(request.Method, request.URL)
		startTime := time.Now()
		response, err := c.client.Do(request)
		c.logFinishedRequest(request.Method, request.URL, time.Now().Sub(startTime), response)
		encapsulatedErr := c.readResponse(response, err, request.URL, responseBody)

		if c.shouldRetry(encapsulatedErr) {
			return &retry.RetryableError{Err: encapsulatedErr}
		}
		return encapsulatedErr
	})
}

func (c *Client) logNewRequest(method string, url *url.URL) {
	if c.logging {
		log.Printf("Outgoing request to [%s] [%s] \n", method, url)
	}
}

func (c *Client) logFinishedRequest(method string, url *url.URL, elapsed time.Duration, response *corehttp.Response) {
	if !c.logging {
		return
	}
	if response != nil && response.StatusCode >= 400 {
		log.Printf("Outgoing request to [%s] [%s] failed with status [%d] in [%s] \n", method, url, response.StatusCode, elapsed.String())
	} else {
		log.Printf("Outgoing request to [%s] [%s] completed in [%s] \n", method, url, elapsed.String())
	}
}

func (c *Client) readResponse(response *corehttp.Response, err error, url *url.URL, responseBody interface{}) error {
	if err != nil || response == nil {
		return &ClientError{Message: "network error", Url: url, Err: err, IsRetryable: true}
	}

	buffer, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return &ClientError{Message: "io error", Url: url, Err: err, IsRetryable: false}
	}

	if response.StatusCode >= 400 {
		return &ClientHttpError{Url: url, StatusCode: response.StatusCode, ResponseBody: buffer, IsRetryable: response.StatusCode >= 500}
	}

	if responseBody == nil {
		return nil
	}
	err = json.Unmarshal(buffer, responseBody)
	if err != nil {
		return &ClientError{Message: "parsing error", Url: url, Err: err, IsRetryable: false}
	}
	return nil
}

func (c *Client) shouldRetry(err error) bool {
	switch err.(type) {
	case *ClientError:
		var clientError *ClientError
		errors.As(err, &clientError)
		return clientError.IsRetryable
	case *ClientHttpError:
		var clientError *ClientHttpError
		errors.As(err, &clientError)
		return clientError.IsRetryable
	case nil:
		return false
	default:
		return true
	}
}
