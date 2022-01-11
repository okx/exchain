package ante_test

import (
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

func (suite *AnteTestSuite) TestWrappedTxSignatureRecover() {
	suite.ctx = suite.ctx.WithBlockHeight(1)

	addr1, priv1 := newTestAddrKey()
	addr2, _ := newTestAddrKey()

	acc1 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr1)
	_ = acc1.SetCoins(newTestCoins())
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc1)

	acc2 := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr2)
	_ = acc2.SetCoins(newTestCoins())
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc2)

	// require a valid Ethereum tx to pass
	to := ethcmn.BytesToAddress(addr2.Bytes())
	amt := big.NewInt(32)
	gas := big.NewInt(20)
	ethMsg := evmtypes.NewMsgEthereumTx(0, &to, amt, 22000, gas, []byte("test"))

	tx, err := newTestEthTx(suite.ctx, ethMsg, priv1)
	suite.Require().NoError(err)

	// TODO: should refactor this function to build new wrapped tx generating logic

	//message, err := suite.app.Codec().MarshalBinaryBare(tx)
	suite.Require().NoError(err)
	//signatue, err := suite.nodePriv.Sign(message)
	suite.Require().NoError(err)
	//wrapped, err := NewWrappedTx(tx, signatue, suite.nodePub.Bytes())

	requireValidTx(suite.T(), suite.anteHandler, suite.ctx, tx, false)
}
