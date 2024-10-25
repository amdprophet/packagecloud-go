package types

type Packages []Package

type Package struct {
	Name               string `json:"name"`
	DistroVersion      string `json:"distro_version"`
	CreatedAt          string `json:"created_at"`
	Version            string `json:"version"`
	Release            string `json:"release"`
	Epoch              int    `json:"epoch"`
	Scope              string `json:"scope"`
	Private            bool   `json:"private"`
	Type               string `json:"type"`
	Filename           string `json:"filename"`
	UploaderName       string `json:"uploader_name"`
	Indexed            bool   `json:"indexed"`
	Sha256Sum          string `json:"sha256sum"`
	RepositoryHtmlUrl  string `json:"repository_html_url"`
	PackageUrl         string `json:"package_url"`
	DownloadsDetailUrl string `json:"downloads_detail_url"`
	DownloadsSeriesUrl string `json:"downloads_series_url"`
	DownloadsCountUrl  string `json:"downloads_count_url"`
	PackageHtmlUrl     string `json:"package_html_url"`
	DownloadUrl        string `json:"download_url"`
	PromoteUrl         string `json:"promote_url"`
	DestroyUrl         string `json:"destroy_url"`
}
