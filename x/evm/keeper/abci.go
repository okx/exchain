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
	if req.Header.GetHeight() == tmtypes.GetMarsHeight() {
		migrateDataInMarsHeight(ctx, k)
	}

	if req.Header.LastBlockId.GetHash() == nil || req.Header.GetHeight() < 1 {
		return
	}

	// Gas costs are handled within msg handler so costs should be ignored
	ctx.SetGasMeter(sdk.NewInfiniteGasMeter())

	// Set the hash -> height and height -> hash mapping.
	currentHash := req.Hash
	lastHash := req.Header.LastBlockId.GetHash()
	height := req.Header.GetHeight() - 1

	blockHash := common.BytesToHash(currentHash)
	k.SetHeightHash(ctx, uint64(height), common.BytesToHash(lastHash))
	k.SetBlockHash(ctx, lastHash, height)

	// reset counters that are used on CommitStateDB.Prepare
	if !ctx.IsTraceTx() {
		k.Bloom = big.NewInt(0)
		k.TxCount = 0
		k.LogSize = 0
		k.LogsManages.Reset()
		k.Bhash = blockHash

		//that can make sure latest block has been committed
		k.UpdatedAccount = k.UpdatedAccount[:0]
		k.EvmStateDb = types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
		k.EvmStateDb.StartPrefetcher("evm")
		k.Watcher.NewHeight(uint64(req.Header.GetHeight()), blockHash, req.Header)
	}
}

// EndBlock updates the accounts and commits state objects to the KV Store, while
// deleting the empty ones. It also sets the bloom filers for the request block to
// the store. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func (k *Keeper) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	// Gas costs are handled within msg handler so costs should be ignored
	ctx.SetGasMeter(sdk.NewInfiniteGasMeter())

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
	}

	k.UpdateInnerBlockData()

	k.Commit(ctx)

	return []abci.ValidatorUpdate{}
}

// migrateDataInMarsHeight migrates data from evm store to param store
// 0. chainConfig
// 1. white address list
// 2. blocked addresses list
func migrateDataInMarsHeight(ctx sdk.Context, k *Keeper) {
	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	newStore := k.paramSpace.CustomKVStore(ctx)

	// 0. migrate chainConfig
	config, _ := k.GetChainConfig(ctx)
	newStore.Set(types.KeyPrefixChainConfig, k.cdc.MustMarshalBinaryBare(config))

	// 1、migrate white list
	whiteList := csdb.GetContractDeploymentWhitelist()
	for i := 0; i < len(whiteList); i++ {
		newStore.Set(types.GetContractDeploymentWhitelistMemberKey(whiteList[i]), []byte(""))
	}

	// 2.1、deploy blocked list
	blockedList := csdb.GetContractBlockedList()
	for i := 0; i < len(blockedList); i++ {
		newStore.Set(types.GetContractBlockedListMemberKey(blockedList[i]), []byte(""))
	}

	// 2.2、migrate blocked method list
	methodBlockedList := csdb.GetContractMethodBlockedList()
	for i := 0; i < len(methodBlockedList); i++ {
		if !methodBlockedList[i].IsAllMethodBlocked() {
			types.SortContractMethods(methodBlockedList[i].BlockMethods)
			value := k.cdc.MustMarshalJSON(methodBlockedList[i].BlockMethods)
			sortedValue := sdk.MustSortJSON(value)
			newStore.Set(types.GetContractBlockedListMemberKey(methodBlockedList[i].Address), sortedValue)
		}
	}
}
