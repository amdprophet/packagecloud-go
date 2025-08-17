package packagecloud

import (
	"fmt"
	"sync"

	"github.com/amdprophet/packagecloud-go/types"
)

type ListVersionsOptions struct {
	// Repo is the repository to list versions for.
	Repo Repo

	// PackageName is the name of the package to list versions for.
	PackageName string

	// Filter can be used to search by package type.
	// (RPMs, Debs, DSCs, Gem, Python, Node).
	// Ignored when Dist is present.
	Filter string

	// Dist is the name of the distribution that the package is in.
	// (i.e. ubuntu, el/6)
	Dist string

	// Arch is the architecture of the packages. (i.e. x86_64, arm64, amd64).
	// Alpine/RPM/Debian only.
	Arch string

	// PerPage is the number of packages to return from the results set. If
	// nothing passed the default is 30.
	PerPage string
}

func (o ListVersionsOptions) SearchOptions() SearchOptions {
	return SearchOptions{
		RepoUser: o.Repo.User,
		RepoName: o.Repo.Name,
		Filter:   o.Filter,
		Dist:     o.Dist,
		Arch:     o.Arch,
		PerPage:  o.PerPage,
	}
}

func (o ListVersionsOptions) Validate() error {
	if err := o.Repo.Validate(); err != nil {
		return fmt.Errorf("repository validation failed: %w", err)
	}
	if o.PackageName == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	return nil
}

func (c *Client) ListVersions(options ListVersionsOptions) (types.PackageVersions, error) {
	versions := types.PackageVersions{}
	mu := &sync.RWMutex{}

	callback := func(streamPackages types.PackageFragments) {
		mu.Lock()
		for _, pkg := range streamPackages {
			if pkg.Name == options.PackageName {
				key := pkg.Version
				if pkg.Type == "rpm" {
					key = fmt.Sprintf("%s-%s", key, pkg.Release)
				}
				versions[key]++
			}
		}
		mu.Unlock()
	}

	if options.Filter != "" || options.Dist != "" || options.Arch != "" {
		searchOptions := options.SearchOptions()

		if err := c.SearchStream(searchOptions, callback); err != nil {
			return nil, err
		}
	} else {
		if err := c.ListPackagesStream(options.Repo, callback); err != nil {
			return nil, err
		}
	}

	return versions, nil
}

func (c *Client) LatestVersion(options ListVersionsOptions) (string, error) {
	versions, err := c.ListVersions(options)
	if err != nil {
		return "", err
	}

	return versions.LatestVersion()
}

func (c *Client) PreviousVersion(options ListVersionsOptions) (string, error) {
	versions, err := c.ListVersions(options)
	if err != nil {
		return "", err
	}

	return versions.PreviousVersion()
}
