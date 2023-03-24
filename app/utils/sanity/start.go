package sanity

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/okex/exchain/app/config"
	apptype "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/consensus"
	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/types"
	db "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/infura"
)

// CheckStart check start command's flags. if user set conflict flags return error.
// the conflicts flags are:
// --fast-query      conflict with --pruning=nothing
// --enable-preruntx conflict with --download-delta
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
//     --dynamic-gp-mode=2
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
// --node-mode=archive(--pruning=nothing) conflicts with --fast-query

var (
	startDependentElems = []dependentPair{
		{ // if infura.FlagEnable=true , watcher.FlagFastQuery must be set to true
			config:       boolItem{name: infura.FlagEnable, expect: true},
			reliedConfig: boolItem{name: watcher.FlagFastQuery, expect: true},
		},
	}
	// conflicts flags
	startConflictElems = []conflictPair{
		// --fast-query      conflict with --pruning=nothing
		{
			configA: boolItem{name: watcher.FlagFastQuery, expect: true},
			configB: stringItem{name: server.FlagPruning, expect: cosmost.PruningOptionNothing},
		},
		// --enable-preruntx conflict with --download-delta
		{
			configA: boolItem{name: consensus.EnablePrerunTx, expect: true},
			configB: boolItem{name: types.FlagDownloadDDS, expect: true},
		},
		// --multi-cache conflict with --download-delta
		{
			configA: boolItem{name: sdk.FlagMultiCache, expect: true},
			configB: boolItem{name: types.FlagDownloadDDS, expect: true},
		},
		{
			configA: stringItem{name: apptype.FlagNodeMode, expect: string(apptype.RpcNode)},
			configB: stringItem{name: server.FlagPruning, expect: cosmost.PruningOptionNothing},
		},
		// --node-mode=archive(--pruning=nothing) conflicts with --fast-query
		{
			configA: stringItem{name: apptype.FlagNodeMode, expect: string(apptype.ArchiveNode)},
			configB: boolItem{name: watcher.FlagFastQuery, expect: true},
		},
		{
			configA: stringItem{name: apptype.FlagNodeMode, expect: string(apptype.RpcNode)},
			configB: boolItem{name: config.FlagEnablePGU, expect: true},
		},
		{
			configA: stringItem{name: apptype.FlagNodeMode, expect: string(apptype.ArchiveNode)},
			configB: boolItem{name: config.FlagEnablePGU, expect: true},
		},
		{
			configA: stringItem{name: apptype.FlagNodeMode, expect: string(apptype.InnertxNode)},
			configB: boolItem{name: config.FlagEnablePGU, expect: true},
		},
	}

	checkRangeItems = []rangeItem{
		{
			enumRange: []int{int(state.DeliverTxsExecModeSerial), state.DeliverTxsExecModeParallel},
			name:      state.FlagDeliverTxsExecMode,
		},
	}
)

// CheckStart check start command.If it has conflict pair above. then return the conflict error
func CheckStart(ctx *server.Context) error {
	if viper.GetBool(FlagDisableSanity) {
		return nil
	}
	for _, v := range startDependentElems {
		if err := v.check(); err != nil {
			return err
		}
	}
	for _, v := range startConflictElems {
		if err := v.check(); err != nil {
			return err
		}
	}

	for _, v := range checkRangeItems {
		if err := v.checkRange(); err != nil {
			return err
		}
	}

	rocksDBMisspelling(ctx)
	return nil
}

func rocksDBMisspelling(ctx *server.Context) {
	//A copy of all the constant variables indicating rocksDB option parameters
	//in exchain/libs/tm-db/rocksdb.go
	const editDistanceThreshold = 2
	var minIndex int
	rocksDBConst := []string{
		"block_size",
		"block_cache",
		"statistics",
		"max_open_files",
		"allow_mmap_reads",
		"allow_mmap_writes",
		"unordered_write",
		"pipelined_write",
	}
	params := parseOptParams(viper.GetString(db.FlagRocksdbOpts))
	if params == nil {
		return
	}
	for _, str := range rocksDBConst {
		delete(params, str)
	}
	if len(params) != 0 {
		for inputOpt, _ := range params {
			optEditDistance := make([]int, len(rocksDBConst))
			minDistance := 20
			for i, expectOpt := range rocksDBConst {
				optEditDistance[i] = editDistance(inputOpt, expectOpt)
				if optEditDistance[i] < minDistance {
					minDistance = optEditDistance[i]
					minIndex = i
				}
			}
			if minDistance <= editDistanceThreshold {
				ctx.Logger.Info(fmt.Sprintf("%s %s failed to set rocksDB parameters, expect %s", db.FlagRocksdbOpts, inputOpt, rocksDBConst[minIndex]))
			} else {
				ctx.Logger.Info(fmt.Sprintf("%s %s failed to set rocksDB parameters, invalid parameter", db.FlagRocksdbOpts, inputOpt))
			}
		}
	}

	return
}

func parseOptParams(params string) map[string]struct{} {
	if len(params) == 0 {
		return nil
	}

	opts := make(map[string]struct{})
	for _, s := range strings.Split(params, ",") {
		opt := strings.Split(s, "=")
		if len(opt) != 2 {
			panic("Invalid options parameter, like this 'block_size=4kb,statistics=true")
		}
		opts[strings.TrimSpace(opt[0])] = struct{}{}
	}
	return opts
}

func editDistance(s, t string) int {
	m := len(s)
	n := len(t)

	if m == 0 {
		return n
	}

	if n == 0 {
		return m
	}

	// 创建二维切片
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	// 初始化第一行和第一列
	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	// 计算编辑距离
	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				insert := d[i][j-1] + 1
				delete := d[i-1][j] + 1
				substitute := d[i-1][j-1] + 1
				if insert <= delete && insert <= substitute {
					d[i][j] = insert
				} else if delete <= insert && delete <= substitute {
					d[i][j] = delete
				} else {
					d[i][j] = substitute
				}
			}
		}
	}

	return d[m][n]
}
