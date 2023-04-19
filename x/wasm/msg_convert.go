package wasm

import (
	"encoding/json"
	"errors"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/wasm/types"
)

var (
	ErrCheckSignerFail = errors.New("check signer fail")
	ErrNotFindHandle   = errors.New("not find handle")
)

func init() {
	RegisterConvert()
}

func RegisterConvert() {
	baseapp.RegisterCmHandleV1("wasm/MsgStoreCode", baseapp.NewCMHandleV1(ConvertMsgStoreCode))
	baseapp.RegisterCmHandleV1("wasm/MsgInstantiateContract", baseapp.NewCMHandleV1(ConvertMsgInstantiateContract))
	baseapp.RegisterCmHandleV1("wasm/MsgExecuteContract", baseapp.NewCMHandleV1(ConvertMsgExecuteContract))
	baseapp.RegisterCmHandleV1("wasm/MsgMigrateContract", baseapp.NewCMHandleV1(ConvertMsgMigrateContract))
	baseapp.RegisterCmHandleV1("wasm/MsgUpdateAdmin", baseapp.NewCMHandleV1(ConvertMsgUpdateAdmin))
}

func ConvertMsgStoreCode(data []byte, signers []sdk.AccAddress, height int64) (sdk.Msg, error) {
	if !tmtypes.HigherThanVenus6(height) {
		return nil, ErrNotFindHandle
	}
	newMsg := types.MsgStoreCode{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return nil, err
	}
	err = newMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return &newMsg, nil
}

func ConvertMsgInstantiateContract(data []byte, signers []sdk.AccAddress, height int64) (sdk.Msg, error) {
	if !tmtypes.HigherThanVenus6(height) {
		return nil, ErrNotFindHandle
	}
	newMsg := types.MsgInstantiateContract{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return nil, err
	}
	err = newMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return &newMsg, nil
}

func ConvertMsgExecuteContract(data []byte, signers []sdk.AccAddress, height int64) (sdk.Msg, error) {
	if !tmtypes.HigherThanVenus6(height) {
		return nil, ErrNotFindHandle
	}
	newMsg := types.MsgExecuteContract{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return nil, err
	}
	err = newMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return &newMsg, nil
}

func ConvertMsgMigrateContract(data []byte, signers []sdk.AccAddress, height int64) (sdk.Msg, error) {
	if !tmtypes.HigherThanVenus6(height) {
		return nil, ErrNotFindHandle
	}
	newMsg := types.MsgMigrateContract{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return nil, err
	}
	err = newMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return &newMsg, nil
}

func ConvertMsgUpdateAdmin(data []byte, signers []sdk.AccAddress, height int64) (sdk.Msg, error) {
	if !tmtypes.HigherThanVenus6(height) {
		return nil, ErrNotFindHandle
	}
	newMsg := types.MsgUpdateAdmin{}
	err := json.Unmarshal(data, &newMsg)
	if err != nil {
		return nil, err
	}
	err = newMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	if ok := common.CheckSignerAddress(signers, newMsg.GetSigners()); !ok {
		return nil, ErrCheckSignerFail
	}
	return &newMsg, nil
}
