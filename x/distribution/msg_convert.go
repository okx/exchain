package distribution

import (
	"encoding/json"
	"errors"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/distribution/types"
)

var (
	ErrCheckSignerFail = errors.New("check signer fail")
)

func init() {
	RegisterConvert()
}

func RegisterConvert() {
	//baseapp.Register(types.ModuleName, "set-withdraw-addr", ConvertSetWithdrawAddressMsg)
	baseapp.Register(types.ModuleName, "withdraw-rewards", ConvertWithdrawDelegatorRewardMsg)
	//baseapp.Register(types.ModuleName, "withdraw-rewards-commission", ConvertWithdrawValidatorCommissionMsg)
}

func ConvertSetWithdrawAddressMsg(data []byte, signers []sdk.AccAddress) (sdk.Msg, error) {
	newMsg := types.MsgSetWithdrawAddress{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return newMsg, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return newMsg, nil
}

func ConvertWithdrawDelegatorRewardMsg(data []byte, signers []sdk.AccAddress) (sdk.Msg, error) {
	newMsg := types.MsgWithdrawDelegatorReward{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return newMsg, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return newMsg, nil
}

func ConvertWithdrawValidatorCommissionMsg(data []byte, signers []sdk.AccAddress) (sdk.Msg, error) {
	newMsg := types.MsgWithdrawValidatorCommission{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return newMsg, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return newMsg, nil
}
