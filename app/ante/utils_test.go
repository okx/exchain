package ante_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/simapp/helpers"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	ibc_tx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/okex/exchain/libs/ibc-go/testing/mock"
	helpers2 "github.com/okex/exchain/libs/ibc-go/testing/simapp/helpers"

	"github.com/stretchr/testify/suite"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"

	"github.com/okex/exchain/app"
	ante "github.com/okex/exchain/app/ante"
	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	okexchain "github.com/okex/exchain/app/types"
	evmtypes "github.com/okex/exchain/x/evm/types"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
)

type AnteTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.OKExChainApp
	anteHandler sdk.AnteHandler
}

func (suite *AnteTestSuite) SetupTest() {
	checkTx := false
	chainId := "okexchain-3"

	suite.app = app.Setup(checkTx)
	suite.app.Codec().RegisterConcrete(&sdk.TestMsg{}, "test/TestMsg", nil)

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: chainId, Time: time.Now().UTC()})
	suite.app.EvmKeeper.SetParams(suite.ctx, evmtypes.DefaultParams())

	suite.anteHandler = ante.NewAnteHandler(suite.app.AccountKeeper, suite.app.EvmKeeper, suite.app.SupplyKeeper, nil, suite.app.WasmHandler, suite.app.IBCKeeper)

	err := okexchain.SetChainId(chainId)
	suite.Nil(err)

	appconfig.RegisterDynamicConfig(suite.app.Logger())
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

func newTestEthTx(ctx sdk.Context, msg *evmtypes.MsgEthereumTx, priv tmcrypto.PrivKey) (sdk.Tx, error) {
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

func newTxConfig() client.TxConfig {
	interfaceRegistry := types2.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	return ibc_tx.NewTxConfig(marshaler, ibc_tx.DefaultSignModes)
}

func mockIbcTx(accNum, seqNum []uint64, priv tmcrypto.PrivKey, chainId string, addr sdk.AccAddress) *sdk.Tx {
	txConfig := newTxConfig()
	packet := channeltypes.NewPacket([]byte(mock.MockPacketData), 1,
		"transfer", "channel-0",
		"transfer", "channel-1",
		clienttypes.NewHeight(1, 0), 0)
	msgs := []ibcmsg.Msg{channeltypes.NewMsgRecvPacket(packet, []byte("proof"), clienttypes.NewHeight(0, 1), addr.String())}
	ibcTx, err := helpers2.GenTx(
		txConfig,
		msgs,
		sdk.CoinAdapters{sdk.NewCoinAdapter(sdk.DefaultIbcWei, sdk.NewIntFromBigInt(big.NewInt(0)))},
		helpers.DefaultGenTxGas,
		chainId,
		accNum,
		seqNum,
		1,
		priv,
	)
	if err != nil {
		return nil
	}
	return &ibcTx
}
