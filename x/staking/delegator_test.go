package staking

import (
	"fmt"
	"testing"
	"time"

	"github.com/okex/okexchain/x/common"

	"github.com/okex/okexchain/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidatorSMProxyDelegationSmoke(t *testing.T) {
	common.InitConfig()
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	proxyOriginTokens := MaxDelegatedToken
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		delegatorDepositAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},
		proxyBindAction{bAction, ValidDelegator2, ProxiedDelegator},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ProxiedDelegator}},
		proxyUnBindAction{bAction, ValidDelegator1},

		delegatorRegProxyAction{bAction, ProxiedDelegator, false},
	}

	expZeroDec := sdk.ZeroDec()
	expProxiedToken1 := DelegatedToken1
	expProxiedToken2 := expProxiedToken1.Add(DelegatedToken2)
	prxBindChecker1 := andChecker{[]actResChecker{
		queryDelegatorProxyCheck(ValidDelegator1, false, true, &expZeroDec, &ProxiedDelegator, nil),
		queryDelegatorProxyCheck(ValidDelegator2, false, false, &expZeroDec, nil, nil),
		queryDelegatorProxyCheck(ProxiedDelegator, true, false, &expProxiedToken1,
			nil, []sdk.AccAddress{ValidDelegator1}),
	}}

	prxBindChecker2 := andChecker{[]actResChecker{
		queryDelegatorProxyCheck(ValidDelegator1, false, true, &expZeroDec, &ProxiedDelegator, nil),
		queryDelegatorProxyCheck(ValidDelegator2, false, true, &expZeroDec, &ProxiedDelegator, nil),
		queryDelegatorProxyCheck(ProxiedDelegator, true, false, &expProxiedToken2,
			nil, []sdk.AccAddress{ValidDelegator1, ValidDelegator2}),
	}}

	proxyAddSharesChecker3 := andChecker{[]actResChecker{
		queryDelegatorProxyCheck(ValidDelegator1, false, true, &expZeroDec, &ProxiedDelegator, nil),
		queryDelegatorProxyCheck(ValidDelegator2, false, true, &expZeroDec, &ProxiedDelegator, nil),
		queryDelegatorProxyCheck(ProxiedDelegator, true, false, &expProxiedToken2, nil, nil),

		queryDelegatorCheck(ValidDelegator1, true, []sdk.ValAddress{}, nil, &DelegatedToken1, nil),
		queryDelegatorCheck(ProxiedDelegator, true, []sdk.ValAddress{startUpValidator.GetOperator()}, nil, &proxyOriginTokens, nil),
		querySharesToCheck(startUpValidator.GetOperator(), 1, []sdk.AccAddress{ProxiedDelegator}),
	}}

	prxUnbindChecker4 := andChecker{[]actResChecker{
		queryDelegatorProxyCheck(ValidDelegator1, false, false, &expZeroDec, nil, nil),
		queryDelegatorProxyCheck(ProxiedDelegator, true, false, &DelegatedToken2, nil, nil),
		validatorDelegatorShareIncreased(false),
		delegatorAddSharesInvariantCheck(),
		nonNegativePowerInvariantCustomCheck(),
		positiveDelegatorInvariantCheck(),
		moduleAccountInvariantsCustomCheck(),
	}}

	actionsAndChecker := []actResChecker{
		nil,
		queryDelegatorCheck(ProxiedDelegator, true, nil, &expZeroDec, &proxyOriginTokens, &expZeroDec),
		queryDelegatorProxyCheck(ProxiedDelegator, true, false, &expZeroDec, nil, nil),
		prxBindChecker1.GetChecker(),
		prxBindChecker2.GetChecker(),
		noErrorInHandlerResult(false),
		proxyAddSharesChecker3.GetChecker(),
		prxUnbindChecker4.GetChecker(),
		queryDelegatorProxyCheck(ProxiedDelegator, false, false, &expZeroDec, nil, nil),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.printParticipantSnapshot(t)
	smTestCase.Run(t)
}

func TestDelegator(t *testing.T) {

	common.InitConfig()
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	zeroDec := sdk.ZeroDec()
	tokenPerTime := sdk.NewDec(8000)
	inputActions := []IAction{
		createValidatorAction{bAction, nil},

		// send delegate messages
		delegatorDepositAction{bAction, ValidDelegator1, tokenPerTime, "testtoken"},
		delegatorDepositAction{bAction, ValidDelegator1, tokenPerTime, sdk.DefaultBondDenom},
		delegatorDepositAction{bAction, ValidDelegator1, tokenPerTime, sdk.DefaultBondDenom},
		delegatorDepositAction{bAction, ValidDelegator1, tokenPerTime.MulInt64(10), sdk.DefaultBondDenom},
		endBlockAction{bAction},

		// send add shares messages
		delegatorsAddSharesAction{bAction, false, false, 0, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, false, 0, []sdk.AccAddress{ValidDelegator2}},
		delegatorsAddSharesAction{bAction, false, false, 1, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, false, int(params.MaxValsToAddShares) + 1, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, true, false, 0, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		// [E] repeat addshares to previous validator
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, true, 1, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},

		// send withdraw message
		delegatorWithdrawAction{bAction, ValidDelegator2, sdk.ZeroDec(), "testtoken"},
		delegatorWithdrawAction{bAction, ValidDelegator2, sdk.ZeroDec(), sdk.DefaultBondDenom},
		delegatorWithdrawAction{bAction, ValidDelegator1, tokenPerTime.MulInt64(2), sdk.DefaultBondDenom},
		delegatorWithdrawAction{bAction, ValidDelegator1, tokenPerTime.QuoInt64(2), sdk.DefaultBondDenom},
		delegatorWithdrawAction{bAction, ValidDelegator1, tokenPerTime.QuoInt64(2), sdk.DefaultBondDenom},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},

		// add shares after dlg.share == 0, expect failed
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		// check delegate response
		noErrorInHandlerResult(false),
		//  1. ValidDelegator1 delegate 8000 okt, success
		noErrorInHandlerResult(true),
		//  2. ValidDelegator1 delegate 8000 okt again, fail coz no balance
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		validatorStatusChecker(sdk.Bonded.String()),

		// check adding shares response
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),
		// [E]
		queryDelegatorCheck(ValidDelegator1, true, nil, nil, nil, nil),
		noErrorInHandlerResult(false),
		queryValidatorCheck(sdk.Bonded, false, nil, nil, nil),

		// check withdraw response
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		//   1. ValidDelegator1 UnBond 4000okt, success
		noErrorInHandlerResult(true),
		//   2. ValidDelegator1 UnBond 4000okt, success
		noErrorInHandlerResult(true),
		nil,
		//   3. Check ValidDelegator1's not exists any more
		queryDelegatorCheck(ValidDelegator1, false, []sdk.ValAddress{}, &zeroDec, &zeroDec, nil),

		// check adding shares after dlg.share == 0
		noErrorInHandlerResult(false),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.Run(t)

}

func TestProxy(t *testing.T) {
	common.InitConfig()
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	proxyOriginTokens := MaxDelegatedToken
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},

		// failed to register & unregister
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},
		delegatorRegProxyAction{bAction, ProxiedDelegator, false},
		delegatorDepositAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},

		// successfully regiester
		// delegate again
		delegatorDepositAction{bAction, ProxiedDelegator, MaxDelegatedToken, sdk.DefaultBondDenom},
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},

		// bind
		proxyBindAction{bAction, ValidDelegator1, InvalidDelegator},
		proxyBindAction{bAction, ValidDelegator1, ValidDelegator2},
		proxyBindAction{bAction, InvalidDelegator, ProxiedDelegator},
		proxyBindAction{bAction, ProxiedDelegator, ProxiedDelegator},
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},
		proxyBindAction{bAction, ValidDelegator2, ProxiedDelegator},

		// [E] delegator bind the same proxy again
		proxyBindAction{bAction, ValidDelegator2, ProxiedDelegator},

		// add shares
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ProxiedDelegator}},
		delegatorsAddSharesAction{bAction, true, true, 0, []sdk.AccAddress{ProxiedDelegator}},

		// [E] delegator add shares to the same validator again
		delegatorsAddSharesAction{bAction, true, true, 0, []sdk.AccAddress{ProxiedDelegator}},

		// redeposit & rewithdraw
		delegatorDepositAction{bAction, ValidDelegator1, DelegatedToken1, sdk.DefaultBondDenom},
		delegatorWithdrawAction{bAction, ValidDelegator2, DelegatedToken2, sdk.DefaultBondDenom},

		// unbind
		proxyUnBindAction{bAction, InvalidDelegator},
		proxyUnBindAction{bAction, ProxiedDelegator},
		proxyUnBindAction{bAction, ValidDelegator1},

		// [E] ProxiedDelegator unbind again
		proxyUnBindAction{bAction, ValidDelegator1},

		// unregister
		delegatorRegProxyAction{bAction, ValidDelegator1, false},
		delegatorRegProxyAction{bAction, ProxiedDelegator, false},
	}

	delegatorsChecker := andChecker{[]actResChecker{
		queryDelegatorCheck(ValidDelegator1, true, nil, nil, nil, nil),
		queryDelegatorCheck(ValidDelegator2, true, nil, nil, nil, nil),
		queryDelegatorCheck(ProxiedDelegator, true, nil, nil, nil, nil),
	}}

	actionsAndChecker := []actResChecker{
		nil,
		validatorStatusChecker(sdk.Unbonded.String()),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),

		// register result
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(false),

		// bind
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),

		// [E] bind
		delegatorsChecker.GetChecker(),

		// add shares
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),

		// [E] add shares
		delegatorsChecker.GetChecker(),

		// redeposit & rewithdraw
		queryDelegatorCheck(ValidDelegator1, true, nil, nil, nil, nil),
		queryDelegatorCheck(ValidDelegator2, false, nil, nil, nil, nil),

		// unbind
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),

		// [E] unbind
		noErrorInHandlerResult(false),

		// unregister result
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.printParticipantSnapshot(t)
	smTestCase.Run(t)
}

func TestRebindProxy(t *testing.T) {
	common.InitConfig()
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	proxyOriginTokens := MaxDelegatedToken
	ProxiedDelegatorAlternative := ValidDelegator2
	zeroDec := sdk.ZeroDec()
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},

		// register proxy
		delegatorDepositAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},

		// register another proxy
		delegatorDepositAction{bAction, ProxiedDelegatorAlternative, proxyOriginTokens, sdk.DefaultBondDenom},
		delegatorRegProxyAction{bAction, ProxiedDelegatorAlternative, true},
		endBlockAction{bAction},

		// bind proxy
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},

		// vote validator
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ProxiedDelegator}},
		endBlockAction{bAction},

		// rebind to an alternative proxy
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegatorAlternative},
		endBlockAction{bAction},
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegatorAlternative},
	}

	ProxyChecker1 := andChecker{[]actResChecker{
		noErrorInHandlerResult(true),
		queryDelegatorProxyCheck(ProxiedDelegator, true, false,
			&DelegatedToken1, nil, []sdk.AccAddress{ValidDelegator1}),
		queryDelegatorProxyCheck(ProxiedDelegatorAlternative, true, false,
			&zeroDec, nil, nil),
	}}

	voteActionChecker := andChecker{checkers: []actResChecker{
		noErrorInHandlerResult(true),
		queryDelegatorProxyCheck(ProxiedDelegator, true, false,
			&DelegatedToken1, nil, []sdk.AccAddress{ValidDelegator1}),
		queryDelegatorProxyCheck(ProxiedDelegatorAlternative, true, false,
			&zeroDec, nil, nil),
		validatorDelegatorShareIncreased(true),
	}}

	ProxyChecker2 := andChecker{[]actResChecker{
		noErrorInHandlerResult(true),
		queryProxyCheck(ProxiedDelegator, true, sdk.ZeroDec()),
		queryProxyCheck(ProxiedDelegatorAlternative, true, DelegatedToken1),

		queryDelegatorProxyCheck(ProxiedDelegator, true, false,
			&zeroDec, nil, nil),
		queryDelegatorProxyCheck(ProxiedDelegatorAlternative, true, false,
			&DelegatedToken1, nil, nil),

		validatorDelegatorShareIncreased(false),
	}}

	actionsAndChecker := []actResChecker{
		nil,
		validatorStatusChecker(sdk.Unbonded.String()),
		// register proxy
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),

		// register another proxy
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),
		nil,

		// bind proxy
		ProxyChecker1.GetChecker(),

		// vote check
		voteActionChecker.GetChecker(),
		nil,

		// rebind to an alternative proxy
		ProxyChecker2.GetChecker(),
		nil,
		// rebind ProxiedDelegatorAlternative
		queryDelegatorProxyCheck(ProxiedDelegatorAlternative, true, false,
			&DelegatedToken1, nil, nil),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.printParticipantSnapshot(t)
	smTestCase.Run(t)
}

func TestLimitedProxy(t *testing.T) {
	common.InitConfig()
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	proxyOriginTokens := MaxDelegatedToken
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},

		// register proxy
		delegatorDepositAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},

		// bind proxy
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},

		// register proxy without unbinding
		delegatorRegProxyAction{bAction, ValidDelegator1, true},
	}

	actionsAndChecker := []actResChecker{
		nil,
		validatorStatusChecker(sdk.Unbonded.String()),
		// register proxy
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),
		// bind proxy
		noErrorInHandlerResult(true),
		// register proxy without unbinding (failed)
		noErrorInHandlerResult(false),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.printParticipantSnapshot(t)
	smTestCase.Run(t)

}

//
// Context: create 2 delegator(d1,d2) + 1 proxy(p) + 1 validator(v)
// Operation Group:
//          setup: v(create), p(deposit), d2(addShares to v)
//          step1: p  5(regProxy, addShare(v), bind(p), unbind(p), withdrawSome)
//			step2: d1 4(bind(p), addShare(v), unbind(p), withdrawSome)
//			step3: d1 4(deposit, bind(p), unbind(p), withdrawSome)
//          step4: p  5(regProxy, addShare(v), bind(p), unbind(p), withdrawSome)
//          teardown: p(unReg), v(destroy)
//          case possibilities: 1 * 1 * 5 * 5 * 4 * 4 = 400
//          iterate all the possibilities to run delegatorConstraintCheck and validatorConstrainCheck
//
func TestDelegatorProxyValidatorConstraints4Steps(t *testing.T) {
	common.InitConfig()
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300
	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)
	startUpStatus := baseValidatorStatus{startUpValidator}
	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{}

	step1Actions := IActions{
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ProxiedDelegator}},
		delegatorWithdrawAction{bAction, ProxiedDelegator, sdk.OneDec(), sdk.DefaultBondDenom},
		proxyBindAction{bAction, ProxiedDelegator, ProxiedDelegator},
		proxyUnBindAction{bAction, ProxiedDelegator},
	}

	step2Actions := IActions{
		delegatorDepositAction{bAction, ValidDelegator1, DelegatedToken1, sdk.DefaultBondDenom},
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		proxyUnBindAction{bAction, ValidDelegator1},
		delegatorWithdrawAction{bAction, ValidDelegator1, sdk.OneDec(), sdk.DefaultBondDenom},
	}

	step3Actions := IActions{
		delegatorDepositAction{bAction, ValidDelegator1, DelegatedToken1, sdk.DefaultBondDenom},
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},
		proxyUnBindAction{bAction, ValidDelegator1},
		delegatorWithdrawAction{bAction, ValidDelegator1, sdk.OneDec(), sdk.DefaultBondDenom},
	}

	step4Actions := IActions{
		delegatorRegProxyAction{bAction, ProxiedDelegator, true},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ProxiedDelegator}},
		proxyBindAction{bAction, ProxiedDelegator, ProxiedDelegator},
		proxyUnBindAction{bAction, ProxiedDelegator},
		delegatorWithdrawAction{bAction, ProxiedDelegator, sdk.OneDec(), sdk.DefaultBondDenom},
	}

	for s1 := 0; s1 < len(step1Actions); s1++ {
		for s2 := 0; s2 < len(step2Actions); s2++ {
			for s3 := 0; s3 < len(step3Actions); s3++ {
				for s4 := 0; s4 < len(step4Actions); s4++ {
					inputActions := []IAction{
						createValidatorAction{bAction, nil},
						delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator2}},
						delegatorDepositAction{bAction, ProxiedDelegator, MaxDelegatedToken, sdk.DefaultBondDenom},
						step1Actions[s1],
						step2Actions[s2],
						step3Actions[s3],
						step4Actions[s4],
						delegatorRegProxyAction{bAction, ProxiedDelegator, false},
						destroyValidatorAction{bAction},
					}

					actionsAndChecker, caseName := generateActionsAndCheckers(inputActions, 3)

					t.Logf("============================================== indexes:[%d,%d,%d,%d]  %s ==============================================", s1, s2, s3, s4, caseName)
					_, _, mk := CreateTestInput(t, false, SufficientInitPower)
					smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
					smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
					smTestCase.printParticipantSnapshot(t)
					smTestCase.Run(t)
					t.Log("============================================================================================")
				}

			}
		}
	}
}

type IActions []IAction
type ResCheckers []actResChecker

func defaultConstraintChecker(checkVa bool, checkD2 bool) actResChecker {

	var checkerGroup ResCheckers
	checkerGroup = []actResChecker{
		queryDelegatorCheck(ValidDelegator1, true, nil, nil, nil, nil),
		queryDelegatorCheck(ProxiedDelegator, true, nil, nil, nil, nil),
	}

	if checkD2 {
		checkerGroup = append(checkerGroup, queryDelegatorCheck(ValidDelegator2, true, nil, nil, nil, nil))
	}
	if checkVa {
		checkerGroup = append(checkerGroup, validatorCheck(StartUpValidatorAddr))
	}

	checker := andChecker{checkerGroup}
	return checker.GetChecker()
}

func generateActionsAndCheckers(stepActions IActions, skipActCnt int) (ResCheckers, string) {
	checkers := ResCheckers{}
	caseName := ""
	for i := 0; i < len(stepActions); i++ {
		a := stepActions[i]
		caseName = caseName + fmt.Sprintf("step_%d_%s#$", i, a.desc())
		if i >= skipActCnt {
			checkVa := a.desc() != "destroyVa"
			checkD2 := a.desc() != "dlg2WithdrawAll"
			checkers = append(checkers, defaultConstraintChecker(checkVa, checkD2))

		} else {
			checkers = append(checkers, nil)
		}
	}
	return checkers, caseName
}

//
// Context: create 1 delegator(d) + 1 proxy(p) + 1 validator(v)
// Operation Group:
//          step1: v 1(create)
//          step2: p 7(deposit, regProxy, unregProxy, bind(p), unbind(p), addShare, withdraw)
//			step3: d 5(deposit, bind(p), unbind(p), addShare(v), withdraw)
//          step4: v 2(nil, destroy)
//          step5: d 5(nil, deposit, withdraw, bind, unbind)
//          step6: p 5(nil, deposit, withdraw, bind, unbind)
//          step7: v 2(nil, destroy)
//          possibilities: 1 * 5 * 7 * 2 * 5 * 5 * 2 = 3500
//          iterate all the possibilities an run delegatorConstraintCheck and validatorConstrainCheck
//
func TestDelegatorProxyValidatorShares7Steps(t *testing.T) {

}
