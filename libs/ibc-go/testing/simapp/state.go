package simapp

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	//	tmjson "github.com/okex/exchain/libs/tendermint/libs/json"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	//stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"

	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	stakingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	//banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	banktypes "github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	//	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	//simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	simtypes "github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	simappparams "github.com/okex/exchain/libs/ibc-go/testing/simapp/params"
)

// AppStateFn returns the initial application state using a genesis or the simulation parameters.
// It panics if the user provides files for both of them.
// If a file is not given for the genesis or the sim params, it creates a randomized one.
func AppStateFn(cdc codec.CodecProxy, simManager *module.SimulationManager) simtypes.AppStateFn {
	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config,
	) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {

		if FlagGenesisTimeValue == 0 {
			genesisTimestamp = simtypes.RandTimestamp(r)
		} else {
			genesisTimestamp = time.Unix(FlagGenesisTimeValue, 0)
		}

		chainID = config.ChainID
		switch {
		case config.ParamsFile != "" && config.GenesisFile != "":
			panic("cannot provide both a genesis file and a params file")

		case config.GenesisFile != "":
			// override the default chain-id from simapp to set it later to the config
			genesisDoc, accounts := AppStateFromGenesisFileFn(r, cdc, config.GenesisFile)

			if FlagGenesisTimeValue == 0 {
				// use genesis timestamp if no custom timestamp is provided (i.e no random timestamp)
				genesisTimestamp = genesisDoc.GenesisTime
			}

			appState = genesisDoc.AppState
			chainID = genesisDoc.ChainID
			simAccs = accounts

		case config.ParamsFile != "":
			appParams := make(simtypes.AppParams)
			bz, err := ioutil.ReadFile(config.ParamsFile)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bz, &appParams)
			if err != nil {
				panic(err)
			}
			appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)

		default:
			appParams := make(simtypes.AppParams)
			appState, simAccs = AppStateRandomizedFn(simManager, r, cdc, accs, genesisTimestamp, appParams)
		}

		rawState := make(map[string]json.RawMessage)
		err := json.Unmarshal(appState, &rawState)
		if err != nil {
			panic(err)
		}

		stakingStateBz, ok := rawState[stakingtypes.ModuleName]
		if !ok {
			panic("staking genesis state is missing")
		}

		stakingState := new(stakingtypes.GenesisState)
		err = cdc.GetCdc().UnmarshalJSON(stakingStateBz, stakingState)
		if err != nil {
			panic(err)
		}
		// compute not bonded balance
		// notBondedTokens := sdk.ZeroInt()
		// for _, val := range stakingState.Validators {
		//	if val.Status != sdk.Unbonded {
		//		continue
		//	}
		//	notBondedTokens = notBondedTokens.Add(val.GetTokens())
		// }
		//	notBondedCoins := sdk.NewCoin(stakingState.Params.BondDenom, notBondedTokens)
		// edit bank state to make it have the not bonded pool tokens
		bankStateBz, ok := rawState[banktypes.ModuleName]
		// TODO(fdymylja/jonathan): should we panic in this case
		if !ok {
			panic("bank genesis state is missing")
		}
		bankState := banktypes.NewGenesisState(true)
		//		bankState := new(banktypes.GenesisState)
		err = cdc.GetCdc().UnmarshalJSON(bankStateBz, bankState)
		if err != nil {
			panic(err)
		}
		// todo our genesis bank do not have balance infos
		//		bankState.Balances = append(bankState.Balances, banktypes.Balance{
		//			Address: authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String(),
		//			Coins:   sdk.NewCoins(notBondedCoins),
		//		})

		// change appState back
		rawState[stakingtypes.ModuleName] = cdc.GetCdc().MustMarshalJSON(stakingState)
		rawState[banktypes.ModuleName] = cdc.GetCdc().MustMarshalJSON(bankState)

		// replace appstate
		appState, err = json.Marshal(rawState)
		if err != nil {
			panic(err)
		}
		return appState, simAccs, chainID, genesisTimestamp
	}
}

// AppStateRandomizedFn creates calls each module's GenesisState generator function
// and creates the simulation params
func AppStateRandomizedFn(
	simManager *module.SimulationManager, r *rand.Rand, cdc codec.CodecProxy,
	accs []simtypes.Account, genesisTimestamp time.Time, appParams simtypes.AppParams,
) (json.RawMessage, []simtypes.Account) {
	numAccs := int64(len(accs))
	genesisState := NewDefaultGenesisState(cdc.GetProtocMarshal())

	// generate a random amount of initial stake coins and a random initial
	// number of bonded accounts
	var initialStake, numInitiallyBonded int64
	appParams.GetOrGenerate(
		cdc.GetCdc(), simappparams.StakePerAccount, &initialStake, r,
		func(r *rand.Rand) { initialStake = r.Int63n(1e12) },
	)
	appParams.GetOrGenerate(
		cdc.GetCdc(), simappparams.InitiallyBondedValidators, &numInitiallyBonded, r,
		func(r *rand.Rand) { numInitiallyBonded = int64(r.Intn(300)) },
	)

	if numInitiallyBonded > numAccs {
		numInitiallyBonded = numAccs
	}

	fmt.Printf(
		`Selected randomly generated parameters for simulated genesis:
{
  stake_per_account: "%d",
  initially_bonded_validators: "%d"
}
`, initialStake, numInitiallyBonded,
	)

	simState := &module.SimulationState{
		AppParams:    appParams,
		Cdc:          cdc.GetCdc(),
		Rand:         r,
		GenState:     genesisState,
		Accounts:     accs,
		InitialStake: initialStake,
		NumBonded:    numInitiallyBonded,
		GenTimestamp: genesisTimestamp,
	}

	simManager.GenerateGenesisStates(simState)

	appState, err := json.Marshal(genesisState)
	if err != nil {
		panic(err)
	}

	return appState, accs
}

// AppStateFromGenesisFileFn util function to generate the genesis AppState
// from a genesis.json file.
func AppStateFromGenesisFileFn(r io.Reader, cdc codec.CodecProxy, genesisFile string) (tmtypes.GenesisDoc, []simtypes.Account) {
	// NOTE: Tendermint uses a custom JSON decoder for GenesisDoc
	genesis, err := tmtypes.GenesisDocFromFile(genesisFile)
	if err != nil {
		panic(err)
	}

	var appState GenesisState
	err = json.Unmarshal(genesis.AppState, &appState)
	if err != nil {
		panic(err)
	}

	var authGenesis authtypes.GenesisState
	if appState[authtypes.ModuleName] != nil {
		cdc.GetCdc().MustUnmarshalJSON(appState[authtypes.ModuleName], &authGenesis)
	}

	newAccs := make([]simtypes.Account, len(authGenesis.Accounts))
	for i, acc := range authGenesis.Accounts {
		// Pick a random private key, since we don't know the actual key
		// This should be fine as it's only used for mock Tendermint validators
		// and these keys are never actually used to sign by mock Tendermint.
		privkeySeed := make([]byte, 15)
		if _, err := r.Read(privkeySeed); err != nil {
			panic(err)
		}

		privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)

		// create simulator accounts
		simAcc := simtypes.Account{PrivKey: privKey, PubKey: privKey.PubKey(), Address: acc.GetAddress()}
		newAccs[i] = simAcc
	}

	return *genesis, newAccs
}
