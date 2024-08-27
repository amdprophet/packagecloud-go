package packagecloud

import "errors"

var (
	ErrPackageAlreadyExists = errors.New("package already exists")
	ErrPaymentRequired      = errors.New("payment required")
)
