package account

import (
	"accountapi-client/http"
	"accountapi-client/retry"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net/url"
	"os"
	"testing"
	"time"
)

var (
	ApiPath         = os.Getenv("ACCOUNT_API_URL")
	ValidAttributes = Attributes{
		Country:                 "PL",
		BaseCurrency:            "PLN",
		BankId:                  "400300",
		AccountNumber:           "112233",
		Bic:                     "NWBKGB22",
		Iban:                    "PL87109024022347282157329766",
		CustomerId:              "22331",
		AccountClassification:   "Personal",
		JointAccount:            true,
		AccountMatchingOptOut:   true,
		SecondaryIdentification: "lets check",
		BankIdCode:              "GBDSC",
	}
)

func TestNewClient(t *testing.T) {

	testCases := []struct {
		Timeout       time.Duration
		Retries       retry.RetriesConfig
		ExpectedError error
	}{
		{
			Retries:       retry.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
			Timeout:       time.Nanosecond * 0,
			ExpectedError: http.TimeoutZeroError,
		},
		{
			Retries:       retry.RetriesConfig{MaxRetries: -1, Delay: time.Millisecond, Factor: 2},
			Timeout:       time.Second,
			ExpectedError: retry.MaxRetriesZeroError,
		},
	}

	for _, testCase := range testCases {

		t.Logf("Given valid ClientConfig retries=%+v timeout=%s", testCase.Retries, testCase.Timeout)
		config := ClientConfig{
			RetriesConfig: &testCase.Retries,
			Timeout:       testCase.Timeout,
			Url:           &url.URL{},
		}

		t.Logf("When creating Client")
		client, err := NewClient(config)

		t.Logf("Should return '%s' error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Nil(t, client)
	}
}

func TestClientCreate(t *testing.T) {
	t.Logf("Given valid HTTP client")
	client := initClient()

	t.Logf("And given valid create request")
	request := validCreateAccountRequest()

	t.Logf("When creating account")
	response, err := client.Create(context.Background(), request)

	t.Logf("Should not return any errors")
	assert.NoError(t, err)
	assert.Equal(t, response, &CreateAccountResponse{Account: &Account{
		Type:           "accounts",
		Id:             request.Account.Id,
		OrganisationId: request.Account.OrganisationId,
		CreatedOn:      response.Account.CreatedOn,
		ModifiedOn:     response.Account.ModifiedOn,
		Attributes:     &ValidAttributes,
	}})

	deleteAccount(client, response.Account)
}

func TestClientCreateWithValidationErrors(t *testing.T) {
	testCases := []struct {
		Request       *CreateAccountRequest
		ExpectedError error
	}{
		{
			Request: &CreateAccountRequest{Account: &Account{
				Id:             "",
				OrganisationId: generateUuid(),
				Attributes:     &Attributes{Country: "PL"}}},
			ExpectedError: &ValidationError{Message: "id cannot be empty"},
		},
		{
			Request: &CreateAccountRequest{Account: &Account{
				Id:             "aa",
				OrganisationId: generateUuid(),
				Attributes:     &Attributes{Country: "PL"}}},
			ExpectedError: &ValidationError{Message: "id has to be UUID V1"},
		},
		{
			Request: &CreateAccountRequest{Account: &Account{
				Id:             generateUuid(),
				OrganisationId: "",
				Attributes:     &Attributes{Country: "PL"}}},
			ExpectedError: &ValidationError{Message: "organisationId cannot be empty"},
		},
		{
			Request: &CreateAccountRequest{Account: &Account{
				Id:             generateUuid(),
				OrganisationId: "a",
				Attributes:     &Attributes{Country: "PL"}}},
			ExpectedError: &ValidationError{Message: "organisationId has to be UUID V1"},
		},
		{
			Request: &CreateAccountRequest{Account: &Account{
				Id:             generateUuid(),
				OrganisationId: generateUuid(),
				Attributes:     &Attributes{}}},
			ExpectedError: &ValidationError{Message: "country cannot be empty"},
		},
	}
	for _, testCase := range testCases {
		t.Logf("Given valid HTTP client")
		client := initClient()

		t.Logf("When creating account with invalid request")
		response, err := client.Create(context.Background(), testCase.Request)

		t.Logf("Should return %s error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Nil(t, response)
	}
}

func TestClientCreateWithInValidBody(t *testing.T) {
	t.Logf("Given valid HTTP client")
	client := initClient()

	t.Logf("And given in valid create request")
	request := CreateAccountRequest{Account: &Account{
		Type:           "accounts",
		Id:             generateUuid(),
		OrganisationId: generateUuid(),
		Attributes:     &Attributes{Country: "XXX"},
	}}

	t.Logf("When creating account")
	response, err := client.Create(context.Background(), &request)

	t.Logf("Should return 400 HTTP error")
	assert.Error(t, err)
	var expectedError *http.ClientHttpError
	assert.True(t, errors.As(err, &expectedError))
	assert.Equal(t, expectedError.StatusCode, 400)
	assert.Nil(t, response)
}

func TestClientFetch(t *testing.T) {
	t.Logf("Given valid HTTP client")
	client := initClient()

	t.Logf("And given new account")
	createResponse, _ := client.Create(context.Background(), validCreateAccountRequest())

	t.Logf("And given valid fetch request")
	fetchRequest := FetchAccountRequest{Id: createResponse.Account.Id}

	t.Logf("When fetching account")
	response, err := client.Fetch(context.Background(), &fetchRequest)

	t.Logf("Should not return any errors")
	assert.NoError(t, err)
	assert.Equal(t, response, &FetchAccountResponse{Account: &Account{
		Type:           "accounts",
		Id:             createResponse.Account.Id,
		OrganisationId: createResponse.Account.OrganisationId,
		CreatedOn:      response.Account.CreatedOn,
		ModifiedOn:     response.Account.ModifiedOn,
		Attributes:     &ValidAttributes,
	}})

	deleteAccount(client, response.Account)
}

func TestClientFetchWithValidationError(t *testing.T) {
	testCases := []struct {
		Request       *FetchAccountRequest
		ExpectedError error
	}{
		{
			Request: &FetchAccountRequest{
				Id: "",
			},
			ExpectedError: &ValidationError{Message: "id cannot be empty"},
		},
		{
			Request: &FetchAccountRequest{
				Id: "aa",
			},
			ExpectedError: &ValidationError{Message: "id has to be UUID V1"},
		},
	}

	for _, testCase := range testCases {
		t.Logf("Given valid HTTP client")
		client := initClient()

		t.Logf("When fetching account")
		response, err := client.Fetch(context.Background(), testCase.Request)

		t.Logf("Should return %s error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Nil(t, response)
	}
}

func TestClientFetchWithInValidId(t *testing.T) {
	t.Logf("Given valid HTTP client")
	client := initClient()

	t.Logf("And given valid fetch request")
	fetchRequest := FetchAccountRequest{Id: generateUuid()}

	t.Logf("When fetching account")
	response, err := client.Fetch(context.Background(), &fetchRequest)

	t.Logf("Should return 404 HTTP error")
	assert.Error(t, err)
	var expectedError *http.ClientHttpError
	assert.True(t, errors.As(err, &expectedError))
	assert.Equal(t, expectedError.StatusCode, 404)
	assert.Nil(t, response)
}

func TestClientDelete(t *testing.T) {
	t.Logf("Given valid HTTP client")
	client := initClient()

	t.Logf("And given new account")
	createResponse, _ := client.Create(context.Background(), validCreateAccountRequest())

	t.Logf("And given valid delete request")
	deleteRequest := DeleteAccountRequest{Id: createResponse.Account.Id, Version: createResponse.Account.Version}

	t.Logf("When deleting account")
	err := client.Delete(context.Background(), &deleteRequest)

	t.Logf("Should not return any errors")
	assert.NoError(t, err)
}

func TestClientDeleteWithValidationErrors(t *testing.T) {
	testCases := []struct {
		Request       *DeleteAccountRequest
		ExpectedError error
	}{
		{
			Request: &DeleteAccountRequest{
				Id:      "",
				Version: 0,
			},
			ExpectedError: &ValidationError{Message: "id cannot be empty"},
		},
		{
			Request: &DeleteAccountRequest{
				Id:      "aa",
				Version: 0,
			},
			ExpectedError: &ValidationError{Message: "id has to be UUID V1"},
		},
		{
			Request: &DeleteAccountRequest{
				Id:      generateUuid(),
				Version: -1,
			},
			ExpectedError: &ValidationError{Message: "version has to be larger than zero"},
		},
	}
	for _, testCase := range testCases {
		t.Logf("Given valid HTTP client")
		client := initClient()

		t.Logf("When deleting account")
		err := client.Delete(context.Background(), testCase.Request)

		t.Logf("Should return %s error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
	}

}

func TestClientList(t *testing.T) {
	t.Logf("Given valid HTTP client")
	client := initClient()

	t.Logf("And given new accounts")
	createResponse1, _ := client.Create(context.Background(), validCreateAccountRequest())
	createResponse2, _ := client.Create(context.Background(), validCreateAccountRequest())

	t.Logf("And given valid list request")
	listRequest := ListAccountsRequest{PageNumber: 0, PageSize: 100}

	t.Logf("When listing account")
	response, err := client.List(context.Background(), &listRequest)

	t.Logf("Should not return any errors")
	assert.NoError(t, err)
	assert.Equal(t, response, &ListAccountResponse{Accounts: []*Account{
		{
			Type:           "accounts",
			Id:             createResponse1.Account.Id,
			OrganisationId: createResponse1.Account.OrganisationId,
			CreatedOn:      createResponse1.Account.CreatedOn,
			ModifiedOn:     createResponse1.Account.ModifiedOn,
			Attributes:     &ValidAttributes,
		},
		{
			Type:           "accounts",
			Id:             createResponse2.Account.Id,
			OrganisationId: createResponse2.Account.OrganisationId,
			CreatedOn:      createResponse2.Account.CreatedOn,
			ModifiedOn:     createResponse2.Account.ModifiedOn,
			Attributes:     &ValidAttributes,
		},
	}})

	deleteAccount(client, response.Accounts[0])
	deleteAccount(client, response.Accounts[1])
}

func TestClientListWithValidationErrors(t *testing.T) {
	testCases := []struct {
		Request       *ListAccountsRequest
		ExpectedError error
	}{
		{
			Request: &ListAccountsRequest{
				PageNumber: -1,
				PageSize:   1,
			},
			ExpectedError: &ValidationError{Message: "pageNumber has to be larger than zero"},
		},
		{
			Request: &ListAccountsRequest{
				PageNumber: 1,
				PageSize:   -1,
			},
			ExpectedError: &ValidationError{Message: "pageSize has to be larger than zero"},
		},
	}
	for _, testCase := range testCases {
		t.Logf("Given valid HTTP client")
		client := initClient()

		t.Logf("When listing account")
		response, err := client.List(context.Background(), testCase.Request)

		t.Logf("Should %s error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Nil(t, response)
	}
}

func TestClientListWithPagination(t *testing.T) {
	t.Logf("Given valid HTTP client")
	client := initClient()

	t.Logf("And given new 120 accounts")
	var createResponses []*CreateAccountResponse
	for range makeRange(120) {
		createResponse, _ := client.Create(context.Background(), validCreateAccountRequest())
		createResponses = append(createResponses, createResponse)
	}

	t.Logf("When listing page 0")
	listRequest1 := ListAccountsRequest{PageNumber: 0, PageSize: 100}
	response1, err1 := client.List(context.Background(), &listRequest1)

	t.Logf("Should return 100 records")
	assert.NoError(t, err1)
	assert.Equal(t, 100, len(response1.Accounts))

	t.Logf("When listing page 1")
	listRequest2 := ListAccountsRequest{PageNumber: 1, PageSize: 100}
	response2, err2 := client.List(context.Background(), &listRequest2)

	t.Logf("Should return 20 records")
	assert.NoError(t, err2)
	assert.Equal(t, 20, len(response2.Accounts))

	for _, createResponse := range createResponses {
		account := *createResponse.Account
		deleteAccount(client, &account)
	}
}

func initClient() *Client {
	apiPath, err := url.Parse(ApiPath)
	if err != nil {
		log.Fatalf("Invalid URL %s %s", ApiPath, err)
	}
	config := ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url:     apiPath,
		RetriesConfig: &retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.5,
		},
	}

	client, _ := NewClient(config)
	return client

}

func deleteAccount(client *Client, account *Account) {
	err := client.Delete(context.Background(), &DeleteAccountRequest{Id: account.Id, Version: account.Version})
	if err != nil {
		log.Fatal(err)
	}

}

func validCreateAccountRequest() *CreateAccountRequest {
	return &CreateAccountRequest{Account: &Account{
		Type:           "accounts",
		Id:             generateUuid(),
		OrganisationId: generateUuid(),
		Attributes:     &ValidAttributes,
	}}
}

func makeRange(max int) []int {
	a := make([]int, max)
	for i := range a {
		a[i] = i
	}
	return a
}

func generateUuid() string {
	out, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}
