package staking

import (
	"testing"
	"time"

	"github.com/okex/okchain/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidatorSMProxyDelegationSmoke(t *testing.T) {
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
		newDelegatorAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},
		baseProxyRegAction{bAction, ProxiedDelegator, true},
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},
		proxyBindAction{bAction, ValidDelegator2, ProxiedDelegator},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ProxiedDelegator}},
		proxyUnBindAction{bAction, ValidDelegator1},

		baseProxyRegAction{bAction, ProxiedDelegator, false},
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

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	zeroDec := sdk.ZeroDec()
	tokenPerTime := sdk.NewDec(8000)
	inputActions := []IAction{
		createValidatorAction{bAction, nil},

		// send delegate messages
		newDelegatorAction{bAction, ValidDelegator1, tokenPerTime, "testtoken"},
		newDelegatorAction{bAction, ValidDelegator1, tokenPerTime, sdk.DefaultBondDenom},
		newDelegatorAction{bAction, ValidDelegator1, tokenPerTime, sdk.DefaultBondDenom},
		newDelegatorAction{bAction, ValidDelegator1, tokenPerTime.MulInt64(10), sdk.DefaultBondDenom},
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
		delegatorUnbondAction{bAction, ValidDelegator2, sdk.ZeroDec(), "testtoken"},
		delegatorUnbondAction{bAction, ValidDelegator2, sdk.ZeroDec(), sdk.DefaultBondDenom},
		delegatorUnbondAction{bAction, ValidDelegator1, tokenPerTime.MulInt64(2), sdk.DefaultBondDenom},
		delegatorUnbondAction{bAction, ValidDelegator1, tokenPerTime.QuoInt64(2), sdk.DefaultBondDenom},
		delegatorUnbondAction{bAction, ValidDelegator1, tokenPerTime.QuoInt64(2), sdk.DefaultBondDenom},
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
		baseProxyRegAction{bAction, ProxiedDelegator, true},
		baseProxyRegAction{bAction, ProxiedDelegator, false},
		newDelegatorAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},

		// successfully regiester
		// delegate again
		newDelegatorAction{bAction, ProxiedDelegator, MaxDelegatedToken, sdk.DefaultBondDenom},
		baseProxyRegAction{bAction, ProxiedDelegator, true},
		baseProxyRegAction{bAction, ProxiedDelegator, true},

		// bind
		proxyBindAction{bAction, ValidDelegator1, InvalidDelegator},
		proxyBindAction{bAction, ValidDelegator1, ValidDelegator2},
		proxyBindAction{bAction, InvalidDelegator, ProxiedDelegator},
		proxyBindAction{bAction, ProxiedDelegator, ProxiedDelegator},
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},
		proxyBindAction{bAction, ValidDelegator2, ProxiedDelegator},

		// add shares
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ProxiedDelegator}},
		delegatorsAddSharesAction{bAction, true, true, 0, []sdk.AccAddress{ProxiedDelegator}},

		// redelegate & unbond
		newDelegatorAction{bAction, ValidDelegator1, DelegatedToken1, sdk.DefaultBondDenom},
		delegatorUnbondAction{bAction, ValidDelegator2, DelegatedToken2, sdk.DefaultBondDenom},

		// unbind
		proxyUnBindAction{bAction, InvalidDelegator},
		proxyUnBindAction{bAction, ProxiedDelegator},
		proxyUnBindAction{bAction, ValidDelegator1},

		// unregister
		baseProxyRegAction{bAction, ValidDelegator1, false},
		baseProxyRegAction{bAction, ProxiedDelegator, false},
	}

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

		// add shares
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),

		// redelegate & unbond
		noErrorInHandlerResult(true),
		noErrorInHandlerResult(true),

		// unbind
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(false),
		noErrorInHandlerResult(true),

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

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	proxyOriginTokens := MaxDelegatedToken
	ProxiedDelegatorAlternative := ValidDelegator2
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},

		// register proxy
		newDelegatorAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},
		baseProxyRegAction{bAction, ProxiedDelegator, true},

		// register another proxy
		newDelegatorAction{bAction, ProxiedDelegatorAlternative, proxyOriginTokens, sdk.DefaultBondDenom},
		baseProxyRegAction{bAction, ProxiedDelegatorAlternative, true},
		endBlockAction{bAction},

		// bind proxy
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},

		// rebind to an alternative proxy
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegatorAlternative},
	}

	ProxyChecker1 := andChecker{[]actResChecker{
		noErrorInHandlerResult(true),
		queryProxyCheck(ProxiedDelegator, true, DelegatedToken1),
		queryProxyCheck(ProxiedDelegatorAlternative, true, sdk.ZeroDec()),
	}}

	ProxyChecker2 := andChecker{[]actResChecker{
		noErrorInHandlerResult(true),
		queryProxyCheck(ProxiedDelegator, true, sdk.ZeroDec()),
		queryProxyCheck(ProxiedDelegatorAlternative, true, DelegatedToken1),
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

		// rebind to an alternative proxy
		ProxyChecker2.GetChecker(),
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
		newDelegatorAction{bAction, ProxiedDelegator, proxyOriginTokens, sdk.DefaultBondDenom},
		baseProxyRegAction{bAction, ProxiedDelegator, true},

		// bind proxy
		proxyBindAction{bAction, ValidDelegator1, ProxiedDelegator},

		// register proxy without unbinding
		baseProxyRegAction{bAction, ValidDelegator1, true},
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
