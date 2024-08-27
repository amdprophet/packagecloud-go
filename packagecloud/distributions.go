package packagecloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/amdprophet/packagecloud-go/util"
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

	if resp.StatusCode == 402 {
		var jsonErrs map[string][]string
		json.Unmarshal(body, &jsonErrs)
		if len(jsonErrs) == 1 {
			if errMsgs, ok := jsonErrs["error"]; ok {
				if len(errMsgs) == 1 && util.SliceContainsString(errMsgs, "payment required") {
					return nil, ErrPaymentRequired
				}
			}
		}
		return nil, fmt.Errorf("api responded with error: %s", string(body))
	}

	return body, nil
}
