package server

import (
	"fmt"
	"strings"

	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"

	"github.com/spf13/viper"

	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	tmiavl "github.com/okx/okbchain/libs/iavl"
	iavlcfg "github.com/okx/okbchain/libs/iavl/config"
)

// GetPruningOptionsFromFlags parses command flags and returns the correct
// PruningOptions. If a pruning strategy is provided, that will be parsed and
// returned, otherwise, it is assumed custom pruning options are provided.
func GetPruningOptionsFromFlags() (types.PruningOptions, error) {
	strategy := strings.ToLower(viper.GetString(FlagPruning))

	switch strategy {
	case types.PruningOptionNothing:
		tmiavl.EnablePruningHistoryState = false
		tmiavl.CommitIntervalHeight = 1
		mpt.TrieCommitGap = 1
		iavlcfg.DynamicConfig.SetCommitGapHeight(1)
		return types.NewPruningOptionsFromString(strategy), nil

	case types.PruningOptionDefault, types.PruningOptionEverything:
		return types.NewPruningOptionsFromString(strategy), nil

	case types.PruningOptionCustom:
		opts := types.NewPruningOptions(
			viper.GetUint64(FlagPruningKeepRecent),
			viper.GetUint64(FlagPruningKeepEvery), viper.GetUint64(FlagPruningInterval),
			viper.GetUint64(FlagPruningMaxWsNum),
		)

		if err := opts.Validate(); err != nil {
			return opts, fmt.Errorf("invalid custom pruning options: %w", err)
		}

		mpt.TrieDirtyDisabled = opts.KeepEvery == 1

		return opts, nil

	default:
		return store.PruningOptions{}, fmt.Errorf("unknown pruning strategy %s", strategy)
	}
}
