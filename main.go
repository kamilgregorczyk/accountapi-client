package main

import (
	"accountapi-client/http/retry"
	"accountapi-client/account"
	"context"
	"log"
	"net/url"
	"time"
)

func main() {
	accountClient, err := account.NewClient(account.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: url.URL{
			Scheme: "https",
			Host:   "account.raspicluster.pl"},
		RetriesConfig: retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	items, err := accountClient.GetItems(context.Background())
	if err != nil {
		log.Panicf(err.Error())
	} else {
		log.Printf("Items: %+v", items)
		item, err := accountClient.GetItem(context.Background(), items[0].Id)
		if err != nil {
			log.Panicf(err.Error())
		} else {
			log.Printf("Item: %+v", item)
		}

		item2, err := accountClient.CreateItem(context.Background(), account.CreateInventory{Name: "aa", Description: "cc"})
		if err != nil {
			log.Panicf(err.Error())
		} else {
			log.Printf("Item: %+v", item2)
		}
	}
}
