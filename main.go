package main

import (
	"accountapi-client/http/retry"
	"accountapi-client/inventory"
	"context"
	"log"
	"net/url"
	"time"
)

func main() {
	inventoryClient, err := inventory.NewClient(inventory.ClientConfig{
		Timeout: time.Second,
		Logging: true,
		Url: url.URL{
			Scheme: "https",
			Host:   "inventory.raspicluster.pl"},
		RetriesConfig: retry.RetriesConfig{
			MaxRetries: 3,
			Delay:      time.Second,
			Factor:     1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	items, err := inventoryClient.GetItems(context.Background())
	if err != nil {
		log.Panicf(err.Error())
	} else {
		log.Printf("Items: %+v", items)
		item, err := inventoryClient.GetItem(context.Background(), items[0].Id)
		if err != nil {
			log.Panicf(err.Error())
		} else {
			log.Printf("Item: %+v", item)
		}

		item2, err := inventoryClient.CreateItem(context.Background(), inventory.CreateInventory{Name: "aa", Description: "cc"})
		if err != nil {
			log.Panicf(err.Error())
		} else {
			log.Printf("Item: %+v", item2)
		}
	}
}
