package global

import "strings"

func IsTiKv(rootDir string) bool {
	return strings.Contains(rootDir, "127.0.0.1")
}
