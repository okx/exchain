package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// build the origin tx ante handlers
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

// build the origin evm tx ante handlers
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

// when at the wrapped tx mode
// should use this function to build a light ante handlers chain
// only check the account nonce and transaction gas
func buildLightStdtxAnteHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		authante.NewSetUpContextDecorator(),
		NewAccountSetupDecorator(ak),
		authante.NewSetPubKeyDecorator(ak),
		authante.NewValidateBasicDecorator(),
		authante.NewSigVerificationDecorator(ak),
		authante.NewIncrementSequenceDecorator(ak),
	)
}

// when at the wrapped tx mode
// should use this function to build a light ante handlers chain
// only check the account nonce and transaction gas
func buildLightEvmTxAnteHandler(ak auth.AccountKeeper, evmKeeper EVMKeeper, sk types.SupplyKeeper, validateMsgHandler ValidateMsgHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewEthSetupContextDecorator(),
		NewEthSigVerificationDecorator(),
		NewAccountVerificationDecorator(ak, evmKeeper),
		NewNonceVerificationDecorator(ak),
		NewIncrementSenderSequenceDecorator(ak),
	)
}
