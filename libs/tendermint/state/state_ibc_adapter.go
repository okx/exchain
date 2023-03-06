package state

import (
	"github.com/okx/okbchain/libs/tendermint/version"
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

var ibcStateVersion = Version{
	Consensus: version.Consensus{
		Block: version.IBCBlockProtocol,
		App:   0,
	},
	Software: version.TMCoreSemVer,
}

func GetStateVersion(h int64) Version {
	return ibcStateVersion
}
