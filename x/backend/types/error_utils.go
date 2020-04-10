package types

import (
	"runtime/debug"
	"strings"
)

// nolint
type ErrorsMerged struct {
	errors []error
}

// nolint
func (em ErrorsMerged) Error() string {

	errStrs := []string{}
	for _, e := range em.errors {
		errStrs = append(errStrs, e.Error())
	}

	return strings.Join(errStrs, "; ")
}

// NewErrorsMerged plenty of errors into a single error.
func NewErrorsMerged(args ...error) error {

	filtered := []error{}
	for _, e := range args {
		if e != nil {
			filtered = append(filtered, e)
		}
	}

	if len(filtered) > 0 {
		return ErrorsMerged{errors: filtered}
	}
	return nil

}

// nolint
func PrintStackIfPanic() {
	r := recover()
	if r != nil {
		debug.PrintStack()
	}
}
