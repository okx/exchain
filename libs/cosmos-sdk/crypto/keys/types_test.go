package keys

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"

	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func Test_writeReadLedgerInfo(t *testing.T) {
	var tmpKey secp256k1.PubKeySecp256k1
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A")
	copy(tmpKey[:], bz)

	lInfo := newLedgerInfo("some_name", tmpKey, *hd.NewFundraiserParams(5, sdk.CoinType, 1), Secp256k1)
	require.Equal(t, TypeLedger, lInfo.GetType())

	path, err := lInfo.GetPath()
	require.NoError(t, err)
	require.Equal(t, "m/44'/118'/5'/0/1", path.String())
	require.Equal(t,
		"cosmospub1addwnpepqddddqg2glc8x4fl7vxjlnr7p5a3czm5kcdp4239sg6yqdc4rc2r5wmxv8p",
		sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, lInfo.GetPubKey()))

	// Serialize and restore
	serialized := marshalInfo(lInfo)
	restoredInfo, err := unmarshalInfo(serialized)
	require.NoError(t, err)
	require.NotNil(t, restoredInfo)

	// Check both keys match
	require.Equal(t, lInfo.GetName(), restoredInfo.GetName())
	require.Equal(t, lInfo.GetType(), restoredInfo.GetType())
	require.Equal(t, lInfo.GetPubKey(), restoredInfo.GetPubKey())

	restoredPath, err := restoredInfo.GetPath()
	require.NoError(t, err)
	require.Equal(t, path, restoredPath)
}
