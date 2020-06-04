package types

import (
	"encoding/json"
	"net/url"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxWasmSize = 500 * 1024

	// MaxLabelSize is the longest label that can be used when Instantiating a contract
	MaxLabelSize = 128

	// BuildTagRegexp is a docker image regexp.
	// We only support max 128 characters, with at least one organization name (subset of all legal names).
	//
	// Details from https://docs.docker.com/engine/reference/commandline/tag/#extended-description :
	//
	// An image name is made up of slash-separated name components (optionally prefixed by a registry hostname).
	// Name components may contain lowercase characters, digits and separators.
	// A separator is defined as a period, one or two underscores, or one or more dashes. A name component may not start or end with a separator.
	//
	// A tag name must be valid ASCII and may contain lowercase and uppercase letters, digits, underscores, periods and dashes.
	// A tag name may not start with a period or a dash and may contain a maximum of 128 characters.
	BuildTagRegexp = "^[a-z0-9][a-z0-9._-]*[a-z0-9](/[a-z0-9][a-z0-9._-]*[a-z0-9])+:[a-zA-Z0-9_][a-zA-Z0-9_.-]*$"

	MaxBuildTagSize = 128
)

type MsgStoreCode struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
	// WASMByteCode can be raw or gzip compressed
	WASMByteCode []byte `json:"wasm_byte_code" yaml:"wasm_byte_code"`
	// Source is a valid absolute HTTPS URI to the contract's source code, optional
	Source string `json:"source" yaml:"source"`
	// Builder is a valid docker image name with tag, optional
	Builder string `json:"builder" yaml:"builder"`
}

func (msg MsgStoreCode) Route() string {
	return RouterKey
}

func (msg MsgStoreCode) Type() string {
	return "store-code"
}

func (msg MsgStoreCode) ValidateBasic() sdk.Error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdk.ErrUnknownRequest(err.Error())
	}

	if len(msg.WASMByteCode) == 0 {
		return sdk.ErrUnknownRequest("empty wasm code")
	}

	if len(msg.WASMByteCode) > MaxWasmSize {
		return sdk.ErrUnknownRequest("wasm code too large")
	}

	if msg.Source != "" {
		u, err := url.Parse(msg.Source)
		if err != nil {
			return sdk.ErrUnknownRequest("source should be a valid url")
		}
		if !u.IsAbs() {
			return sdk.ErrUnknownRequest("source should be an absolute url")
		}
		if u.Scheme != "https" {
			return sdk.ErrUnknownRequest("source must use https")
		}
	}

	return validateBuilder(msg.Builder)
}

func (msg MsgStoreCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func validateBuilder(buildTag string) sdk.Error {
	if len(buildTag) > MaxBuildTagSize {
		return sdk.ErrUnknownRequest("builder tag longer than 128 characters")
	}

	if buildTag != "" {
		ok, err := regexp.MatchString(BuildTagRegexp, buildTag)
		if err != nil || !ok {
			return sdk.ErrUnknownRequest("invalid tag supplied for builder")
		}
	}

	return nil
}

type MsgInstantiateContract struct {
	Sender    sdk.AccAddress  `json:"sender" yaml:"sender"`
	Code      uint64          `json:"code_id" yaml:"code_id"`
	Label     string          `json:"label" yaml:"label"`
	InitMsg   json.RawMessage `json:"init_msg" yaml:"init_msg"`
	InitFunds sdk.Coins       `json:"init_funds" yaml:"init_funds"`
}

func (msg MsgInstantiateContract) Route() string {
	return RouterKey
}

func (msg MsgInstantiateContract) Type() string {
	return "instantiate"
}

func (msg MsgInstantiateContract) ValidateBasic() sdk.Error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdk.ErrUnknownRequest(err.Error())
	}

	if msg.Code == 0 {
		return sdk.ErrUnknownRequest("code_id is required")
	}
	if msg.Label == "" {
		return sdk.ErrUnknownRequest("label is required")
	}
	if len(msg.Label) > MaxLabelSize {
		return sdk.ErrUnknownRequest("label cannot be longer than 128 characters")
	}

	if msg.InitFunds.IsAnyNegative() {
		return sdk.ErrInvalidCoins("negative SentFunds")
	}
	return nil
}

func (msg MsgInstantiateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgInstantiateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgExecuteContract struct {
	Sender    sdk.AccAddress  `json:"sender" yaml:"sender"`
	Contract  sdk.AccAddress  `json:"contract" yaml:"contract"`
	Msg       json.RawMessage `json:"msg" yaml:"msg"`
	SentFunds sdk.Coins       `json:"sent_funds" yaml:"sent_funds"`
}

func (msg MsgExecuteContract) Route() string {
	return RouterKey
}

func (msg MsgExecuteContract) Type() string {
	return "execute"
}

func (msg MsgExecuteContract) ValidateBasic() sdk.Error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdk.ErrUnknownRequest(err.Error())
	}
	if err := sdk.VerifyAddressFormat(msg.Contract); err != nil {
		return sdk.ErrUnknownRequest(err.Error())
	}

	if msg.SentFunds.IsAnyNegative() {
		return sdk.ErrInvalidCoins("negative SentFunds")
	}
	return nil
}

func (msg MsgExecuteContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgExecuteContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
