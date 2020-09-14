package v0_11

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/dex/legacy/v0_10"
)

func Migrate(oldGenState v0_10.GenesisState) GenesisState {
	params := Params{
		ListFee:                oldGenState.Params.ListFee,
		TransferOwnershipFee:   oldGenState.Params.TransferOwnershipFee,
		DelistMaxDepositPeriod: oldGenState.Params.DelistMaxDepositPeriod,
		DelistMinDeposit:       oldGenState.Params.DelistMinDeposit,
		DelistVotingPeriod:     oldGenState.Params.DelistVotingPeriod,
		WithdrawPeriod:         oldGenState.Params.WithdrawPeriod,
		RegisterOperatorFee:    sdk.NewDecCoinFromDec(common.NativeToken, sdk.ZeroDec()),
	}

	operatorMap := make(map[string]struct{})
	var operators DEXOperators
	var maxTokenPairID uint64 = 0
	for _, pair := range oldGenState.TokenPairs {
		if pair.ID > maxTokenPairID {
			maxTokenPairID = pair.ID
		}

		if _, exist := operatorMap[pair.Owner.String()]; exist {
			continue
		}
		operatorMap[pair.Owner.String()] = struct{}{}
		operator := DEXOperator{
			Address:            pair.Owner,
			HandlingFeeAddress: pair.Owner,
			Website:            "",
			InitHeight:         pair.BlockHeight,
			TxHash:             "",
		}
		operators = append(operators, operator)
	}

	return GenesisState{
		Params:         params,
		TokenPairs:     oldGenState.TokenPairs,
		WithdrawInfos:  oldGenState.WithdrawInfos,
		ProductLocks:   oldGenState.ProductLocks,
		Operators:      operators,
		MaxTokenPairID: maxTokenPairID,
	}
}
