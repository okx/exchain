package simapp_test

import (
	"testing"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/simapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"

	"github.com/stretchr/testify/require"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
)

func TestSimGenesisAccountValidate(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	vestingStart := time.Now().UTC()

	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 1000))
	baseAcc := authtypes.NewBaseAccount(addr, nil, pubkey, 0, 0)
	require.NoError(t, baseAcc.SetCoins(coins))

	testCases := []struct {
		name    string
		sga     simapp.SimGenesisAccount
		wantErr bool
	}{
		{
			"valid basic account",
			simapp.SimGenesisAccount{
				BaseAccount: baseAcc,
			},
			false,
		},
		{
			"invalid basic account with mismatching address/pubkey",
			simapp.SimGenesisAccount{
				BaseAccount: authtypes.NewBaseAccount(addr, nil, secp256k1.GenPrivKey().PubKey(), 0, 0),
			},
			true,
		},
		{
			"valid basic account with module name",
			simapp.SimGenesisAccount{
				BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(crypto.AddressHash([]byte("testmod"))), nil, nil, 0, 0),
				ModuleName:  "testmod",
			},
			false,
		},
		{
			"valid basic account with invalid module name/pubkey pair",
			simapp.SimGenesisAccount{
				BaseAccount: baseAcc,
				ModuleName:  "testmod",
			},
			true,
		},
		{
			"valid basic account with valid vesting attributes",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins,
				StartTime:       vestingStart.Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			false,
		},
		{
			"valid basic account with invalid vesting end time",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins,
				StartTime:       vestingStart.Add(2 * time.Hour).Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			true,
		},
		{
			"valid basic account with invalid original vesting coins",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins.Add(coins...),
				StartTime:       vestingStart.Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.wantErr, tc.sga.Validate() != nil)
		})
	}
}

func TestSimGenesisAccountRLP(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	vestingStart := time.Now().UTC()

	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 1000))
	baseAcc := authtypes.NewBaseAccount(addr, nil, pubkey, 0, 0)
	require.NoError(t, baseAcc.SetCoins(coins))

	sga := simapp.SimGenesisAccount{
		BaseAccount:     baseAcc,
		OriginalVesting: coins.Add(coins...),
		StartTime:       vestingStart.Unix(),
		EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
	}

	data, err := sga.RLPEncodeToBytes()
	require.NoError(t, err)

	var sgacc simapp.SimGenesisAccount
	err = sgacc.RLPDecodeBytes(data)
	require.NoError(t, err)

	require.Equal(t, sgacc.ModuleName, sga.ModuleName)
	require.Equal(t, len(sgacc.ModulePermissions), len(sga.ModulePermissions))
}