package util

import (
	"os"
	"strconv"
	"strings"
)

// EnvBool reads a boolean environment variable.
// Returns (value, true) when the variable is set; (false, false) when unset.
func EnvBool(key string) (bool, bool) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return false, false
	}
	b, err := strconv.ParseBool(v)
	if err == nil {
		return b, true
	}
	switch strings.ToLower(v) {
	case "yes", "on":
		return true, true
	case "no", "off":
		return false, true
	default:
		return false, false
	}
}
