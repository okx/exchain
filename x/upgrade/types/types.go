package types

import (
	"fmt"

	"github.com/okex/okexchain/x/common/proto"
)

// VersionInfo is the struct of version info
type VersionInfo struct {
	UpgradeInfo proto.AppUpgradeConfig `json:"upgrade_info"`
	Success     bool                   `json:"success"`
}

// NewVersionInfo creates a new instance of NewVersionInfo
func NewVersionInfo(upgradeConfig proto.AppUpgradeConfig, success bool) VersionInfo {
	return VersionInfo{
		upgradeConfig,
		success,
	}
}

// QueryVersion is designed for version query
type QueryVersion struct {
	Ver uint64 `json:"version"`
}

// NewQueryVersion creates a new instance of QueryVersion
func NewQueryVersion(ver uint64) QueryVersion {
	return QueryVersion{
		Ver: ver,
	}
}

// String returns a human readable string representation of QueryVersion
func (qv QueryVersion) String() string {
	return fmt.Sprintf("The query version is %d", qv.Ver)
}
