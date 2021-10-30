package types

import (
	"encoding/hex"
	"testing"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	FixPK   = ed25519.GenPrivKeyFromSecret([]byte{1}).PubKey()
	FixAddr = FixPK.Address()
)

func TestGetValidatorPowerRank(t *testing.T) {
	valAddr1 := sdk.ValAddress(FixAddr)
	emptyDesc := Description{}
	val1 := NewValidator(valAddr1, pk1, emptyDesc, DefaultMinSelfDelegation)
	val1.DelegatorShares = sdk.ZeroDec()
	val2, val3, val4 := val1, val1, val1
	val2.DelegatorShares = sdk.OneDec()
	val3.DelegatorShares = sdk.OneDec().MulInt64(10)
	val4.DelegatorShares = sdk.OneDec().MulInt64(1 << 16)

	tests := []struct {
		validator Validator
		wantHex   string
	}{
		{val1, "2300000000000000009c288ede7df62742fc3b7d0962045a8cef0f79f6"},
		{val2, "2300000000000000019c288ede7df62742fc3b7d0962045a8cef0f79f6"},
		{val3, "23000000000000000a9c288ede7df62742fc3b7d0962045a8cef0f79f6"},
		{val4, "2300000000000100009c288ede7df62742fc3b7d0962045a8cef0f79f6"},
	}
	for i, tt := range tests {
		got := hex.EncodeToString(getValidatorPowerRank(tt.validator))

		assert.Equal(t, tt.wantHex, got, "Keys did not match on test case %d", i)
	}
}
