package dataworks

import "github.com/go-resty/resty/v2"

type Client struct {
	Endpint    string
	httpClient *resty.Client
}

func NewClient(endpoint string) *Client {
	httpClient := resty.New()

	return &Client{
		Endpint:    endpoint,
		httpClient: httpClient,
	}
}
