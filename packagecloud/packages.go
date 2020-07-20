package packagecloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

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

func (c *Client) PushPackage(options PushPackageOptions) ([]byte, error) {
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

	packagesURL := fmt.Sprintf(packagesPath, options.RepoUser, options.RepoName)
	req, err := c.newRequest("POST", packagesURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	if resp.StatusCode == 422 {
		var jsonErrs map[string][]string
		json.Unmarshal(body, &jsonErrs)
		if len(jsonErrs) == 1 {
			if errMsgs, ok := jsonErrs["filename"]; ok {
				if len(errMsgs) == 1 && util.SliceContainsString(errMsgs, "has already been taken") {
					return nil, ErrPackageAlreadyExists
				}
			}
		}
		return nil, fmt.Errorf("api responded with error: %s", string(body))
	}

	return nil, nil
}
