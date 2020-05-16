package v0_9

import (
	"fmt"

	v09dex "github.com/okex/okchain/x/dex/legacy/v0_9"
	v034distr "github.com/okex/okchain/x/distribution/legacy/v0_34"
	v036distr "github.com/okex/okchain/x/distribution/legacy/v0_36"
	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
	v09gov "github.com/okex/okchain/x/gov/legacy/v0_9"
	v08order "github.com/okex/okchain/x/order/legacy/v0_8"
	v09order "github.com/okex/okchain/x/order/legacy/v0_9"
	v09params "github.com/okex/okchain/x/params/legacy/v0_9"
	v034staking "github.com/okex/okchain/x/staking/legacy/v0_34"
	v036staking "github.com/okex/okchain/x/staking/legacy/v0_36"
	v08token "github.com/okex/okchain/x/token/legacy/v0_8"
	v09token "github.com/okex/okchain/x/token/legacy/v0_9"
	v08upgrade "github.com/okex/okchain/x/upgrade/legacy/v0_8"
	v09upgrade "github.com/okex/okchain/x/upgrade/legacy/v0_9"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v034auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v0_34"
	v036auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v0_36"
	v036crisis "github.com/cosmos/cosmos-sdk/x/crisis"
	v034genAccounts "github.com/cosmos/cosmos-sdk/x/genaccounts/legacy/v0_34"
	v036genAccounts "github.com/cosmos/cosmos-sdk/x/genaccounts/legacy/v0_36"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	v034gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v0_34"
	v036mint "github.com/cosmos/cosmos-sdk/x/mint"
	v036supply "github.com/cosmos/cosmos-sdk/x/supply/legacy/v0_36"
	"github.com/tendermint/tendermint/crypto"
)

const (
	notBondedPoolName = "not_bonded_tokens_pool"
	bondedPoolName    = "bonded_tokens_pool"
	feeCollectorName  = "fee_collector"
	mintModuleName    = "mint"

	basic   = "basic"
	minter  = "minter"
	burner  = "burner"
	staking = "staking"
)

// Migrate migrates exported state from v0.34 to a v0.36 genesis state
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v08Codec := codec.New()
	codec.RegisterCrypto(v08Codec)

	v09Codec := codec.New()
	codec.RegisterCrypto(v09Codec)
	v08gov.RegisterCodec(v08Codec)
	v09gov.RegisterCodec(v09Codec)

	// migrate genesis accounts state
	if appState[v034genAccounts.ModuleName] != nil {
		var genAccs v034genAccounts.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v034genAccounts.ModuleName], &genAccs)

		var authGenState v034auth.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v034auth.ModuleName], &authGenState)

		var govGenState v08gov.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v08gov.ModuleName], &govGenState)

		var distrGenState v034distr.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v034distr.ModuleName], &distrGenState)

		var stakingGenState v034staking.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v034staking.ModuleName], &stakingGenState)

		var tokenGenState v08token.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v08token.ModuleName], &tokenGenState)

		delete(appState, v034genAccounts.ModuleName) // delete old key in case the name changed
		var deposits []v034gov.DepositWithMetadata
		for _, deposit := range govGenState.Deposits {
			deposits = append(deposits, v034gov.DepositWithMetadata{
				ProposalID: deposit.ProposalID,
				Deposit: v034gov.Deposit{
					ProposalID: deposit.ProposalID,
					Depositor:  deposit.Depositor,
					Amount:     deposit.Amount,
				},
			})
		}
		lockCoins := sdk.DecCoins{}
		for _, lock := range tokenGenState.LockCoins {
			lockCoins = lockCoins.Add(lock.Coins)
		}

		appState[v036genAccounts.ModuleName] = v09Codec.MustMarshalJSON(
			genAccountMigrate(
				genAccs, authGenState.CollectedFees, distrGenState.FeePool.CommunityPool, lockCoins, deposits,
				stakingGenState.Validators, stakingGenState.UnbondingDelegations, distrGenState.OutstandingRewards,
				stakingGenState.Params.BondDenom, v036distr.ModuleName, v09gov.ModuleName, v09token.ModuleName,
			),
		)
	}

	// migrate auth state
	if appState[v034auth.ModuleName] != nil {
		var authGenState v034auth.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v034auth.ModuleName], &authGenState)

		delete(appState, v034auth.ModuleName) // delete old key in case the name changed
		appState[v036auth.ModuleName] = v09Codec.MustMarshalJSON(v036auth.Migrate(authGenState))
	}

	// migrate gov state
	var govGenState v08gov.GenesisState
	if appState[v08gov.ModuleName] != nil {
		v08Codec.MustUnmarshalJSON(appState[v08gov.ModuleName], &govGenState)

		delete(appState, v08gov.ModuleName) // delete old key in case the name changed
		appState[v09gov.ModuleName] = v09Codec.MustMarshalJSON(v09gov.Migrate(govGenState))
	}

	// migrate params state
	appState[v09params.ModuleName] = v09Codec.MustMarshalJSON(v09params.Migrate(govGenState.Params))

	// migrate upgrade state
	if appState[v08upgrade.ModuleName] != nil {
		var upgradeGenState v08upgrade.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v08upgrade.ModuleName], &upgradeGenState)

		delete(appState, v08upgrade.ModuleName) // delete old key in case the name changed
		appState[v09upgrade.ModuleName] = v09Codec.MustMarshalJSON(v09upgrade.Migrate(upgradeGenState, govGenState.Params))
	}

	// migrate token state and dex state
	if appState[v08token.ModuleName] != nil {
		var tokenGenState v08token.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v08token.ModuleName], &tokenGenState)

		delete(appState, v08token.ModuleName) // delete old key in case the name changed
		appState[v09token.ModuleName] = v09Codec.MustMarshalJSON(v09token.Migrate(tokenGenState, govGenState.Params))
		appState[v09dex.ModuleName] = v09Codec.MustMarshalJSON(v09dex.Migrate(tokenGenState, govGenState.Params))
	}

	// migrate order state
	if appState[v08order.ModuleName] != nil {
		var orderGenState v08order.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v08order.ModuleName], &orderGenState)

		delete(appState, v08order.ModuleName) // delete old key in case the name changed
		appState[v09order.ModuleName] = v09Codec.MustMarshalJSON(v09order.Migrate(orderGenState))
	}

	// migrate distribution state
	if appState[v034distr.ModuleName] != nil {
		var distrGenState v034distr.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v034distr.ModuleName], &distrGenState)

		delete(appState, v034distr.ModuleName) // delete old key in case the name changed
		appState[v036distr.ModuleName] = v09Codec.MustMarshalJSON(v036distr.Migrate(distrGenState))
	}

	// migrate staking state
	if appState[v034staking.ModuleName] != nil {
		var stakingGenState v034staking.GenesisState
		v08Codec.MustUnmarshalJSON(appState[v034staking.ModuleName], &stakingGenState)

		delete(appState, v034staking.ModuleName) // delete old key in case the name changed
		appState[v036staking.ModuleName] = v09Codec.MustMarshalJSON(v036staking.Migrate(stakingGenState))
	}

	// migrate supply state
	appState[v036supply.ModuleName] = v09Codec.MustMarshalJSON(v036supply.EmptyGenesisState())
	// migrate crisis state
	appState[v036crisis.ModuleName] = v09Codec.MustMarshalJSON(v036crisis.DefaultGenesisState())

	appState[v036mint.ModuleName] = v09Codec.MustMarshalJSON(v036mint.DefaultGenesisState())

	return appState
}

// Migrate accepts exported genesis state from v0.34 and migrates it to v0.36 genesis state
// It deletes the governance base accounts and creates the new module accounts
// The remaining accounts are updated to the new GenesisAccount type from 0.36
func genAccountMigrate(
	oldGenState v034genAccounts.GenesisState, fees sdk.Coins, communityPool sdk.DecCoins, lockCoins sdk.Coins,
	deposits []v034gov.DepositWithMetadata, vals v034staking.Validators, ubds []v034staking.UnbondingDelegation,
	valOutRewards []v034distr.ValidatorOutstandingRewardsRecord, bondDenom, distrModuleName, govModuleName,
	tokenModuleName string) v036genAccounts.GenesisState {

	depositedCoinsAccAddr := sdk.AccAddress(crypto.AddressHash([]byte("govDepositedCoins")))
	burnedDepositCoinsAccAddr := sdk.AccAddress(crypto.AddressHash([]byte("govBurnedDepositCoins")))
	dexListDepositedCoinsAccAddr := sdk.AccAddress(crypto.AddressHash([]byte("govDexListDepositedCoins")))

	bondedAmt := sdk.ZeroInt()
	notBondedAmt := sdk.ZeroInt()

	// remove the two previous governance base accounts for deposits and burned
	// coins from rejected proposals add six new module accounts:
	// distribution, gov, mint, fee collector, bonded and not bonded pool
	var (
		newGenState      v036genAccounts.GenesisState
		govCoins         sdk.Coins
		dexListDeposited = sdk.Coins{}
		extraAccounts    = 7
	)

	for _, acc := range oldGenState {
		switch {
		case acc.Address.Equals(depositedCoinsAccAddr):
			// remove gov deposits base account
			govCoins = acc.Coins
			extraAccounts--

		case acc.Address.Equals(burnedDepositCoinsAccAddr):
			// remove gov burned deposits base account
			extraAccounts--

		case acc.Address.Equals(dexListDepositedCoinsAccAddr):
			// remove gov dexList deposits base account
			dexListDeposited = acc.Coins
			extraAccounts--

		default:
			newGenState = append(
				newGenState,
				v036genAccounts.NewGenesisAccount(
					acc.Address, acc.Coins, acc.Sequence,
					acc.OriginalVesting, acc.DelegatedFree, acc.DelegatedVesting,
					acc.StartTime, acc.EndTime, "", []string{},
				),
			)
		}
	}

	var expDeposits sdk.Coins
	for _, deposit := range deposits {
		expDeposits = expDeposits.Add(deposit.Deposit.Amount)
	}

	if !expDeposits.IsEqual(govCoins) {
		panic(
			fmt.Sprintf(
				"pre migration deposit base account coins ≠ stored deposits coins (%s ≠ %s)",
				expDeposits.String(), govCoins.String(),
			),
		)
	}

	// get staking module accounts coins
	for _, validator := range vals {
		switch validator.Status {
		case sdk.Bonded:
			bondedAmt = bondedAmt.Add(validator.MinSelfDelegation)

		case sdk.Unbonding, sdk.Unbonded:
			notBondedAmt = notBondedAmt.Add(validator.Tokens)

		default:
			panic("invalid validator status")
		}
	}

	for _, ubd := range ubds {
		for _, entry := range ubd.Entries {
			notBondedAmt = notBondedAmt.Add(entry.Balance)
		}
	}
	//bondedCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, bondedAmt))
	//bondedCoins := sdk.NewDecCoins(sdk.NewDecCoinsFromDec(bondDenom, sdk.NewDecFromBigInt(bondedAmt.BigInt())))
	bondedCoins := sdk.NewDecCoins(sdk.NewDecCoinsFromDec(bondDenom, sdk.NewDecFromBigIntWithPrec(bondedAmt.BigInt(), 3)))
	notBondedCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, notBondedAmt))

	// get distr module account coins
	var distrDecCoins sdk.DecCoins
	for _, reward := range valOutRewards {
		distrDecCoins = distrDecCoins.Add(reward.OutstandingRewards)
	}

	distrCoins, _ := distrDecCoins.Add(communityPool).TruncateDecimal()

	// get module account addresses
	feeCollectorAddr := sdk.AccAddress(crypto.AddressHash([]byte(feeCollectorName)))
	govAddr := sdk.AccAddress(crypto.AddressHash([]byte(govModuleName)))
	bondedAddr := sdk.AccAddress(crypto.AddressHash([]byte(bondedPoolName)))
	notBondedAddr := sdk.AccAddress(crypto.AddressHash([]byte(notBondedPoolName)))
	distrAddr := sdk.AccAddress(crypto.AddressHash([]byte(distrModuleName)))
	mintAddr := sdk.AccAddress(crypto.AddressHash([]byte(mintModuleName)))
	tokenAddr := sdk.AccAddress(crypto.AddressHash([]byte(tokenModuleName)))

	// create module genesis accounts
	feeCollectorModuleAcc := v036genAccounts.NewGenesisAccount(
		feeCollectorAddr, fees, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, feeCollectorName, []string{basic},
	)
	govModuleAcc := v036genAccounts.NewGenesisAccount(
		govAddr, govCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, govModuleName, []string{burner},
	)
	distrModuleAcc := v036genAccounts.NewGenesisAccount(
		distrAddr, distrCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, distrModuleName, []string{basic},
	)
	bondedModuleAcc := v036genAccounts.NewGenesisAccount(
		bondedAddr, bondedCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, bondedPoolName, []string{burner, staking},
	)
	notBondedModuleAcc := v036genAccounts.NewGenesisAccount(
		notBondedAddr, notBondedCoins, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, notBondedPoolName, []string{burner, staking},
	)
	mintModuleAcc := v036genAccounts.NewGenesisAccount(
		mintAddr, sdk.Coins{}, 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, mintModuleName, []string{minter},
	)
	tokenModuleAcc := v036genAccounts.NewGenesisAccount(
		tokenAddr, dexListDeposited.Add(lockCoins), 0,
		sdk.Coins{}, sdk.Coins{}, sdk.Coins{},
		0, 0, tokenModuleName, []string{minter, burner},
	)

	newGenState = append(
		newGenState,
		[]v036genAccounts.GenesisAccount{
			feeCollectorModuleAcc, govModuleAcc, distrModuleAcc,
			bondedModuleAcc, notBondedModuleAcc, mintModuleAcc, tokenModuleAcc,
		}...,
	)

	// verify the total number of accounts is correct
	if len(newGenState) != len(oldGenState)+extraAccounts {
		panic(
			fmt.Sprintf(
				"invalid total number of genesis accounts; got: %d, expected: %d",
				len(newGenState), len(oldGenState)+extraAccounts),
		)
	}

	return newGenState
}
