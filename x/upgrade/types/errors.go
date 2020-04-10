package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/okex/okchain/x/gov/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "upgrade"

	CodeInvalidMsgType        sdk.CodeType = 100
	CodeUnSupportedMsgType    sdk.CodeType = 101
	CodeUnknownRequest        sdk.CodeType = sdk.CodeUnknownRequest
	CodeNotCurrentProposal    sdk.CodeType = 102
	CodeNotValidator          sdk.CodeType = 103
	CodeDoubleSwitch          sdk.CodeType = 104
	CodeNoUpgradeConfig       sdk.CodeType = 105
	CodeInvalidUpgradeParams  sdk.CodeType = 107
	CodeInvalidSoftWareDescri sdk.CodeType = 108
	CodeInvalidVersion        sdk.CodeType = 109
	CodeSwitchPeriodInProcess sdk.CodeType = 110
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidMsgType:
		return "Invalid msg type"
	case CodeUnSupportedMsgType:
		return "Current version software doesn't support the msg type"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func NewError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func ErrZeroSwitchHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, govTypes.CodeInvalidHeight, fmt.Sprintf("Protocol switchHeight in AppUpgradeProposal cannot be 0."))
}

func ErrInvalidUpgradeThreshold(codespace sdk.CodespaceType, Threshold sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidUpgradeParams, fmt.Sprintf("Invalid Upgrade Threshold( "+Threshold.String()+" ) should be [0.75, 1)."))
}

func ErrInvalidLength(codespace sdk.CodespaceType, descriptor string, got, max int) sdk.Error {
	msg := fmt.Sprintf("bad length for %v, got length %v, max is %v", descriptor, got, max)
	return sdk.NewError(codespace, CodeInvalidSoftWareDescri, msg)
}

func ErrInvalidVersion(codespace sdk.CodespaceType, version uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVersion, fmt.Sprintf("Version [%v] in AppUpgradeProposal isn't valid.", version))
}

func ErrInvalidSwitchHeight(codespace sdk.CodespaceType, blockHeight uint64, switchHeight uint64) sdk.Error {
	return sdk.NewError(codespace, govTypes.CodeInvalidHeight, fmt.Sprintf("Protocol switchHeight [%v] in AppUpgradeProposal isn't large than current block height [%v].", switchHeight, blockHeight))
}

func ErrSwitchPeriodInProcess(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSwitchPeriodInProcess, fmt.Sprintf("App Upgrade Switch Period is in process."))
}
