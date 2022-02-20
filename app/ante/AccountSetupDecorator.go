package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/keeper"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"time"
)

// AccountSetupDecorator sets an account to state if it's not stored already. This only applies for MsgEthermint.
type AccountSetupDecorator struct {
	ak auth.AccountKeeper
}

// NewAccountSetupDecorator creates a new AccountSetupDecorator instance
func NewAccountSetupDecorator(ak auth.AccountKeeper) AccountSetupDecorator {
	return AccountSetupDecorator{
		ak: ak,
	}
}

// AnteHandle sets an account for MsgEthermint (evm) if the sender is registered.
// NOTE: Since the account is set without any funds, the message execution will
// fail if the validator requires a minimum fee > 0.
func (asd AccountSetupDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	pinAnte(ctx.AnteTracer(), "9-AccountSetupDecorator")
	time.Sleep(900*time.Millisecond)

	msgs := tx.GetMsgs()

	time.Sleep(3*time.Second)


	if len(msgs) == 0 {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "no messages included in transaction")
	}

	for _, msg := range msgs {
		if msgEthermint, ok := msg.(evmtypes.MsgEthermint); ok {
			setupAccount(asd.ak, ctx, msgEthermint.From)
		}
	}

	return next(ctx, tx, simulate)
}

func setupAccount(ak keeper.AccountKeeper, ctx sdk.Context, addr sdk.AccAddress) {
	acc := ak.GetAccount(ctx, addr)
	if acc != nil {
		return
	}

	acc = ak.NewAccountWithAddress(ctx, addr)
	ak.SetAccount(ctx, acc)
}