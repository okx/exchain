package cli

import (
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
)

type (
	stakingMsgBuildingHelpers = genutilcli.StakingMsgBuildingHelpers
)

var (
	// nolint
	ValidateGenesisCmd = genutilcli.ValidateGenesisCmd
)
