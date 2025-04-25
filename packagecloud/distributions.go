package packagecloud

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/amdprophet/packagecloud-go/types"
)

const (
	distributionsPath = "/api/v1/distributions.json"
)

func (c *Client) GetDistributions() (types.PackageTypes, error) {
	distributionsURL, err := url.Parse(distributionsPath)
	if err != nil {
		return nil, fmt.Errorf("this is a bug, failed to parse relative url: %s", err)
	}

	endpoint := c.getURL(distributionsURL)

	resp, err := c.apiRequest("GET", endpoint.String(), nil, "application/json")
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %s", err)
	}

	var packageTypes types.PackageTypes
	if err := json.Unmarshal(resp.Body, &packageTypes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal distributions json: %s", err)
	}

	return packageTypes, nil
}
