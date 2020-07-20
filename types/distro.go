package types

import "fmt"

type PackageTypes map[string][]Distro

type Distro struct {
	DisplayName string          `json:"display_name"`
	IndexName   string          `json:"index_name"`
	Versions    []DistroVersion `json:"versions"`
}

type DistroVersion struct {
	ID            int    `json:"id"`
	DisplayName   string `json:"display_name"`
	IndexName     string `json:"index_name"`
	VersionNumber string `json:"version_number"`
}

func GetDistroID(distros []Distro, packageType, name, version string) (int, error) {
	var matchingDistro *Distro
	for _, distro := range distros {
		if distro.IndexName == name {
			matchingDistro = &distro
			break
		}
	}
	if matchingDistro == nil {
		return -1, fmt.Errorf("distro was not found for given name: %s", name)
	}

	for _, v := range matchingDistro.Versions {
		if v.IndexName == version {
			return v.ID, nil
		}
	}

	return -1, fmt.Errorf("distro version was not found for given name and version: %s/%s", name, version)
}
