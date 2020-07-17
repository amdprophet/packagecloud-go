package packagecloud

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	listDistributionsPath = "/api/v1/distributions.json"
)

type GetClientFn func() (*Client, error)

type Config struct {
	ServiceURL string `json:"url"`
	Token      string `json:"token"`
	Verbose    bool   `json:"verbose"`
}

func (c Config) Validate() error {
	if c.Token == "" {
		return errors.New("token must not be empty")
	}

	if _, err := url.Parse(c.ServiceURL); err != nil {
		return fmt.Errorf("invalid url: %s", err)
	}

	return nil
}

type Client struct {
	config *Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: &config,
	}
}

func (c *Client) newRequest(method string, path string) (*http.Request, error) {
	baseURL, _ := url.Parse(c.config.ServiceURL)
	relativeURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relative url: %s")
	}
	requestURL := baseURL.ResolveReference(relativeURL)

	fmt.Println("url:", requestURL)
	//req := http.NewRequest(method )

	return nil, nil
}

func (c *Client) ListDistributions() {
	_, err := c.newRequest("GET", listDistributionsPath)
	if err != nil {
		// no-op for now
	}
}
