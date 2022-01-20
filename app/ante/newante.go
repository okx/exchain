package ante

//import (
//	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
//	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
//	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
//	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
//	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
//	"github.com/okex/exchain/libs/tendermint/libs/log"
//	"sync"
//)
//
//var logger anteLogger
//var loggerOnce sync.Once
//func SetLogger(l log.Logger) {
//	loggerOnce.Do(func() {
//		logger.Logger = l.With("module", "ante")
//	})
//}
//
//type anteLogger struct {
//	log.Logger
//}
//
//func (l anteLogger) Info(msg string, keyvals ...interface{}) {
//	if l.Logger == nil {
//		return
//	}
//	l.Logger.Info(msg, keyvals...)
//}
//
//// NewAnteHandler returns an ante handler responsible for attempting to route an
//// Ethereum or SDK transaction to an internal ante handler for performing
//// transaction-level processing (e.g. fee payment, signature verification) before
//// being passed onto it's respective handler.
//func NewAnteHandler4Wtx(ak auth.AccountKeeper, evmKeeper EVMKeeper,
//	sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
//	return func(
//		ctx sdk.Context, tx sdk.Tx, sim bool,
//	) (newCtx sdk.Context, err error) {
//		var anteHandler sdk.AnteHandler
//
//		stdTxAnteHandler := sdk.ChainAnteDecorators(
//			authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
//			NewAccountSetupDecorator(ak),
//			NewAccountBlockedVerificationDecorator(evmKeeper), //account blocked check AnteDecorator
//			authante.NewMempoolFeeDecorator(),
//			authante.NewValidateBasicDecorator(),
//			authante.NewValidateMemoDecorator(ak),
//			authante.NewConsumeGasForTxSizeDecorator(ak),
//			authante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
//			authante.NewValidateSigCountDecorator(ak),
//			authante.NewDeductFeeDecorator(ak, sk),
//			authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
//			authante.NewSigVerificationDecorator(ak),
//			authante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
//			NewValidateMsgHandlerDecorator(validateMsgHandler),
//		)
//
//		evmTxAnteHandler := sdk.ChainAnteDecorators(
//			NewEthSetupContextDecorator(), // outermost AnteDecorator. EthSetUpContext must be called first
//			NewGasLimitDecorator(evmKeeper),
//			NewEthMempoolFeeDecorator(evmKeeper),
//			authante.NewValidateBasicDecorator(),
//			NewEthSigVerificationDecorator(),
//			NewAccountBlockedVerificationDecorator(evmKeeper), //account blocked check AnteDecorator
//			NewAccountVerificationDecorator(ak, evmKeeper),
//			NewNonceVerificationDecorator(ak),
//			NewEthGasConsumeDecorator(ak, sk, evmKeeper),
//			NewIncrementSenderSequenceDecorator(ak), // innermost AnteDecorator.
//		)
//
//		switch tx.GetType() {
//		case sdk.StdTxType:
//			logger.Info("ante StdTx")
//			anteHandler = stdTxAnteHandler
//		case sdk.EvmTxType:
//			logger.Info("ante MsgEthereumTx")
//			anteHandler = evmTxAnteHandler
//		case sdk.WrappedTxType:
//			logger.Info("ante WrappedTx")
//			anteHandler = func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
//				return wrappedTxAnteHandler(ctx, tx, sim, tx.GetPayloadTx(), stdTxAnteHandler, evmTxAnteHandler)
//			}
//		default:
//			logger.Info("invalid transaction type: %T", tx)
//			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
//		}
//
//		return anteHandler(ctx, tx, sim)
//	}
//}
//
//func wrappedTxAnteHandler(ctx sdk.Context, wtx sdk.Tx, sim bool,
//	payloadTx sdk.Tx, stdTxAnteHandler, evmTxAnteHandler sdk.AnteHandler) (newCtx sdk.Context, err error) {
//
//	// 1. try wrapped tx AnteHandler
//	wtxAnteHandler := sdk.ChainAnteDecorators(
//		authante.NewNodeSignatureDecorator(logger.Logger),
//	)
//
//	newCtx, err = wtxAnteHandler(ctx, wtx, sim)
//	if err == nil {
//		return
//	}
//
//	logger.Info("Wrapped tx anteHandler failed", "err", err)
//
//	// 2. try payload tx AnteHandler
//	var payloadAnteHandler sdk.AnteHandler
//	switch payloadTx.GetType() {
//	case sdk.StdTxType:
//		payloadAnteHandler = stdTxAnteHandler
//	case sdk.EvmTxType:
//		payloadAnteHandler = evmTxAnteHandler
//	default:
//		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid payload transaction type: %T", payloadTx)
//	}
//	newCtx, err = payloadAnteHandler(newCtx, payloadTx, sim)
//	logger.Info("Payload tx anteHandler", "err", err)
//
//	return
//}
