package sanity

import (
	"fmt"
	apptype "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/tendermint/consensus"
	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/cobra"
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
	conflictElems = []conflictPair{
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

type conflictPair struct {
	configA item
	configB item
}

// checkConflict check configA vs configB by cmd and viper
// if both configA and configB is set by user,
// and the value is equal to the conflicts value then complain it
func (cp *conflictPair) checkConflict(cmd *cobra.Command) error {
	if checkUserSetFlag(cmd, cp.configA.label()) &&
		checkUserSetFlag(cmd, cp.configB.label()) {
		if cp.configA.check() &&
			cp.configB.check() {
			return fmt.Errorf(" %v conflict with %v", cp.configA.verbose(), cp.configB.verbose())
		}
	}

	return nil
}

func CheckStart(cmd *cobra.Command) error {
	for _, v := range conflictElems {
		if err := v.checkConflict(cmd); err != nil {
			return err
		}
	}

	return nil
}

// checkUserSetFlag If the user set the value (or if left to default)
func checkUserSetFlag(cmd *cobra.Command, inFlag string) bool {
	flag := cmd.Flags().Lookup(inFlag)
	return flag != nil && flag.Changed
}
