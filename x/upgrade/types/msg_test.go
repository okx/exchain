package types

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

// just 4 test
/////////////////////////////////////////////////////////////
var (
	pubkey = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	addr   = sdk.AccAddress(pubkey.Address())
)

func TestMsgUpgradeConfig_ValidateBasic(t *testing.T) {
	msg := NewMsgUpgradeConfig(1, 1, 100, "http://web.abc", addr)
	require.NoError(t, msg.ValidateBasic())

	msg = NewMsgUpgradeConfig(1, 1, 100, "http://web.abc", sdk.AccAddress{})
	require.Error(t, msg.ValidateBasic())
}

func TestMsgUpgradeConfig_GetSignBytes(t *testing.T) {
	msg := NewMsgUpgradeConfig(1, 1, 100, "http://web.abc", addr)
	require.NotEqual(t, len(msg.GetSignBytes()), 0)
}

func TestMsgUpgradeConfig_GetSigners(t *testing.T) {
	msg := NewMsgUpgradeConfig(1, 1, 100, "http://web.abc", addr)
	require.NotEqual(t, len(msg.GetSigners()), 0)
}

func TestMsgUpgradeConfig_Route(t *testing.T) {
	msg := NewMsgUpgradeConfig(1, 1, 100, "http://web.abc", addr)
	require.NotEqual(t, len(msg.Route()), 0)
}

func TestMsgUpgradeConfig_Type(t *testing.T) {
	msg := NewMsgUpgradeConfig(1, 1, 100, "http://web.abc", addr)
	require.NotEqual(t, len(msg.Type()), 0)
}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}
/////////////////////////////////////////////////////////////
