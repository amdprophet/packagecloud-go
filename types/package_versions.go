package types

import (
	"errors"
	"fmt"
	"sort"

	semver "github.com/Masterminds/semver/v3"
)

type PackageVersions map[string]int

func (p PackageVersions) SemanticVersions() ([]*semver.Version, error) {
	versions := []*semver.Version{}

	for version := range p {
		v, err := semver.NewVersion(version)
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %s", err)
		}
		versions = append(versions, v)
	}

	return versions, nil
}

func (p PackageVersions) ReverseSorted() ([]*semver.Version, error) {
	versions, err := p.SemanticVersions()
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, errors.New("no versions available")
	}

	// Sort versions in descending order
	sort.Sort(sort.Reverse(semver.Collection(versions)))

	return versions, nil
}

func (p PackageVersions) LatestVersion() (string, error) {
	versions, err := p.ReverseSorted()
	if err != nil {
		return "", err
	}

	// Return the latest version as a string
	return versions[0].String(), nil
}

func (p PackageVersions) PreviousVersion() (string, error) {
	versions, err := p.ReverseSorted()
	if err != nil {
		return "", err
	}

	// Return the second latest version as a string
	return versions[1].String(), nil
}
