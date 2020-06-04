package staking

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/keeper"
	"github.com/okex/okchain/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// dummy addresses used for testing
var (
	Addrs = keeper.Addrs
	PKs   = keeper.PKs

	addrDels = []sdk.AccAddress{
		Addrs[0],
		Addrs[1],
		Addrs[2],
		Addrs[3],
	}

	addrVals = []sdk.ValAddress{
		sdk.ValAddress(Addrs[4]),
		sdk.ValAddress(Addrs[5]),
		sdk.ValAddress(Addrs[6]),
		sdk.ValAddress(Addrs[7]),
		sdk.ValAddress(Addrs[8]),
	}

	StartUpValidatorAddr   = addrVals[0]
	StartUpValidatorPubkey = PKs[0]

	MostPowerfulVaAddr = addrVals[len(addrVals)-1]
	MostPowerfulVaPub  = PKs[len(PKs)-1]

	InvalidDelegator    = addrDels[0]
	ValidDelegator1     = addrDels[1]
	ValidDelegator2     = addrDels[2]
	ProxiedDelegator    = addrDels[3]
	SufficientInitPower = int64(10000)
	MaxDelegatedToken   = sdk.NewDec(4096)
	DefaultMSD          = types.DefaultMinSelfDelegation
	VotesFromDefaultMSD = sdk.OneDec()
	DelegatedToken1     = VotesFromDefaultMSD.MulInt64(1024)
	DelegatedToken2     = VotesFromDefaultMSD.MulInt64(2048)
)

var (
	CreateTestInput             = keeper.CreateTestInput
	ValidatorByPowerIndexExists = keeper.ValidatorByPowerIndexExists
	NewTestMsgCreateValidator   = keeper.NewTestMsgCreateValidator
	SimpleCheckValidator        = keeper.SimpleCheckValidator
)

// --------------------------------------------------------------
// Test Interfaces of Validator State Machine

// IValidatorStatus shows the action of validator status
type IValidatorStatus interface {
	getValidator() Validator
	getStatus() string
	desc() string
}

// IAction shows the action of a role in staking test
type IAction interface {
	apply(ctx sdk.Context, vaStatus IValidatorStatus, result *ActionResultCtx)
}

type ActionResultCtx struct {
	txMsgResult      *sdk.Result
	errorResult      error
	endBlockerResult abci.ValidatorUpdates
	isBlkHeightInc   bool
	params           types.Params
	tc               *basicStakingSMTestCase
	t                *testing.T
}

func (ar *ActionResultCtx) String() string {
	return fmt.Sprintf("TxMsgResult: %+v, ErrorResult: %+v, EndBlockResult: %+v",
		ar.txMsgResult, ar.errorResult, ar.endBlockerResult)
}

type actResChecker func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, realRes *ActionResultCtx) bool

// baseValidatorStatus is an implementation of Validator State Machine
type baseValidatorStatus struct {
	va Validator
}

func (s baseValidatorStatus) getValidator() Validator {
	return s.va
}

func (s baseValidatorStatus) getStatus() string {
	return s.getValidator().GetStatus().String()
}

func (s baseValidatorStatus) desc() string {
	if s.getValidator().GetOperator() != nil {
		return s.getValidator().String()
	}
	return "Validator's destroyed or not initialized"
}

type baseAction struct {
	mStkKeeper keeper.MockStakingKeeper
}

type createValidatorAction struct {
	baseAction
	newVal IValidatorStatus
}

func (a createValidatorAction) apply(ctx sdk.Context, expVaStatus IValidatorStatus, resultCtx *ActionResultCtx) {

	vaStatus := expVaStatus
	if a.newVal != nil {
		vaStatus = a.newVal
	}

	val := vaStatus.getValidator()
	resultCtx.t.Logf("====> Apply createValidatorAction[%d], addr:%s, msd: %s, maxVA: %d\n",
		ctx.BlockHeight(), val.OperatorAddress, val.MinSelfDelegation, resultCtx.params.MaxValidators)

	msgCreateValidator := NewTestMsgCreateValidator(val.OperatorAddress, val.ConsPubKey, val.MinSelfDelegation)
	if err := msgCreateValidator.ValidateBasic(); err != nil {
		panic(err)
	}
	handler := NewHandler(a.mStkKeeper.Keeper)

	msgResponse := handler(ctx, msgCreateValidator)

	if resultCtx != nil {
		resultCtx.txMsgResult = &msgResponse
		resultCtx.isBlkHeightInc = false

		validator, found := resultCtx.tc.mockKeeper.Keeper.GetValidator(ctx, val.OperatorAddress)
		if !found {
			panic("failed to create a validator")
		}
		resultCtx.t.Logf("     ==>>> CreateValidator Result: %s msd: %s, votes: %s\n",
			validator.OperatorAddress, validator.MinSelfDelegation, validator.DelegatorShares)
	}

}

type otherMostPowerfulValidatorEnter struct {
	baseAction
}

func (a otherMostPowerfulValidatorEnter) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {

	val := vaStatus.getValidator()

	resultCtx.t.Logf("====> Apply otherMostPowerfulValidatorEnter[%d], msd: %s\n",
		ctx.BlockHeight(), val.MinSelfDelegation)

	newValidator := NewValidator(MostPowerfulVaAddr, MostPowerfulVaPub, Description{})

	newVaStatus := baseValidatorStatus{newValidator}
	cva := createValidatorAction{a.baseAction, nil}
	cva.apply(ctx, newVaStatus, resultCtx)

	// increase the voting power by voting
	handler := NewHandler(resultCtx.tc.mockKeeper.Keeper)
	handler(ctx, NewMsgVote(ValidDelegator2, []sdk.ValAddress{newValidator.OperatorAddress}))

	// get the info of powerful validator
	validator, found := resultCtx.tc.mockKeeper.Keeper.GetValidator(ctx, newValidator.OperatorAddress)
	if !found {
		panic("failed to create a validator")
	}
	resultCtx.t.Logf("     ==>>> OtherMostPowerfulValidatorEnter Result: %s msd: %s, votes: %s\n",
		validator.OperatorAddress, validator.MinSelfDelegation, validator.DelegatorShares)

}

type destroyValidatorAction struct {
	baseAction
}

func (a destroyValidatorAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	val := vaStatus.getValidator()
	resultCtx.t.Logf("====> Apply destroyValidatorAction[%d], msd: %s\n",
		ctx.BlockHeight(), val.MinSelfDelegation)

	msgDestroyValidator := types.NewMsgDestroyValidator(val.OperatorAddress.Bytes())
	if err := msgDestroyValidator.ValidateBasic(); err != nil {
		panic(err)
	}
	handler := NewHandler(a.mStkKeeper.Keeper)

	msgResponse := handler(ctx, msgDestroyValidator)

	if resultCtx != nil {
		resultCtx.txMsgResult = &msgResponse
		resultCtx.isBlkHeightInc = false

		validator, found := resultCtx.tc.mockKeeper.Keeper.GetValidator(ctx, val.OperatorAddress)
		if !found {
			panic("validator is missing")
		}
		resultCtx.t.Logf("     ==>>> DestroyValidator Result: %s msd: %s, votes: %s\n",
			validator.OperatorAddress, validator.MinSelfDelegation, validator.DelegatorShares)
	}
}

type jailValidatorAction struct {
	baseAction
}

func (a jailValidatorAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	val := vaStatus.getValidator()
	resultCtx.t.Logf("====> Apply jailValidatorAction[%d], msd: %s\n",
		ctx.BlockHeight(), val.MinSelfDelegation)

	// No Response here
	a.mStkKeeper.Jail(ctx, val.GetConsAddr())
	a.mStkKeeper.AppendAbandonedValidatorAddrs(ctx, val.GetConsAddr())
	if resultCtx != nil {
		resultCtx.isBlkHeightInc = false

		validator, found := resultCtx.tc.mockKeeper.Keeper.GetValidator(ctx, val.OperatorAddress)
		if !found {
			panic("validator is missing")
		}
		resultCtx.t.Logf("     ==>>> JailValidator Result: %s msd: %s, votes: %s\n",
			validator.OperatorAddress, validator.MinSelfDelegation, validator.DelegatorShares)
	}
}

type unJailValidatorAction struct {
	baseAction
}

func (a unJailValidatorAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	val := vaStatus.getValidator()
	resultCtx.t.Logf("====> Apply unJailValidatorAction[%d], msd: %s\n",
		ctx.BlockHeight(), val.MinSelfDelegation)

	a.mStkKeeper.Unjail(ctx, val.GetConsAddr())
	if resultCtx != nil {
		resultCtx.isBlkHeightInc = false

		validator, found := resultCtx.tc.mockKeeper.Keeper.GetValidator(ctx, val.OperatorAddress)
		if !found {
			panic("validator is missing")
		}
		resultCtx.t.Logf("     ==>>> UnJailValidator Result: %s msd: %s, votes: %s\n",
			validator.OperatorAddress, validator.MinSelfDelegation, validator.DelegatorShares)
	}

}

type waitUntilUnbondingTimeExpired struct {
	baseAction
}

func (a waitUntilUnbondingTimeExpired) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply waitUntilUnbondingTimeExpired[%d], msd: %s\n",
		ctx.BlockHeight(), vaStatus.getValidator().GetMinSelfDelegation().String())

	time.Sleep(resultCtx.params.UnbondingTime + time.Millisecond)
}

type endBlockAction struct {
	baseAction
}

func (action endBlockAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	r := EndBlocker(ctx, action.mStkKeeper.Keeper)
	if resultCtx != nil {
		resultCtx.t.Logf("====> Apply endBlockAction[%d]\n", ctx.BlockHeight())
		resultCtx.endBlockerResult = r
		resultCtx.isBlkHeightInc = true
	}
}

type newDelegatorAction struct {
	baseAction
	dlgAddr   sdk.AccAddress
	dlgAmount sdk.Dec
	dlgDenom  string
}

func (action newDelegatorAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply newDelegatorAction[%d], dlgAddr: %s, dlgAmount: %s, dlgDenon: %s\n",
		ctx.BlockHeight(), action.dlgAddr.String(), action.dlgAmount.String(), action.dlgDenom)
	handler := NewHandler(resultCtx.tc.mockKeeper.Keeper)
	coins := sdk.NewDecCoinFromDec(action.dlgDenom, action.dlgAmount)
	msgDelegate := NewMsgDelegate(action.dlgAddr, coins)
	if err := msgDelegate.ValidateBasic(); err != nil {
		panic(err)
	}

	res := handler(ctx, msgDelegate)

	newDlg, _ := resultCtx.tc.mockKeeper.Keeper.GetDelegator(ctx, action.dlgAddr)
	resultCtx.t.Logf("     ==>>> NewDelegatorInfo :%s, resOK: %+v, info: %+v \n", action.dlgAddr.String(), res.IsOK(), newDlg)
	if resultCtx != nil {
		resultCtx.txMsgResult = &res
	}
}

type delegatorsVoteAction struct {
	baseAction
	voteOnVaSet   bool
	voteOnStartup bool
	voteOnFakes   int
	delegators    []sdk.AccAddress
}

func (action delegatorsVoteAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply delegatorsVoteAction[%d]\n", ctx.BlockHeight())

	handler := NewHandler(action.mStkKeeper.Keeper)

	var vaAddrs []sdk.ValAddress
	if action.voteOnStartup {
		vaAddrs = append(vaAddrs, vaStatus.getValidator().OperatorAddress)
	}

	if action.voteOnVaSet {
		for _, v := range resultCtx.tc.originVaSet {
			vaAddrs = append(vaAddrs, v.getValidator().OperatorAddress)
		}
	}

	for i := 0; i < action.voteOnFakes; i++ {
		vaAddrs = append(vaAddrs, sdk.ValAddress(fmt.Sprintf("fakeAddr%d", i)))
	}

	if len(action.delegators) == 0 {
		for _, v := range resultCtx.tc.originDlgSet {
			action.delegators = append(action.delegators, v.DelegatorAddress)
		}
	}

	for _, d := range action.delegators {
		resultCtx.t.Logf("     ==>>> Delegator: %s vote to Validators: %s\n", d.String(), vaAddrs)
		voteMsg := NewMsgVote(d, vaAddrs)

		res := handler(ctx, voteMsg)
		if resultCtx != nil {
			resultCtx.txMsgResult = &res
		}
	}

	if resultCtx != nil {
		resultCtx.isBlkHeightInc = false
	}
}

type delegatorUnbondAction struct {
	baseAction
	dlgAddr     sdk.AccAddress
	unbondToken sdk.Dec
	tokenDenom  string
}

func (action delegatorUnbondAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply delegatorUnbondAction [%d]\n", ctx.BlockHeight())

	handler := NewHandler(action.mStkKeeper.Keeper)
	coins := sdk.NewDecCoinFromDec(action.tokenDenom, action.unbondToken)

	msg := NewMsgUndelegate(action.dlgAddr, coins)
	res := handler(ctx, msg)
	if resultCtx != nil {
		resultCtx.txMsgResult = &res
	}

	newDlg, _ := resultCtx.tc.mockKeeper.Keeper.GetDelegator(ctx, action.dlgAddr)
	resultCtx.t.Logf("     ==>>> DelegatorUnbonded Result: %s unbond: %s, resOK: %+v, info: %+v \n", msg.DelegatorAddress, coins.String(), res.IsOK(), newDlg)
}

type delegatorsUnBondAction struct {
	baseAction
	allDelegatorDoUnbound bool
	unbondAllTokens       bool
}

func (action delegatorsUnBondAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply delegatorsUnBondAction[%d]\n", ctx.BlockHeight())

	maxIdx := len(resultCtx.tc.originDlgSet) - 1
	if !action.allDelegatorDoUnbound {
		maxIdx = len(resultCtx.tc.originDlgSet)/2 - 1
	}

	counter := 0
	for _, v := range resultCtx.tc.originDlgSet {
		if counter > maxIdx {
			break
		}

		dlg, _ := action.mStkKeeper.Keeper.GetDelegator(ctx, v.DelegatorAddress)

		coins := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, dlg.Tokens)
		if !action.unbondAllTokens {
			coins = sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, dlg.Tokens.QuoInt64(2))
		}

		subAction := delegatorUnbondAction{action.baseAction,
			dlg.DelegatorAddress, coins.Amount, coins.Denom}
		subAction.apply(ctx, vaStatus, resultCtx)

		counter++
		resultCtx.tc.originDlgSet[v.DelegatorAddress.String()] = &dlg
	}

	if resultCtx != nil {
		resultCtx.isBlkHeightInc = false
	}
}

type baseProxyRegAction struct {
	baseAction
	proxyAddr sdk.AccAddress
	doReg     bool
}

func (action baseProxyRegAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply baseProxyRegAction[%d] ProxyAddress: %s, DoRegister: %+v\n",
		ctx.BlockHeight(), action.proxyAddr, action.doReg)

	handler := NewHandler(action.mStkKeeper.Keeper)
	msg := types.NewMsgRegProxy(action.proxyAddr, action.doReg)
	if err := msg.ValidateBasic(); err != nil {
		panic(err)
	}

	res := handler(ctx, msg)

	if resultCtx != nil {
		resultCtx.txMsgResult = &res
	}
}

type proxyBindAction struct {
	baseAction
	dlgAddr   sdk.AccAddress
	proxyAddr sdk.AccAddress
}

func (action proxyBindAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply proxyBindAction[%d]\n", ctx.BlockHeight())
	msg := types.NewMsgBindProxy(action.dlgAddr, action.proxyAddr)
	handler := NewHandler(action.mStkKeeper.Keeper)
	res := handler(ctx, msg)

	if resultCtx != nil {
		resultCtx.txMsgResult = &res
	}
}

type proxyUnBindAction struct {
	baseAction
	dlgAddr sdk.AccAddress
}

func (action proxyUnBindAction) apply(ctx sdk.Context, vaStatus IValidatorStatus, resultCtx *ActionResultCtx) {
	resultCtx.t.Logf("====> Apply proxyUnBindAction[%d]\n", ctx.BlockHeight())
	msg := types.NewMsgUnbindProxy(action.dlgAddr)
	handler := NewHandler(action.mStkKeeper.Keeper)
	res := handler(ctx, msg)
	if resultCtx != nil {
		resultCtx.txMsgResult = &res
	}
}

// --------------------------------------------------------------
// Implementation of actResChecker Checker

func validatorStatusChecker(expStatus string) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		return assert.NotNil(t, beforeStatus) &&
			assert.NotNil(t, afterStatus) &&
			assert.EqualValues(t, expStatus, afterStatus.getStatus(), beforeStatus.desc(), afterStatus.desc())
	}
}

func validatorDelegatorShareLeft(expectLeft bool) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		return assert.NotNil(t, afterStatus) &&
			assert.True(t, expectLeft && afterStatus.getValidator().GetDelegatorShares().GT(sdk.ZeroDec()) ||
				!expectLeft && afterStatus.getValidator().GetDelegatorShares().Equal(sdk.ZeroDec()),
				afterStatus.desc())
	}
}

func validatorKickedOff(expectKickedOff bool) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		b1 := assert.NotNil(t, beforeStatus) && assert.NotNil(t, afterStatus)
		b2 := b1 && expectKickedOff && afterStatus.getValidator().IsJailed()
		b3 := b1 && !expectKickedOff && !afterStatus.getValidator().IsJailed()

		return b2 || b3
	}
}

func validatorRemoved(expectRemoved bool) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		b1 := expectRemoved && assert.True(t, afterStatus.getValidator().GetOperator() == nil)
		b2 := !expectRemoved && assert.True(t, afterStatus.getValidator().GetOperator() != nil)
		//resultCtx.tc.printParticipantSnapshot()

		return b1 || b2
	}

}

func validatorDelegatorShareIncreased(expectIncr bool) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		b1 := assert.NotNil(t, beforeStatus) && assert.NotNil(t, afterStatus)

		beforeDS := beforeStatus.getValidator().GetDelegatorShares()
		afterDS := afterStatus.getValidator().GetDelegatorShares()

		b2 := b1 && expectIncr && assert.True(t, afterDS.GT(beforeDS),
			fmt.Sprintf("beforeDS: %s, afterDS: %s", beforeDS.String(), afterDS.String()))
		b3 := b1 && !expectIncr && assert.False(t, afterDS.GT(beforeDS),
			fmt.Sprintf("beforeDS: %s, afterDS: %s", beforeDS.String(), afterDS.String()))

		return b2 || b3
	}
}

func noErrorInHandlerResult(expectNoError bool) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		b1 := assert.NotNil(t, resultCtx.txMsgResult, resultCtx)
		if expectNoError {
			b1 = b1 && assert.True(t, resultCtx.txMsgResult.IsOK(), resultCtx.txMsgResult)
		} else {
			b1 = b1 && assert.False(t, resultCtx.txMsgResult.IsOK(), resultCtx.txMsgResult)
		}

		return b1
	}
}

// --------------------------------------------------------------
// Implementation of actResChecker Checker for queries

func queryValidatorCheck(expStatus sdk.BondStatus, expJailed bool, expDS *sdk.Dec, expMsd *sdk.Dec, expUnbdHght *int64) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		q := keeper.NewQuerier(resultCtx.tc.mockKeeper.Keeper)
		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)

		basicParams := types.NewQueryValidatorParams(afterStatus.getValidator().OperatorAddress)
		bz, _ := amino.MarshalJSON(basicParams)
		res, err := q(ctx, []string{types.QueryValidator}, abci.RequestQuery{Data: bz})
		require.Nil(t, err)

		validator := types.Validator{}
		require.NoError(t, amino.UnmarshalJSON(res, &validator))

		b1 := assert.Equal(t, validator.GetStatus(), expStatus, validator.Standardize().String())
		b2 := assert.Equal(t, validator.IsJailed(), expJailed, validator.Standardize().String())

		b3 := true
		if expDS != nil {
			b3 = assert.Equal(t, *expDS, validator.GetDelegatorShares(), validator.Standardize().String())
		}

		b4 := true
		if expMsd != nil {
			b4 = assert.Equal(t, *expMsd, validator.GetMinSelfDelegation(), validator.Standardize().String())
		}

		b5 := true
		if expUnbdHght != nil {
			b5 = assert.Equal(t, *expUnbdHght, validator.UnbondingHeight, validator.Standardize().String())
		}

		return b1 && b2 && b3 && b4 && b5
	}
}

func queryProxyCheck(proxyAddr sdk.AccAddress, isProxy bool, totalDelTokens sdk.Dec) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)
		q := keeper.NewQuerier(resultCtx.tc.mockKeeper.Keeper)

		cdc := ModuleCdc

		queryDlgParams := types.NewQueryDelegatorParams(proxyAddr)
		bz := cdc.MustMarshalJSON(queryDlgParams)
		res, sdkErr := q(ctx, []string{types.QueryDelegator}, abci.RequestQuery{Data: bz})
		if sdkErr != nil {
			panic(fmt.Sprintf("failed. Proxy %s is missing", proxyAddr))
		}

		var proxy Delegator
		cdc.MustUnmarshalJSON(res, &proxy)

		require.Equal(t, isProxy, proxy.IsProxy)
		require.True(t, totalDelTokens.Equal(proxy.TotalDelegatedTokens))
		return true
	}
}

func queryDelegatorCheck(dlgAddr sdk.AccAddress, expExist bool, expVAs []sdk.ValAddress, expShares *sdk.Dec,
	expToken *sdk.Dec, expUnbondingToken *sdk.Dec) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)
		q := keeper.NewQuerier(resultCtx.tc.mockKeeper.Keeper)

		cdc := ModuleCdc

		queryDlgParams := types.NewQueryDelegatorParams(dlgAddr)
		bz := cdc.MustMarshalJSON(queryDlgParams)
		res, sdkErr := q(ctx, []string{types.QueryDelegator}, abci.RequestQuery{Data: bz})
		found := sdkErr == nil

		b1 := assert.True(t, found == expExist)
		b2, b3, b4, b5 := true, true, true, true

		if expExist {

			dlg := Delegator{}
			_ = cdc.UnmarshalJSON(res, &dlg)

			resultCtx.tc.originDlgSet[dlgAddr.String()] = &dlg
			// check voted validators
			b2 = true
			delegatorStr := fmt.Sprintf("%+v", dlg)
			if len(expVAs) > 0 {
				cnt := 0
				for _, exp := range expVAs {
					for _, real := range dlg.ValidatorAddresses {
						if real.Equals(exp) {
							cnt++
							break
						}
					}
				}

				b2 = assert.Equal(t, len(expVAs), cnt, delegatorStr)
			}

			if expShares != nil {
				b3 = assert.Equal(t, *expShares, dlg.GetLastVotes(), delegatorStr)
			}

			if expToken != nil {
				b4 = assert.Equal(t, *expToken, dlg.Tokens, delegatorStr)
			}

		}

		if expUnbondingToken != nil {

			// query for the undelegation info
			basicParams := types.NewQueryDelegatorParams(dlgAddr)
			bz, err := cdc.MarshalJSON(basicParams)
			require.NoError(t, err)
			res, sdkErr := q(ctx, []string{types.QueryUnbondingDelegation}, abci.RequestQuery{Data: bz})
			if expUnbondingToken.Equal(sdk.ZeroDec()) && sdkErr != nil {
				b5 = assert.True(t, sdkErr.Code() == 102, sdkErr.Error())
			} else {
				require.NoError(t, sdkErr)
				unDelegationInfo := types.DefaultUndelegation()
				require.NoError(t, cdc.UnmarshalJSON(res, &unDelegationInfo))
				b5 = assert.Equal(t, *expUnbondingToken, unDelegationInfo.Quantity, unDelegationInfo.String())
			}

		}

		return b1 && b2 && b3 && b4 && b5
	}
}

// queryDelegatorProxyCheck returns the callback function for the querier if delegator proxy
// TODO: proxyVotes not found in codes.
func queryDelegatorProxyCheck(dlgAddr sdk.AccAddress, expIsProxy bool, expHasProxy bool,
	expTotalDlgTokens *sdk.Dec, expBindedToProxy *sdk.AccAddress, expBindedDelegators []sdk.AccAddress) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {

		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)

		//query delegator from keeper directly
		dlg, found := resultCtx.tc.mockKeeper.Keeper.GetDelegator(ctx, dlgAddr)
		require.True(t, found)

		b1 := assert.Equal(t, expIsProxy, dlg.IsProxy)
		b2 := assert.Equal(t, expHasProxy, dlg.HasProxy())
		b3 := true
		if expTotalDlgTokens != nil {
			b3 = assert.Equal(t, expTotalDlgTokens.String(), dlg.TotalDelegatedTokens.String())
		}

		var b4 bool
		if expBindedToProxy != nil {
			b4 = assert.Equal(t, *expBindedToProxy, dlg.ProxyAddress)
		} else {
			b4 = dlg.ProxyAddress == nil
		}

		b5 := true
		if len(expBindedDelegators) > 0 {
			q := NewQuerier(resultCtx.tc.mockKeeper.Keeper)
			para := types.NewQueryDelegatorParams(dlgAddr)
			bz, _ := types.ModuleCdc.MarshalJSON(para)
			data, err := q(ctx, []string{types.QueryProxy}, abci.RequestQuery{Data: bz})
			require.NoError(t, err)

			realProxiedDelegators := []sdk.AccAddress{}
			require.NoError(t, ModuleCdc.UnmarshalJSON(data, &realProxiedDelegators))

			b5 = assert.Equal(t, len(expBindedDelegators), len(realProxiedDelegators))
			if b5 {
				cnt := 0
				for _, e := range expBindedDelegators {
					for _, r := range realProxiedDelegators {
						if r.Equals(e) {
							cnt++
							continue
						}
					}
				}
				b5 = assert.Equal(t, len(expBindedDelegators), cnt)
			}

		}

		r := b1 && b2 && b3 && b4 && b5
		if !r {
			resultCtx.tc.printParticipantSnapshot(resultCtx.t)
		}

		return r
	}
}

func queryAllValidatorCheck(expStatuses []sdk.BondStatus, expStatusCnt []int) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {

		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)
		q := keeper.NewQuerier(resultCtx.tc.mockKeeper.Keeper)
		cdc := ModuleCdc

		//if expStatuses == nil && expStatusCnt == nil {
		//	return true
		//}

		require.True(t, len(expStatusCnt) == len(expStatuses), expStatusCnt, expStatuses)

		for i := 0; i < len(expStatuses); i++ {

			params := types.NewQueryValidatorsParams(1, 100, expStatuses[i].String())
			bz, _ := cdc.MarshalJSON(params)
			res, err := q(ctx, []string{types.QueryValidators}, abci.RequestQuery{Data: bz})

			require.Nil(t, err)
			vals := Validators{}
			require.NoError(t, cdc.UnmarshalJSON(res, &vals))
			require.Equal(t, expStatusCnt[i], len(vals))
		}

		return true
	}
}

func queryVotesToCheck(valAddr sdk.ValAddress, expVoterCnt int, expVoters []sdk.AccAddress) actResChecker {

	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {

		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)
		q := keeper.NewQuerier(resultCtx.tc.mockKeeper.Keeper)
		cdc := ModuleCdc

		params := types.NewQueryValidatorParams(valAddr)
		bz, _ := cdc.MarshalJSON(params)

		res, e := q(ctx, []string{types.QueryValidatorVotes}, abci.RequestQuery{Data: bz})
		require.Nil(t, e, e)

		var sharesResponses types.SharesResponses
		err := cdc.UnmarshalJSON(res, &sharesResponses)
		b1 := assert.Nil(t, err, err) &&
			assert.Equal(t, expVoterCnt, len(sharesResponses))

		b2 := true
		if b1 && expVoters != nil && len(expVoters) > 0 {

			cnt := 0
			for _, exp := range expVoters {
				for _, real := range sharesResponses {
					if real.DelAddr.Equals(exp) {
						cnt++
						break
					}
				}
			}

			b2 = assert.Equal(t, len(expVoters), cnt, expVoters, sharesResponses.String())

		}

		return b1 && b2
	}

}

func queryPoolCheck(expBonded *sdk.Dec, expUnbonded *sdk.Dec) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {

		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)
		q := keeper.NewQuerier(resultCtx.tc.mockKeeper.Keeper)
		cdc := ModuleCdc

		res, e := q(ctx, []string{types.QueryPool}, abci.RequestQuery{})
		require.Nil(t, e, e)

		pool := types.Pool{}
		require.NoError(t, cdc.UnmarshalJSON(res, &pool))
		require.NotNil(t, pool.String())

		b1 := true
		if expBonded != nil {
			b1 = assert.Equal(t, *expBonded, pool.BondedTokens)
		}

		b2 := true
		if expUnbonded != nil {
			b2 = assert.Equal(t, *expUnbonded, pool.NotBondedTokens)
		}

		stkKeeper := resultCtx.tc.mockKeeper.Keeper
		totalBonded := stkKeeper.TotalBondedTokens(ctx)
		bonedRatio := stkKeeper.BondedRatio(ctx)
		require.True(t, totalBonded.GT(sdk.ZeroDec()))
		// bonedRatio will be equals to Zero when there is only msd in the pool
		require.True(t, bonedRatio.GTE(sdk.ZeroDec()))

		return b1 && b2

	}
}

func baseInVariantCheck(t *testing.T, invariant sdk.Invariant, resultCtx *ActionResultCtx) bool {
	ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)
	msg, broken := invariant(ctx)
	if broken {
		t.Error(msg)
	}
	return !broken

}

func delegatorVotesInvariantCheck() actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		invariant := keeper.DelegatorVotesInvariant(resultCtx.tc.mockKeeper.Keeper)
		return baseInVariantCheck(t, invariant, resultCtx)
	}
}

func positiveDelegatorInvariantCheck() actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		invariant := keeper.PositiveDelegatorInvariant(resultCtx.tc.mockKeeper.Keeper)
		return baseInVariantCheck(t, invariant, resultCtx)
	}
}

func nonNegativePowerInvariantCustomCheck() actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		invariant := keeper.NonNegativePowerInvariantCustom(resultCtx.tc.mockKeeper.Keeper)
		return baseInVariantCheck(t, invariant, resultCtx)
	}
}

func moduleAccountInvariantsCustomCheck() actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		invariant := keeper.ModuleAccountInvariantsCustom(resultCtx.tc.mockKeeper.Keeper)
		return baseInVariantCheck(t, invariant, resultCtx)
	}
}

func getLatestGenesisValidatorCheck(expGenValCnt int) actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		ctx := getNewContext(resultCtx.tc.mockKeeper.MountedStore, resultCtx.tc.currentHeight)
		genVals := GetLatestGenesisValidator(ctx, resultCtx.tc.mockKeeper.Keeper)
		ok := assert.NotNil(t, genVals)
		ok = ok && assert.Equal(t, expGenValCnt, len(genVals), genVals)
		return ok
	}
}

type andChecker struct {
	checkers []actResChecker
}

func (o *andChecker) GetChecker() actResChecker {
	return func(t *testing.T, beforeStatus, afterStatus IValidatorStatus, resultCtx *ActionResultCtx) bool {
		for _, chk := range o.checkers {
			if !chk(t, beforeStatus, afterStatus, resultCtx) {
				return false
			}
		}
		return true
	}
}

// --------------------------------------------------------------
// Validator State Machine TestCase

type basicStakingSMTestCase struct {
	mockKeeper       keeper.MockStakingKeeper
	stkParams        types.Params
	startUpVaStatus  IValidatorStatus
	sequenceActions  []IAction
	actionsResChecks []actResChecker
	currentHeight    int64
	originDlgSet     map[string]*Delegator
	originVaSet      []IValidatorStatus
	test             *testing.T
}

func newValidatorSMTestCase(mk keeper.MockStakingKeeper, params types.Params, startUpStatus IValidatorStatus,
	inputActions []IAction, actionsResCheckers []actResChecker, t *testing.T) basicStakingSMTestCase {

	tc := basicStakingSMTestCase{
		mk,
		params,
		startUpStatus,
		inputActions,
		actionsResCheckers,
		0,
		nil,
		[]IValidatorStatus{},
		t,
	}

	tc.originDlgSet = make(map[string]*Delegator, 10)

	//initialization
	stkKeeper := mk.Keeper
	ctx := getNewContext(mk.MountedStore, tc.currentHeight)
	stkKeeper.SetParams(ctx, tc.stkParams)

	return tc

}

func getNewContext(ms store.MultiStore, height int64) sdk.Context {
	header := abci.Header{ChainID: keeper.TestChainID, Height: height}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	ctx = ctx.WithBlockTime(time.Now())

	return ctx
}

func (tc *basicStakingSMTestCase) SetupValidatorSetAndDelegatorSet(maxValidator int) {

	ctx := getNewContext(tc.mockKeeper.MountedStore, tc.currentHeight)
	bAction := baseAction{tc.mockKeeper}
	var lastStatus IValidatorStatus
	for i := 0; i < maxValidator; i++ {
		v := NewValidator(addrVals[i+1], PKs[i+1], Description{})

		lastStatus = baseValidatorStatus{v}
		result := ActionResultCtx{}
		result.params = tc.stkParams
		result.t = tc.test
		result.tc = tc
		createValidatorAction{bAction, nil}.apply(ctx, lastStatus, &result)
		tc.originVaSet = append(tc.originVaSet, lastStatus)
	}

	// two delegators
	handler := NewHandler(tc.mockKeeper.Keeper)

	handler(ctx, NewMsgDelegate(ValidDelegator1, sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, DelegatedToken1)))
	delegator1, _ := tc.mockKeeper.Keeper.GetDelegator(ctx, ValidDelegator1)
	tc.originDlgSet[delegator1.DelegatorAddress.String()] = &delegator1

	handler(ctx, NewMsgDelegate(ValidDelegator2, sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, DelegatedToken2)))
	delegator2, _ := tc.mockKeeper.Keeper.GetDelegator(ctx, ValidDelegator2)
	tc.originDlgSet[delegator2.DelegatorAddress.String()] = &delegator2

	endBlockAction{bAction}.apply(ctx, lastStatus, nil)
	tc.currentHeight++
}

func (tc *basicStakingSMTestCase) printParticipantSnapshot(t *testing.T) {
	ctx := getNewContext(tc.mockKeeper.MountedStore, tc.currentHeight)

	allVas := tc.mockKeeper.Keeper.GetAllValidators(ctx)
	t.Logf("        ==> Debug Validator Set & Delegators info ")
	for _, v := range allVas {
		t.Logf("Va: %s, Status: %s, Msd: %s,  DS: %s\n", v.GetOperator().String(), v.GetStatus().String(),
			v.GetMinSelfDelegation().String(), v.GetDelegatorShares().String())
	}

	for _, d := range tc.originDlgSet {
		latestDlg, _ := tc.mockKeeper.Keeper.GetDelegator(ctx, d.DelegatorAddress)
		t.Logf("Dlg: %s, VoteTo: %s, BondedToken: %s, GotShare: %s, IsProxy: %+v, HasProxy: %+v, TotalDS: %s \n",
			latestDlg.DelegatorAddress.String(), latestDlg.ValidatorAddresses, latestDlg.Tokens.String(),
			latestDlg.Shares.String(), latestDlg.IsProxy, latestDlg.HasProxy(), latestDlg.TotalDelegatedTokens.String())
	}

}

func (tc *basicStakingSMTestCase) Run(t *testing.T) {

	stkKeeper := tc.mockKeeper.Keeper
	ctx := getNewContext(tc.mockKeeper.MountedStore, tc.currentHeight)
	stkKeeper.SetParams(ctx, tc.stkParams)

	if len(tc.sequenceActions) != len(tc.actionsResChecks) {
		panic(fmt.Sprintf("length of seqenceActions(%d) & resultChecker(%d) should be the same", len(tc.sequenceActions), len(tc.actionsResChecks)))
	}

	//1. enter validator status and actions loop
	beforeStatus := tc.startUpVaStatus
	for i := 0; i < len(tc.sequenceActions); i++ {
		action := tc.sequenceActions[i]

		check := tc.actionsResChecks[i]
		resultCtx := ActionResultCtx{}
		resultCtx.params = tc.stkParams
		resultCtx.tc = tc
		resultCtx.t = t

		action.apply(ctx, beforeStatus, &resultCtx)

		afterValidator, _ := stkKeeper.GetValidator(ctx, tc.startUpVaStatus.getValidator().OperatorAddress)
		afterStatus := baseValidatorStatus{afterValidator}

		if check != nil {
			cr := check(t, beforeStatus, afterStatus, &resultCtx)

			if !cr {
				tc.printParticipantSnapshot(t)
			}

			require.True(t, cr, fmt.Sprintf("====ActionRound: %d\n", i+1),
				tc.stkParams, beforeStatus.desc(), afterStatus.desc(), resultCtx.String())

		}

		if resultCtx.isBlkHeightInc {
			tc.currentHeight++
			resultCtx.isBlkHeightInc = false
		}

		ctx = getNewContext(tc.mockKeeper.MountedStore, tc.currentHeight)
		beforeStatus = afterStatus
	}
}
