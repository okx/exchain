package baseapp

import (
	"strings"
)

// Supported ABCI Query prefixes
const (
	QueryPathApp    = "app"
	QueryPathCustom = "custom"
	QueryPathP2P    = "p2p"
	QueryPathStore  = "store"
)

// SplitABCIQueryPath splits a string path using the delimiter '/'.
//
// e.g. "this/is/funny" becomes []string{"this", "is", "funny"}
func SplitABCIQueryPath(requestPath string) (path []string) {
	path = strings.Split(requestPath, "/")

	// first element is empty string
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	return path
}
