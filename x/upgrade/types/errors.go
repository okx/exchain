package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/okex/okchain/x/gov/types"
)

// const
const (
	DefaultCodespace          sdk.CodespaceType = "upgrade"
	CodeUnknownRequest                          = sdk.CodeUnknownRequest
	CodeInvalidMsgType        sdk.CodeType      = 100
	CodeUnSupportedMsgType    sdk.CodeType      = 101
	CodeNotCurrentProposal    sdk.CodeType      = 102
	CodeNotValidator          sdk.CodeType      = 103
	CodeDoubleSwitch          sdk.CodeType      = 104
	CodeNoUpgradeConfig       sdk.CodeType      = 105
	CodeInvalidUpgradeParams  sdk.CodeType      = 107
	CodeInvalidSoftWareDescri sdk.CodeType      = 108
	CodeInvalidVersion        sdk.CodeType      = 109
	CodeSwitchPeriodInProcess sdk.CodeType      = 110
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidMsgType:
		return "Invalid msg type"
	case CodeUnSupportedMsgType:
		return "current version software doesn't support the msg type"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// NewError returns a new sdk.Error with a specific msg
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

// ErrInvalidVersion returns an error when the version is invalid
func ErrInvalidVersion(codespace sdk.CodespaceType, version uint64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVersion,
		fmt.Sprintf("failed. version [%v] in AppUpgradeProposal is invalid", version))
}

// ErrInvalidSwitchHeight returns an error when the switch height for upgrade is invalid
func ErrInvalidSwitchHeight(codespace sdk.CodespaceType, blockHeight uint64, switchHeight uint64) sdk.Error {
	return sdk.NewError(codespace, govTypes.CodeInvalidHeight,
		fmt.Sprintf("failed. protocol switchHeight [%v] in AppUpgradeProposal isn't large than current block height [%v]",
			switchHeight, blockHeight))
}

// ErrSwitchPeriodInProcess returns an error when the UpgradeConfig has existed
func ErrSwitchPeriodInProcess(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSwitchPeriodInProcess, "failed. app upgrade switch period is in process")
}

func errZeroSwitchHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, govTypes.CodeInvalidHeight,
		fmt.Sprintf("failed. protocol switchHeight in AppUpgradeProposal isn't allowed to be 0"))
}

func errInvalidUpgradeThreshold(codespace sdk.CodespaceType, Threshold sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidUpgradeParams,
		fmt.Sprintf("failed. invalid Upgrade Threshold( "+Threshold.String()+" ) should be [0.75, 1)"))
}

func errInvalidLength(codespace sdk.CodespaceType, descriptor string, got, max int) sdk.Error {
	msg := fmt.Sprintf("failed. bad length for %v, got length %v, max is %v", descriptor, got, max)
	return sdk.NewError(codespace, CodeInvalidSoftWareDescri, msg)
}
