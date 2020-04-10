package staking

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidatorMultiCreates(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = 1
	params.MinSelfDelegationLimit = sdk.NewDec(2)

	validMsd := sdk.NewDec(3)
	invalidMsd := sdk.OneDec()

	startUpValidator := NewValidator(addrVals[0], PKs[0], Description{})
	startUpValidator.MinSelfDelegation = validMsd
	expectDelegatorShares := validMsd

	startUpStatus := baseValidatorStatus{startUpValidator}

	invalidVal := NewValidator(addrVals[1], PKs[1], Description{})
	invalidVal.MinSelfDelegation = invalidMsd
	invalidVaStatus := baseValidatorStatus{invalidVal}

	bAction := baseAction{mk}
	inputActions := []IAction{
		createValidatorAction{bAction, invalidVaStatus},
		createValidatorAction{bAction, startUpStatus},
		endBlockAction{bAction},
	}

	actionsAndChecker := []actResChecker{
		noErrorInHandlerResult(false),
		validatorStatusChecker(sdk.Unbonded.String()),
		queryValidatorCheck(sdk.Bonded, false, &expectDelegatorShares, &validMsd, nil),
	}

	smTestCase := newValidatorSMTestCase(mk, params, startUpStatus, inputActions, actionsAndChecker)
	smTestCase.Run(t)
}

func TestValidatorSM1Create2Destroy3Create(t *testing.T) {

	_, _, mk := CreateTestInput(t, false, SufficientInitPower)

	params := DefaultParams()
	params.MaxValidators = 1
	params.Epoch = 1
	params.UnbondingTime = time.Millisecond * 300

	startUpValidator := NewValidator(addrVals[0], PKs[0], Description{})
	startUpValidator.MinSelfDelegation = InitMsd2000
	expectDelegatorShares := InitMsd2000

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
		queryValidatorCheck(sdk.Bonded, false, &expectDelegatorShares, &InitMsd2000, nil),
		noErrorInHandlerResult(true),
		validatorKickedOff(true),
		nil,
		nil,
		nil,
		queryValidatorCheck(sdk.Unbonded, false, &expectDelegatorShares, &InitMsd2000, nil),
		queryValidatorCheck(sdk.Bonded, false, &expectDelegatorShares, &InitMsd2000, nil),
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
