package packagecloud

import (
	"fmt"
	"strings"
)

func isEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}

func bugPanic(msg string) {
	// This function is used to panic with a message indicating a bug in the code.
	// It should only be called when the code is expected to never reach this point
	// due to a logic error or incorrect usage.
	panic(fmt.Sprintf("BUG: %s", msg))
}
