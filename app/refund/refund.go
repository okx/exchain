package refund

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/x/auth/refund"

	ethermint "github.com/okex/okexchain/app/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	evmtypes "github.com/okex/okexchain/x/evm/types"
)

type EVMKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
}

func NewGasRefundHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (err error) {
		var gasRefundHandler sdk.GasRefundHandler
		switch tx.(type) {
		case evmtypes.MsgEthereumTx:
			gasRefundHandler = EthGasRefundDecorator(ak, evmKeeper, sk)
		default:
			return nil
		}
		return gasRefundHandler(ctx, tx)
	}
}

type EthGasRefundHandler struct {
	ak        auth.AccountKeeper
	sk        types.SupplyKeeper
	evmKeeper EVMKeeper
}

func (egrh EthGasRefundHandler) GasRefundHandle(ctx sdk.Context, tx sdk.Tx) (err error) {

	currentGasMeter := ctx.GasMeter()
	TempGasMeter := sdk.NewInfiniteGasMeter()
	ctx = ctx.WithGasMeter(TempGasMeter)

	defer func() {
		ctx = ctx.WithGasMeter(currentGasMeter)
	}()

	gasLimit := currentGasMeter.Limit()
	gasUsed := currentGasMeter.GasConsumed()

	if gasUsed >= gasLimit {
		return nil
	}

	msgEthTx, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}

	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return err
	}
	_, err = msgEthTx.VerifySig(chainIDEpoch)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "signature verification failed: %s", err.Error())
	}

	// sender address should be in the tx cache from the previous AnteHandle call
	address := msgEthTx.From()
	if address.Empty() {
		panic("sender address cannot be empty")
	}

	gasLeft := new(big.Int).Sub(new(big.Int).SetUint64(gasLimit), new(big.Int).SetUint64(gasUsed))
	gasRefund := new(big.Int).Mul(msgEthTx.Data.Price, gasLeft)

	evmDenom := egrh.evmKeeper.GetParams(ctx).EvmDenom
	feeAmt := sdk.NewCoins(
		sdk.NewCoin(evmDenom, sdk.NewDecFromBigIntWithPrec(gasRefund, sdk.Precision)),
	)

	err = refund.RefundFees(egrh.sk, ctx, address, feeAmt)
	if err != nil {
		return err
	}

	return nil
}

func EthGasRefundDecorator(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {

	egrh := EthGasRefundHandler{
		ak:        ak,
		sk:        sk,
		evmKeeper: evmKeeper,
	}

	return func(ctx sdk.Context, tx sdk.Tx) (err error) {
		return egrh.GasRefundHandle(ctx, tx)
	}
}
