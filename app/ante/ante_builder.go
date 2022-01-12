package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// nolint
func buildOriginStdtxAnteHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
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
}

func buildOriginEvmTxAnteHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
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
}

func buildLightStdtxAnteHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		authante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewAccountSetupDecorator(ak),
		authante.NewValidateBasicDecorator(),
		authante.NewConsumeGasForTxSizeDecorator(ak),
		authante.NewDeductFeeDecorator(ak, sk),
		authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		authante.NewSigVerificationDecorator(ak),
		authante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}

func buildLightEvmTxAnteHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
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
}
