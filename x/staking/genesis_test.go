package staking

import (
	"fmt"
	"testing"
	"time"

	"github.com/okex/okchain/x/common"

	"github.com/cosmos/cosmos-sdk/x/supply"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/exported"
	"github.com/okex/okchain/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

//func tokensFromTendermintPower(power int64) sdk.Int {
//	return sdk.NewInt(power).Mul(sdk.PowerReduction)
//}

func TestInitGenesis(t *testing.T) {
	ctx, _, mKeeper := CreateTestInput(t, false, 1000)
	keeper := mKeeper.Keeper
	supplyKeeper := mKeeper.SupplyKeeper
	clearNotBondedPool(t, ctx, supplyKeeper)
	valTokens := int64(1)

	params := keeper.GetParams(ctx)
	validators := make([]Validator, 2)
	// initialize the validators
	validators[0].OperatorAddress = sdk.ValAddress(Addrs[0])
	validators[0].ConsPubKey = PKs[0]
	validators[0].Description = types.NewDescription("hoop", "", "", "")
	validators[0].Status = sdk.Bonded
	validators[0].DelegatorShares = sdk.NewDec(valTokens)
	validators[0].MinSelfDelegation = sdk.OneDec()
	validators[1].OperatorAddress = sdk.ValAddress(Addrs[1])
	validators[1].ConsPubKey = PKs[1]
	validators[1].Description = types.NewDescription("bloop", "", "", "")
	validators[1].Status = sdk.Bonded
	validators[1].DelegatorShares = sdk.NewDec(valTokens)
	validators[1].MinSelfDelegation = sdk.OneDec()

	delegators := make([]Delegator, 1)
	delegators[0] = types.NewDelegator(Addrs[10])

	unbondingDelegations := make([]types.UndelegationInfo, 1)
	unbondingDelegations[0] = types.NewUndelegationInfo(Addrs[11], sdk.NewDec(100), time.Now().Add(time.Minute*10))

	proxysDelegator := make([]ProxyDelegatorKeyExported, 2)
	proxysDelegator[0].DelAddr = Addrs[2]
	proxysDelegator[0].ProxyAddr = Addrs[3]
	proxysDelegator[1].DelAddr = Addrs[4]
	proxysDelegator[1].ProxyAddr = Addrs[5]

	votes := make([]types.VotesExported, 2)
	votes[0] = types.NewVoteExported(Addrs[6], validators[0].OperatorAddress, sdk.NewDec(1))
	votes[1] = types.NewVoteExported(Addrs[7], validators[1].OperatorAddress, sdk.NewDec(1))

	genesisState := NewGenesisState(params, validators, delegators)
	genesisState.ProxyDelegatorKeys = proxysDelegator
	genesisState.Votes = votes
	genesisState.UnbondingDelegations = unbondingDelegations
	// for the token sum check in staking init-genesis
	//coinsToModuleAcc := sdk.DecCoins{sdk.NewDecCoinFromDec("okt", sdk.NewDec(100))}.ToCoins()
	coinsToModuleAcc := sdk.NewCoin(common.NativeToken, sdk.NewInt(100)).ToCoins()
	_ = supplyKeeper.SendCoinsFromAccountToModule(ctx, Addrs[11], types.NotBondedPoolName, coinsToModuleAcc)
	vals := InitGenesis(ctx, keeper, nil, supplyKeeper, genesisState)

	actualGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, genesisState.Params, actualGenesis.Params)
	require.Equal(t, genesisState.Delegators, actualGenesis.Delegators)
	require.EqualValues(t, keeper.GetAllValidators(ctx).Export(), actualGenesis.Validators)
	require.True(t, actualGenesis.Exported)
	require.Equal(t, genesisState.ProxyDelegatorKeys, actualGenesis.ProxyDelegatorKeys)

	// now make sure the validators are bonded and intra-tx counters are correct
	resVal, found := keeper.GetValidator(ctx, sdk.ValAddress(Addrs[0]))
	require.True(t, found)
	require.Equal(t, sdk.Bonded, resVal.Status)

	resVal, found = keeper.GetValidator(ctx, sdk.ValAddress(Addrs[1]))
	require.True(t, found)
	require.Equal(t, sdk.Bonded, resVal.Status)

	abcivals := make([]abci.ValidatorUpdate, len(vals))
	for i, val := range validators {
		abcivals[i] = val.ABCIValidatorUpdateByVotes()
	}
	require.EqualValues(t, abcivals, vals)

	newCtx, _, newMKeeper := CreateTestInput(t, false, 1000)
	newKeeper := newMKeeper.Keeper
	newSupplyKeeper := newMKeeper.SupplyKeeper
	clearNotBondedPool(t, newCtx, newSupplyKeeper)
	_ = newSupplyKeeper.SendCoinsFromAccountToModule(newCtx, Addrs[11], types.NotBondedPoolName, coinsToModuleAcc)
	InitGenesis(newCtx, newKeeper, nil, newSupplyKeeper, actualGenesis)
	// 0x11
	require.Equal(t, actualGenesis.LastValidatorPowers[0].Power, newKeeper.GetLastValidatorPower(newCtx, actualGenesis.LastValidatorPowers[0].Address))
	require.Equal(t, actualGenesis.LastValidatorPowers[1].Power, newKeeper.GetLastValidatorPower(newCtx, actualGenesis.LastValidatorPowers[1].Address))
	// 0x12
	totalPower := sdk.ZeroInt()
	for _, v := range actualGenesis.LastValidatorPowers {
		totalPower = totalPower.AddRaw(v.Power)
	}
	require.Equal(t, totalPower, newKeeper.GetLastTotalPower(newCtx))
	// 0x21
	resVal, found = newKeeper.GetValidator(newCtx, actualGenesis.Validators[0].OperatorAddress)
	require.True(t, found)
	require.Equal(t, actualGenesis.Validators[0].Import(), resVal)
	resVal, found = newKeeper.GetValidator(newCtx, actualGenesis.Validators[1].OperatorAddress)
	require.True(t, found)
	require.Equal(t, actualGenesis.Validators[1].Import(), resVal)
	// 0x22
	resVal, found = newKeeper.GetValidatorByConsAddr(newCtx,
		sdk.GetConsAddress(sdk.MustGetConsPubKeyBech32(actualGenesis.Validators[0].ConsPubKey)))
	require.True(t, found)
	require.Equal(t, actualGenesis.Validators[0].Import(), resVal)
	resVal, found = newKeeper.GetValidatorByConsAddr(newCtx,
		sdk.GetConsAddress(sdk.MustGetConsPubKeyBech32(actualGenesis.Validators[1].ConsPubKey)))
	require.True(t, found)
	require.Equal(t, actualGenesis.Validators[1].Import(), resVal)
	// 0x23
	newKeeper.IterateBondedValidatorsByPower(newCtx, func(index int64, validator exported.ValidatorI) (stop bool) {
		require.Equal(t, actualGenesis.Validators[index].Import(), validator)
		return false
	})
	// 0x51
	require.Equal(t, types.NewVoteResponse(actualGenesis.Votes[0].VoterAddress, actualGenesis.Votes[0].Votes),
		newKeeper.GetValidatorVotes(newCtx, actualGenesis.Votes[0].ValidatorAddress)[0])
	require.Equal(t, types.NewVoteResponse(actualGenesis.Votes[1].VoterAddress, actualGenesis.Votes[1].Votes),
		newKeeper.GetValidatorVotes(newCtx, actualGenesis.Votes[1].ValidatorAddress)[0])
	// 0x52
	delegator, ok := newKeeper.GetDelegator(newCtx, actualGenesis.Delegators[0].DelegatorAddress)
	require.True(t, ok)
	require.Equal(t, actualGenesis.Delegators[0], delegator)
	// 0x53
	unbondingDelegator, ok := newKeeper.GetUndelegating(newCtx, actualGenesis.UnbondingDelegations[0].DelegatorAddress)
	require.True(t, ok)
	require.Equal(t, actualGenesis.UnbondingDelegations[0], unbondingDelegator)
	// 0x54
	newKeeper.IterateKeysBeforeCurrentTime(newCtx, time.Now().Add(time.Hour),
		func(index int64, key []byte) (stop bool) {
			oldTime, delAddr := types.SplitCompleteTimeWithAddrKey(key)
			require.Equal(t, actualGenesis.UnbondingDelegations[index].CompletionTime, oldTime)
			require.Equal(t, actualGenesis.UnbondingDelegations[index].DelegatorAddress, delAddr)
			return false
		})
	// 0x55
	newKeeper.IterateProxy(newCtx, nil, false,
		func(index int64, delAddr, proxyAddr sdk.AccAddress) (stop bool) {
			require.Equal(t, actualGenesis.ProxyDelegatorKeys[index], types.NewProxyDelegatorKeyExported(delAddr, proxyAddr))
			return false
		})
	// 0x60
	// Will never happen in InitGenesis

	exportGenesis := ExportGenesis(newCtx, newKeeper)
	require.Equal(t, actualGenesis.Params, exportGenesis.Params)
	require.Equal(t, actualGenesis.Delegators, exportGenesis.Delegators)
	require.EqualValues(t, newKeeper.GetAllValidators(newCtx).Export(), exportGenesis.Validators)
	require.True(t, exportGenesis.Exported)
	require.Equal(t, actualGenesis.ProxyDelegatorKeys, exportGenesis.ProxyDelegatorKeys)

	exportGenesis.Validators[0].UnbondingCompletionTime = time.Now()
	exportGenesis.Validators[0].Status = sdk.Unbonding
	newCtx, _, newMKeeper = CreateTestInput(t, false, 1000)
	newKeeper = newMKeeper.Keeper
	newSupplyKeeper = newMKeeper.SupplyKeeper
	clearNotBondedPool(t, newCtx, newSupplyKeeper)
	_ = newSupplyKeeper.SendCoinsFromAccountToModule(newCtx, Addrs[11], types.NotBondedPoolName, coinsToModuleAcc)
	InitGenesis(newCtx, newKeeper, nil, newSupplyKeeper, exportGenesis)
	// 0x43
	require.Equal(t, []sdk.ValAddress{exportGenesis.Validators[0].OperatorAddress}, newKeeper.GetValidatorQueueTimeSlice(newCtx, exportGenesis.Validators[0].UnbondingCompletionTime))
}

func clearNotBondedPool(t *testing.T, ctx sdk.Context, supplyKeeper supply.Keeper) {
	notBondedPool := supplyKeeper.GetModuleAccount(ctx, types.NotBondedPoolName)
	zeroCoins := sdk.NewCoins(sdk.NewInt64Coin(common.NativeToken, 0))
	require.NoError(t, notBondedPool.SetCoins(zeroCoins))
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
}

func TestInitGenesisLargeValidatorSet(t *testing.T) {
	size := 200
	require.True(t, size > 100)

	ctx, _, mKeeper := CreateTestInput(t, false, 1000)
	keeper := mKeeper.Keeper
	supplyKeeper := mKeeper.SupplyKeeper
	params := keeper.GetParams(ctx)
	delegators := []Delegator{}
	validators := make([]Validator, size)

	for i := range validators {
		validators[i] = NewValidator(sdk.ValAddress(Addrs[i]),
			PKs[i], NewDescription(fmt.Sprintf("#%d", i), "", "", ""))

		validators[i].Status = sdk.Bonded

		tokens := int64(1)
		if i < 100 {
			tokens = int64(2)
		}
		//validators[i].Tokens = tokens
		validators[i].DelegatorShares = sdk.NewDec(tokens)
	}

	genesisState := NewGenesisState(params, validators, delegators)
	vals := InitGenesis(ctx, keeper, nil, supplyKeeper, genesisState)

	require.EqualValues(t, len(vals), keeper.GetParams(ctx).MaxValidators)
}

func TestValidateGenesis(t *testing.T) {
	genValidators := make([]Validator, 1, 5)
	pk := ed25519.GenPrivKey().PubKey()
	generatedValDesc := NewDescription("", "", "", "")
	genValidators[0] = NewValidator(sdk.ValAddress(pk.Address()), pk, generatedValDesc)
	genValidators[0].DelegatorShares = sdk.OneDec()
	genValidatorExported := genValidators[0].Export()

	tests := []struct {
		name    string
		mutate  func(*GenesisState)
		wantErr bool
	}{
		{"default", func(*GenesisState) {}, false},
		// validate genesis validators
		{"duplicate validator", func(data *types.GenesisState) {
			(*data).Validators = []ValidatorExport{genValidatorExported, genValidatorExported}
		}, true},
		{"no delegator shares", func(data *types.GenesisState) {
			(*data).Validators = []ValidatorExport{genValidatorExported}
			(*data).Validators[0].DelegatorShares = sdk.ZeroDec()
		}, true},
		{"jailed and bonded validator", func(data *types.GenesisState) {
			(*data).Validators = []ValidatorExport{genValidatorExported}
			(*data).Validators[0].Jailed = true
			(*data).Validators[0].Status = sdk.Bonded
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genesisState := types.DefaultGenesisState()
			tt.mutate(&genesisState)
			if tt.wantErr {
				assert.Error(t, ValidateGenesis(genesisState))
			} else {
				assert.NoError(t, ValidateGenesis(genesisState))
			}
		})
	}
}
