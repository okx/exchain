package ante

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/innertx"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

// EVMKeeper defines the expected keeper interface used on the Eth AnteHandler
type EVMKeeper interface {
	innertx.InnerTxKeeper
	GetParams(ctx sdk.Context) evmtypes.Params
	IsAddressBlocked(ctx sdk.Context, addr sdk.AccAddress) bool
}

// NewGasLimitDecorator creates a new GasLimitDecorator.
func NewGasLimitDecorator(evm EVMKeeper) GasLimitDecorator {
	return GasLimitDecorator{
		evm: evm,
	}
}

type GasLimitDecorator struct {
	evm EVMKeeper
}

func (g GasLimitDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	pinAnte(ctx.AnteTracer(), "GasLimitDecorator")

	currentGasMeter := ctx.GasMeter() // avoid race
	infGasMeter := sdk.GetReusableInfiniteGasMeter()
	ctx.SetGasMeter(infGasMeter)
	if tx.GetGas() > g.evm.GetParams(ctx).MaxGasLimitPerTx {
		ctx.SetGasMeter(currentGasMeter)
		sdk.ReturnInfiniteGasMeter(infGasMeter)
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrTxTooLarge, "too large gas limit, it must be less than %d", g.evm.GetParams(ctx).MaxGasLimitPerTx)
	}

	ctx.SetGasMeter(currentGasMeter)
	sdk.ReturnInfiniteGasMeter(infGasMeter)
	return next(ctx, tx, simulate)
}
