package proto

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

// keys
var (
	upgradeConfigKey     = []byte("upgrade_config")
	currentVersionKey    = []byte("current_version")
	lastFailedVersionKey = []byte("last_failed_version")
	cdc                  = codec.New()
)
