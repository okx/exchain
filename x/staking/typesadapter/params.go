package typesadapter

import (
	"github.com/okex/exchain/x/staking/types"
)

func (p *Params) From(pp types.Params) {
	p.MaxValidators = uint32(pp.MaxValidators)
	p.UnbondingTime = pp.UnbondingTime
	p.HistoricalEntries = pp.HistoricalEntries
}
