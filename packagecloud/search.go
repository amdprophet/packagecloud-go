package packagecloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/amdprophet/packagecloud-go/util"
)

const (
	// user, repo
	searchPath = "/api/v1/repos/%s/%s/search.json"
)

type SearchOptions struct {
	RepoUser string
	RepoName string
	Query    string
	Filter   string
	Dist     string
}

func (c *Client) Search(options SearchOptions) ([]byte, error) {
	searchURL := fmt.Sprintf(searchPath, options.RepoUser, options.RepoName)

	req, err := c.newRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	query := req.URL.Query()
	query.Add("per_page", "250")

	if options.Query != "" {
		query.Add("q", options.Query)
	}

	if options.Filter != "" {
		query.Add("filter", options.Filter)
	}

	if options.Dist != "" {
		if query.Has("filter") {
			query.Del("filter")
		}
		query.Add("dist", options.Dist)
	}

	req.URL.RawQuery = query.Encode()

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
