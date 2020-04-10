package cases

import (
	"os"
)

func GetBackendDBDir() string {
	gopath := os.Getenv("GOPATH")
	dir := gopath + "/src/github.com/okex/okchain/x/backend/cases"
	return dir
}
