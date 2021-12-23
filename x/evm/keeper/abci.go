package keeper

import (
	"math/big"

	"github.com/okex/exchain/x/evm/watcher"

	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"github.com/ethereum/go-ethereum/common"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/x/evm/types"
)

// BeginBlock sets the block hash -> block height map for the previous block height
// and resets the Bloom filter and the transaction count to 0.
func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	if req.Header.LastBlockId.GetHash() == nil || req.Header.GetHeight() < 1 {
		return
	}

	// Gas costs are handled within msg handler so costs should be ignored
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	// Set the hash -> height and height -> hash mapping.
	currentHash := req.Hash
	lastHash := req.Header.LastBlockId.GetHash()
	height := req.Header.GetHeight() - 1

	k.SetHeightHash(ctx, uint64(height), common.BytesToHash(lastHash))
	k.SetBlockHash(ctx, lastHash, height)
	k.InitInnerBlock(common.BytesToHash(currentHash).Hex())

	// reset counters that are used on CommitStateDB.Prepare
	k.Bloom = big.NewInt(0)
	k.TxCount = 0
	k.LogSize = 0
	k.LogsManages = NewLogManager()
	k.Bhash = common.BytesToHash(currentHash)

	//that can make sure latest block has been committed
	k.Watcher.NewHeight(uint64(req.Header.GetHeight()), common.BytesToHash(currentHash), req.Header)
	k.Watcher.ExecuteDelayEraseKey()

	k.UpdatedAccount = k.UpdatedAccount[:0]
	k.EvmStateDb = types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
}

// EndBlock updates the accounts and commits state objects to the KV Store, while
// deleting the empty ones. It also sets the bloom filers for the request block to
// the store. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func (k Keeper) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	// Gas costs are handled within msg handler so costs should be ignored
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	// set the block bloom filter bytes to store
	bloom := ethtypes.BytesToBloom(k.Bloom.Bytes())
	k.SetBlockBloom(ctx, req.Height, bloom)

	if types.GetEnableBloomFilter() {
		// the hash of current block is stored when executing BeginBlock of next block.
		// so update section in the next block.
		if indexer := types.GetIndexer(); indexer != nil {
			if types.GetIndexer().IsProcessing() {
				// notify new height
				go func() {
					indexer.NotifyNewHeight(ctx)
				}()
			} else {
				interval := uint64(req.Height - tmtypes.GetStartBlockHeight())
				if interval >= (indexer.GetValidSections()+1)*types.BloomBitsBlocks {
					go types.GetIndexer().ProcessSection(ctx, k, interval, k.Watcher.GetBloomDataPoint())
				}
			}
		}
	}

	if watcher.IsWatcherEnabled() && k.Watcher.IsFirstUse() {
		store := ctx.KVStore(k.storeKey)
		iteratorBlockedList := sdk.KVStorePrefixIterator(store, types.KeyPrefixContractBlockedList)
		defer iteratorBlockedList.Close()
		for ; iteratorBlockedList.Valid(); iteratorBlockedList.Next() {
			vaule := iteratorBlockedList.Value()
			if len(vaule) == 0 {
				k.Watcher.SaveContractBlockedListItem(iteratorBlockedList.Key()[1:])
			} else {
				k.Watcher.SaveContractMethodBlockedListItem(iteratorBlockedList.Key()[1:], vaule)
			}
		}

		iteratorDeploymentWhitelist := sdk.KVStorePrefixIterator(store, types.KeyPrefixContractDeploymentWhitelist)
		defer iteratorDeploymentWhitelist.Close()
		for ; iteratorDeploymentWhitelist.Valid(); iteratorDeploymentWhitelist.Next() {
			k.Watcher.SaveContractDeploymentWhitelistItem(iteratorDeploymentWhitelist.Key()[1:])
		}

		k.Watcher.Used()
	}

	if watcher.IsWatcherEnabled() {
		params := k.GetParams(ctx)
		k.Watcher.SaveParams(params)

		k.Watcher.SaveBlock(bloom)
		k.Watcher.Commit()
	}

	k.UpdateInnerBlockData()

	k.EvmStateDb.WithContext(ctx).Commit(true)

	return []abci.ValidatorUpdate{}
}
