package staking

import (
	"testing"
	"time"

	"github.com/okex/okexchain/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidatorMultiCreates(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = 1
	params.Epoch = 1

	startUpValidator := NewValidator(addrVals[0], PKs[0], Description{}, types.DefaultMinSelfDelegation)
	startUpStatus := baseValidatorStatus{startUpValidator}

	invalidVal := NewValidator(addrVals[1], PKs[1], Description{}, types.DefaultMinSelfDelegation)
	invalidVaStatus := baseValidatorStatus{invalidVal}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, startUpStatus},
		createValidatorAction{bAction, invalidVaStatus},
		delegatorsAddSharesAction{bAction, false, true, 0, []sdk.AccAddress{ValidDelegator1}},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		noErrorInHandlerResult(false),
		validatorDelegatorShareIncreased(true),
		validatorStatusChecker(sdk.Bonded.String()),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker, t)
	smTestCase.SetupValidatorSetAndDelegatorSet(int(params.MaxValidators))
	smTestCase.Run(t)
}

func TestValidatorSM1Create2Destroy3Create(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = 1
	params.Epoch = 1
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(addrVals[0], PKs[0], Description{}, types.DefaultMinSelfDelegation)

	startUpStatus := baseValidatorStatus{startUpValidator}
	recreateValStatus := baseValidatorStatus{startUpValidator}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, nil},
		endBlockAction{bAction},
		destroyValidatorAction{bAction},
		endBlockAction{bAction},
		waitUntilUnbondingTimeExpired{bAction},
		endBlockAction{bAction},
		endBlockAction{bAction},
		createValidatorAction{bAction, &recreateValStatus},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &SharesFromDefaultMSD, &DefaultMSD, nil),
		noErrorInHandlerResult(true),
		validatorKickedOff(true),
		nil,
		nil,
		nil,
		queryValidatorCheck(sdk.Unbonded, false, &SharesFromDefaultMSD, &DefaultMSD, nil),
		queryValidatorCheck(sdk.Bonded, false, &SharesFromDefaultMSD, &DefaultMSD, nil),
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
