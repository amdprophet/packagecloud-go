package util

import "testing"

func TestSliceContainsString(t *testing.T) {
	slice := []string{"'my-package_0.1.0_amd64.deb' has already been taken"}
	result := SliceContainsString(slice, "has already been taken")
	if !result {
		t.Error("slice did not contain string")
	}
}
