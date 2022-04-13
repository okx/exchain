package ibctesting

import (
	"encoding/json"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"testing"

	//cryptocodec "github.com/okex/exchain/app/crypto/ethsecp256k1"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"

	//authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"

	capabilitykeeper "github.com/okex/exchain/libs/cosmos-sdk/x/capability/keeper"
	stakingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	stakingkeeper "github.com/okex/exchain/x/staking"
	"github.com/stretchr/testify/require"

	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/ibc-go/modules/core/keeper"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp"
)

var DefaultTestingAppInit func() (TestingApp, map[string]json.RawMessage) = SetupTestingApp

type TestingApp interface {
	abci.Application

	// ibc-go additions
	GetBaseApp() *bam.BaseApp
	GetStakingKeeper() stakingkeeper.Keeper
	GetIBCKeeper() *keeper.Keeper
	GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper
	//GetTxConfig() client.TxConfig

	// Implemented by SimApp
	AppCodec() *codec.CodecProxy

	// Implemented by BaseApp
	LastCommitID() sdk.CommitID
	LastBlockHeight() int64
}

func SetupTestingApp() (TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	//encCdc := simapp.MakeTestEncodingConfig()
	app := simapp.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 5)
	return app, simapp.NewDefaultGenesisState(nil)
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authexported.GenesisAccount, balances ...sdk.Coins) TestingApp {
	app, genesisState := DefaultTestingAppInit()
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().GetCdc().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.NewInt(1000000)

	for _, val := range valSet.Validators {
		//pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		//require.NoError(t, err)
		//pkAny, err := codectypes.NewAnyWithValue(pk)
		//require.NoError(t, err)
		validator := stakingtypes.Validator{
			//OperatorAddress:   sdk.ValAddress(val.Address).String(),
			//ConsensusPubkey:   pkAny,
			Jailed: false,
			//Status:            stakingtypes.Bonded,
			Tokens:          bondAmt,
			DelegatorShares: sdk.OneDec(),
			Description:     stakingtypes.Description{},
			UnbondingHeight: int64(0),
			//UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

	}
	// set validators and delegations
	//stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	//genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	// todo bank genesis state file
	//	totalSupply := sdk.NewCoins()
	//	for _, b := range balances {
	// add genesis acc tokens and delegated tokens to total supply
	//		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))...)
	//	}

	// add bonded amount to bonded pool module account
	// balances = append(balances, banktypes.Balance{
	// 	Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
	// 	Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	// })

	// update total supply
	// bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	// genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Validators: []abci.ValidatorUpdate{},
			//ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes: stateBytes,
		},
	)

	// commit genesis changes
	// app.Commit()
	// app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
	// 	Height:             app.LastBlockHeight() + 1,
	// 	AppHash:            app.LastCommitID().Hash,
	// 	ValidatorsHash:     valSet.Hash(),
	// 	NextValidatorsHash: valSet.Hash(),
	// }})

	return app
}
