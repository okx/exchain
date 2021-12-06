package wrap

import (
	"github.com/ethereum/go-ethereum/rlp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrapAccountRLP(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := types.NewBaseAccount(addr, nil, pubkey, 0, 0)

	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 246)}
	seq := uint64(7)

	// set everything on the account
	err := baseAcc.SetSequence(seq)
	require.Nil(t, err)
	err = baseAcc.SetCoins(someCoins)
	require.Nil(t, err)

	rst, err := rlp.EncodeToBytes(baseAcc)
	require.Nil(t, err)

	var baAcc WrapAccount
	err = rlp.DecodeBytes(rst, &baAcc)
	require.Nil(t, err)

	require.Equal(t, baseAcc.Address, baAcc.RealAcc.GetAddress())
}