package cli

import (
	"encoding/hex"
	"strings"
	"testing"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestGetValsSet(t *testing.T) {
	pubKeys := []crypto.PubKey{
		newPubKey("0000000000000000000000000000000000000000000000000000000000000000"),
		newPubKey("1111111111111111111111111111111111111111111111111111111111111111"),
		newPubKey("2222222222222222222222222222222222222222222222222222222222222222"),
		newPubKey("3333333333333333333333333333333333333333333333333333333333333333"),
	}

	accAddrs := []sdk.AccAddress{
		sdk.AccAddress(pubKeys[0].Address()),
		sdk.AccAddress(pubKeys[1].Address()),
		sdk.AccAddress(pubKeys[2].Address()),
		sdk.AccAddress(pubKeys[3].Address()),
	}

	var valAddrsStr []string
	var expectedValAddrs []sdk.ValAddress
	for i := 0; i < 4; i++ {
		valAddr := sdk.ValAddress(accAddrs[i])
		expectedValAddrs = append(expectedValAddrs, valAddr)
		valAddrsStr = append(valAddrsStr, valAddr.String())
	}

	arg := strings.Join(valAddrsStr, ",")

	valAddrs, err := getValsSet(arg)
	require.NoError(t, err)
	require.Equal(t, expectedValAddrs, valAddrs)
}

func newPubKey(pubKey string) (res crypto.PubKey) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}
	var pubKeyEd25519 ed25519.PubKeyEd25519
	copy(pubKeyEd25519[:], pubKeyBytes[:])
	return pubKeyEd25519
}
