package keeper

//import (
//	"errors"
//	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
//	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
//	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
//	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
//	"github.com/okex/exchain/x/wasm/types"
//	"github.com/okex/exchain/x/wasm/watcher"
//)
//
//type WatchKeeper struct {
//	watchDB *watcher.Watcher
//	keeper  *Keeper
//}
//
//func NewWatchKeeper(watchdb *watcher.Watcher, keeper *Keeper) *WatchKeeper {
//	return &WatchKeeper{
//		watchDB: watchdb,
//		keeper:  keeper,
//	}
//}
//
//func (wk *WatchKeeper) GetContractHistory(ctx sdk.Context, contractAddr sdk.AccAddress) []types.ContractCodeHistoryEntry {
//	return wk.keeper.GetContractHistory(ctx, contractAddr)
//	r := make([]types.ContractCodeHistoryEntry, 0)
//	prefixStore := prefix.NewStore(wk.watchDB, types.GetContractCodeHistoryElementPrefix(contractAddr))
//	iter := prefixStore.Iterator(nil, nil)
//	defer iter.Close()
//
//	for ; iter.Valid(); iter.Next() {
//		var e types.ContractCodeHistoryEntry
//		wk.keeper.cdc.GetProtocMarshal().MustUnmarshal(iter.Value(), &e)
//		r = append(r, e)
//	}
//	return r
//	return nil
//}
//
//func (wk *WatchKeeper) QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
//	contractInfo, codeInfo, prefixStore, err := wk.contractInstance(ctx, contractAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	// treat it as pinned
//	gasRegister := NewDefaultWasmGasRegister()
//	smartQuerySetupCosts := gasRegister.InstantiateContractCosts(wk.IsPinnedCode(ctx, contractInfo.CodeID), len(req))
//	ctx.GasMeter().ConsumeGas(smartQuerySetupCosts, "Loading CosmWasm module: query")
//
//	// prepare querier
//	bankWatcher := NewWatchBankKeeper(watcher.NewWatcher(), wk.keeper.bankKeeper)
//	queryPlugin := DefaultQueryPlugins(bankWatcher, nil, wk.keeper.queryRouter, wk)
//	querier := NewQueryHandler(ctx, queryPlugin, contractAddr, gasRegister)
//
//	env := types.NewEnv(ctx, contractAddr)
//	queryResult, gasUsed, qErr := wk.keeper.wasmVM.Query(codeInfo.CodeHash, env, req, prefixStore, cosmwasmAPI, querier, NewMultipliedGasMeter(ctx.GasMeter(), wk.keeper.gasRegister), wk.keeper.runtimeGasForContract(ctx), costJSONDeserialization)
//	wk.keeper.consumeRuntimeGas(ctx, gasUsed)
//	if qErr != nil {
//		return nil, sdkerrors.Wrap(types.ErrQueryFailed, qErr.Error())
//	}
//	return queryResult, nil
//}
//
//func (wk *WatchKeeper) QueryRaw(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) []byte {
//	if key == nil {
//		return nil
//	}
//	prefixStoreKey := types.GetContractStorePrefix(contractAddress)
//	prefixStore := prefix.NewStore(wk.watchDB, prefixStoreKey)
//	return prefixStore.Get(key)
//}
//
//func (wk *WatchKeeper) HasContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool {
//	store := wk.watchDB
//	return store.Has(types.GetContractAddressKey(contractAddress))
//}
//
//func (wk *WatchKeeper) GetContractInfo(_ sdk.Context, contractAddress sdk.AccAddress) *types.ContractInfo {
//	store := wk.watchDB
//	var contract types.ContractInfo
//	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
//	if contractBz == nil {
//		return nil
//	}
//	wk.keeper.cdc.GetProtocMarshal().MustUnmarshal(contractBz, &contract)
//	return &contract
//}
//
//func (wk *WatchKeeper) IterateContractInfo(ctx sdk.Context, cb func(sdk.AccAddress, types.ContractInfo) bool) {
//}
//
//func (wk *WatchKeeper) IterateContractsByCode(ctx sdk.Context, codeID uint64, cb func(address sdk.AccAddress) bool) {
//
//}
//
//func (wk *WatchKeeper) IterateContractState(ctx sdk.Context, contractAddress sdk.AccAddress, cb func(key, value []byte) bool) {
//
//}
//
//func (wk *WatchKeeper) GetCodeInfo(ctx sdk.Context, codeID uint64) *types.CodeInfo {
//	store := wk.watchDB
//	var codeInfo types.CodeInfo
//	codeInfoBz := store.Get(types.GetCodeKey(codeID))
//	if codeInfoBz == nil {
//		return nil
//	}
//	wk.keeper.cdc.GetProtocMarshal().MustUnmarshal(codeInfoBz, &codeInfo)
//	return &codeInfo
//}
//
//func (wk *WatchKeeper) IterateCodeInfos(ctx sdk.Context, cb func(uint64, types.CodeInfo) bool) {
//}
//
//func (wk *WatchKeeper) GetByteCode(ctx sdk.Context, codeID uint64) ([]byte, error) {
//	codeInfo := wk.GetCodeInfo(ctx, codeID)
//	return wk.keeper.wasmVM.GetCode(codeInfo.CodeHash)
//}
//
//func (wk *WatchKeeper) IsPinnedCode(_ sdk.Context, codeID uint64) bool {
//	store := wk.watchDB
//	return store.Has(types.GetPinnedCodeIndexPrefix(codeID))
//}
//
//func (wk *WatchKeeper) contractInstance(_ sdk.Context, contractAddress sdk.AccAddress) (types.ContractInfo, types.CodeInfo, types.StoreAdapter, error) {
//	// TODO
//	// set watch db when set to state db
//	store := wk.watchDB
//
//	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
//	if contractBz == nil {
//		return types.ContractInfo{}, types.CodeInfo{}, types.StoreAdapter{}, sdkerrors.Wrap(types.ErrNotFound, "contract")
//	}
//	var contractInfo types.ContractInfo
//	wk.keeper.cdc.GetProtocMarshal().MustUnmarshal(contractBz, &contractInfo)
//
//	codeInfoBz := store.Get(types.GetCodeKey(contractInfo.CodeID))
//	if codeInfoBz == nil {
//		return contractInfo, types.CodeInfo{}, types.StoreAdapter{}, sdkerrors.Wrap(types.ErrNotFound, "code info")
//	}
//	var codeInfo types.CodeInfo
//	wk.keeper.cdc.GetProtocMarshal().MustUnmarshal(codeInfoBz, &codeInfo)
//	prefixStoreKey := types.GetContractStorePrefix(contractAddress)
//	prefixStore := prefix.NewStore(store, prefixStoreKey)
//	return contractInfo, codeInfo, types.NewStoreAdapter(prefixStore), nil
//}
//
//type WatchBankKeeper struct {
//	db     *watcher.Watcher
//	keeper types.BankViewKeeper
//}
//
//func NewWatchBankKeeper(w *watcher.Watcher, keeper types.BankViewKeeper) *WatchBankKeeper {
//	return &WatchBankKeeper{
//		db:     w,
//		keeper: keeper,
//	}
//}
//
//func (wbk *WatchBankKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
//
//	return nil
//}
//
//func (wbk *WatchBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
//	return sdk.Coin{}
//}
//
//type WatchQueryHandler struct {
//	Ctx         sdk.Context
//	Plugins     WasmVMQueryHandler
//	Caller      sdk.AccAddress
//	gasRegister GasRegister
//}
//
//func NewWatchQueryHandler(ctx sdk.Context, vmQueryHandler WasmVMQueryHandler, caller sdk.AccAddress, gasRegister GasRegister) *WatchQueryHandler {
//	return &WatchQueryHandler{
//		Ctx:         ctx,
//		Plugins:     vmQueryHandler,
//		Caller:      caller,
//		gasRegister: gasRegister,
//	}
//}
//
//func (w *WatchQueryHandler) Query(request wasmvmtypes.QueryRequest, gasLimit uint64) ([]byte, error) {
//	// set a limit for a subCtx
//	sdkGas := w.gasRegister.FromWasmVMGas(gasLimit)
//	// discard all changes/ events in subCtx by not committing the cached context
//	subCtx := sdk.Context{}
//	subCtx.SetGasMeter(sdk.NewGasMeter(sdkGas))
//
//	// make sure we charge the higher level context even on panic
//	defer func() {
//		w.Ctx.GasMeter().ConsumeGas(subCtx.GasMeter().GasConsumed(), "contract sub-query")
//	}()
//
//	res, err := w.Plugins.HandleQuery(subCtx, w.Caller, request)
//	if err == nil {
//		// short-circuit, the rest is dealing with handling existing errors
//		return res, nil
//	}
//
//	// special mappings to system error (which are not redacted)
//	var noSuchContract *types.ErrNoSuchContract
//	if ok := errors.As(err, &noSuchContract); ok {
//		err = wasmvmtypes.NoSuchContract{Addr: noSuchContract.Addr}
//	}
//
//	// Issue #759 - we don't return error string for worries of non-determinism
//	return nil, redactError(err)
//}
//
//func (w *WatchQueryHandler) GasConsumed() uint64 {
//	return w.Ctx.GasMeter().GasConsumed()
//}
