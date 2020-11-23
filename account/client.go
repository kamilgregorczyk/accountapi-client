package account

import (
	"accountapi-client/http"
	"accountapi-client/http/retry"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ClientConfig struct {
	Timeout       time.Duration
	Logging       bool
	Url           url.URL
	RetriesConfig retry.RetriesConfig
}

type Client struct {
	Url    url.URL
	Client *http.Client
}

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

func (c *Client) Create(ctx context.Context, account *Account) (*Account, error) {
	path, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/organisation/accounts", c.Url.String()))
	if err != nil {
		return nil, err
	}
	var accountData Data
	err = c.Client.Post(ctx, path, &Data{account}, &accountData)
	return accountData.Data, err
}

func (c *Client) Fetch(ctx context.Context, id string) (*Account, error) {
	path, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/organisation/accounts/%s", c.Url.String(), id))
	if err != nil {
		return nil, err
	}

	var accountData Data
	err = c.Client.Get(ctx, path, &accountData)
	return accountData.Data, err
}

func (c *Client) Delete(ctx context.Context, id string, version int64) error {
	path, err := url.ParseRequestURI(fmt.Sprintf("%s/v1/organisation/accounts/%s", c.Url.String(), id))
	if err != nil {
		return err
	}

	query := path.Query()
	query.Set("version", strconv.FormatInt(version, 10))
	path.RawQuery = query.Encode()
	err = c.Client.Delete(ctx, path)
	return err
}
