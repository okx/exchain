package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestSplitBorrowedKey(t *testing.T) {
	pk := ed25519.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pk.Address())
	product := "test-product"

	key := GetBorrowedKey(addr, product)
	splitAddr := SplitBorrowedKey(key, product)
	require.True(t, addr.Equals(splitAddr))
}
