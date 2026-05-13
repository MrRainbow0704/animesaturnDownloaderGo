package version

import "strings"

var version = "v0.0.0"

func Get() string {
	return version
}

func IsDev() bool {
	return strings.Contains(Get(), "dev")
}
