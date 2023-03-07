package typesadapter

import (
	"github.com/okx/okbchain/x/staking/types"
)

func (p *Params) From(pp types.Params) {
	p.MaxValidators = uint32(pp.MaxValidators)
	p.UnbondingTime = pp.UnbondingTime
	if pp.HistoricalEntries == 0 {
		p.HistoricalEntries = 10000
	} else {
		p.HistoricalEntries = pp.HistoricalEntries
	}
}
