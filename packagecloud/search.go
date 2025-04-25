package packagecloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/amdprophet/packagecloud-go/types"
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

func (c *Client) Search(options SearchOptions, per_page string) (types.PackageFragments, error) {
	var packages types.PackageFragments
	var mu = &sync.RWMutex{}

	if err := c.SearchStream(options, per_page, func(streamPackages types.PackageFragments) {
		mu.Lock()
		packages = append(packages, streamPackages...)
		mu.Unlock()
	}); err != nil {
		return nil, err
	}

	return packages, nil
}

func (c *Client) SearchStream(options SearchOptions, per_page string, fn func(types.PackageFragments)) error {
	searchPath := fmt.Sprintf(searchPath, options.RepoUser, options.RepoName)
	searchURL, err := url.Parse(searchPath)
	if err != nil {
		return fmt.Errorf("this is a bug, failed to parse relative url: %s", err)
	}

	query := searchURL.Query()
	query.Add("per_page", per_page)

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

	searchURL.RawQuery = query.Encode()
	endpoint := c.getURL(searchURL)

	return c.paginatedRequest("GET", endpoint.String(), nil, "application/json", func(bytes []byte) error {
		var packages types.PackageFragments
		if err := json.Unmarshal(bytes, &packages); err != nil {
			return &UnmarshalError{
				Data: bytes,
				Err:  err,
			}
		}
		fn(packages)
		return nil
	})
}
