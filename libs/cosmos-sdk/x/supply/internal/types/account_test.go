package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/crypto/sr25519"

	yaml "gopkg.in/yaml.v2"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"

	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"

	"github.com/stretchr/testify/require"
)

func TestModuleAccountMarshalYAML(t *testing.T) {
	name := "test"
	moduleAcc := NewEmptyModuleAccount(name, Minter, Burner, Staking)
	bs, err := yaml.Marshal(moduleAcc)
	require.NoError(t, err)

	want := "|\n  address: cosmos1n7rdpqvgf37ktx30a2sv2kkszk3m7ncmg5drhe\n  coins: []\n  public_key: \"\"\n  account_number: 0\n  sequence: 0\n  name: test\n  permissions:\n  - minter\n  - burner\n  - staking\n"
	require.Equal(t, want, string(bs))
}

func TestHasPermissions(t *testing.T) {
	name := "test"
	macc := NewEmptyModuleAccount(name, Staking, Minter, Burner)
	cases := []struct {
		permission string
		expectHas  bool
	}{
		{Staking, true},
		{Minter, true},
		{Burner, true},
		{"other", false},
	}

	for i, tc := range cases {
		hasPerm := macc.HasPermission(tc.permission)
		if tc.expectHas {
			require.True(t, hasPerm, "test case #%d", i)
		} else {
			require.False(t, hasPerm, "test case #%d", i)
		}
	}
}

func TestValidate(t *testing.T) {
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	baseAcc := authtypes.NewBaseAccount(addr, sdk.Coins{}, nil, 0, 0)
	tests := []struct {
		name   string
		acc    authexported.GenesisAccount
		expErr error
	}{
		{
			"valid module account",
			NewEmptyModuleAccount("test"),
			nil,
		},
		{
			"invalid name and address pair",
			NewModuleAccount(baseAcc, "test"),
			fmt.Errorf("address %s cannot be derived from the module name 'test'", addr),
		},
		{
			"empty module account name",
			NewModuleAccount(baseAcc, "    "),
			errors.New("module account name cannot be blank"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acc.Validate()
			require.Equal(t, tt.expErr, err)
		})
	}
}

func TestModuleAccountJSON(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 5))
	baseAcc := authtypes.NewBaseAccount(addr, coins, nil, 10, 50)
	acc := NewModuleAccount(baseAcc, "test", "burner")

	bz, err := json.Marshal(acc)
	require.NoError(t, err)

	bz1, err := acc.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz1), string(bz))

	var a ModuleAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, acc.String(), a.String())
}

func TestModuleAccountAmino(t *testing.T) {
	cdc := codec.New()
	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	RegisterCodec(cdc)
	cdc.RegisterInterface((*tmcrypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(sr25519.PubKeySr25519{},
		sr25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoName, nil)

	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	accounts := []ModuleAccount{
		{
			authtypes.NewBaseAccount(
				addr,
				sdk.NewCoins(sdk.Coin{"heco", sdk.Dec{big.NewInt(1)}}),
				pubkey,
				1,
				1,
			),
			"name",
			[]string{"perm1", "perm2"},
		},
		{
			authtypes.NewBaseAccount(
				addr,
				sdk.NewCoins(sdk.NewInt64Coin("test", 5), sdk.NewInt64Coin("ok", 100000)),
				pubkey,
				9098,
				1000,
			),
			"name",
			[]string{"perm1", "perm2"},
		},
		{
			authtypes.NewBaseAccount(
				nil,
				nil,
				nil,
				0,
				0,
			),
			"",
			[]string{""},
		},
		{
			authtypes.NewBaseAccount(
				nil,
				nil,
				nil,
				0,
				0,
			),
			"",
			nil,
		},
	}

	for _, acc := range accounts {
		var iacc authexported.Account = acc
		bz, err := cdc.MarshalBinaryBare(iacc)
		require.NoError(t, err)

		var accountExpect authexported.Account
		err = cdc.UnmarshalBinaryBare(bz, &accountExpect)
		require.NoError(t, err)

		var account authexported.Account
		v, err := cdc.UnmarshalBinaryBareWithRegisteredUbmarshaller(bz, &account)
		require.NoError(t, err)
		accountActual, ok := v.(authexported.Account)
		require.True(t, ok)
		require.EqualValues(t, accountExpect, accountActual)

		nbz, err := cdc.MarshalBinaryBareWithRegisteredMarshaller(iacc)
		require.NoError(t, err)
		require.EqualValues(t, bz, nbz)
	}
}
