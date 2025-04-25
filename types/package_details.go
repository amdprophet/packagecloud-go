package types

type PackageDetails struct {
	// Name is the name of the package.
	Name string `json:"name"`

	// DistroVersion is the distro version for the package.
	DistroVersion string `json:"distro_version"`

	// Architecture is the architecture of the package.
	Architecture string `json:"architecture"`

	// Repository is the name of the repository for the package.
	Repository string `json:"repository"`

	// Size is the size (in bytes) of the package. Returned as a String for
	// JavaScript support.
	Size string `json:"size"`

	// Summary is the summary of the package (if available).
	Summary string `json:"summary"`

	// Filename is the filename of the package.
	Filename string `json:"filename"`

	// Description is the description of the package (if available).
	Description string `json:"description"`

	// MD5Sum is the MD5 checksum for the package (if available).
	MD5Sum string `json:"md5sum"`

	// SHA1Sum is the SHA1 checksum for the package (if available).
	SHA1Sum string `json:"sha1sum"`

	// SHA256Sum is the SHA256 checksum for the package (if available).
	SHA256Sum string `json:"sha256sum"`

	// SHA512Sum is the SHA512 checksum for the package (if available).
	SHA512Sum string `json:"sha512sum"`

	// Private specifies whether or not the package is in a private repository.
	Private bool `json:"private"`

	// UploaderName is the name of the uploader for the package.
	UploaderName string `json:"uploader_name"`

	// CreatedAt is a timestamp of when the package was uploaded.
	CreatedAt string `json:"created_at"`

	// Licenses is a list of licenses for this package (if available).
	Licenses []string `json:"licenses"`

	// Version is the version of the package.
	Version string `json:"version"`

	// Release is the release of the package (if available).
	Release string `json:"release"`

	// Epoch is the epoch of the package (if available).
	Epoch int `json:"epoch"`

	// Indexed specifies whether or not this package has been indexed.
	Indexed bool `json:"indexed"`

	// Scope is the scope of the package (if available).
	Scope string `json:"scope"`

	// RepositoryHTMLURL is the HTML URL of the repository.
	RepositoryHTMLURL string `json:"repository_html_url"`

	// VersionsURL is the API URL of other versions of this package.
	VersionsURL string `json:"versions_url"`

	// PromoteURL is the API URL that can be used to promote this package.
	PromoteURL string `json:"promote_url"`

	// TotalDownloadsCount is the number of times this package has been
	// downloaded.
	TotalDownloadsCount int `json:"total_downloads_count"`

	// DownloadURL is the URL to download this package.
	DownloadURL string `json:"download_url"`

	// PackageHTMLURL is the HTML URL for this package.
	PackageHTMLURL string `json:"package_html_url"`

	// DownloadsDetailURL is the URL to get access log details for package
	// downloads.
	DownloadsDetailURL string `json:"downloads_detail_url"`

	// DownloadSeriesURL is the URL to get time series data for package
	// downloads.
	DownloadsSeriesURL string `json:"downloads_series_url"`

	// DownloadsCountURL is the URL to get the total number of package
	// downloads.
	DownloadsCountURL string `json:"downloads_count_url"`

	// DestroyURL is the URL for the HTTP DELETE request to destroy this
	// package.
	DestroyURL string `json:"destroy_url"`

	// SelfURL is the API URL for this response.
	SelfURL string `json:"self_url"`
}
