package main

import (
	"accountapi-client/account"
	"accountapi-client/retry"
	"context"
	"github.com/google/uuid"
	"log"
	"net/url"
	"time"
)

func main() {
	accountClient, err := account.NewClient(account.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: url.URL{
			Scheme: "http",
			Host:   "localhost:8080"},
		RetriesConfig: retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1.5,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	accountId, _ := uuid.NewUUID()
	organisationId, _ := uuid.NewUUID()
	createResponse, err := accountClient.Create(context.Background(), &account.CreateAccountRequest{Account: &account.Account{
		Type:           "accounts",
		Id:             accountId.String(),
		OrganisationId: organisationId.String(),
		Attributes: &account.Attributes{
			Country:      "PL",
			BaseCurrency: "PLN",
			BankId:       "400300",
			BankIdCode:   "GBDSC",
			Bic:          "NWBKGB22",
		},
	}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", createResponse.Account)

	fetchResponse, err := accountClient.Fetch(context.Background(), &account.FetchAccountRequest{Id: createResponse.Account.Id})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", fetchResponse.Account)

	err = accountClient.Delete(context.Background(), &account.DeleteAccountRequest{Id: fetchResponse.Account.Id, Version: fetchResponse.Account.Version})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Account should be deleted, lets check...")

	_, err = accountClient.Fetch(context.Background(), &account.FetchAccountRequest{Id: fetchResponse.Account.Id})
	if err != nil {
		log.Println("Account does not exist!")
	}

	listResponse, err := accountClient.List(context.Background(), &account.ListAccountsRequest{PageSize: 1000000, PageNumber: 0})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d", len(listResponse.Accounts))

}
