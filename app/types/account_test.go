package types_test

import (
	"encoding/json"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"testing"

	"github.com/stretchr/testify/suite"

	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"

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
		StateRoot: ethcmn.Hash{},
	}
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
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

func (suite *AccountTestSuite) TestEthAccountRLP() {
	suite.SetupTest()
	//data, err := suite.account.RLPEncodeToBytes()
	//suite.Require().NoError(err, "fail to use rlp to encode ethAccount")

	data, err := rlp.EncodeToBytes(suite.account)
	suite.Require().NoError(err, "fail to use rlp to encode ethAccount")

	var acc types.EthAccount
	err = acc.RLPDecodeBytes(data)
	suite.Require().NoError(err, "fail to use rlp to decode ethAccount")

	suite.Require().Equal(suite.account.CodeHash, acc.CodeHash)
}