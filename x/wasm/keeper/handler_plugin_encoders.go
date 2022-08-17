package keeper

import (
	"encoding/json"
	"fmt"
	ibcadapter "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	bank "github.com/okex/exchain/libs/cosmos-sdk/x/bank"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	ibcclienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"

	"github.com/okex/exchain/x/wasm/types"
)

type BankEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.BankMsg) ([]ibcadapter.Msg, error)
type CustomEncoder func(sender sdk.AccAddress, msg json.RawMessage) ([]ibcadapter.Msg, error)
type DistributionEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.DistributionMsg) ([]ibcadapter.Msg, error)
type StakingEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.StakingMsg) ([]ibcadapter.Msg, error)
type StargateEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.StargateMsg) ([]ibcadapter.Msg, error)
type WasmEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.WasmMsg) ([]ibcadapter.Msg, error)
type IBCEncoder func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]ibcadapter.Msg, error)

type MessageEncoders struct {
	Bank   func(sender sdk.AccAddress, msg *wasmvmtypes.BankMsg) ([]ibcadapter.Msg, error)
	Custom func(sender sdk.AccAddress, msg json.RawMessage) ([]ibcadapter.Msg, error)
	//Distribution func(sender sdk.AccAddress, msg *wasmvmtypes.DistributionMsg) ([]sdk.Msg, error)
	//IBC          func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]sdk.Msg, error)
	//Staking      func(sender sdk.AccAddress, msg *wasmvmtypes.StakingMsg) ([]sdk.Msg, error)
	Stargate func(sender sdk.AccAddress, msg *wasmvmtypes.StargateMsg) ([]ibcadapter.Msg, error)
	Wasm     func(sender sdk.AccAddress, msg *wasmvmtypes.WasmMsg) ([]ibcadapter.Msg, error)
	//Gov          func(sender sdk.AccAddress, msg *wasmvmtypes.GovMsg) ([]sdk.Msg, error)
}

func DefaultEncoders(unpacker codectypes.AnyUnpacker, portSource types.ICS20TransferPortSource) MessageEncoders {
	return MessageEncoders{
		Bank:   EncodeBankMsg,
		Custom: NoCustomMsg,
		//Distribution: EncodeDistributionMsg,
		//IBC:          EncodeIBCMsg(portSource),
		//Staking:      EncodeStakingMsg,
		Stargate: EncodeStargateMsg(unpacker),
		Wasm:     EncodeWasmMsg,
		//Gov:          EncodeGovMsg,
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
	//if o.Distribution != nil {
	//	e.Distribution = o.Distribution
	//}
	//if o.IBC != nil {
	//	e.IBC = o.IBC
	//}
	//if o.Staking != nil {
	//	e.Staking = o.Staking
	//}
	if o.Stargate != nil {
		e.Stargate = o.Stargate
	}
	if o.Wasm != nil {
		e.Wasm = o.Wasm
	}
	//if o.Gov != nil {
	//	e.Gov = o.Gov
	//}
	return e
}

func (e MessageEncoders) Encode(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]ibcadapter.Msg, error) {
	switch {
	case msg.Bank != nil:
		return e.Bank(contractAddr, msg.Bank)
	case msg.Custom != nil:
		return e.Custom(contractAddr, msg.Custom)
	//case msg.Distribution != nil:
	//	return e.Distribution(contractAddr, msg.Distribution)
	//case msg.IBC != nil:
	//	return e.IBC(ctx, contractAddr, contractIBCPortID, msg.IBC)
	//case msg.Staking != nil:
	//	return e.Staking(contractAddr, msg.Staking)
	case msg.Stargate != nil:
		return e.Stargate(contractAddr, msg.Stargate)
	case msg.Wasm != nil:
		return e.Wasm(contractAddr, msg.Wasm)
		//case msg.Gov != nil:
		//	return EncodeGovMsg(contractAddr, msg.Gov)
	}
	return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Wasm")
}

func EncodeBankMsg(sender sdk.AccAddress, msg *wasmvmtypes.BankMsg) ([]ibcadapter.Msg, error) {
	if msg.Send == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Bank")
	}
	if len(msg.Send.Amount) == 0 {
		return nil, nil
	}
	toSend, err := ConvertWasmCoinsToSdkCoins(msg.Send.Amount)
	if err != nil {
		return nil, err
	}

	sdkMsg := bank.MsgSendAdapter{
		FromAddress: sender.String(),
		ToAddress:   msg.Send.ToAddress,
		Amount:      toSend,
	}
	return []ibcadapter.Msg{&sdkMsg}, nil
}

func NoCustomMsg(sender sdk.AccAddress, msg json.RawMessage) ([]ibcadapter.Msg, error) {
	return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "custom variant not supported")
}

//func EncodeDistributionMsg(sender sdk.AccAddress, msg *wasmvmtypes.DistributionMsg) ([]ibcadapter.Msg, error) {
//	switch {
//	case msg.SetWithdrawAddress != nil:
//		withDrawAddress, err := sdk.AccAddressFromBech32(msg.SetWithdrawAddress.Address)
//		if err != nil {
//			return nil, err
//		}
//		setMsg := distributiontypes.MsgSetWithdrawAddress{
//			DelegatorAddress: sender,
//			WithdrawAddress:  withDrawAddress,
//		}
//		return []ibcadapter.Msg{&setMsg}, nil
//	case msg.WithdrawDelegatorReward != nil:
//		validatorAddress, err := sdk.AccAddressFromBech32(msg.WithdrawDelegatorReward.Validator)
//		if err != nil {
//			return nil, err
//		}
//		withdrawMsg := distributiontypes.MsgWithdrawDelegatorReward{
//			DelegatorAddress: sender,
//			ValidatorAddress: sdk.ValAddress(validatorAddress),
//		}
//		return []ibcadapter.Msg{&withdrawMsg}, nil
//	default:
//		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Distribution")
//	}
//}

//func EncodeStakingMsg(sender sdk.AccAddress, msg *wasmvmtypes.StakingMsg) ([]ibcadapter.Msg, error) {
//	switch {
//	case msg.Delegate != nil:
//		coin, err := ConvertWasmCoinToSdkCoin(msg.Delegate.Amount)
//		if err != nil {
//			return nil, err
//		}
//		validatorAddress, err := sdk.AccAddressFromBech32(msg.Delegate.Validator)
//		if err != nil {
//			return nil, err
//		}
//		sdkMsg := stakingtypes.MsgDelegate{
//			DelegatorAddress: sender,
//			ValidatorAddress: sdk.ValAddress(validatorAddress),
//			Amount:           coin,
//		}
//		return []ibcadapter.Msg{&sdkMsg}, nil
//
//	case msg.Redelegate != nil:
//		coin, err := ConvertWasmCoinToSdkCoin(msg.Redelegate.Amount)
//		if err != nil {
//			return nil, err
//		}
//		srcValidatorAddress, err := sdk.AccAddressFromBech32(msg.Redelegate.SrcValidator)
//		if err != nil {
//			return nil, err
//		}
//		dstValidatorAddress, err := sdk.AccAddressFromBech32(msg.Redelegate.DstValidator)
//		if err != nil {
//			return nil, err
//		}
//		sdkMsg := stakingtypes.MsgBeginRedelegate{
//			DelegatorAddress:    sender,
//			ValidatorSrcAddress: sdk.ValAddress(srcValidatorAddress),
//			ValidatorDstAddress: sdk.ValAddress(dstValidatorAddress),
//			Amount:              coin,
//		}
//		return []sdk.Msg{&sdkMsg}, nil
//	case msg.Undelegate != nil:
//		coin, err := ConvertWasmCoinToSdkCoin(msg.Undelegate.Amount)
//		if err != nil {
//			return nil, err
//		}
//		dstValidatorAddress, err := sdk.AccAddressFromBech32(msg.Undelegate.Validator)
//		if err != nil {
//			return nil, err
//		}
//		sdkMsg := stakingtypes.MsgUndelegate{
//			DelegatorAddress: sender,
//			ValidatorAddress: sdk.ValAddress(dstValidatorAddress),
//			Amount:           coin,
//		}
//		return []sdk.Msg{&sdkMsg}, nil
//	default:
//		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Staking")
//	}
//}

func EncodeStargateMsg(unpacker codectypes.AnyUnpacker) StargateEncoder {
	return func(sender sdk.AccAddress, msg *wasmvmtypes.StargateMsg) ([]ibcadapter.Msg, error) {
		any := codectypes.Any{
			TypeUrl: msg.TypeURL,
			Value:   msg.Value,
		}
		var sdkMsg ibcadapter.Msg
		if err := unpacker.UnpackAny(&any, &sdkMsg); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("Cannot unpack proto message with type URL: %s", msg.TypeURL))
		}
		if err := codectypes.UnpackInterfaces(sdkMsg, unpacker); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("UnpackInterfaces inside msg: %s", err))
		}
		return []ibcadapter.Msg{sdkMsg}, nil
	}
}

func EncodeWasmMsg(sender sdk.AccAddress, msg *wasmvmtypes.WasmMsg) ([]ibcadapter.Msg, error) {
	switch {
	case msg.Execute != nil:
		coins, err := ConvertWasmCoinsToSdkCoins(msg.Execute.Funds)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgExecuteContract{
			Sender:   sender.String(),
			Contract: msg.Execute.ContractAddr,
			Msg:      msg.Execute.Msg,
			Funds:    coins,
		}
		return []ibcadapter.Msg{&sdkMsg}, nil
	case msg.Instantiate != nil:
		coins, err := ConvertWasmCoinsToSdkCoins(msg.Instantiate.Funds)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgInstantiateContract{
			Sender: sender.String(),
			CodeID: msg.Instantiate.CodeID,
			Label:  msg.Instantiate.Label,
			Msg:    msg.Instantiate.Msg,
			Admin:  msg.Instantiate.Admin,
			Funds:  coins,
		}
		return []ibcadapter.Msg{&sdkMsg}, nil
	case msg.Migrate != nil:
		sdkMsg := types.MsgMigrateContract{
			Sender:   sender.String(),
			Contract: msg.Migrate.ContractAddr,
			CodeID:   msg.Migrate.NewCodeID,
			Msg:      msg.Migrate.Msg,
		}
		return []ibcadapter.Msg{&sdkMsg}, nil
	case msg.ClearAdmin != nil:
		sdkMsg := types.MsgClearAdmin{
			Sender:   sender.String(),
			Contract: msg.ClearAdmin.ContractAddr,
		}
		return []ibcadapter.Msg{&sdkMsg}, nil
	case msg.UpdateAdmin != nil:
		sdkMsg := types.MsgUpdateAdmin{
			Sender:   sender.String(),
			Contract: msg.UpdateAdmin.ContractAddr,
			NewAdmin: msg.UpdateAdmin.Admin,
		}
		return []ibcadapter.Msg{&sdkMsg}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Wasm")
	}
}

func EncodeIBCMsg(portSource types.ICS20TransferPortSource) func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]ibcadapter.Msg, error) {
	return func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]ibcadapter.Msg, error) {
		switch {
		case msg.CloseChannel != nil:
			return []ibcadapter.Msg{&channeltypes.MsgChannelCloseInit{
				PortId:    PortIDForContract(sender),
				ChannelId: msg.CloseChannel.ChannelID,
				Signer:    sender.String(),
			}}, nil
		case msg.Transfer != nil:
			amount, err := ConvertWasmCoinToSdkCoin(msg.Transfer.Amount)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "amount")
			}
			msg := &ibctransfertypes.MsgTransfer{
				SourcePort:       portSource.GetPort(ctx),
				SourceChannel:    msg.Transfer.ChannelID,
				Token:            amount,
				Sender:           sender.String(),
				Receiver:         msg.Transfer.ToAddress,
				TimeoutHeight:    ConvertWasmIBCTimeoutHeightToCosmosHeight(msg.Transfer.Timeout.Block),
				TimeoutTimestamp: msg.Transfer.Timeout.Timestamp,
			}
			return []ibcadapter.Msg{msg}, nil
		default:
			return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "Unknown variant of IBC")
		}
	}
}

//func EncodeGovMsg(sender sdk.AccAddress, msg *wasmvmtypes.GovMsg) ([]ibcadapter.Msg, error) {
//	var option govtypes.VoteOption
//	switch msg.Vote.Vote {
//	case wasmvmtypes.Yes:
//		option = govtypes.OptionYes
//	case wasmvmtypes.No:
//		option = govtypes.OptionNo
//	case wasmvmtypes.NoWithVeto:
//		option = govtypes.OptionNoWithVeto
//	case wasmvmtypes.Abstain:
//		option = govtypes.OptionAbstain
//	}
//	vote := &govtypes.MsgVote{
//		ProposalID: msg.Vote.ProposalId,
//		Voter:      sender,
//		Option:     option,
//	}
//	return []ibcadapter.Msg{vote}, nil
//}

// ConvertWasmIBCTimeoutHeightToCosmosHeight converts a wasmvm type ibc timeout height to ibc module type height
func ConvertWasmIBCTimeoutHeightToCosmosHeight(ibcTimeoutBlock *wasmvmtypes.IBCTimeoutBlock) ibcclienttypes.Height {
	if ibcTimeoutBlock == nil {
		return ibcclienttypes.NewHeight(0, 0)
	}
	return ibcclienttypes.NewHeight(ibcTimeoutBlock.Revision, ibcTimeoutBlock.Height)
}

// ConvertWasmCoinsToSdkCoins converts the wasm vm type coins to sdk type coins
func ConvertWasmCoinsToSdkCoins(coins []wasmvmtypes.Coin) (sdk.CoinAdapters, error) {
	var toSend sdk.CoinAdapters
	for _, coin := range coins {
		c, err := ConvertWasmCoinToSdkCoin(coin)
		if err != nil {
			return nil, err
		}
		toSend = append(toSend, c)
	}
	return toSend, nil
}

// ConvertWasmCoinToSdkCoin converts a wasm vm type coin to sdk type coin
func ConvertWasmCoinToSdkCoin(coin wasmvmtypes.Coin) (sdk.CoinAdapter, error) {
	amount, ok := sdk.NewIntFromString(coin.Amount)
	if !ok {
		return sdk.CoinAdapter{}, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, coin.Amount+coin.Denom)
	}
	r := sdk.CoinAdapter{
		Denom:  coin.Denom,
		Amount: amount,
	}
	if err := sdk.ValidateDenom(coin.Denom); err != nil {
		return sdk.CoinAdapter{}, err
	}
	if r.IsNegative() {
		return sdk.CoinAdapter{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, coin.Amount+coin.Denom)
	}
	return r, nil
}
