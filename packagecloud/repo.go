package packagecloud

import (
	"errors"
	"fmt"
	"strings"
)

type Repo struct {
	User string
	Name string
}

func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.User, r.Name)
}

func NewRepo(user, name string) Repo {
	return Repo{
		User: user,
		Name: name,
	}
}

func NewRepoFromString(s string) (Repo, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || isEmptyString(parts[0]) || isEmptyString(parts[1]) {
		return Repo{}, errors.New("repo must be in the format 'user/repo'")
	}
	return Repo{
		User: parts[0],
		Name: parts[1],
	}, nil
}

func (r Repo) Validate() error {
	if isEmptyString(r.User) {
		return fmt.Errorf("invalid repository: %s, user cannot be empty", r.String())
	}
	if isEmptyString(r.Name) {
		return fmt.Errorf("invalid repository: %s, name cannot be empty", r.String())
	}
	return nil
}
