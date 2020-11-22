package http

import (
	"accountapi-client/http/retry"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

var validClientConfig = ClientConfig{
	Retries: retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
	Timeout: time.Second,
	Headers: Headers{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	},
}

type DummyRequest struct {
	Title string `json:"title"`
}

type DummyResponse struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func TestNewClientWithValidConfig(t *testing.T) {
	config := validClientConfig

	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("When creating Client")
	client, err := NewClient(config)

	t.Logf("Should not return any errors")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithInValidConfig(t *testing.T) {
	testCases := []struct {
		Timeout       time.Duration
		Headers       Headers
		Retries       retry.RetriesConfig
		ExpectedError error
	}{
		{
			Retries:       retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
			Timeout:       time.Nanosecond * 0,
			Headers:       Headers{},
			ExpectedError: TimeoutZeroError,
		},
		{
			Retries:       retry.RetriesConfig{MaxRetries: -1, Delay: time.Millisecond, Factor: 2},
			Timeout:       time.Second,
			Headers:       Headers{},
			ExpectedError: retry.MaxRetriesZeroError,
		},
	}

	for _, testCase := range testCases {

		t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", testCase.Retries, testCase.Timeout, testCase.Headers)
		config := ClientConfig{
			Retries: testCase.Retries,
			Timeout: testCase.Timeout,
			Headers: testCase.Headers,
		}

		t.Logf("When creating Client")
		client, err := NewClient(config)

		t.Logf("Should return '%s' error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Nil(t, client)
	}
}

func TestClient_Get(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, DummyResponse{Title: "Jan", Id: 1}))

	defer server.Close()

	t.Logf("When calling GET")
	var dummyResponse DummyResponse
	err := client.Get(context.Background(), server.URL, &dummyResponse)

	t.Logf("Should return DummyResponse")
	assert.NoError(t, err)
	assert.Equal(t, DummyResponse{Title: "Jan", Id: 1}, dummyResponse)
	assert.Equal(t, 1, callCount["/"])
}

func TestClient_GetWithHttpError(t *testing.T) {
	t.Logf("Given HTTP server")
	callCount := make(map[string]int)
	mux := http.NewServeMux()
	mux.Handle("/400", requestHandler(400, &callCount))
	mux.Handle("/404", requestHandler(404, &callCount))
	mux.Handle("/500", requestHandler(500, &callCount))
	mux.Handle("/503", requestHandler(503, &callCount))
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		StatusCode    int
		CallCount     int
		Url           string
		ExpectedError error
	}{
		{StatusCode: 400, CallCount: 1, Url: server.URL + "/400", ExpectedError: &ClientHttpError{Url: server.URL + "/400", StatusCode: 400}},
		{StatusCode: 404, CallCount: 1, Url: server.URL + "/404", ExpectedError: &ClientHttpError{Url: server.URL + "/404", StatusCode: 404}},
		{StatusCode: 500, CallCount: 4, Url: server.URL + "/500", ExpectedError: &ClientHttpError{Url: server.URL + "/500", StatusCode: 500}},
		{StatusCode: 503, CallCount: 4, Url: server.URL + "/503", ExpectedError: &ClientHttpError{Url: server.URL + "/503", StatusCode: 503}},
	}

	for _, testCase := range testCases {
		config := validClientConfig
		t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

		t.Logf("And given Client")
		client, _ := NewClient(config)

		t.Logf("When calling GET")
		var dummyResponse DummyResponse
		err := client.Get(context.Background(), testCase.Url, &dummyResponse)

		t.Logf("Should return ClientHttpError with statusCode %d", testCase.StatusCode)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Equal(t, DummyResponse{}, dummyResponse)
		assert.Equal(t, testCase.CallCount, callCount[fmt.Sprintf("/%d", testCase.StatusCode)])

	}
}

func TestClient_GetWithResponseBodyParsingError(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, "aa"))

	defer server.Close()

	t.Logf("When calling GET")
	var dummyResponse DummyResponse
	err := client.Get(context.Background(), server.URL, &dummyResponse)

	t.Logf("Should return ClientError with parsing message")
	assert.EqualError(t, err, (&ClientError{Message: "parsing error", Url: server.URL, Err: &json.UnmarshalTypeError{
		Value:  "string",
		Type:   reflect.TypeOf(dummyResponse),
		Offset: 4,
	}}).Error())
	assert.Equal(t, DummyResponse{}, dummyResponse)
	assert.Equal(t, 1, callCount["/"])
}

func TestClient_GetWithDialError(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, "aa"))

	defer server.Close()

	t.Logf("When calling GET")
	var dummyResponse DummyResponse
	err := client.Get(context.Background(), "http://localhost:3322/asdasd", &dummyResponse)

	t.Logf("Should return ClientError with network message")
	var expectedError *ClientError
	assert.True(t, errors.As(err, &expectedError))
	assert.Equal(t, "network error", expectedError.Message)
	assert.Equal(t, DummyResponse{}, dummyResponse)
	assert.Equal(t, 0, callCount["/"])
}

func TestClient_Delete(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, DummyResponse{Title: "Jan", Id: 1}))

	defer server.Close()

	t.Logf("When calling DELETE")
	var dummyResponse DummyResponse
	err := client.Delete(context.Background(), server.URL, &dummyResponse)

	t.Logf("Should return DummyResponse")
	assert.NoError(t, err)
	assert.Equal(t, DummyResponse{Title: "Jan", Id: 1}, dummyResponse)
	assert.Equal(t, 1, callCount["/"])
}

func TestClient_DeleteWithHttpError(t *testing.T) {
	t.Logf("Given HTTP server")
	callCount := make(map[string]int)
	mux := http.NewServeMux()
	mux.Handle("/400", requestHandler(400, &callCount))
	mux.Handle("/404", requestHandler(404, &callCount))
	mux.Handle("/500", requestHandler(500, &callCount))
	mux.Handle("/503", requestHandler(503, &callCount))
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		StatusCode    int
		CallCount     int
		Url           string
		ExpectedError error
	}{
		{StatusCode: 400, CallCount: 1, Url: server.URL + "/400", ExpectedError: &ClientHttpError{Url: server.URL + "/400", StatusCode: 400}},
		{StatusCode: 404, CallCount: 1, Url: server.URL + "/404", ExpectedError: &ClientHttpError{Url: server.URL + "/404", StatusCode: 404}},
		{StatusCode: 500, CallCount: 4, Url: server.URL + "/500", ExpectedError: &ClientHttpError{Url: server.URL + "/500", StatusCode: 500}},
		{StatusCode: 503, CallCount: 4, Url: server.URL + "/503", ExpectedError: &ClientHttpError{Url: server.URL + "/503", StatusCode: 503}},
	}

	for _, testCase := range testCases {
		config := validClientConfig
		t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

		t.Logf("And given Client")
		client, _ := NewClient(config)

		t.Logf("When calling DELETE")
		var dummyResponse DummyResponse
		err := client.Delete(context.Background(), testCase.Url, &dummyResponse)

		t.Logf("Should return ClientHttpError with statusCode %d", testCase.StatusCode)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Equal(t, DummyResponse{}, dummyResponse)
		assert.Equal(t, testCase.CallCount, callCount[fmt.Sprintf("/%d", testCase.StatusCode)])

	}
}

func TestClient_DeleteWithResponseBodyParsingError(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, "aa"))

	defer server.Close()

	t.Logf("When calling DELETE")
	var dummyResponse DummyResponse
	err := client.Delete(context.Background(), server.URL, &dummyResponse)

	t.Logf("Should return ClientError with parsing message")
	assert.EqualError(t, err, (&ClientError{Message: "parsing error", Url: server.URL, Err: &json.UnmarshalTypeError{
		Value:  "string",
		Type:   reflect.TypeOf(dummyResponse),
		Offset: 4,
	}}).Error())
	assert.Equal(t, DummyResponse{}, dummyResponse)
	assert.Equal(t, 1, callCount["/"])
}

func TestClient_DeleteWithDialError(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, "aa"))

	defer server.Close()

	t.Logf("When calling DELETE")
	var dummyResponse DummyResponse
	err := client.Delete(context.Background(), "http://localhost:3322/asdasd", &dummyResponse)

	t.Logf("Should return ClientError with network message")
	var expectedError *ClientError
	assert.True(t, errors.As(err, &expectedError))
	assert.Equal(t, "network error", expectedError.Message)
	assert.Equal(t, DummyResponse{}, dummyResponse)
	assert.Equal(t, 0, callCount["/"])
}

func TestClient_Post(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, DummyResponse{Title: "Jan", Id: 1}))

	defer server.Close()

	t.Logf("When calling POST with request")
	var dummyResponse DummyResponse
	err := client.Post(context.Background(), server.URL, &DummyRequest{Title: "Jan"}, &dummyResponse)

	t.Logf("Should return DummyResponse")
	assert.NoError(t, err)
	assert.Equal(t, DummyResponse{Title: "Jan", Id: 1}, dummyResponse)
	assert.Equal(t, 1, callCount["/"])
}

func TestClient_PostWithHttpError(t *testing.T) {
	t.Logf("Given HTTP server")
	callCount := make(map[string]int)
	mux := http.NewServeMux()
	mux.Handle("/400", requestHandler(400, &callCount))
	mux.Handle("/404", requestHandler(404, &callCount))
	mux.Handle("/500", requestHandler(500, &callCount))
	mux.Handle("/503", requestHandler(503, &callCount))
	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		StatusCode    int
		CallCount     int
		Url           string
		ExpectedError error
	}{
		{StatusCode: 400, CallCount: 1, Url: server.URL + "/400", ExpectedError: &ClientHttpError{Url: server.URL + "/400", StatusCode: 400}},
		{StatusCode: 404, CallCount: 1, Url: server.URL + "/404", ExpectedError: &ClientHttpError{Url: server.URL + "/404", StatusCode: 404}},
		{StatusCode: 500, CallCount: 4, Url: server.URL + "/500", ExpectedError: &ClientHttpError{Url: server.URL + "/500", StatusCode: 500}},
		{StatusCode: 503, CallCount: 4, Url: server.URL + "/503", ExpectedError: &ClientHttpError{Url: server.URL + "/503", StatusCode: 503}},
	}

	for _, testCase := range testCases {
		config := validClientConfig
		t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

		t.Logf("And given Client")
		client, _ := NewClient(config)

		t.Logf("When calling POST")
		var dummyResponse DummyResponse
		err := client.Post(context.Background(), testCase.Url, &DummyRequest{Title: "Jan"}, &dummyResponse)

		t.Logf("Should return ClientHttpError with statusCode %d", testCase.StatusCode)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Equal(t, DummyResponse{}, dummyResponse)
		assert.Equal(t, testCase.CallCount, callCount[fmt.Sprintf("/%d", testCase.StatusCode)])

	}
}

func TestClient_PostWithResponseBodyParsingError(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, "aa"))

	defer server.Close()

	t.Logf("When calling POST")
	var dummyResponse DummyResponse
	err := client.Post(context.Background(), server.URL, &DummyRequest{Title: "Jan"}, &dummyResponse)

	t.Logf("Should return ClientError with parsing message")
	assert.EqualError(t, err, (&ClientError{Message: "parsing error", Url: server.URL, Err: &json.UnmarshalTypeError{
		Value:  "string",
		Type:   reflect.TypeOf(dummyResponse),
		Offset: 4,
	}}).Error())
	assert.Equal(t, DummyResponse{}, dummyResponse)
	assert.Equal(t, 1, callCount["/"])
}

func TestClient_PostWithRequestBodyParsingError(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, DummyResponse{Title: "Jan", Id: 1}))

	defer server.Close()

	t.Logf("When calling POST")
	var dummyResponse DummyResponse
	err := client.Post(context.Background(), server.URL, func() {}, &dummyResponse)

	t.Logf("Should return ClientError with parsing message")
	assert.EqualError(t, err, (&ClientError{Message: "body parse error",
		Url: server.URL,
		Err: &json.UnsupportedTypeError{Type: reflect.TypeOf(func() {})},
	}).Error())
	assert.Equal(t, DummyResponse{}, dummyResponse)
	assert.Equal(t, 0, callCount["/"])
}

func TestClient_PostWithDialError(t *testing.T) {
	config := validClientConfig
	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("And given Client")
	client, _ := NewClient(config)

	t.Logf("And HTTP server returning 200 status")
	callCount := make(map[string]int)
	server := httptest.NewServer(requestHandlerWithBody(200, &callCount, "aa"))

	defer server.Close()

	t.Logf("When calling POST on unknown URL")
	var dummyResponse DummyResponse
	err := client.Post(context.Background(), "http://localhost:3322/asdasd", &DummyRequest{Title: "Jan"}, &dummyResponse)

	t.Logf("Should return ClientError with network message")
	var expectedError *ClientError
	assert.True(t, errors.As(err, &expectedError))
	assert.Equal(t, "network error", expectedError.Message)
	assert.Equal(t, DummyResponse{}, dummyResponse)
	assert.Equal(t, 0, callCount["/"])
}

func requestHandler(statusCode int, callCount *map[string]int) http.HandlerFunc {
	return requestHandlerWithBody(statusCode, callCount, nil)
}

func requestHandlerWithBody(statusCode int, callCount *map[string]int, responseBody interface{}) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		(*callCount)[req.RequestURI] = (*callCount)[req.RequestURI] + 1
		res.WriteHeader(statusCode)
		if responseBody != nil {
			js, _ := json.Marshal(responseBody)
			res.Write(js)
		}
	}
}
