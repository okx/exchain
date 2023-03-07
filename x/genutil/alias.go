package genutil

import (
	"github.com/okx/okbchain/x/genutil/types"

	sdkgenutil "github.com/okx/okbchain/libs/cosmos-sdk/x/genutil"
	sdkgenutiltypes "github.com/okx/okbchain/libs/cosmos-sdk/x/genutil/types"
)

// const
const (
	ModuleName = types.ModuleName
)

type (
	// GenesisState is the type alias of the one in cmsdk
	GenesisState = types.GenesisState
	// InitConfig is the type alias of the one in cmsdk
	InitConfig = sdkgenutil.InitConfig
	// GenesisAccountsIterator is the type alias of the one in cmsdk
	GenesisAccountsIterator = sdkgenutiltypes.GenesisAccountsIterator
)

var (
	// nolint
	ModuleCdc                           = types.ModuleCdc
	GenesisStateFromGenFile             = sdkgenutil.GenesisStateFromGenFile
	NewGenesisState                     = sdkgenutil.NewGenesisState
	SetGenesisStateInAppState           = sdkgenutil.SetGenesisStateInAppState
	InitializeNodeValidatorFiles        = sdkgenutil.InitializeNodeValidatorFiles
	ExportGenesisFileWithTime           = sdkgenutil.ExportGenesisFileWithTime
	NewInitConfig                       = sdkgenutil.NewInitConfig
	ValidateGenesis                     = types.ValidateGenesis
	GenesisStateFromGenDoc              = sdkgenutil.GenesisStateFromGenDoc
	SetGenTxsInAppGenesisState          = sdkgenutil.SetGenTxsInAppGenesisState
	ExportGenesisFile                   = sdkgenutil.ExportGenesisFile
	InitializeNodeValidatorFilesByIndex = sdkgenutil.InitializeNodeValidatorFilesByIndex
)
