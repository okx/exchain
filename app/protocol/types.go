package protocol

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/okex/okchain/x/gov"
	"github.com/okex/okchain/x/token"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/okex/okchain/x/staking"

	"github.com/cosmos/cosmos-sdk/x/crisis"

	"github.com/okex/okchain/x/backend"
	distr "github.com/okex/okchain/x/distribution"
	"github.com/okex/okchain/x/stream"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Protocol interface {
	GetVersion() uint64
	ExportAppStateAndValidators(ctx sdk.Context) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error)

	// load base installation 4 each protolcol
	LoadContext()
	Init()
	GetCodec() *codec.Codec

	// gracefully stop okchaind
	CheckStopped()

	// setter
	SetLogger(log log.Logger) Protocol
	SetParent(parent Parent) Protocol

	//getter
	GetParent() Parent

	// get specific keeper
	GetBackendKeeper() backend.Keeper
	GetStreamKeeper() stream.Keeper
	GetCrisisKeeper() crisis.Keeper
	GetStakingKeeper() staking.Keeper
	GetDistrKeeper() distr.Keeper
	GetSlashingKeeper() slashing.Keeper
	GetTokenKeeper() token.Keeper

	// fit 4 cm36
	GetKVStoreKeysMap() map[string]*sdk.KVStoreKey
	GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey
	ExportGenesis(ctx sdk.Context) map[string]json.RawMessage
}

//----------------------
// baseapp interface
//----------------------
type Parent interface {
	DeliverTx(abci.RequestDeliverTx) abci.ResponseDeliverTx
	PushInitChainer(initChainer sdk.InitChainer)
	PushBeginBlocker(beginBlocker sdk.BeginBlocker)
	PushEndBlocker(endBlocker sdk.EndBlocker)
	PushAnteHandler(ah sdk.AnteHandler)
	SetRouter(router sdk.Router, queryRouter sdk.QueryRouter)
}

// types of exported
type ExportedAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.DecCoins   `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`

	/* vesting account fields */
	// total vesting coins upon initialization
	OriginalVesting sdk.Coins `json:"original_vesting"`
	// delegated vested coins at time of delegation
	DelegatedFree sdk.Coins `json:"delegated_free"`
	// delegated vesting coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting"`
	// vesting start time (UNIX Epoch time)
	StartTime int64 `json:"start_time"`
	// vesting end time (UNIX Epoch time)
	EndTime int64 `json:"end_time"`
}

func NewExportedAccount(acc exported.Account) ExportedAccount {
	exportedAcc := ExportedAccount{
		Address:       acc.GetAddress(),
		Coins:         sdk.NewDecCoins(acc.GetCoins()),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	vestAcc, ok := acc.(auth.VestingAccount)
	if ok {
		exportedAcc.OriginalVesting = vestAcc.GetOriginalVesting()
		exportedAcc.DelegatedFree = vestAcc.GetDelegatedFree()
		exportedAcc.DelegatedVesting = vestAcc.GetDelegatedVesting()
		exportedAcc.StartTime = vestAcc.GetStartTime()
		exportedAcc.EndTime = vestAcc.GetEndTime()
	}

	return exportedAcc
}

// state
type ExportState struct {
	AuthData auth.GenesisState `json:"auth"`
	BankData bank.GenesisState `json:"bank"`
	Accounts []ExportedAccount `json:"accounts"`
	GovData  gov.GenesisState  `json:"gov"`
}
