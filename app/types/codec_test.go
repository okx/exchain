package types

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/okex/exchain/libs/tendermint/crypto/sr25519"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

func unmarshalEthAccountFromAminoWithTypePrefix(data []byte) (*EthAccount, error) {
	var typePrefix = []byte{0x4c, 0x96, 0xdf, 0xce}
	if 0 != bytes.Compare(typePrefix, data[0:4]) {
		return nil, errors.New("type error")
	}
	data = data[4:]

	var dataLen uint64 = 0
	account := &EthAccount{}

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, _, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return nil, err
		}
		data = data[1:]

		var n int
		dataLen, n, _ = amino.DecodeUvarint(data)

		data = data[n:]
		subData := data[:dataLen]

		switch pos {
		case 1:
			baseAccount, err := authtypes.UnmarshalBaseAccountFromAmino(subData)
			if err != nil {
				return nil, err
			}
			account.BaseAccount = baseAccount
		case 2:
			account.CodeHash = subData
		}
	}
	return account, nil
}

func TestAccountAmino(t *testing.T) {
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

	// ak := mock.NewAddrKeys(addr, pubKey, privKey)
	balance := sdk.NewCoins(NewPhotonCoin(sdk.OneInt()), sdk.Coin{"heco", sdk.Dec{big.NewInt(1)}})
	testAccount := EthAccount{
		BaseAccount: auth.NewBaseAccount(addr, balance, pubKey, 1, 1),
		CodeHash:    ethcrypto.Keccak256(nil),
	}

	data, err := cdc.MarshalBinaryBare(&testAccount)
	if err != nil {
		t.Fatal("marshal error")
	}

	var account exported.Account

	err = cdc.UnmarshalBinaryBare(data, &account)
	if err != nil {
		t.Fatal("unmarshal error")
	}

	ethAccount, err := unmarshalEthAccountFromAminoWithTypePrefix(data)
	require.NoError(t, err)

	var account2 exported.Account
	v, err := cdc.UnmarshalBinaryBareWithRegisteredUbmarshaller(data, &account2)
	require.NoError(t, err)
	account2, ok := v.(exported.Account)
	require.True(t, ok)

	require.EqualValues(t, &testAccount, ethAccount)
	require.EqualValues(t, &testAccount, account)
	require.EqualValues(t, &testAccount, account2)
}

func BenchmarkAccountAmino(b *testing.B) {
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

	// ak := mock.NewAddrKeys(addr, pubKey, privKey)
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

	b.Run("direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = unmarshalEthAccountFromAminoWithTypePrefix(data)
		}
	})

	b.Run("amino-direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var account exported.Account
			_, _ = cdc.UnmarshalBinaryBareWithRegisteredUbmarshaller(data, &account)
		}
	})
}
