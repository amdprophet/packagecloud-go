package packagecloud

import (
	"errors"
	"fmt"
)

var (
	ErrPackageAlreadyExists = errors.New("package already exists")
	ErrPaymentRequired      = errors.New("payment required")
)

type UnmarshalError struct {
	Data []byte
	Err  error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("failed to unmarshal data: %s\n%s", e.Err, string(e.Data))
}
