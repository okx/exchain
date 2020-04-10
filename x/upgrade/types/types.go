package types

import (
	"fmt"

	"github.com/okex/okchain/x/common/proto"
)

type VersionInfo struct {
	UpgradeInfo proto.AppUpgradeConfig `json:"upgrade_info"`
	Success     bool                   `json:"success"`
}

func NewVersionInfo(upgradeConfig proto.AppUpgradeConfig, success bool) VersionInfo {
	return VersionInfo{
		upgradeConfig,
		success,
	}
}

type QueryVersion struct {
	Ver uint64 `json:"version"`
}

func NewQueryVersion(ver uint64) QueryVersion {
	return QueryVersion{
		Ver: ver,
	}
}

func (qv QueryVersion) String() string {
	return fmt.Sprintf("The query version is %d", qv.Ver)
}
