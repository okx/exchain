package keeper

import (
	"encoding/json"
	"fmt"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/okex/okchain/x/wasm/internal/types"
)

type MessageHandler struct {
	router   sdk.Router
	encoders MessageEncoders
}

func NewMessageHandler(router sdk.Router, customEncoders *MessageEncoders) MessageHandler {
	encoders := DefaultEncoders().Merge(customEncoders)
	return MessageHandler{
		router:   router,
		encoders: encoders,
	}
}

type BankEncoder func(sender sdk.AccAddress, msg *wasmTypes.BankMsg) ([]sdk.Msg, error)
type CustomEncoder func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error)
type StakingEncoder func(sender sdk.AccAddress, msg *wasmTypes.StakingMsg) ([]sdk.Msg, error)
type WasmEncoder func(sender sdk.AccAddress, msg *wasmTypes.WasmMsg) ([]sdk.Msg, error)

type MessageEncoders struct {
	Bank    BankEncoder
	Custom  CustomEncoder
	Staking StakingEncoder
	Wasm    WasmEncoder
}

func DefaultEncoders() MessageEncoders {
	return MessageEncoders{
		Bank:    EncodeBankMsg,
		Custom:  NoCustomMsg,
		Staking: EncodeStakingMsg,
		Wasm:    EncodeWasmMsg,
	}
}

func (e MessageEncoders) Merge(o *MessageEncoders) MessageEncoders {
	if o == nil {
		return e
	}
	if o.Bank != nil {
		e.Bank = o.Bank
	}
	if o.Custom != nil {
		e.Custom = o.Custom
	}
	if o.Staking != nil {
		e.Staking = o.Staking
	}
	if o.Wasm != nil {
		e.Wasm = o.Wasm
	}
	return e
}

func (e MessageEncoders) Encode(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, error) {
	switch {
	case msg.Bank != nil:
		return e.Bank(contractAddr, msg.Bank)
	case msg.Custom != nil:
		return e.Custom(contractAddr, msg.Custom)
	case msg.Staking != nil:
		return e.Staking(contractAddr, msg.Staking)
	case msg.Wasm != nil:
		return e.Wasm(contractAddr, msg.Wasm)
	}
	return nil, types.ErrInvalidMsg("Unknown variant of Wasm")
}

func EncodeBankMsg(sender sdk.AccAddress, msg *wasmTypes.BankMsg) ([]sdk.Msg, error) {
	if msg.Send == nil {
		return nil, types.ErrInvalidMsg("Unknown variant of Bank")
	}
	if len(msg.Send.Amount) == 0 {
		return nil, nil
	}
	fromAddr, stderr := sdk.AccAddressFromBech32(msg.Send.FromAddress)
	if stderr != nil {
		return nil, sdk.ErrInvalidAddress(msg.Send.FromAddress)
	}
	toAddr, stderr := sdk.AccAddressFromBech32(msg.Send.ToAddress)
	if stderr != nil {
		return nil, sdk.ErrInvalidAddress(msg.Send.ToAddress)
	}
	toSend, err := convertWasmCoinsToSdkCoins(msg.Send.Amount)
	if err != nil {
		return nil, err
	}
	sdkMsg := bank.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      toSend,
	}
	return []sdk.Msg{sdkMsg}, nil
}

func NoCustomMsg(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	return nil, types.ErrInvalidMsg("Custom variant not supported")
}

func EncodeStakingMsg(sender sdk.AccAddress, msg *wasmTypes.StakingMsg) ([]sdk.Msg, error) {
	if msg.Delegate != nil {
		validator, err := sdk.ValAddressFromBech32(msg.Delegate.Validator)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(msg.Delegate.Validator)
		}
		coin, err := convertWasmCoinToSdkCoin(msg.Delegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := staking.MsgDelegate{
			DelegatorAddress: sender,
			ValidatorAddress: validator,
			Amount:           coin,
		}
		return []sdk.Msg{sdkMsg}, nil
	}
	if msg.Redelegate != nil {
		src, err := sdk.ValAddressFromBech32(msg.Redelegate.SrcValidator)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(msg.Redelegate.SrcValidator)
		}
		dst, err := sdk.ValAddressFromBech32(msg.Redelegate.DstValidator)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(msg.Redelegate.DstValidator)
		}
		coin, err := convertWasmCoinToSdkCoin(msg.Redelegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := staking.MsgBeginRedelegate{
			DelegatorAddress:    sender,
			ValidatorSrcAddress: src,
			ValidatorDstAddress: dst,
			Amount:              coin,
		}
		return []sdk.Msg{sdkMsg}, nil
	}
	if msg.Undelegate != nil {
		validator, err := sdk.ValAddressFromBech32(msg.Undelegate.Validator)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(msg.Undelegate.Validator)
		}
		coin, err := convertWasmCoinToSdkCoin(msg.Undelegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := staking.MsgUndelegate{
			DelegatorAddress: sender,
			ValidatorAddress: validator,
			Amount:           coin,
		}
		return []sdk.Msg{sdkMsg}, nil
	}
	if msg.Withdraw != nil {
		var err error
		rcpt := sender
		if len(msg.Withdraw.Recipient) != 0 {
			rcpt, err = sdk.AccAddressFromBech32(msg.Withdraw.Recipient)
			if err != nil {
				return nil, sdk.ErrInvalidAddress(msg.Withdraw.Recipient)
			}
		}
		validator, err := sdk.ValAddressFromBech32(msg.Withdraw.Validator)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(msg.Withdraw.Validator)
		}
		setMsg := distribution.MsgSetWithdrawAddress{
			DelegatorAddress: sender,
			WithdrawAddress:  rcpt,
		}
		withdrawMsg := distribution.MsgWithdrawDelegatorReward{
			DelegatorAddress: sender,
			ValidatorAddress: validator,
		}
		return []sdk.Msg{setMsg, withdrawMsg}, nil
	}
	return nil, types.ErrInvalidMsg("Unknown variant of Staking")
}

func EncodeWasmMsg(sender sdk.AccAddress, msg *wasmTypes.WasmMsg) ([]sdk.Msg, error) {
	if msg.Execute != nil {
		contractAddr, err := sdk.AccAddressFromBech32(msg.Execute.ContractAddr)
		if err != nil {
			return nil, sdk.ErrInvalidAddress(msg.Execute.ContractAddr)
		}
		coins, err := convertWasmCoinsToSdkCoins(msg.Execute.Send)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgExecuteContract{
			Sender:    sender,
			Contract:  contractAddr,
			Msg:       msg.Execute.Msg,
			SentFunds: coins,
		}
		return []sdk.Msg{sdkMsg}, nil
	}
	if msg.Instantiate != nil {
		coins, err := convertWasmCoinsToSdkCoins(msg.Instantiate.Send)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgInstantiateContract{
			Sender: sender,
			Code:   msg.Instantiate.CodeID,
			// TODO: add this to CosmWasm
			Label:     fmt.Sprintf("Auto-created by %s", sender),
			InitMsg:   msg.Instantiate.Msg,
			InitFunds: coins,
		}
		return []sdk.Msg{sdkMsg}, nil
	}
	return nil, types.ErrInvalidMsg("Unknown variant of Wasm")
}

func (h MessageHandler) Dispatch(ctx sdk.Context, contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) error {
	sdkMsgs, err := h.encoders.Encode(contractAddr, msg)
	if err != nil {
		return err
	}
	for _, sdkMsg := range sdkMsgs {
		if err := h.handleSdkMessage(ctx, contractAddr, sdkMsg); err != nil {
			return err
		}
	}
	return nil
}

func (h MessageHandler) handleSdkMessage(ctx sdk.Context, contractAddr sdk.Address, msg sdk.Msg) error {
	// make sure this account can send it
	for _, acct := range msg.GetSigners() {
		if !acct.Equals(contractAddr) {
			return sdk.ErrUnauthorized("contract doesn't have permission")
		}
	}

	// find the handler and execute it
	handler := h.router.Route(msg.Route())
	if handler == nil {
		return sdk.ErrUnknownRequest(msg.Route())
	}
	res := handler(ctx, msg)
	if !res.IsOK() {
		return sdk.ErrInternal("fail to execute: " + msg.Route())
	}
	// redispatch all events, (type sdk.EventTypeMessage will be filtered out in the handler)
	ctx.EventManager().EmitEvents(res.Events)

	return nil
}

func convertWasmCoinsToSdkCoins(coins []wasmTypes.Coin) (sdk.Coins, error) {
	var toSend sdk.Coins
	for _, coin := range coins {
		c, err := convertWasmCoinToSdkCoin(coin)
		if err != nil {
			return nil, err
		}
		toSend = append(toSend, c)
	}
	return toSend, nil
}

func convertWasmCoinToSdkCoin(coin wasmTypes.Coin) (sdk.Coin, error) {
	amount, ok := sdk.NewIntFromString(coin.Amount)
	if !ok {
		return sdk.Coin{}, sdk.ErrInvalidCoins(coin.Amount + coin.Denom)
	}
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: amount.ToDec(),
	}, nil
}
