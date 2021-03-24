package gas

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	ethermint "github.com/okex/okexchain/app/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
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

	// fetch sender account from signature
	senderAcc, err := auth.GetSignerAcc(ctx, egh.ak, address)
	if err != nil {
		return err
	}

	if senderAcc == nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"sender account %s (%s) is nil", common.BytesToAddress(address.Bytes()), address,
		)
	}

	gasUsed := currentGasMeter.GasConsumed()
	gasCost := new(big.Int).Mul(msgEthTx.Data.Price, new(big.Int).SetUint64(gasUsed))
	evmDenom := egh.evmKeeper.GetParams(ctx).EvmDenom
	feeAmt := sdk.NewCoins(
		sdk.NewCoin(evmDenom, sdk.NewDecFromBigIntWithPrec(gasCost, sdk.Precision)),
	)

	err = auth.DeductFees(egh.sk, ctx, senderAcc, feeAmt, false)
	if err != nil {
		return err
	}

	ctx = ctx.WithGasMeter(currentGasMeter)

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
	gasUsed := currentGasMeter.GasConsumed()

	for i, fee := range fees {
		gasPrice := new(big.Int).Div(fee.Amount.BigInt(), new(big.Int).SetUint64(gas))
		gasConsumed := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasUsed))
		gasFee := sdk.NewCoin(fee.Denom, sdk.NewDecFromBigIntWithPrec(gasConsumed, sdk.Precision))
		gasFees[i] = gasFee
	}

	err = auth.DeductFees(cgh.supplyKeeper, ctx, feePayerAcc, gasFees, false)
	if err != nil {
		return err
	}

	ctx = ctx.WithGasMeter(currentGasMeter)

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
