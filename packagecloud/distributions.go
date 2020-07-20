package packagecloud

import (
	"fmt"
	"io/ioutil"
)

const (
	distributionsPath = "/api/v1/distributions.json"
)

func (c *Client) GetDistributions() ([]byte, error) {
	req, err := c.newRequest("GET", distributionsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	return body, nil
}
