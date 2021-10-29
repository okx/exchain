package types_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"math/big"

	"errors"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/sr25519"
	"testing"

	"github.com/stretchr/testify/suite"

	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/types"
)

func init() {
	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
}

type AccountTestSuite struct {
	suite.Suite

	account *types.EthAccount
}

func (suite *AccountTestSuite) SetupTest() {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	balance := sdk.NewCoins(types.NewPhotonCoin(sdk.OneInt()))
	baseAcc := auth.NewBaseAccount(addr, balance, pubkey, 10, 50)
	suite.account = &types.EthAccount{
		BaseAccount: baseAcc,
		CodeHash:    []byte{1, 2},
	}
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

func parsePosAndType(data byte) (pos int, aminoType amino.Typ3) {
	aminoType = amino.Typ3(data & 0x07)
	pos = int(data) >> 3
	return
}

func unmarshalEthAccountFromAmino(data []byte) (*types.EthAccount, error) {
	var typePrefix = []byte{0x4c, 0x96, 0xdf, 0xce}
	if 0 != bytes.Compare(typePrefix, data[0:4]) {
		return nil, errors.New("type error")
	}
	data = data[4:]

	var dataLen uint64 = 0
	account := &types.EthAccount{}

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, _ := parsePosAndType(data[0])
		data = data[1:]

		var n int
		dataLen, n, _ = amino.DecodeUvarint(data)

		data = data[n:]
		subData := data[:dataLen]

		switch pos {
		case 1:
			baseAccount, err := unmarshalBaseAccountFromAmino(subData)
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

func unmarshalBaseAccountFromAmino(data []byte) (*auth.BaseAccount, error) {
	var dataLen uint64 = 0
	var subData []byte
	account := &auth.BaseAccount{}

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, aminoType := parsePosAndType(data[0])
		data = data[1:]

		if aminoType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			account.Address = make([]byte, len(subData), len(subData))
			copy(account.Address, subData)
			// account.Address = subData
		case 2:
			coin, err := unmarshalCoinFromAmino(subData)
			if err != nil {
				return nil, err
			}
			account.Coins = append(account.Coins, coin)
		case 3:
			pubkey, err := unmarshalPubKeyFromAmino(subData)
			if err != nil {
				return nil, err
			}
			account.PubKey = pubkey
		case 4:
			uvarint, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}
			account.AccountNumber = uvarint
			dataLen = uint64(n)
		case 5:
			uvarint, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return nil, err
			}
			account.Sequence = uvarint
			dataLen = uint64(n)
		}
	}
	return account, nil
}

func unmarshalCoinFromAmino(data []byte) (coin sdk.DecCoin, err error) {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]

		if len(data) <= 0 {
			break
		}

		pos, aminoType := parsePosAndType(data[0])
		data = data[1:]

		if aminoType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			coin.Denom = string(subData)
		case 2:
			amt := big.NewInt(0)
			err = amt.UnmarshalText(subData)
			if err != nil {
				return
			}
			coin.Amount = sdk.Dec{
				amt,
			}
		}
	}
	return
}

var typePubKeySecp256k1Prefix = []byte{0xeb, 0x5a, 0xe9, 0x87}

func unmarshalPubKeyFromAmino(data []byte) (tmcrypto.PubKey, error) {
	if 0 == bytes.Compare(typePubKeySecp256k1Prefix, data[0:4]) {
		if data[4] != 33 {
			return nil, errors.New("pubkey secp256k1 size error")
		}
		data = data[5:]
		pubKey := secp256k1.PubKeySecp256k1{}
		copy(pubKey[:], data)
		return pubKey, nil
	}
	panic("not implement")
}

func TestAccountAmino(t *testing.T) {
	cdc := codec.New()
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	types.RegisterCodec(cdc)

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
	balance := sdk.NewCoins(types.NewPhotonCoin(sdk.OneInt()), sdk.Coin{"heco", sdk.Dec{big.NewInt(1)}})
	testAccount := types.EthAccount{
		BaseAccount: auth.NewBaseAccount(addr, balance, pubKey, 1, 1),
		CodeHash: ethcrypto.Keccak256(nil),
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

	ethAccount, err := unmarshalEthAccountFromAmino(data)
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
	types.RegisterCodec(cdc)

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
	balance := sdk.NewCoins(types.NewPhotonCoin(sdk.OneInt()))
	testAccount := types.EthAccount{
		BaseAccount: auth.NewBaseAccount(addr, balance, pubKey, 1, 1),
		CodeHash: ethcrypto.Keccak256(nil),
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
			_, _ = unmarshalEthAccountFromAmino(data)
		}
	})

	b.Run("amino-direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var account exported.Account
			_, _ = cdc.TryUnmarshalBinaryBareInterfaceWithRegisteredUbmarshaller(data,&account)
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
		{"positive diff", types.NativeToken, sdk.Coins{}, sdk.OneInt()},
		{"zero diff, same coin", types.NativeToken, sdk.NewCoins(types.NewPhotonCoin(sdk.ZeroInt())), sdk.ZeroInt()},
		{"zero diff, other coin", sdk.DefaultBondDenom, sdk.NewCoins(types.NewPhotonCoin(sdk.ZeroInt())), sdk.ZeroInt()},
		{"negative diff", types.NativeToken, sdk.NewCoins(types.NewPhotonCoin(sdk.NewInt(10))), sdk.NewInt(1)},
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

	var a types.EthAccount
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
	types.SetBech32Prefixes(config)

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

	res := new(types.EthAccount)
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

	res = new(types.EthAccount)
	err = res.UnmarshalJSON([]byte(jsonAcc))
	suite.Require().NoError(err)
	suite.Require().Equal(suite.account.Address.String(), res.Address.String())

	jsonAcc = fmt.Sprintf(
		`{"address":"","eth_address":"","coins":[{"denom":"aphoton","amount":"1"}],"public_key":"%s","account_number":10,"sequence":50,"code_hash":"0102"}`,
		bech32pubkey,
	)

	res = new(types.EthAccount)
	err = res.UnmarshalJSON([]byte(jsonAcc))
	suite.Require().Error(err, "should fail if both address are empty")

	// test that the sdk.AccAddress is populated from the hex address
	jsonAcc = fmt.Sprintf(
		`{"address": "%s","eth_address":"0x0000000000000000000000000000000000000000","coins":[{"denom":"aphoton","amount":"1"}],"public_key":"%s","account_number":10,"sequence":50,"code_hash":"0102"}`,
		suite.account.Address.String(), bech32pubkey,
	)

	res = new(types.EthAccount)
	err = res.UnmarshalJSON([]byte(jsonAcc))
	suite.Require().Error(err, "should fail if addresses mismatch")
}
