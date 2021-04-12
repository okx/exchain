package cases

import (
	"os"
)

// GetBackendDBDir return the path "$GOPATH/src/github.com/okex/exchain/x/backend/cases"
func GetBackendDBDir() string {
	gopath := os.Getenv("GOPATH")
	dir := gopath + "/src/github.com/okex/exchain/x/backend/cases"
	return dir
}
