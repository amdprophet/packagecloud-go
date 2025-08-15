package packagecloud

import (
	"errors"
	"fmt"
	"strings"
)

type Distro struct {
	Name    string
	Version string
}

func (d Distro) String() string {
	return fmt.Sprintf("%s/%s", d.Name, d.Version)
}

func NewDistro(name, version string) Distro {
	return Distro{
		Name:    name,
		Version: version,
	}
}

func NewDistroFromString(s string) (Distro, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || isEmptyString(parts[0]) || isEmptyString(parts[1]) {
		return Distro{}, errors.New("distro must be in the format 'name/version'")
	}
	return Distro{
		Name:    parts[0],
		Version: parts[1],
	}, nil
}

func (d Distro) Validate() error {
	if isEmptyString(d.Name) {
		return fmt.Errorf("invalid distro: %s, name cannot be empty", d.String())
	}
	if isEmptyString(d.Version) {
		return fmt.Errorf("invalid distro: %s, version cannot be empty", d.String())
	}
	return nil
}
