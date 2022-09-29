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
	// TODO after merged the pr withdraw-all-rewards, modify it
	baseapp.RegisterCmHandle(types.ModuleName, "withdraw-rewards", ConvertWithdrawDelegatorRewardMsg)
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
