package types

import (
	"math/big"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/okex/exchain/libs/tendermint/crypto/sr25519"
	"github.com/stretchr/testify/require"
)

func TestEthAccountAmino(t *testing.T) {
	cdc := codec.New()
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	RegisterCodec(cdc)

	cdc.RegisterInterface((*tmcrypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(sr25519.PubKeySr25519{},
		sr25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoName, nil)

	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())

	accounts := []EthAccount{
		{},
		{
			auth.NewBaseAccount(
				addr,
				sdk.NewCoins(NewPhotonCoin(sdk.OneInt()), sdk.Coin{"heco", sdk.Dec{big.NewInt(1)}}),
				pubKey,
				1,
				1,
			),
			ethcrypto.Keccak256(nil),
		},
		{
			auth.NewBaseAccount(
				addr,
				sdk.NewCoins(NewPhotonCoin(sdk.ZeroInt()), sdk.Coin{"heco", sdk.Dec{big.NewInt(0)}}),
				pubKey,
				0,
				0,
			),
			ethcrypto.Keccak256(nil),
		},
		{
			auth.NewBaseAccount(
				nil,
				nil,
				nil,
				0,
				0,
			),
			ethcrypto.Keccak256(nil),
		},
		{
			BaseAccount: &auth.BaseAccount{},
		},
	}

	for _, testAccount := range accounts {
		data, err := cdc.MarshalBinaryBare(&testAccount)
		if err != nil {
			t.Fatal("marshal error")
		}

		var accountFromAmino exported.Account

		err = cdc.UnmarshalBinaryBare(data, &accountFromAmino)
		if err != nil {
			t.Fatal("unmarshal error")
		}

		var accountFromUnmarshaller exported.Account
		v, err := cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(data, &accountFromUnmarshaller)
		require.NoError(t, err)
		accountFromUnmarshaller, ok := v.(exported.Account)
		require.True(t, ok)

		require.EqualValues(t, accountFromAmino, accountFromUnmarshaller)

		var ethAccount EthAccount
		err = ethAccount.UnmarshalFromAmino(cdc, data[4:])
		require.NoError(t, err)
		require.EqualValues(t, accountFromAmino, &ethAccount)

		dataFromMarshaller, err := cdc.MarshalBinaryBareWithRegisteredMarshaller(&testAccount)
		require.NoError(t, err)
		require.EqualValues(t, data, dataFromMarshaller)

		dataFromMarshaller, err = ethAccount.MarshalToAmino(cdc)
		if dataFromMarshaller == nil {
			dataFromMarshaller = []byte{}
		}
		require.Equal(t, data[4:], dataFromMarshaller)
	}
}

func BenchmarkEthAccountAminoUnmarshal(b *testing.B) {
	cdc := codec.New()
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	RegisterCodec(cdc)

	cdc.RegisterInterface((*tmcrypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(sr25519.PubKeySr25519{},
		sr25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoName, nil)

	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())

	balance := sdk.NewCoins(NewPhotonCoin(sdk.OneInt()))
	testAccount := EthAccount{
		BaseAccount: auth.NewBaseAccount(addr, balance, pubKey, 1, 1),
		CodeHash:    ethcrypto.Keccak256(nil),
	}

	data, _ := cdc.MarshalBinaryBare(&testAccount)

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("amino", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var account exported.Account
			_ = cdc.UnmarshalBinaryBare(data, &account)
		}
	})

	b.Run("unmarshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var account exported.Account
			_, _ = cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(data, &account)
		}
	})
}

func BenchmarkEthAccountAminoMarshal(b *testing.B) {
	cdc := codec.New()
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	RegisterCodec(cdc)

	cdc.RegisterInterface((*tmcrypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(sr25519.PubKeySr25519{},
		sr25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoName, nil)

	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())

	balance := sdk.NewCoins(NewPhotonCoin(sdk.OneInt()))
	testAccount := EthAccount{
		BaseAccount: auth.NewBaseAccount(addr, balance, pubKey, 1, 1),
		CodeHash:    ethcrypto.Keccak256(nil),
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.Run("amino", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, _ := cdc.MarshalBinaryBare(&testAccount)
			_ = data
		}
	})

	b.Run("marshaller", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, _ := cdc.MarshalBinaryBareWithRegisteredMarshaller(&testAccount)
			_ = data
		}
	})
}
