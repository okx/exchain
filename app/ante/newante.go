package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"sync"

	evmtypes "github.com/okex/exchain/x/evm/types"
)

var logger anteLogger
var loggerOnce sync.Once
func SetLogger(l log.Logger) {
	loggerOnce.Do(func() {
		logger.Logger = l.With("module", "ante")
	})
}

type anteLogger struct {
	log.Logger
}

func (l anteLogger) Info(msg string, keyvals ...interface{}) {
	if l.Logger == nil {
		return
	}
	l.Logger.Info(msg, keyvals...)
}

// NewAnteHandler returns an ante handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal ante handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewAnteHandler4Wtx(ak auth.AccountKeeper, evmKeeper EVMKeeper,
	sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		stdTxAnteHandler := sdk.ChainAnteDecorators(
			authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
			NewAccountSetupDecorator(ak),
			NewAccountBlockedVerificationDecorator(evmKeeper), //account blocked check AnteDecorator
			authante.NewMempoolFeeDecorator(),
			authante.NewValidateBasicDecorator(),
			authante.NewValidateMemoDecorator(ak),
			authante.NewConsumeGasForTxSizeDecorator(ak),
			authante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
			authante.NewValidateSigCountDecorator(ak),
			authante.NewDeductFeeDecorator(ak, sk),
			authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
			authante.NewSigVerificationDecorator(ak),
			authante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
			NewValidateMsgHandlerDecorator(validateMsgHandler),
		)

		evmTxAnteHandler := sdk.ChainAnteDecorators(
			NewEthSetupContextDecorator(), // outermost AnteDecorator. EthSetUpContext must be called first
			NewGasLimitDecorator(evmKeeper),
			NewEthMempoolFeeDecorator(evmKeeper),
			authante.NewValidateBasicDecorator(),
			NewEthSigVerificationDecorator(),
			NewAccountBlockedVerificationDecorator(evmKeeper), //account blocked check AnteDecorator
			NewAccountVerificationDecorator(ak, evmKeeper),
			NewNonceVerificationDecorator(ak),
			NewEthGasConsumeDecorator(ak, sk, evmKeeper),
			NewIncrementSenderSequenceDecorator(ak), // innermost AnteDecorator.
		)

		switch txType := tx.(type) {
		case auth.StdTx:
			logger.Info("ante auth.StdTx")
			anteHandler = stdTxAnteHandler
		case evmtypes.MsgEthereumTx:
			logger.Info("ante MsgEthereumTx")

			anteHandler = evmTxAnteHandler
		case auth.WrappedTx:
			logger.Info("ante auth.WrappedTx")
			anteHandler = func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
				return checkTxAnteHandler(ctx, tx, sim, txType.Tx, stdTxAnteHandler, evmTxAnteHandler)
			}
		default:
			logger.Info("invalid transaction type: %T", tx)
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}
}

func checkTxAnteHandler(ctx sdk.Context, tx sdk.Tx, sim bool, payloadTx sdk.Tx, stdTxAnteHandler, evmTxAnteHandler sdk.AnteHandler) (newCtx sdk.Context, err error) {

	var payloadAnteHandler sdk.AnteHandler
	logger.Info("ante checkTxAnteHandler")

	switch payloadTx.(type) {
	case auth.StdTx:
		payloadAnteHandler = stdTxAnteHandler
	case evmtypes.MsgEthereumTx:
		payloadAnteHandler = evmTxAnteHandler
	default:
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid payload transaction type: %T", payloadTx)
	}

	chkTxAnteHandler := sdk.ChainAnteDecorators(
		authante.NewNodeSignatureDecorator(logger.Logger),
	)

	newCtx, err = chkTxAnteHandler(ctx, tx, sim)
	if err != nil {
		logger.Info("Wrapped tx anteHandler failed", "err", err)
		newCtx, err = payloadAnteHandler(newCtx, payloadTx, sim)
		logger.Info("Payload tx anteHandler", "err", err)
	}

	return newCtx, err
}

