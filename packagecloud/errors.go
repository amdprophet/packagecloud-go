package packagecloud

import (
	"errors"
	"fmt"
)

var (
	ErrPackageAlreadyExists = errors.New("package already exists")
	ErrPaymentRequired      = errors.New("payment required")
)

type MissingSearchOptionsError struct{}

func (e *MissingSearchOptionsError) Error() string {
	return "one or more of the query, filter, dist, and/or arch flags must be specified"
}

type UnmarshalError struct {
	Data []byte
	Err  error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("failed to unmarshal data: %s\n%s", e.Err, string(e.Data))
}

type MissingOptionError struct {
	Field string
}

func (e *MissingOptionError) Error() string {
	return fmt.Sprintf("missing required option: %s", e.Field)
}
