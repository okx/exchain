package staking

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidatorSMCreateValidator(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = 1

	startUpValidator := NewValidator(addrVals[0], PKs[0], Description{})
	expectDelegatorShares := SharesFromDefaultMSD

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &expectDelegatorShares, &DefaultMSD, nil),
		getLatestGenesisValidatorCheck(1),
	}

	smTestCase := basicStakingSMTestCase{
		mk,
		params,
		startUpStatus,
		inputActions,
		actionsAndChecker,
		0,
		nil,
		nil,
		t,
	}

	smTestCase.Run(t)
}

func TestValidatorSMCreateValidatorWithValidatorSet(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 3
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)
}

func TestValidatorSMNormalFullLifeCircle(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		// ensure the validator in the val set
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		endBlockAction{bAction},

		destroyValidatorAction{bAction},

		endBlockAction{bAction},
		// clear the shares on the startUpValidator
		delegatorWithdrawAction{bAction, ValidDelegator1, DelegatedToken1, sdk.DefaultBondDenom},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	destroyChecker := andChecker{[]actResChecker{
		validatorDelegatorShareLeft(true),
		validatorKickedOff(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),

		// destroyValidatorAction checker
		destroyChecker.GetChecker(),

		validatorStatusChecker(sdk.Unbonding.String()),
		validatorDelegatorShareLeft(false),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorRemoved(true),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)

}

func TestValidatorSMEvilFullLifeCircle(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		delegatorsAddSharesAction{bAction, true, false, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator2}},
		endBlockAction{bAction},
		jailValidatorAction{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	jailedChecker := andChecker{[]actResChecker{
		validatorDelegatorShareLeft(true),
		validatorKickedOff(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Bonded.String()),
		jailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)
}

func TestValidatorSMEvilFullLifeCircleWithUnjail(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		delegatorsAddSharesAction{bAction, true, false, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator2}},
		endBlockAction{bAction},
		jailValidatorAction{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
		unJailValidatorAction{bAction},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	jailedChecker := andChecker{[]actResChecker{
		validatorDelegatorShareLeft(true),
		validatorKickedOff(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}}

	unJailedChecker := andChecker{[]actResChecker{
		validatorDelegatorShareLeft(true),
		validatorKickedOff(false),
		validatorStatusChecker(sdk.Unbonding.String()),
	}}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Bonded.String()),
		jailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		unJailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Bonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)
}

func TestValidatorSMEvilFullLifeCircleWithUnjail2(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		delegatorsAddSharesAction{bAction, true, false, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator2}},
		endBlockAction{bAction},
		jailValidatorAction{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
		unJailValidatorAction{bAction},
		endBlockAction{bAction},
	}

	jailedChecker := andChecker{[]actResChecker{
		validatorDelegatorShareLeft(true),
		validatorKickedOff(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Bonded.String()),
		jailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorKickedOff(false),
		validatorStatusChecker(sdk.Bonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)
}

func TestValidatorSMEpochRotate(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		endBlockAction{bAction},
		otherMostPowerfulValidatorEnter{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		// startUpValidator created
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),

		// more powerful validator enter
		validatorStatusChecker(sdk.Bonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),

		// entering a new epoch
		// startUpValidator fail to keep validator's position
		validatorStatusChecker(sdk.Unbonding.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)

}

func TestValidatorSMReRankPowerIndex(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}

	addSharesChecker := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Unbonded.String()),
	}}

	withdrawChecker := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(false),
	}}

	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, true, false, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator2}},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorWithdrawAction{bAction, ValidDelegator2, DelegatedToken2, sdk.DefaultBondDenom},
		endBlockAction{bAction},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		addSharesChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),
		withdrawChecker.GetChecker(),
		validatorStatusChecker(sdk.Bonded.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)

}

// the following case should be designed again carefully, focus on:
// 0. multi-voting (5 validators & 2delegator)
// 1. validator msd, delegatorshares, add shares
// 2. check delegator's token, unbonded token, shares
func TestValidatorSMMultiVoting(t *testing.T) {

	ctx, _, mk := CreateTestInput(t, false, SufficientInitPower)
	clearNotBondedPool(t, ctx, mk.SupplyKeeper)

	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	expZeroDec := sdk.ZeroDec()
	expValsBondedToken := DefaultMSD.MulInt64(int64(len(fullVaSet)))
	expDlgGrpBondedToken := DelegatedToken1.Add(DelegatedToken2)
	expAllBondedToken := expValsBondedToken.Add(expDlgGrpBondedToken)
	startUpCheck := andChecker{[]actResChecker{
		queryPoolCheck(&expAllBondedToken, &expZeroDec),
		noErrorInHandlerResult(true),
	}}

	// after delegator in group finish adding shares, do following check
	addSharesChecker := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Unbonded.String()),
		queryDelegatorCheck(ValidDelegator2, true, fullVaSet, nil, &DelegatedToken2, &expZeroDec),
		queryAllValidatorCheck([]sdk.BondStatus{sdk.Unbonded, sdk.Bonded, sdk.Unbonding}, []int{1, 4, 0}),
		querySharesToCheck(startUpStatus.getValidator().OperatorAddress, 1, []sdk.AccAddress{ValidDelegator2}),
		queryPoolCheck(&expAllBondedToken, &expZeroDec),
		noErrorInHandlerResult(true),
	}}

	// All Deleagtor Unbond half of the delegation
	expDlgBondedTokens1 := DelegatedToken1.QuoInt64(2)
	expDlgUnbondedToken1 := expDlgBondedTokens1
	expDlgBondedTokens2 := DelegatedToken2.QuoInt64(2)
	expDlgUnbondedToken2 := expDlgBondedTokens2
	expAllUnBondedToken1 := expDlgUnbondedToken1.Add(expDlgUnbondedToken2)
	expAllBondedToken1 := DefaultMSD.MulInt64(int64(len(fullVaSet))).Add(expDlgBondedTokens1).Add(expDlgBondedTokens2)
	withdrawChecker1 := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		queryDelegatorCheck(ValidDelegator1, true, originVaSet, nil, &expDlgBondedTokens1, &expDlgUnbondedToken1),
		queryDelegatorCheck(ValidDelegator2, true, fullVaSet, nil, &expDlgBondedTokens2, &expDlgUnbondedToken2),
		queryAllValidatorCheck([]sdk.BondStatus{sdk.Unbonded, sdk.Bonded, sdk.Unbonding}, []int{1, 4, 0}),
		querySharesToCheck(startUpStatus.getValidator().OperatorAddress, 1, []sdk.AccAddress{ValidDelegator2}),
		queryPoolCheck(&expAllBondedToken1, &expAllUnBondedToken1),
	}}

	// All Deleagtor Unbond the delegation left
	expDlgGrpUnbonded2 := expZeroDec
	expAllBondedToken2 := DefaultMSD.MulInt64(int64(len(fullVaSet)))
	withdrawChecker2 := andChecker{[]actResChecker{
		// cannot find unbonding token in GetUnbonding info
		queryDelegatorCheck(ValidDelegator1, false, []sdk.ValAddress{}, nil, &expZeroDec, nil),
		queryDelegatorCheck(ValidDelegator2, false, []sdk.ValAddress{}, nil, &expZeroDec, nil),
		queryAllValidatorCheck([]sdk.BondStatus{sdk.Unbonded, sdk.Bonded, sdk.Unbonding}, []int{1, 4, 0}),
		querySharesToCheck(startUpStatus.getValidator().OperatorAddress, 0, []sdk.AccAddress{}),
		queryPoolCheck(&expAllBondedToken2, &expDlgGrpUnbonded2),
	}}

	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, true, false, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, true, true, 0, []sdk.AccAddress{ValidDelegator2}},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorsWithdrawAction{bAction, true, false},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorsWithdrawAction{bAction, true, true},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		startUpCheck.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorDelegatorShareIncreased(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		addSharesChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		nil,
		withdrawChecker1.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		nil,
		withdrawChecker2.GetChecker(),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.printParticipantSnapshot(t)
	smTestCase.Run(t)

}

// 1. apply DestroyValidator action on a Bonded VA x
// 2. Wait for an Unbonded VA x
// 3. Then delegator unbond all tokens to withdraw shares from VA x
// 4. VA is removed
func TestValidatorSMDestroyValidatorUnbonding2UnBonded2Removed(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 1
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, false, true, 0, nil},
		endBlockAction{bAction},
		destroyValidatorAction{bAction},
		endBlockAction{bAction},

		// first unbonding time pass, delegator shares left, validator unbonding --> unbonded
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},

		// delegators unbond all tokens back, validator has no msd & delegator shares now, delegator removed
		delegatorsWithdrawAction{bAction, true, true},
	}

	expZeroDec := sdk.ZeroDec()
	dlgAddSharesCheck1 := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(true),
		validatorRemoved(false),
		validatorDelegatorShareLeft(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}}

	afterUnbondingTimeExpiredCheck1 := andChecker{[]actResChecker{
		validatorRemoved(false),
		validatorStatusChecker(sdk.Unbonded.String()),
	}}

	dlgUnbondCheck2 := andChecker{[]actResChecker{
		noErrorInHandlerResult(true),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorRemoved(true),
		queryDelegatorCheck(ValidDelegator1, false, nil, nil, &expZeroDec, nil),
	}}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &SharesFromDefaultMSD, &startUpValidator.MinSelfDelegation, nil),
		dlgAddSharesCheck1.GetChecker(),
		nil,
		queryValidatorCheck(sdk.Bonded, true, nil, &expZeroDec, nil),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		afterUnbondingTimeExpiredCheck1.GetChecker(),
		dlgUnbondCheck2.GetChecker(),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.printParticipantSnapshot(t)
	smTestCase.Run(t)
}

// 1. apply DestroyValidator action on a Bonded VA x
// 2. Then delegator unbond all tokens to withdraw shares from VA x
// 3. Wait for an Unbonded VA x
// 4. VA is removed
func TestValidatorSMDestroyValidatorUnbonding2Removed(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 1
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		delegatorsAddSharesAction{bAction, false, true, 0, nil},
		endBlockAction{bAction},
		destroyValidatorAction{bAction},
		endBlockAction{bAction},

		// delegators unbond all tokens back, validator has no msd & delegator shares now, delegator removed
		delegatorsWithdrawAction{bAction, true, true},

		// second unbonding time pass, no delegator shares left, unbonding --> validator removed
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	//expZeroInt := sdk.ZeroInt()
	expZeroDec := sdk.ZeroDec()
	dlgAddSharesCheck1 := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(true),
		validatorRemoved(false),
		validatorDelegatorShareLeft(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}}

	dlgUnbondCheck2 := andChecker{[]actResChecker{
		noErrorInHandlerResult(true),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorRemoved(false),
		queryDelegatorCheck(ValidDelegator1, false, nil, nil, &expZeroDec, nil),
	}}

	afterUnbondingTimeExpiredCheck1 := andChecker{[]actResChecker{
		validatorRemoved(true),
	}}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &SharesFromDefaultMSD, &startUpValidator.MinSelfDelegation, nil),
		dlgAddSharesCheck1.GetChecker(),
		nil,
		queryValidatorCheck(sdk.Bonded, true, nil, &expZeroDec, nil),
		validatorStatusChecker(sdk.Unbonding.String()),
		dlgUnbondCheck2.GetChecker(),
		queryValidatorCheck(sdk.Unbonding, true, nil, &expZeroDec, nil),
		afterUnbondingTimeExpiredCheck1.GetChecker(),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.printParticipantSnapshot(t)
	smTestCase.Run(t)
}
