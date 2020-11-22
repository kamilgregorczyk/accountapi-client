package main

import (
	"accountapi-client/account"
	"accountapi-client/http/retry"
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

	accountId, err := uuid.NewUUID()
	account, err := accountClient.Create(context.Background(), &account.Account{
		Type:           "accounts",
		Id:             accountId.String(),
		OrganisationId: "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		Attributes: &account.AccountAttributes{
			Country:      "PL",
			BaseCurrency: "PLN",
			BankId:       "400300",
			BankIdCode:   "GBDSC",
			Bic:          "NWBKGB22",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", account)

	account, err = accountClient.Fetch(context.Background(), accountId.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", account)

	err = accountClient.Delete(context.Background(), accountId.String(), account.Version)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Account should be deleted, lets check...")

	_, err = accountClient.Fetch(context.Background(), accountId.String())
	if err != nil {
		log.Println("Account does not exist!")
	}
}
