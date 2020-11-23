package http

import (
	"errors"
	"fmt"
	"net/url"
)

// Errors thrown by NewClient when ConfigClient has errors
var (
	TimeoutZeroError = errors.New("timeout has to be larger than 0ms")
)

// Throw by the Client on unexpected non-http related issues like parsing, dialing or tls handshake issues
type ClientError struct {
	Url     *url.URL
	Message string
	Err     error
	IsRetryable bool
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("failed to call %s due to %s: %s", e.Url, e.Message, e.Err)
}

func (e *ClientError) Unwrap() error {
	return e.Err
}

// Throw by the Client on server-side http errors, it is returned on anything beyond or equal to HTTP-400
type ClientHttpError struct {
	Url          *url.URL
	StatusCode   int
	ResponseBody []byte
	IsRetryable bool
}

func (e *ClientHttpError) Error() string {
	return fmt.Sprintf("failed to call %s due to HTTP error %d with body `%s`", e.Url, e.StatusCode, e.ResponseBody)
}
