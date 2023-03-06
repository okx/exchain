package keeper

import (
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

func TestHooks(t *testing.T) {
	//for test mock staking keeper hooks
	ctx, _, mkeeper := CreateTestInput(t, false, 0)
	keeper := mkeeper.Keeper
	valsOld := createVals(ctx, 4, keeper)
	vals := []sdk.ValAddress{valsOld[0].GetOperator(), valsOld[1].GetOperator()}

	//mock staking keeper hooks execute an empty statement
	keeper.AfterValidatorCreated(ctx, valsOld[0].GetOperator())
	keeper.BeforeValidatorModified(ctx, valsOld[0].GetOperator())
	keeper.AfterValidatorRemoved(ctx, valsOld[0].GetConsAddr(), valsOld[0].GetOperator())
	keeper.AfterValidatorBonded(ctx, valsOld[0].GetConsAddr(), valsOld[0].GetOperator())
	keeper.AfterValidatorBeginUnbonding(ctx, valsOld[0].GetConsAddr(), valsOld[0].GetOperator())
	keeper.AfterValidatorDestroyed(ctx, valsOld[0].GetConsAddr(), valsOld[0].GetOperator())
	keeper.BeforeDelegationCreated(ctx, addrDels[0], vals)
	keeper.BeforeDelegationSharesModified(ctx, addrDels[0], vals)
	keeper.BeforeDelegationRemoved(ctx, addrDels[0], valsOld[0].GetOperator())
	keeper.AfterDelegationModified(ctx, addrDels[0], vals)
}
