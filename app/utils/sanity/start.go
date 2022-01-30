package sanity

import (
	apptype "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/tendermint/consensus"
	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/watcher"
)

// CheckStart check start command's flags. if user set conflict flags return error.
// the conflicts flags are:
// --fast-query      conflict with --paralleled-tx=true
// --fast-query      conflict with --pruning=nothing
// --enable-preruntx conflict with --download-delta
// --upload-delta    conflict with --paralleled-tx=true
//
// based the conflicts above and node-mode below
// --node-mode=rpc manage the following flags:
//     --disable-checktx-mutex=true
//     --disable-query-mutex=true
//     --enable-bloom-filter=true
//     --fast-lru=10000
//     --fast-query=true
//     --iavl-enable-async-commit=true
//     --max-open=20000
//     --mempool.enable_pending_pool=true
//     --cors=*
//
// --node-mode=validator manage the following flags:
//     --disable-checktx-mutex=true
//     --disable-query-mutex=true
//     --enable-dynamic-gp=false
//     --iavl-enable-async-commit=true
//     --iavl-cache-size=10000000
//     --pruning=everything
//
// --node-mode=archive manage the following flags:
//    --pruning=nothing
//    --disable-checktx-mutex=true
//    --disable-query-mutex=true
//    --enable-bloom-filter=true
//    --iavl-enable-async-commit=true
//    --max-open=20000
//    --cors=*
//
// then
// --node-mode=rpc(--fast-query) conflicts with --paralleled-tx=true and --pruning=nothing
// --node-mode=archive(--pruning=nothing) conflicts with --fast-query

var (
	// conflicts flags
	startConflictElems = []conflictPair{
		// --fast-query      conflict with --paralleled-tx=true
		{
			configA: boolItem{name: watcher.FlagFastQuery, value: true},
			configB: boolItem{name: state.FlagParalleledTx, value: true},
		},
		// --fast-query      conflict with --pruning=nothing
		{
			configA: boolItem{name: watcher.FlagFastQuery, value: true},
			configB: stringItem{name: server.FlagPruning, value: cosmost.PruningOptionNothing},
		},
		// --enable-preruntx conflict with --download-delta
		{
			configA: boolItem{name: consensus.EnablePrerunTx, value: true},
			configB: boolItem{name: types.FlagDownloadDDS, value: true},
		},
		// --upload-delta    conflict with --paralleled-tx=true
		{
			configA: boolItem{name: types.FlagUploadDDS, value: true},
			configB: boolItem{name: state.FlagParalleledTx, value: true},
		},
		// --node-mode=rpc(--fast-query) conflicts with --paralleled-tx=true and --pruning=nothing
		{
			configA: stringItem{name: apptype.FlagNodeMode, value: string(apptype.RpcNode)},
			configB: boolItem{name: state.FlagParalleledTx, value: true},
		},
		{
			configA: stringItem{name: apptype.FlagNodeMode, value: string(apptype.RpcNode)},
			configB: stringItem{name: server.FlagPruning, value: cosmost.PruningOptionNothing},
		},
		// --node-mode=archive(--pruning=nothing) conflicts with --fast-query
		{
			configA: stringItem{name: apptype.FlagNodeMode, value: string(apptype.ArchiveNode)},
			configB: boolItem{name: watcher.FlagFastQuery, value: true},
		},
	}
)

// CheckStart check start command.If it has conflict pair above. then return the conflict error
func CheckStart() error {
	for _, v := range startConflictElems {
		if err := v.checkConflict(); err != nil {
			return err
		}
	}

	return nil
}
