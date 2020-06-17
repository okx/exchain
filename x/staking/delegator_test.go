package staking

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidatorSMProxyDelegationSmoke(t *testing.T) {
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

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

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

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
		noErrorInHandlerResult(false),
		nil,

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
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

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

		// redelegate & unbond
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

		// redelegate & unbond
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),

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
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

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

	voteActionChecker := andChecker{checkers:[]actResChecker{
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
	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

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
