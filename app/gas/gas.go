package gas

import (
	"math/big"

	ethermint "github.com/okex/okexchain/app/types"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	evmtypes "github.com/okex/okexchain/x/evm/types"
)

type EVMKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
}

func NewGasHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper) sdk.GasHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (err error) {
		var gasHandler sdk.GasHandler
		switch tx.(type) {
		case auth.StdTx:
			gasHandler = CosmosGasDecorator(ak, sk)
		case evmtypes.MsgEthereumTx:
			gasHandler = EthGasDecorator(ak, evmKeeper, sk)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}
		return gasHandler(ctx, tx, sim)
	}
}

type EthGasHandler struct {
	ak        auth.AccountKeeper
	sk        types.SupplyKeeper
	evmKeeper EVMKeeper
}

func (egh EthGasHandler) GasHandle(ctx sdk.Context, tx sdk.Tx, sim bool) (err error) {

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

	evmDenom := egh.evmKeeper.GetParams(ctx).EvmDenom
	feeAmt := sdk.NewCoins(
		sdk.NewCoin(evmDenom, sdk.NewDecFromBigIntWithPrec(gasRefund, sdk.Precision)),
	)

	err = ante.RefundFees(egh.sk, ctx, address, feeAmt)
	if err != nil {
		return err
	}

	return nil
}

func EthGasDecorator(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper) sdk.GasHandler {

	egh := EthGasHandler{
		ak:        ak,
		sk:        sk,
		evmKeeper: evmKeeper,
	}

	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (err error) {
		return egh.GasHandle(ctx, tx, simulate)
	}
}

type CosmosGasHandler struct {
	ak           keeper.AccountKeeper
	supplyKeeper types.SupplyKeeper
}

func (cgh CosmosGasHandler) GasHandle(ctx sdk.Context, tx sdk.Tx, sim bool) (err error) {

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

	feeTx, ok := tx.(ante.FeeTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	feePayerAcc := cgh.ak.GetAccount(ctx, feePayer)
	if feePayerAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	gas := feeTx.GetGas()
	fees := feeTx.GetFee()
	gasFees := make(sdk.Coins, len(fees))

	for i, fee := range fees {
		gasPrice := new(big.Int).Div(fee.Amount.BigInt(), new(big.Int).SetUint64(gas))
		gasConsumed := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasUsed))
		gasCost := sdk.NewCoin(fee.Denom, sdk.NewDecFromBigIntWithPrec(gasConsumed, sdk.Precision))
		gasRefund := fee.Sub(gasCost)

		gasFees[i] = gasRefund
	}

	err = ante.RefundFees(cgh.supplyKeeper, ctx, feePayerAcc.GetAddress(), gasFees)
	if err != nil {
		return err
	}

	return nil
}

func CosmosGasDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.GasHandler {
	cgh := CosmosGasHandler{
		ak:           ak,
		supplyKeeper: sk,
	}

	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (err error) {
		return cgh.GasHandle(ctx, tx, simulate)
	}
}
