package keeper

import (
	"math/big"
	"sync"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"

	app "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/okx/okbchain/x/evm/types"
	"github.com/okx/okbchain/x/evm/watcher"
	"github.com/okx/okbchain/x/params"
)

const (
	heightCacheLimit = 1024
	hashCacheLimit   = 1024
)

// Keeper wraps the CommitStateDB, allowing us to pass in SDK context while adhering
// to the StateDB interface.
type Keeper struct {
	// Amino codec
	cdc *codec.Codec
	// Store key required for the EVM Prefix KVStore. It is required by:
	// - storing Account's Storage State
	// - storing Account's Code
	// - storing transaction Logs
	// - storing block height -> bloom filter map. Needed for the Web3 API.
	// - storing block hash -> block height map. Needed for the Web3 API.
	storeKey sdk.StoreKey

	// Account Keeper for fetching accounts
	accountKeeper types.AccountKeeper
	paramSpace    types.Subspace
	supplyKeeper  types.SupplyKeeper
	bankKeeper    types.BankKeeper
	govKeeper     GovKeeper
	stakingKeeper types.StakingKeeper

	// Transaction counter in a block. Used on StateSB's Prepare function.
	// It is reset to 0 every block on BeginBlock so there's no point in storing the counter
	// on the KVStore or adding it as a field on the EVM genesis state.
	TxCount int
	Bloom   *big.Int
	Bhash   ethcmn.Hash
	LogSize uint
	Ada     types.DbAdapter

	LogsManages *LogsManager

	// add inner block data
	innerBlockData BlockInnerData

	EvmStateDb     *types.CommitStateDB
	UpdatedAccount []ethcmn.Address

	db ethstate.Database

	startHeight uint64
	triegc      *prque.Prque

	// cache chain config
	cci *chainConfigInfo

	hooks   types.EvmHooks
	logger  log.Logger
	Watcher *watcher.Watcher

	heightCache *lru.Cache // Cache for the most recent block heights
	hashCache   *lru.Cache // Cache for the most recent block hash
}

type chainConfigInfo struct {
	// chainConfig cached chain config
	// nil means invalid the cache, we should cache it again.
	cc *types.ChainConfig

	// gasReduced: cached chain config reduces gas costs.
	// when use cached chain config, we restore the gas cost(gasReduced)
	gasReduced sdk.Gas
}

// NewKeeper generates new evm module keeper
func NewKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace, ak types.AccountKeeper, sk types.SupplyKeeper, bk types.BankKeeper, stk types.StakingKeeper,
	logger log.Logger) *Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	err := initInnerDB()
	if err != nil {
		panic(err)
	}

	if enable := types.GetEnableBloomFilter(); enable {
		db := types.BloomDb()
		types.InitIndexer(db)
	}
	logger = logger.With("module", types.ModuleName)
	heightCache, _ := lru.New(heightCacheLimit)
	hashCache, _ := lru.New(hashCacheLimit)
	// NOTE: we pass in the parameter space to the CommitStateDB in order to use custom denominations for the EVM operations
	k := &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		accountKeeper: ak,
		paramSpace:    paramSpace,
		supplyKeeper:  sk,
		bankKeeper:    bk,
		stakingKeeper: stk,
		TxCount:       0,
		Bloom:         big.NewInt(0),
		LogSize:       0,
		Ada:           types.DefaultPrefixDb{},

		innerBlockData: defaultBlockInnerData(),

		db:             mpt.InstanceOfMptStore(),
		triegc:         prque.New(nil),
		UpdatedAccount: make([]ethcmn.Address, 0),
		cci:            &chainConfigInfo{},
		LogsManages:    NewLogManager(),
		logger:         logger,
		Watcher:        watcher.NewWatcher(logger),
		heightCache:    heightCache,
		hashCache:      hashCache,
	}
	k.Watcher.SetWatchDataManager()
	ak.SetObserverKeeper(k)

	//k.OpenTrie()
	k.EvmStateDb = types.NewCommitStateDB(k.GenerateCSDBParams())
	return k
}

// NewKeeper generates new evm module keeper
func NewSimulateKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace types.Subspace, ak types.AccountKeeper, sk types.SupplyKeeper, bk types.BankKeeper, stk types.StakingKeeper, ada types.DbAdapter,
	logger log.Logger) *Keeper {
	heightCache, _ := lru.New(heightCacheLimit)
	hashCache, _ := lru.New(hashCacheLimit)
	// NOTE: we pass in the parameter space to the CommitStateDB in order to use custom denominations for the EVM operations
	k := &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		accountKeeper:  ak,
		paramSpace:     paramSpace,
		supplyKeeper:   sk,
		bankKeeper:     bk,
		stakingKeeper:  stk,
		TxCount:        0,
		Bloom:          big.NewInt(0),
		LogSize:        0,
		Watcher:        watcher.NewWatcher(nil),
		Ada:            ada,
		UpdatedAccount: make([]ethcmn.Address, 0),

		db: mpt.InstanceOfMptStore(),
		// Optimize memory usage. No need to initialize this variable when simulate tx.
		// triegc:         prque.New(nil),

		cci:         &chainConfigInfo{},
		heightCache: heightCache,
		hashCache:   hashCache,
	}

	//k.OpenTrie()
	k.EvmStateDb = types.NewCommitStateDB(k.GenerateCSDBParams())

	return k
}

// Warning, you need to use pointer object here, for you need to update UpdatedAccount var
func (k *Keeper) OnAccountUpdated(acc auth.Account) {
	if _, ok := acc.(*app.EthAccount); ok {
		k.Watcher.DeleteAccount(acc.GetAddress())
	}

	k.UpdatedAccount = append(k.UpdatedAccount, ethcmn.BytesToAddress(acc.GetAddress().Bytes()))
}

// Logger returns a module-specific logger.
func (k *Keeper) GenerateCSDBParams() types.CommitStateDBParams {
	return types.CommitStateDBParams{
		StoreKey:      k.storeKey,
		ParamSpace:    k.paramSpace,
		AccountKeeper: k.accountKeeper,
		SupplyKeeper:  k.supplyKeeper,
		BankKeeper:    k.bankKeeper,
		Ada:           k.Ada,
		Cdc:           k.cdc,
		DB:            k.db,
	}
}

// GeneratePureCSDBParams generates an instance of csdb params ONLY for store setter and getter
func (k Keeper) GeneratePureCSDBParams() types.CommitStateDBParams {
	return types.CommitStateDBParams{
		StoreKey:   k.storeKey,
		ParamSpace: k.paramSpace,
		Ada:        k.Ada,
		Cdc:        k.cdc,

		DB: k.db,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger
}

func (k Keeper) GetStoreKey() store.StoreKey {
	return k.storeKey
}

// ----------------------------------------------------------------------------
// Block hash mapping functions
// Required by Web3 API.
//  TODO: remove once tendermint support block queries by hash.
// ----------------------------------------------------------------------------

// GetBlockHeight gets block height from block consensus hash
func (k Keeper) GetBlockHeight(ctx sdk.Context, hash ethcmn.Hash) (int64, bool) {
	if cached, ok := k.heightCache.Get(hash.Hex()); ok {
		height := cached.(int64)
		return height, true
	}
	return k.getBlockHashInDiskDB(hash.Bytes())
}

// SetBlockHeight sets the mapping from block consensus hash to block height
func (k Keeper) SetBlockHeight(ctx sdk.Context, hash []byte, height int64) {
	k.setBlockHashInDiskDB(hash, height)
}

// IterateBlockHash iterates all over the block hash in every height
func (k Keeper) IterateBlockHash(ctx sdk.Context, fn func(key []byte, value []byte) (stop bool)) {
	k.iterateBlockHashInDiskDB(fn)
}

// ----------------------------------------------------------------------------
// Epoch Height -> hash mapping functions
// Required by EVM context's GetHashFunc
// ----------------------------------------------------------------------------

// GetHeightHash returns the block header hash associated with a given block height and chain epoch number.
func (k Keeper) GetHeightHash(ctx sdk.Context, height uint64) ethcmn.Hash {
	if cached, ok := k.hashCache.Get(int64(height)); ok {
		hash := cached.(string)
		return ethcmn.HexToHash(hash)
	}
	return types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).GetHeightHash(height)
}

// SetHeightHash sets the block header hash associated with a given height.
func (k Keeper) SetHeightHash(ctx sdk.Context, height uint64, hash ethcmn.Hash) {
	types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx).SetHeightHash(height, hash)
}

// ----------------------------------------------------------------------------
// Block bloom bits mapping functions
// Required by Web3 API.
// ----------------------------------------------------------------------------

// GetBlockBloom gets bloombits from block height
func (k Keeper) GetBlockBloom(ctx sdk.Context, height int64) ethtypes.Bloom {
	return k.getBlockBloomInDiskDB(height)
}

// SetBlockBloom sets the mapping from block height to bloom bits
func (k Keeper) SetBlockBloom(ctx sdk.Context, height int64, bloom ethtypes.Bloom) {
	k.setBlockBloomInDiskDB(height, bloom)
}

// IterateBlockBloom iterates all over the bloom value in every height
func (k Keeper) IterateBlockBloom(ctx sdk.Context, fn func(key []byte, value []byte) (stop bool)) {
	k.iterateBlockBloomInDiskDB(fn)
}

// GetAccountStorage return state storage associated with an account
func (k Keeper) GetAccountStorage(ctx sdk.Context, address ethcmn.Address) (types.Storage, error) {
	storage := types.Storage{}
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	err := csdb.ForEachStorage(address, func(key, value ethcmn.Hash) bool {
		storage = append(storage, types.NewState(key, value))
		return false
	})
	if err != nil {
		return types.Storage{}, err
	}

	return storage, nil
}

// getChainConfig get raw chain config and unmarshal it
func (k *Keeper) getChainConfig(ctx sdk.Context) (types.ChainConfig, bool) {
	// if keeper has cached the chain config, return immediately

	store := k.Ada.NewStore(k.paramSpace.CustomKVStore(ctx), types.KeyPrefixChainConfig)

	// get from an empty key that's already prefixed by KeyPrefixChainConfig
	bz := store.Get([]byte{})
	if len(bz) == 0 {
		return types.ChainConfig{}, false
	}

	var config types.ChainConfig
	// first 4 bytes are type prefix
	// bz len must > 4; otherwise, MustUnmarshalBinaryBare will panic
	if err := config.UnmarshalFromAmino(k.cdc, bz[4:]); err != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &config)
	}

	return config, true
}

// GetChainConfig gets chain config, the result if from cached result, or
// it gains chain config and gas costs from getChainConfig, then
// cache the chain config and gas costs.
func (k Keeper) GetChainConfig(ctx sdk.Context) (types.ChainConfig, bool) {
	// if keeper has cached the chain config, return immediately, and increase gas costs.
	if k.cci.cc != nil {
		ctx.GasMeter().ConsumeGas(k.cci.gasReduced, "cached chain config recover")
		return *k.cci.cc, true
	}

	gasStart := ctx.GasMeter().GasConsumed()
	chainConfig, found := k.getChainConfig(ctx)
	gasStop := ctx.GasMeter().GasConsumed()

	// only cache chain config result when we found it, or try to found again.
	if found {
		k.cci.cc = &chainConfig
		k.cci.gasReduced = gasStop - gasStart
	}

	return chainConfig, found
}

// SetChainConfig sets the mapping from block consensus hash to block height
func (k *Keeper) SetChainConfig(ctx sdk.Context, config types.ChainConfig) {
	store := k.Ada.NewStore(k.paramSpace.CustomKVStore(ctx), types.KeyPrefixChainConfig)

	bz := k.cdc.MustMarshalBinaryBare(config)
	// get to an empty key that's already prefixed by KeyPrefixChainConfig
	store.Set([]byte{}, bz)

	// invalid the chainConfig
	k.cci.cc = nil
}

// SetGovKeeper sets keeper of gov
func (k *Keeper) SetGovKeeper(gk GovKeeper) {
	k.govKeeper = gk
}

var commitStateDBPool = &sync.Pool{
	New: func() interface{} {
		return &types.CommitStateDB{GuFactor: types.DefaultGuFactor}
	},
}

// checks whether the address is blocked
func (k *Keeper) IsAddressBlocked(ctx sdk.Context, addr sdk.AccAddress) bool {
	csdb := commitStateDBPool.Get().(*types.CommitStateDB)
	defer commitStateDBPool.Put(csdb)
	types.ResetCommitStateDB(csdb, k.GenerateCSDBParams(), &ctx)

	// csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	return k.GetParams(ctx).EnableContractBlockedList && csdb.IsContractInBlockedList(addr.Bytes())
}

func (k *Keeper) IsContractInBlockedList(ctx sdk.Context, addr sdk.AccAddress) bool {
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	return csdb.IsContractInBlockedList(addr.Bytes())
}

func (k *Keeper) SetObserverKeeper(infuraKeeper watcher.InfuraKeeper) {
	k.Watcher.InfuraKeeper = infuraKeeper
}

// SetHooks sets the hooks for the EVM module
// It should be called only once during initialization, it panics if called more than once.
func (k *Keeper) SetHooks(hooks types.EvmHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set evm hooks twice")
	}
	k.hooks = hooks

	return k
}

// ResetHooks resets the hooks for the EVM module
func (k *Keeper) ResetHooks() *Keeper {
	k.hooks = nil

	return k
}

// GetHooks gets the hooks for the EVM module
func (k *Keeper) GetHooks() types.EvmHooks {
	return k.hooks
}

// CallEvmHooks delegate the call to the hooks. If no hook has been registered, this function returns with a `nil` error
func (k *Keeper) CallEvmHooks(ctx sdk.Context, st *types.StateTransition, receipt *ethtypes.Receipt) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.PostTxProcessing(ctx, st, receipt)
}

// Add latest block height and hash to lru cache
func (k *Keeper) AddHeightHashToCache(height int64, hash string) {
	k.heightCache.Add(hash, height)
	k.hashCache.Add(height, hash)
}
