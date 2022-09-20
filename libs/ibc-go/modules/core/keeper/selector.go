package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	DefaultFactory = func(higherThan func(h int64) bool, version int, internal IBCServerKeeper) SelectorFactory {
		return func() KeeperSelector {
			return NewDefaultKeepSelector(higherThan, version, internal)
		}
	}
)

type SelectorFactory func() KeeperSelector

type KeeperSelector interface {
	Version() int
	Select(ctx sdk.Context) (IBCServerKeeper, bool)
}

type KeeperSelectors []KeeperSelector

func (k KeeperSelectors) Len() int {
	return len(k)
}

func (k KeeperSelectors) Less(i, j int) bool {
	return k[i].Version() >= k[j].Version()
}

func (k KeeperSelectors) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

var (
	_ KeeperSelector = (*DefaultKeepSelector)(nil)
)

type DefaultKeepSelector struct {
	higherThan func(h int64) bool
	version    int

	internal IBCServerKeeper
}

func (d *DefaultKeepSelector) Version() int {
	return d.version
}

func NewDefaultKeepSelector(higherThan func(h int64) bool, version int, internal IBCServerKeeper) *DefaultKeepSelector {
	return &DefaultKeepSelector{
		higherThan: higherThan,
		version:    version,
		internal:   internal,
	}
}

func (d *DefaultKeepSelector) Select(ctx sdk.Context) (IBCServerKeeper, bool) {
	if d.higherThan(ctx.BlockHeight()) {
		return d.internal, true
	}
	return nil, false
}
