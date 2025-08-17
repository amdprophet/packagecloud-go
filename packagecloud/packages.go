package packagecloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/amdprophet/packagecloud-go/types"
	"github.com/amdprophet/packagecloud-go/util"
)

const (
	// user, repo
	packagesPath = "/api/v1/repos/%s/%s/packages.json"
)

// TODO: packagecloud supports more than this, implement the others
func GetSupportedFileExtensions() []string {
	return []string{
		".deb",
		".rpm",
	}
}

func ValidateFileExtensions(paths []string) error {
	supportedExts := GetSupportedFileExtensions()
	exts := make(map[string]struct{})

	for _, path := range paths {
		ext := filepath.Ext(path)
		if !util.SliceContainsString(supportedExts, ext) {
			return fmt.Errorf("invalid file extension: %s, supported extensions: %s",
				ext, supportedExts)
		}
		exts[ext] = struct{}{}
	}

	if len(exts) > 1 {
		return errors.New("cannot push multiple packages of different types at the same time")
	}

	return nil
}

type PushPackageOptions struct {
	RepoUser string
	RepoName string
	DistroID string
	FilePath string
}

func (c *Client) PushPackage(options PushPackageOptions) (*types.PackageDetails, error) {
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	if err := writer.WriteField("package[distro_version_id]", options.DistroID); err != nil {
		return nil, fmt.Errorf("failed to write form field: %s", err)
	}

	file, err := os.Open(options.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("package[package_file]", filepath.Base(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %s", err)
	}
	io.Copy(part, file)
	writer.Close()

	path := fmt.Sprintf(packagesPath, options.RepoUser, options.RepoName)
	packagesURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("this is a bug, failed to parse relative url: %s", err)
	}

	endpoint := c.getURL(packagesURL)
	contentType := writer.FormDataContentType()

	resp, err := c.apiRequest("POST", endpoint.String(), reqBody, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %s", err)
	}

	var pkg types.PackageDetails
	if err := json.Unmarshal(resp.Body, &pkg); err != nil {
		return nil, &UnmarshalError{
			Data: resp.Body,
			Err:  err,
		}
	}

	return &pkg, nil
}

func (c *Client) ListPackages(repo Repo) (types.PackageFragments, error) {
	var packages types.PackageFragments
	mu := &sync.RWMutex{}

	if err := c.ListPackagesStream(repo, func(streamPackages types.PackageFragments) {
		mu.Lock()
		packages = append(packages, streamPackages...)
		mu.Unlock()
	}); err != nil {
		return nil, err
	}

	return packages, nil
}

func (c *Client) ListPackagesStream(repo Repo, fn func(types.PackageFragments)) error {
	if err := repo.Validate(); err != nil {
		return fmt.Errorf("repository validation failed: %w", err)
	}

	packagesURL, err := url.Parse(fmt.Sprintf(packagesPath, repo.User, repo.Name))
	if err != nil {
		return fmt.Errorf("this is a bug, failed to parse relative url: %s", err)
	}

	endpoint := c.getURL(packagesURL)

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
