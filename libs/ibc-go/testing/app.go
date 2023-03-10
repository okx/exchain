package ibctesting

import (
	"encoding/json"
	staking "github.com/okx/okbchain/x/staking/types"
	"testing"
	"time"

	ibc "github.com/okx/okbchain/libs/ibc-go/modules/core"

	"github.com/okx/okbchain/libs/cosmos-sdk/client"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank"

	//cryptocodec "github.com/okx/okbchain/app/crypto/ethsecp256k1"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"

	//authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"

	capabilitykeeper "github.com/okx/okbchain/libs/cosmos-sdk/x/capability/keeper"
	stakingtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/staking/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
	dbm "github.com/okx/okbchain/libs/tm-db"
	"github.com/okx/okbchain/x/evm"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	stakingkeeper "github.com/okx/okbchain/x/staking"
	"github.com/stretchr/testify/require"

	bam "github.com/okx/okbchain/libs/cosmos-sdk/baseapp"
	"github.com/okx/okbchain/libs/ibc-go/modules/core/keeper"
	"github.com/okx/okbchain/libs/ibc-go/testing/simapp"
)

var DefaultTestingAppInit func() (TestingApp, map[string]json.RawMessage) = SetupTestingApp

// IBC application testing ports

type TestingApp interface {
	abci.Application
	TxConfig() client.TxConfig

	// ibc-go additions
	GetBaseApp() *bam.BaseApp
	GetStakingKeeper() stakingkeeper.Keeper
	GetIBCKeeper() *keeper.Keeper
	GetFacadedKeeper() *ibc.Keeper
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
	return app, simapp.NewDefaultGenesisState()
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(t *testing.T, chainId string, valSet *tmtypes.ValidatorSet, genAccs []authexported.GenesisAccount, balances ...sdk.Coins) TestingApp {
	app, genesisState := DefaultTestingAppInit()
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)

	genesisState[authtypes.ModuleName] = app.AppCodec().GetCdc().MustMarshalJSON(authGenesis)
	var err error
	if err != nil {
		panic("SetupWithGenesisValSet marshal error")
	}
	//var genesisState2 authtypes.GenesisState

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.NewInt(1000000)

	for _, val := range valSet.Validators {
		//pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		//require.NoError(t, err)
		//pkAny, err := codectypes.NewAnyWithValue(pk)
		//require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:         sdk.ValAddress(val.Address),
			ConsPubKey:              val.PubKey,
			Jailed:                  false,
			Status:                  sdk.Bonded,
			Tokens:                  bondAmt,
			DelegatorShares:         sdk.OneDec(),
			Description:             stakingtypes.Description{},
			UnbondingHeight:         int64(0),
			UnbondingCompletionTime: time.Unix(0, 0).UTC(),
			//UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().GetCdc().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		//add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))...)
	}

	balances = append(balances, sdk.Coins{
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000),
	})

	bankGenesis := bank.DefaultGenesisState()
	genesisState[bank.ModuleName] = app.AppCodec().GetCdc().MustMarshalJSON(bankGenesis)

	evmGenesis := evmtypes.DefaultGenesisState()
	evmGenesis.Params.EnableCall = true
	evmGenesis.Params.EnableCreate = true
	genesisState[evm.ModuleName] = app.AppCodec().GetCdc().MustMarshalJSON(evmGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)
	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
			ChainId:         chainId,
		},
	)
	ctx := app.GetBaseApp().NewContext(false, abci.Header{Height: 1, Time: time.Now()})
	app.GetStakingKeeper().SetParams(ctx, staking.DefaultDposParams())
	// commit genesis changes
	app.Commit(abci.RequestCommit{})
	// app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(app.LastBlockHeight() + 1),
		NextValidatorsHash: valSet.Hash(app.LastBlockHeight() + 1),
	}}) //app.Commit(abci.RequestCommit{})

	return app
}
