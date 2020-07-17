package types

type DistroFormats map[string][]Distro

type Distro struct {
	DisplayName string          `json:"display_name"`
	IndexName   string          `json:"index_name"`
	Versions    []DistroVersion `json:"versions"`
}

type DistroVersion struct {
	ID            int64  `json:"id"`
	DisplayName   string `json:"display_name"`
	IndexName     string `json:"index_name"`
	VersionNumber string `json:"version_number"`
}
