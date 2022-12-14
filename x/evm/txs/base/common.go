package base

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types"
)

var commitStateDBPool = &sync.Pool{
	New: func() interface{} {
		return &types.CommitStateDB{}
	},
}

func getCommitStateDB() *types.CommitStateDB {
	return commitStateDBPool.Get().(*types.CommitStateDB)
}

func putCommitStateDB(st *types.CommitStateDB) {
	commitStateDBPool.Put(st)
}

func msg2st(ctx *sdk.Context, k *Keeper, msg *types.MsgEthereumTx, st *types.StateTransition) (reuseCsdb bool, err error) {
	var chainIDEpoch *big.Int
	chainIDEpoch, err = ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return
	}

	var sender common.Address
	// Verify signature and retrieve sender address
	sender, err = getSender(ctx, chainIDEpoch, msg)
	if err != nil {
		return
	}

	txHash := msg.TxHash()
	if len(txHash) == 0 || isVenus(ctx.BlockHeight()) {
		txHash = tmtypes.Tx(ctx.TxBytes()).Hash(ctx.BlockHeight())
	}
	ethHash := common.BytesToHash(txHash)

	st.AccountNonce = msg.Data.AccountNonce
	st.Price = msg.Data.Price
	st.GasLimit = msg.Data.GasLimit
	st.Recipient = msg.Data.Recipient
	st.Amount = msg.Data.Amount
	st.Payload = msg.Data.Payload
	st.ChainID = chainIDEpoch
	st.TxHash = &ethHash
	st.Sender = sender
	st.Simulate = ctx.IsCheckTx()
	st.TraceTx = ctx.IsTraceTx()
	st.TraceTxLog = ctx.IsTraceTxLog()

	if tmtypes.HigherThanMars(ctx.BlockHeight()) && ctx.IsDeliver() {
		st.Csdb = k.EvmStateDb.WithContext(*ctx)
	} else {
		csdb := getCommitStateDB()
		types.ResetCommitStateDB(csdb, k.GenerateCSDBParams(), ctx)
		st.Csdb = csdb
		reuseCsdb = true
	}

	return
}

func getSender(ctx *sdk.Context, chainIDEpoch *big.Int, msg *types.MsgEthereumTx) (sender common.Address, err error) {
	if ctx.IsCheckTx() {
		if from := ctx.From(); len(from) > 0 {
			return common.HexToAddress(from), nil
		}
	}
	err = msg.VerifySig(chainIDEpoch, ctx.BlockHeight())
	if err == nil {
		sender = msg.EthereumAddress()
	}

	return
}

func isVenus(height int64) bool {
	if height <= 0 {
		return false
	}
	lastHeight := height - 1
	if tmtypes.HigherThanVenus(height) && !tmtypes.HigherThanVenus(lastHeight) {
		return true
	}
	return false
}
