package packagecloud

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GetClientFn func() (*Client, error)

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config Config) *Client {
	return &Client{
		config:     &config,
		httpClient: &http.Client{},
	}
}

func (c *Client) getURL(path string) (string, error) {
	baseURL, _ := url.Parse(c.config.ServiceURL)
	relativeURL, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("this is a bug, failed to parse relative url: %s", err)
	}
	requestURL := baseURL.ResolveReference(relativeURL)
	return requestURL.String(), nil
}

func (c *Client) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url, err := c.getURL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.config.Token, "")

	return req, nil
}
