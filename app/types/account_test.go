package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"

	"errors"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/crypto/sr25519"

	"github.com/stretchr/testify/suite"

	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

func init() {
	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
}

type AccountTestSuite struct {
	suite.Suite

	account *EthAccount
}

func (suite *AccountTestSuite) SetupTest() {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	balance := sdk.NewCoins(NewPhotonCoin(sdk.OneInt()))
	baseAcc := auth.NewBaseAccount(addr, balance, pubkey, 10, 50)
	suite.account = &EthAccount{
		BaseAccount: baseAcc,
		CodeHash:    []byte{1, 2},
	}
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

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

		pos, _ := amino.ParseProtoPosAndTypeMustOneByte(data[0])
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
	v, ok := cdc.TryUnmarshalBinaryBareInterfaceWithRegisteredUbmarshaller(data, &account2)
	require.True(t, ok)
	account2, ok = v.(exported.Account)
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
			_, _ = cdc.TryUnmarshalBinaryBareInterfaceWithRegisteredUbmarshaller(data, &account)
		}
	})
}

func (suite *AccountTestSuite) TestEthAccount_Balance() {

	testCases := []struct {
		name         string
		denom        string
		initialCoins sdk.Coins
		amount       sdk.Int
	}{
		{"positive diff", NativeToken, sdk.Coins{}, sdk.OneInt()},
		{"zero diff, same coin", NativeToken, sdk.NewCoins(NewPhotonCoin(sdk.ZeroInt())), sdk.ZeroInt()},
		{"zero diff, other coin", sdk.DefaultBondDenom, sdk.NewCoins(NewPhotonCoin(sdk.ZeroInt())), sdk.ZeroInt()},
		{"negative diff", NativeToken, sdk.NewCoins(NewPhotonCoin(sdk.NewInt(10))), sdk.NewInt(1)},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset values
			suite.account.SetCoins(tc.initialCoins)

			suite.account.SetBalance(tc.denom, sdk.NewDecFromInt(tc.amount))
			suite.Require().Equal(sdk.NewDecFromInt(tc.amount), suite.account.Balance(tc.denom))
		})
	}

}

func (suite *AccountTestSuite) TestEthermintAccountJSON() {
	bz, err := json.Marshal(suite.account)
	suite.Require().NoError(err)

	bz1, err := suite.account.MarshalJSON()
	suite.Require().NoError(err)
	suite.Require().Equal(string(bz1), string(bz))

	var a EthAccount
	suite.Require().NoError(a.UnmarshalJSON(bz))
	suite.Require().Equal(suite.account.String(), a.String())
	suite.Require().Equal(suite.account.PubKey, a.PubKey)
}

func (suite *AccountTestSuite) TestEthermintPubKeyJSON() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	bz := privkey.PubKey().Bytes()

	pubk, err := tmamino.PubKeyFromBytes(bz)
	suite.Require().NoError(err)
	suite.Require().Equal(pubk, privkey.PubKey())
}

func (suite *AccountTestSuite) TestSecpPubKeyJSON() {
	pubkey := secp256k1.GenPrivKey().PubKey()
	bz := pubkey.Bytes()

	pubk, err := tmamino.PubKeyFromBytes(bz)
	suite.Require().NoError(err)
	suite.Require().Equal(pubk, pubkey)
}

func (suite *AccountTestSuite) TestEthermintAccount_String() {
	config := sdk.GetConfig()
	SetBech32Prefixes(config)

	bech32pubkey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, suite.account.PubKey)
	suite.Require().NoError(err)

	accountStr := fmt.Sprintf(`|
  address: %s
  eth_address: %s
  coins:
  - denom: %s
    amount: "1.000000000000000000"
  public_key: %s
  account_number: 10
  sequence: 50
  code_hash: "0102"
`, suite.account.Address, suite.account.EthAddress().String(), sdk.DefaultBondDenom, bech32pubkey)

	suite.Require().Equal(accountStr, suite.account.String())

	i, err := suite.account.MarshalYAML()
	suite.Require().NoError(err)

	var ok bool
	accountStr, ok = i.(string)
	suite.Require().True(ok)
	suite.Require().Contains(accountStr, suite.account.Address.String())
	suite.Require().Contains(accountStr, bech32pubkey)
}

func (suite *AccountTestSuite) TestEthermintAccount_MarshalJSON() {
	bz, err := suite.account.MarshalJSON()
	suite.Require().NoError(err)
	suite.Require().Contains(string(bz), suite.account.EthAddress().String())

	res := new(EthAccount)
	err = res.UnmarshalJSON(bz)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.account, res)

	bech32pubkey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, suite.account.PubKey)
	suite.Require().NoError(err)

	// test that the sdk.AccAddress is populated from the hex address
	jsonAcc := fmt.Sprintf(
		`{"address":"","eth_address":"%s","coins":[{"denom":"aphoton","amount":"1"}],"public_key":"%s","account_number":10,"sequence":50,"code_hash":"0102"}`,
		suite.account.EthAddress().String(), bech32pubkey,
	)

	res = new(EthAccount)
	err = res.UnmarshalJSON([]byte(jsonAcc))
	suite.Require().NoError(err)
	suite.Require().Equal(suite.account.Address.String(), res.Address.String())

	jsonAcc = fmt.Sprintf(
		`{"address":"","eth_address":"","coins":[{"denom":"aphoton","amount":"1"}],"public_key":"%s","account_number":10,"sequence":50,"code_hash":"0102"}`,
		bech32pubkey,
	)

	res = new(EthAccount)
	err = res.UnmarshalJSON([]byte(jsonAcc))
	suite.Require().Error(err, "should fail if both address are empty")

	// test that the sdk.AccAddress is populated from the hex address
	jsonAcc = fmt.Sprintf(
		`{"address": "%s","eth_address":"0x0000000000000000000000000000000000000000","coins":[{"denom":"aphoton","amount":"1"}],"public_key":"%s","account_number":10,"sequence":50,"code_hash":"0102"}`,
		suite.account.Address.String(), bech32pubkey,
	)

	res = new(EthAccount)
	err = res.UnmarshalJSON([]byte(jsonAcc))
	suite.Require().Error(err, "should fail if addresses mismatch")
}
