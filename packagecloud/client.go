package packagecloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/amdprophet/packagecloud-go/util"
	"github.com/peterhellberg/link"
)

type GetClientFn func() (*Client, error)

type APIResponse struct {
	Body      []byte
	LinkGroup link.Group
}

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

func (c *Client) getURL(path *url.URL) *url.URL {
	baseURL, _ := url.Parse(c.config.ServiceURL)
	return baseURL.ResolveReference(path)
}

func (c *Client) apiRequest(method string, url string, payload io.Reader, contentType string) (*APIResponse, error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Set("Accept", contentType)
	req.SetBasicAuth(c.config.Token, "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var httpErr error
	switch resp.StatusCode {
	case 402:
		httpErr = getResponseError(body, "payment required", ErrPaymentRequired)
	case 422:
		httpErr = getResponseError(body, "has already been taken", ErrPackageAlreadyExists)
	}
	if httpErr != nil {
		return nil, fmt.Errorf("api responded with error: %s", string(body))
	}

	return &APIResponse{
		Body:      body,
		LinkGroup: link.ParseResponse(resp),
	}, nil
}

func (c *Client) paginatedRequest(method string, endpoint string, payload io.Reader, contentType string, fn func([]byte) error) error {
	wg := sync.WaitGroup{}
	for {
		resp, err := c.apiRequest(method, endpoint, payload, contentType)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(resp.Body)
		}()

		next, found := resp.LinkGroup["next"]
		if !found {
			break
		}
		endpoint = next.URI
	}
	wg.Wait()

	return nil
}

func getResponseError(bytes []byte, search string, respErr error) error {
	var jsonErrs map[string][]string
	if err := json.Unmarshal(bytes, &jsonErrs); err != nil {
		return fmt.Errorf("failed to parse response body as json: %w", err)
	}
	if len(jsonErrs) == 1 {
		if errMsgs, ok := jsonErrs["error"]; ok {
			if len(errMsgs) == 1 && util.SliceContainsString(errMsgs, search) {
				return respErr
			}
		}
	}
	return nil
}
