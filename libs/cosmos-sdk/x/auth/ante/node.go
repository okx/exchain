package ante

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)


type NodeSignatureDecorator struct{
	logger log.Logger
}

func NewNodeSignatureDecorator(l log.Logger) NodeSignatureDecorator {
	return NodeSignatureDecorator{
		logger: l,
	}
}

func (n NodeSignatureDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// -1 for failure, 1 for success
	ctx = ctx.WithNodeSigVerifyResult(-1)
	var wtx authtypes.WrappedTx
	var ok bool
	if wtx, ok = tx.(authtypes.WrappedTx); !ok {
		return ctx, fmt.Errorf("Invalid WrappedTx")
	}

	inWhitelist := false
	if ctx.NodekeyWhitelist != nil {
		for _, v := range ctx.NodekeyWhitelist {
			if bytes.Compare(v, wtx.NodeKey) == 0 {
				inWhitelist = true
				break
			}
		}
	}

	if !inWhitelist {
		return ctx, fmt.Errorf("The pubkey of wtx is not in the node key whitelist")
	}

	var pubKey ed25519.PubKeyEd25519
	err = pubKey.UnmarshalFromAmino(wtx.NodeKey)
	if err != nil {
		n.logger.Info("Failed to recover node key", "err", err)
		return ctx, err
	}

	if !pubKey.VerifyBytes(wtx.GetPayloadTxBytes(), wtx.Signature) {
		n.logger.Info("Failed to verify payload tx",
			"pubkey", hexutil.Encode(wtx.NodeKey),
			)
		return ctx, err
	}

	ctx = ctx.WithNodeSigVerifyResult(1)

	defer n.logger.Info("NodeSignatureDecorator anteHandle done",
		"pubkey", hexutil.Encode(wtx.NodeKey),
		"tx-type", tx.GetType(),
		)

	return next(ctx, tx, simulate)
}
