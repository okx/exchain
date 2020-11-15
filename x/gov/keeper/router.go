package keeper

import (
	"fmt"
	"regexp"

	"github.com/okex/okexchain/x/gov/types"
)

var (
	_ Router = (*router)(nil)

	isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
)

// Router implements a governance Handler router.
//
// TODO: Use generic router (ref #3976).
type Router interface {
	AddRoute(r string, h types.Handler) (rtr Router)
	HasRoute(r string) bool
	GetRoute(path string) (h types.Handler)
	Seal()
}

type router struct {
	routes map[string]types.Handler
	sealed bool
}

func NewRouter() Router {
	return &router{
		routes: make(map[string]types.Handler),
	}
}

// Seal seals the router which prohibits any subsequent route handlers to be
// added. Seal will panic if called more than once.
func (rtr *router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}
	rtr.sealed = true
}

// AddRoute adds a governance handler for a given path. It returns the Router
// so AddRoute calls can be linked. It will panic if the router is sealed.
func (rtr *router) AddRoute(path string, h types.Handler) Router {
	if rtr.sealed {
		panic("router sealed; cannot add route handler")
	}

	if !isAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.HasRoute(path) {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = h
	return rtr
}

// HasRoute returns true if the router has a path registered or false otherwise.
func (rtr *router) HasRoute(path string) bool {
	return rtr.routes[path] != nil
}

// GetRoute returns a Handler for a given path.
func (rtr *router) GetRoute(path string) types.Handler {
	if !rtr.HasRoute(path) {
		panic(fmt.Sprintf("route \"%s\" does not exist", path))
	}

	return rtr.routes[path]
}
