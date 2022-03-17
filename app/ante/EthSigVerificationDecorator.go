package ante

import (
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// EthSigVerificationDecorator validates an ethereum signature
type EthSigVerificationDecorator struct{}

// NewEthSigVerificationDecorator creates a new EthSigVerificationDecorator
func NewEthSigVerificationDecorator() EthSigVerificationDecorator {
	return EthSigVerificationDecorator{}
}

// AnteHandle validates the signature and returns sender address
func (esvd EthSigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// simulate means 'eth_call' or 'eth_estimateGas', when it means 'eth_estimateGas' we can not 'VerifySig'.so skip here
	if simulate {
		return next(ctx, tx, simulate)
	}
	pinAnte(ctx.AnteTracer(), "EthSigVerificationDecorator")

	msgEthTx, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}

	// parse the chainID from a string to a base-10 integer
	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return ctx, err
	}

	// validate sender/signature and cache the address
	signerSigCache, err := msgEthTx.VerifySig(chainIDEpoch, ctx.BlockHeight(), ctx.TxBytes(), ctx.SigCache())
	if err != nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "signature verification failed: %s", err.Error())
	}

	// update ctx for push signerSigCache
	newCtx = ctx.WithSigCache(signerSigCache)
	// NOTE: when signature verification succeeds, a non-empty signer address can be
	// retrieved from the transaction on the next AnteDecorators.
	return next(newCtx, msgEthTx, simulate)
}