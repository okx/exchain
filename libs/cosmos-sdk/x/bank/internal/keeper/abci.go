package keeper

import (
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/innertx"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, ik innertx.InnerTxKeeper) {
	currentHash := req.Hash
	if ik != nil {
		ik.InitInnerBlock(common.BytesToHash(currentHash).Hex())
	}
}
