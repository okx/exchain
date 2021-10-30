package ante_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/okex/exchain/dependence/cosmos-sdk/server"

	"github.com/stretchr/testify/suite"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/auth"

	"github.com/okex/exchain/app"
	ante "github.com/okex/exchain/app/ante"
	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	okexchain "github.com/okex/exchain/app/types"
	evmtypes "github.com/okex/exchain/x/evm/types"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	abci "github.com/tendermint/tendermint/abci/types"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

type AnteTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.OKExChainApp
	anteHandler sdk.AnteHandler
}

func (suite *AnteTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.app.Codec().RegisterConcrete(&sdk.TestMsg{}, "test/TestMsg", nil)

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "okexchain-3", Time: time.Now().UTC()})
	suite.app.EvmKeeper.SetParams(suite.ctx, evmtypes.DefaultParams())

	suite.anteHandler = ante.NewAnteHandler(suite.app.AccountKeeper, suite.app.EvmKeeper, suite.app.SupplyKeeper, nil)

	appconfig.RegisterDynamicConfig(server.NewDefaultContext())
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func newTestMsg(addrs ...sdk.AccAddress) *sdk.TestMsg {
	return sdk.NewTestMsg(addrs...)
}

func newTestCoins() sdk.Coins {
	return sdk.NewCoins(okexchain.NewPhotonCoinInt64(500000000))
}

func newTestStdFee() auth.StdFee {
	return auth.NewStdFee(220000, sdk.NewCoins(okexchain.NewPhotonCoinInt64(150)))
}

// GenerateAddress generates an Ethereum address.
func newTestAddrKey() (sdk.AccAddress, tmcrypto.PrivKey) {
	privkey, _ := ethsecp256k1.GenerateKey()
	addr := ethcrypto.PubkeyToAddress(privkey.ToECDSA().PublicKey)

	return sdk.AccAddress(addr.Bytes()), privkey
}

func newTestSDKTx(
	ctx sdk.Context, msgs []sdk.Msg, privs []tmcrypto.PrivKey,
	accNums []uint64, seqs []uint64, fee auth.StdFee,
) sdk.Tx {

	sigs := make([]auth.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := auth.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, "")

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:    priv.PubKey(),
			Signature: sig,
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, "")
}

func newTestEthTx(ctx sdk.Context, msg evmtypes.MsgEthereumTx, priv tmcrypto.PrivKey) (sdk.Tx, error) {
	chainIDEpoch, err := okexchain.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, err
	}

	privkey, ok := priv.(ethsecp256k1.PrivKey)
	if !ok {
		return nil, fmt.Errorf("invalid private key type: %T", priv)
	}

	if err := msg.Sign(chainIDEpoch, privkey.ToECDSA()); err != nil {
		return nil, err
	}

	return msg, nil
}
