package state_test

import (
	"os"
	"testing"

	"github.com/okex/exchain/libs/tendermint/types"
)

func TestMain(m *testing.M) {
	types.RegisterMockEvidencesGlobal()
	os.Exit(m.Run())
}
