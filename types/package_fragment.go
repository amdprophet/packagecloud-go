package types

type PackageFragments []PackageFragment

func (p PackageFragments) Indexed() bool {
	for _, pkg := range p {
		if !pkg.Indexed {
			return false
		}
	}
	return true
}

type PackageFragment struct {
	// Name is the name of the package.
	Name string `json:"name"`

	// CreatedAt is a timestamp of when the package was uploaded.
	CreatedAt string `json:"created_at"`

	// DistroVersion is the distro version for the package.
	DistroVersion string `json:"distro_version"`

	// Version is the version of the package.
	Version string `json:"version"`

	// Release is the release of the package (if available).
	Release string `json:"release"`

	// Architecture is the architecture of the package.
	Architecture string `json:"architecture"`

	// Epoch is the epoch of the package (if available).
	Epoch int `json:"epoch"`

	// Scope is the scope of the package (if available).
	Scope string `json:"scope"`

	// Private specifies whether or not the package is in a private repository.
	Private bool `json:"private"`

	// Type is the type of package ("deb", "gem", or "rpm").
	Type string `json:"type"`

	// Filename is the filename of the package.
	Filename string `json:"filename"`

	// UploaderName is the name of the uploader for the package.
	UploaderName string `json:"uploader_name"`

	// Indexed specifies whether or not this package has been indexed.
	Indexed bool `json:"indexed"`

	// RepositoryHTMLURL is the HTML URL of the repository.
	RepositoryHTMLURL string `json:"repository_html_url"`

	// PackageURL is the API URL for this package.
	PackageURL string `json:"package_url"`

	// PackageHTMLURL is the HTML URL for this package.
	PackageHTMLURL string `json:"package_html_url"`

	// DownloadsDetailURL is the URL to get access log details for package
	// downloads.
	DownloadsDetailURL string `json:"downloads_detail_url"`

	// DownloadSeriesURL is the URL to get time series data for package
	// downloads.
	DownloadsSeriesURL string `json:"downloads_series_url"`

	// TotalDownloadsCount is the number of times this package has been
	// downloaded.
	TotalDownloadsCount int `json:"total_downloads_count"`

	// PromoteURL is the API URL that can be used to promote this package.
	PromoteURL string `json:"promote_url"`

	// DestroyURL is the URL for the HTTP DELETE request to destroy this
	// package.
	DestroyURL string `json:"destroy_url"`
}
