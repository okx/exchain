package types

import "github.com/ethereum/go-ethereum/accounts/abi"

const (
	SendToWasmEventName  = "__OKCSendToWasm"
	WasmCalledMethodName = "mintCW20"

	SendToEvmSubMsgName = "__OKCSendToEvm"
	EvmCalledMethodName = "mintERC20"
)

// SendToWasmEventName represent the signature of
// `event __SendToWasmEventName(string wasmAddr,string recipient, string amount)`
var SendToWasmEvent abi.Event

func init() {
	stringType, _ := abi.NewType("string", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)

	SendToWasmEvent = abi.NewEvent(
		SendToWasmEventName,
		SendToWasmEventName,
		false,
		abi.Arguments{
			abi.Argument{
				Name:    "wasmAddr",
				Type:    stringType,
				Indexed: false,
			},
			abi.Argument{
				Name:    "recipient",
				Type:    stringType,
				Indexed: false,
			},
			abi.Argument{
				Name:    "amount",
				Type:    uint256Type,
				Indexed: false,
			},
		},
	)
}
