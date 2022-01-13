package ante_test

import (
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ante "github.com/okex/exchain/app/ante"
	app "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

func buildTestTx(suite *AnteTestSuite) (sdk.Tx, error) {
	suite.ctx = suite.ctx.WithBlockHeight(10)
	addr1, priv1 := newTestAddrKey()
	addr2, _ := newTestAddrKey()
	acc1 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr1)
	_ = acc1.SetCoins(newTestCoins())
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc1)
	acc2 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr2)
	_ = acc2.SetCoins(newTestCoins())
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc2)
	to := ethcmn.BytesToAddress(addr2.Bytes())
	amt := big.NewInt(32)
	gas := big.NewInt(20)
	ethMsg := evmtypes.NewMsgEthereumTx(0, &to, amt, 22000, gas, []byte("test"))
	return newTestEthTx(suite.ctx, ethMsg, priv1)
}

func (suite *AnteTestSuite) TestWrappedTxSignatureRecover() {
	setConfidentKeyListWithCurrent(suite)
	tx, err := buildTestTx(suite)
	suite.Require().NoError(err)
	message, err := suite.app.Codec().MarshalBinaryLengthPrefixed(tx)
	suite.Require().NoError(err)
	signatue, err := suite.nodePriv.Sign(message)
	suite.Require().NoError(err)
	wrapped, err := NewWrappedTx(tx, signatue, suite.nodePub.Bytes())
	wrappedTx := wrapped.(app.WrappedTx)
	suite.Require().NoError(err)
	requireValidTx(suite.T(), suite.anteHandler, suite.ctx, wrappedTx.GetOriginTx(), false)

	suite.Require().NoError(err)
	confident, err := ante.VerifyConfidentTx(message, wrappedTx.Signature, wrappedTx.NodeKey)
	suite.Require().NoError(err)
	suite.Require().Equal(true, confident)
}

func (suite *AnteTestSuite) TestSkipWrappedSignaturePhase() {
	setConfidentKeyList(suite, true)
	ante.SetWrappedTxEffectiveHeight(1000)
	tx, err := buildTestTx(suite)
	suite.Require().NoError(err)
	message, _ := suite.app.Codec().MarshalBinaryLengthPrefixed(tx)
	signature, _ := suite.nodePriv.Sign(message)
	wrapped, _ := NewWrappedTx(tx, signature, suite.nodePub.Bytes())
	newCtx, err := suite.anteHandler(suite.ctx, wrapped, false)
	suite.Require().NoError(err)
	suite.Require().Equal(message, newCtx.ReplaceTx())
}

func (suite *AnteTestSuite) TestWrappedEtherrumTx() {
	setConfidentKeyList(suite, true)
	ante.SetWrappedTxEffectiveHeight(1)
	tx, err := buildTestTx(suite)
	suite.Require().NoError(err)
	message, _ := suite.app.Codec().MarshalBinaryLengthPrefixed(tx)
	suite.ctx = suite.ctx.WithTxBytes(message)
	signature, _ := suite.nodePriv.Sign(message)
	wrapped, _ := NewWrappedTx(tx, signature, suite.nodePub.Bytes())
	newCtx, err := suite.anteHandler(suite.ctx, tx, false)
	suite.Require().NoError(err)
	message, _ = suite.app.Codec().MarshalBinaryLengthPrefixed(wrapped)
	suite.Require().Equal(message, newCtx.ReplaceTx())
}

func (suite *AnteTestSuite) TestGenerateBad() {
	// use another node key to signature and verify with current
}

func (suite *AnteTestSuite) TestGenerateNoWrapper() {
	// confident node signature verify and return true
}

func (suite *AnteTestSuite) TestLightAnteEvm() {
	// test light evm tx ante
}

func (suite *AnteTestSuite) TestLightAnteStd() {
	// test light std tx ante
}
