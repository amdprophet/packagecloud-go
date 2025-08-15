package packagecloud

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/amdprophet/packagecloud-go/types"
)

const (
	// repo user/name, distro name/version, filename
	promotePath = "/api/v1/repos/%s/%s/%s/promote.json"
)

// PromoteByFilename will promote a single package by filename.
func (c *Client) PromoteByFilename(src Repo, dst Repo, distro Distro, filename string) error {
	if err := src.Validate(); err != nil {
		return fmt.Errorf("source repository validation failed: %w", err)
	}
	if err := dst.Validate(); err != nil {
		return fmt.Errorf("destination repository validation failed: %w", err)
	}
	if err := distro.Validate(); err != nil {
		return fmt.Errorf("distro validation failed: %w", err)
	}
	if isEmptyString(filename) {
		return errors.New("filename cannot be empty")
	}

	pkg := types.PackageFragment{
		DistroVersion: distro.String(),
		Filename:      filename,
	}

	pkg.PromoteURL = buildPromoteURL(src, pkg)

	fmt.Println("Promoting package")
	fmt.Printf("  - Source repository:      %s\n", src)
	fmt.Printf("  - Destination repository: %s\n", dst)
	fmt.Printf("  - Filename:               %s\n", pkg.Filename)
	fmt.Printf("  - Distro:                 %s\n", pkg.DistroVersion)
	fmt.Println("")

	if err := c.promote(pkg, dst); err != nil {
		return err
	}

	fmt.Println("Successfully promoted 1 package")

	return nil
}

// PromoteBySearch will search for any packages matching the given search
// options and then promote all matches to the destination repository.
func (c *Client) PromoteBySearch(dst Repo, options SearchOptions) error {
	if err := dst.Validate(); err != nil {
		return fmt.Errorf("destination repository validation failed: %w", err)
	}
	if err := options.Validate(); err != nil {
		return fmt.Errorf("search options validation failed: %w", err)
	}

	src := NewRepo(options.RepoUser, options.RepoName)

	packages, err := c.Search(options)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		fmt.Println("Promoting package")
		fmt.Printf("  - Source repository:      %s\n", src)
		fmt.Printf("  - Destination repository: %s\n", dst)
		fmt.Printf("  - Name:                   %s\n", pkg.Name)
		fmt.Printf("  - Type:                   %s\n", pkg.Type)
		fmt.Printf("  - Version:                %s\n", pkg.Version)
		fmt.Printf("  - Release:                %s\n", pkg.Release)
		fmt.Printf("  - Epoch:                  %d\n", pkg.Epoch)
		fmt.Printf("  - Architecture:           %s\n", pkg.Architecture)
		fmt.Printf("  - Distro:                 %s\n", pkg.DistroVersion)
		fmt.Println("")

		if err := c.promote(pkg, dst); err != nil {
			return err
		}
	}

	fmt.Printf("Successfully promoted %d package(s)\n", len(packages))

	return nil
}

func (c *Client) promote(pkg types.PackageFragment, dst Repo) error {
	// Validate method arguments before proceeding with the promotion. If any
	// of these validations fail, it indicates a bug in the code or an
	// incorrect usage of the method.
	if err := dst.Validate(); err != nil {
		bugPanic("dst repo validation failed: " + err.Error())
	}

	promoteURL, err := url.Parse(pkg.PromoteURL)
	if err != nil {
		msg := fmt.Sprintf("failed to parse relative url: %s", err)
		bugPanic(msg)
	}

	query := promoteURL.Query()
	query.Add("destination", dst.String())
	promoteURL.RawQuery = query.Encode()

	endpoint := c.getURL(promoteURL)

	if _, err := c.apiRequest("POST", endpoint.String(), nil, "application/json"); err != nil {
		return err
	}
	return nil
}

func buildPromoteURL(repo Repo, pkg types.PackageFragment) string {
	// Validate method arguments before generating the promote path. If any of
	// these validations fail, it indicates a bug in the code or an incorrect
	// usage of the method.
	if err := repo.Validate(); err != nil {
		bugPanic("repo validation failed: " + err.Error())
	}
	if pkg.DistroVersion == "" {
		bugPanic("pkg has empty distro version")
	}
	if pkg.Filename == "" {
		bugPanic("pkg has empty filename")
	}

	return fmt.Sprintf(
		promotePath,
		repo.String(),
		pkg.DistroVersion,
		pkg.Filename,
	)
}
