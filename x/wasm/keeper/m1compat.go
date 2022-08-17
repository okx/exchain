package keeper

import (
	"runtime"
	"testing"
)

func SkipIfM1(t *testing.T) {
	if runtime.GOARCH == "arm64" {
		t.Skip("Skipping for M1: Signal Error, Stack Dump")
	}
}
