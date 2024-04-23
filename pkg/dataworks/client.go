package dataworks

import "github.com/go-resty/resty/v2"

type Client struct {
	httpClient *resty.Client
	Endpint    string
}

func NewClient(endpoint string) *Client {
	httpClient := resty.New()

	return &Client{
		Endpint:    endpoint,
		httpClient: httpClient,
	}
}
