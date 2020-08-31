package pulsarclient

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/stream/common"
	"github.com/okex/okchain/x/stream/types"
)

type PulsarData struct {
	Height        int64
	matchResults  []*backend.MatchResult
	newTokenPairs []*dex.TokenPair
}

func NewPulsarData() *PulsarData {
	return &PulsarData{
		matchResults: make([]*backend.MatchResult, 0),
	}
}

func (p *PulsarData) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper, cache *common.Cache) {
	p.Height = ctx.BlockHeight()
	p.matchResults = common.GetMatchResults(ctx, orderKeeper)
	p.newTokenPairs = cache.GetNewTokenPairs()
}

var _ types.IStreamData = (*PulsarData)(nil)

func (p PulsarData) BlockHeight() int64 {
	return p.Height
}

func (p PulsarData) DataType() types.StreamDataKind {
	return types.StreamDataKlineKind
}
