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
	params.MinSelfDelegationLimit = MinSelfDelegationLimited

	startUpValidator := NewValidator(addrVals[0], PKs[0], Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000
	expectDelegatorShares := InitMsd2000

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &expectDelegatorShares, &InitMsd2000, nil),
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
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.Run(t)
}

func TestValidatorSMNormalFullLifeCircle(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
		destroyValidatorAction{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	destroyChecker := andChecker{[]actResChecker{
		validatorDelegatorShareLeft(false),
		validatorKickedOff(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),

		// destroyValidatorAction checker
		destroyChecker.GetChecker(),

		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorRemoved(true),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.Run(t)

}

func TestValidatorSMEvilFullLifeCircle(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
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
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),
		jailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.Run(t)
}

func TestValidatorSMEvilFullLifeCircleWithUnjail(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
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
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),
		jailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		unJailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Bonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.Run(t)
}

func TestValidatorSMEvilFullLifeCircleWithUnjail2(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
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
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),
		jailedChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorKickedOff(false),
		validatorStatusChecker(sdk.Bonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.Run(t)
}

func TestValidatorSMEpochRotate(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
		otherMostPowerfulValidatorEnter{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		// startUpValidator created
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),

		// more powerful validator enter
		validatorStatusChecker(sdk.Bonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),

		// entering a new epoch
		// startUpValidator fail to keep validator's position
		validatorStatusChecker(sdk.Unbonding.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.Run(t)

}

func TestValidatorSMReRankPowerIndex(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = uint16(len(addrVals)) - 1
	params.Epoch = 2
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}

	voteChecker := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Unbonded.String()),
	}}

	undelegateChecker := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(false),
	}}

	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorsVoteAction{bAction, false, true, 0, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorsUnBondAction{bAction, true, true},
		endBlockAction{bAction},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		voteChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Bonded.String()),
		undelegateChecker.GetChecker(),
		validatorStatusChecker(sdk.Bonded.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.Run(t)

}

// the following case should be designed again carefully, focus on:
// 0. multi-voting (5 validators & 2delegator)
// 1. validator msd, delgatorshares, votes
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
	startUpValidator.MinSelfDelegation = InitMsd2000

	startUpStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	expZeroDec := sdk.ZeroDec()
	expVasBondedToken := DefaultValidInitMsd.MulInt64(int64(len(originVaSet))).Add(
		startUpValidator.MinSelfDelegation)
	expDlgGrpBondedToken := MaxDelegatedToken.MulInt64(int64(len(ValidDlgGroup)))
	expAllBondedToken := expVasBondedToken.Add(expDlgGrpBondedToken)
	startUpCheck := andChecker{[]actResChecker{
		queryPoolCheck(&expAllBondedToken, &expZeroDec),
		noErrorInHandlerResult(true),
	}}

	// after delegator in group finish voting, do following check.
	voteChecker := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Unbonded.String()),
		queryDelegatorCheck(ValidDelegator1, true, fullVaSet, nil, &MaxDelegatedToken, &expZeroDec),
		queryAllValidatorCheck([]sdk.BondStatus{sdk.Unbonded, sdk.Bonded, sdk.Unbonding}, []int{1, 4, 0}),
		queryVotesToCheck(startUpStatus.getValidator().OperatorAddress, 2, ValidDlgGroup),
		queryPoolCheck(&expAllBondedToken, &expZeroDec),
		noErrorInHandlerResult(true),
	}}

	// All Deleagtor Unbond 4000 okt * 2 (half of original delegated tokens)
	expDlgBondedTokens1 := MaxDelegatedToken.QuoInt64(2)
	expDlgUnbondedToken1 := expDlgBondedTokens1
	expAllUnBondedToken1 := expDlgUnbondedToken1.MulInt64(2)
	expAllBondedToken1 := expAllBondedToken.Sub(expAllUnBondedToken1)
	undelegateChecker1 := andChecker{[]actResChecker{
		validatorDelegatorShareIncreased(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		queryDelegatorCheck(ValidDelegator1, true, fullVaSet, nil, &expDlgBondedTokens1, &expDlgUnbondedToken1),
		queryAllValidatorCheck([]sdk.BondStatus{sdk.Unbonded, sdk.Bonded, sdk.Unbonding}, []int{1, 4, 0}),
		queryVotesToCheck(startUpStatus.getValidator().OperatorAddress, 2, ValidDlgGroup),
		queryPoolCheck(&expAllBondedToken1, &expAllUnBondedToken1),
	}}

	// All Deleagtor left bonded 4000 okt * 2 (another half of original delegated tokens)
	expDlgGrpUnbonded2 := MaxDelegatedToken.MulInt64(int64(len(ValidDlgGroup)))
	expAllBondedToken2 := expAllBondedToken.Sub(expDlgGrpUnbonded2)
	undelegateChecker2 := andChecker{[]actResChecker{
		// cannot find unbonding token in GetUnbonding info
		queryDelegatorCheck(ValidDelegator1, false, []sdk.ValAddress{}, nil, &expZeroDec, nil),
		queryAllValidatorCheck([]sdk.BondStatus{sdk.Unbonded, sdk.Bonded, sdk.Unbonding}, []int{1, 4, 0}),
		queryVotesToCheck(startUpStatus.getValidator().OperatorAddress, 0, []sdk.AccAddress{}),
		queryPoolCheck(&expAllBondedToken2, &expZeroDec),
	}}

	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorsVoteAction{bAction, true, true, 0, nil},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorsUnBondAction{bAction, true, false},
		endBlockAction{bAction},
		endBlockAction{bAction},
		delegatorsUnBondAction{bAction, true, true},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		startUpCheck.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		voteChecker.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		nil,
		undelegateChecker1.GetChecker(),
		validatorStatusChecker(sdk.Unbonded.String()),
		validatorStatusChecker(sdk.Unbonded.String()),
		nil,
		undelegateChecker2.GetChecker(),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.printParticipantSnapshot()
	smTestCase.Run(t)

}

// 1. apply DestroyValidator action on a Bonded VA x
// 2. Wait for an Unbonded VA x
// 3. Then delegator unbond all tokens to withdraw votes from VA x
// 4. VA is removed
func TestValidatorSMDestroyValidatorUnbonding2UnBonded2Removed(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)
	params := DefaultParams()

	originVaSet := addrVals[1:]
	params.MaxValidators = uint16(len(originVaSet))
	params.Epoch = 1
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(StartUpValidatorAddr, StartUpValidatorPubkey, Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		delegatorsVoteAction{bAction, false, true, 0, nil},
		endBlockAction{bAction},
		destroyValidatorAction{bAction},
		endBlockAction{bAction},

		// first unbonding time pass, delegator shares left, validator unbonding --> unbonded
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},

		// delegators unbond all tokens back, validator has no msd & delegator shares now, delegator removed
		delegatorsUnBondAction{bAction, true, true},
	}

	//expZeroInt := sdk.ZeroInt()
	expZeroDec := sdk.ZeroDec()
	dlgVoteCheck1 := andChecker{[]actResChecker{
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

	expDlgShare := startUpValidator.MinSelfDelegation
	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &expDlgShare, &startUpValidator.MinSelfDelegation, nil),
		dlgVoteCheck1.GetChecker(),
		nil,
		queryValidatorCheck(sdk.Bonded, true, nil, &expZeroDec, nil),
		validatorStatusChecker(sdk.Unbonding.String()),
		validatorStatusChecker(sdk.Unbonding.String()),
		afterUnbondingTimeExpiredCheck1.GetChecker(),
		dlgUnbondCheck2.GetChecker(),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.printParticipantSnapshot()
	smTestCase.Run(t)
}

// 1. apply DestroyValidator action on a Bonded VA x
// 2. Then delegator unbond all tokens to withdraw votes from VA x
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
	startUpValidator.MinSelfDelegation = InitMsd2000.MulInt64(2)

	startUpStatus := baseValidatorStatus{startUpValidator}

	orgValsLen := len(originVaSet)
	fullVaSet := make([]sdk.ValAddress, orgValsLen+1)
	copy(fullVaSet, originVaSet)
	copy(fullVaSet[orgValsLen:], []sdk.ValAddress{startUpStatus.getValidator().GetOperator()})

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		delegatorsVoteAction{bAction, false, true, 0, nil},
		endBlockAction{bAction},
		destroyValidatorAction{bAction},
		endBlockAction{bAction},

		// delegators unbond all tokens back, validator has no msd & delegator shares now, delegator removed
		delegatorsUnBondAction{bAction, true, true},

		// second unbonding time pass, no delegator shares left, unbonding --> validator removed
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
	}

	//expZeroInt := sdk.ZeroInt()
	expZeroDec := sdk.ZeroDec()
	dlgVoteCheck1 := andChecker{[]actResChecker{
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

	expDlgShare := startUpValidator.MinSelfDelegation
	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &expDlgShare, &startUpValidator.MinSelfDelegation, nil),
		dlgVoteCheck1.GetChecker(),
		nil,
		queryValidatorCheck(sdk.Bonded, true, nil, &expZeroDec, nil),
		validatorStatusChecker(sdk.Unbonding.String()),
		dlgUnbondCheck2.GetChecker(),
		queryValidatorCheck(sdk.Unbonding, true, nil, &expZeroDec, nil),
		afterUnbondingTimeExpiredCheck1.GetChecker(),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators), sdk.OneDec().Int64())
	smTestCase.printParticipantSnapshot()
	smTestCase.Run(t)
}
