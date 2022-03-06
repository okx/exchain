package ante

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

type FastAnteHandler func(ctx sdk.Context, ethTx evmtypes.MsgEthereumTx, simulate bool) (newCtx sdk.Context, err error)

type FastETHExecuteDecorator struct {
	ak        auth.AccountKeeper
	sk        types.SupplyKeeper
	evmKeeper EVMKeeper
}

func NewFastETHExecuteDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper, ek EVMKeeper) FastETHExecuteDecorator {
	return FastETHExecuteDecorator{
		ak:        ak,
		sk:        sk,
		evmKeeper: ek,
	}
}

func (feed FastETHExecuteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	trc := ctx.AnteTracer()
	if trc != nil {
		trc.RepeatingPin("FastETHExecuteDecorator")
	}

	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				log := fmt.Sprintf(
					"out of gas in location: %v; gasLimit: %d, gasUsed: %d",
					rType.Descriptor, tx.GetGas(), ctx.GasMeter().GasConsumed(),
				)
				err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, log)
			default:
				panic(r)
			}
		}
	}()

	// simulate means 'eth_call' or 'eth_estimateGas', when it means 'eth_estimateGas' we can not 'VerifySig'.so skip here
	if simulate || ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	if err := tx.ValidateBasic(); err != nil {
		return ctx, err
	}
	return next(ctx, tx, simulate)
}
