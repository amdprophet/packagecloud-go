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
	// RepoUser is the username that the repository to search belongs to.
	RepoUser string

	// RepoName is the name of the repository to search.
	RepoName string

	// Query is a query string to search for package filename. If empty string
	// is passed, all packages are returned.
	Query string

	// Filter can be used to search by package type.
	// (RPMs, Debs, DSCs, Gem, Python, Node).
	// Ignored when Dist is present.
	Filter string

	// Dist is the name of the distribution that the package is in.
	// (i.e. ubuntu, el/6)
	// Overrides Filter.
	Dist string

	// Arch is the architecture of the packages. (i.e. x86_64, arm64, amd64).
	// Alpine/RPM/Debian only.
	Arch string

	// PerPage is the number of packages to return from the results set. If
	// nothing passed the default is 30.
	PerPage string
}

func (o SearchOptions) Validate() error {
	if o.Query == "" && o.Filter == "" && o.Dist == "" && o.Arch == "" {
		return &MissingSearchOptionsError{}
	}
	return nil
}

func (c *Client) Search(options SearchOptions) (types.PackageFragments, error) {
	var packages types.PackageFragments
	var mu = &sync.RWMutex{}

	if err := c.SearchStream(options, func(streamPackages types.PackageFragments) {
		mu.Lock()
		packages = append(packages, streamPackages...)
		mu.Unlock()
	}); err != nil {
		return nil, err
	}

	return packages, nil
}

func (c *Client) SearchStream(options SearchOptions, fn func(types.PackageFragments)) error {
	if err := options.Validate(); err != nil {
		return err
	}

	searchPath := fmt.Sprintf(searchPath, options.RepoUser, options.RepoName)
	searchURL, err := url.Parse(searchPath)
	if err != nil {
		return fmt.Errorf("this is a bug, failed to parse relative url: %s", err)
	}

	query := searchURL.Query()

	if options.Query != "" {
		query.Add("q", options.Query)
	}

	if options.Filter != "" {
		query.Add("filter", options.Filter)
	}

	if options.Dist != "" {
		query.Add("dist", options.Dist)
	}

	if options.Arch != "" {
		query.Add("arch", options.Arch)
	}

	if options.PerPage != "" {
		query.Add("per_page", options.PerPage)
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
