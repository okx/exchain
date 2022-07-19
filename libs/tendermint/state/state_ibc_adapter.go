package state

import (
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/libs/tendermint/version"
)

func (v Version) UpgradeToIBCVersion() Version {
	return Version{
		Consensus: version.Consensus{
			Block: version.IBCBlockProtocol,
			App:   v.Consensus.App,
		},
		Software: v.Software,
	}
}

func (v Version) IsUpgraded() bool {
	return v.Consensus.Block == version.IBCBlockProtocol
}

var ibcStateVersion = Version{
	Consensus: version.Consensus{
		Block: version.IBCBlockProtocol,
		App:   0,
	},
	Software: version.TMCoreSemVer,
}

func GetStateVersion(h int64) Version {
	if types.HigherThanVenus1(h) {
		return ibcStateVersion
	}
	return initStateVersion
}
