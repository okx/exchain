package types

import (
	"encoding/json"
	"strconv"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/wasm/internal/types/bytes"
)

const defaultLRUCacheSize = uint64(0)
const defaultQueryGasLimit = uint64(3000000)

// Model is a struct that holds a KV pair
type Model struct {
	// hex-encode key to read it better (this is often ascii)
	Key bytes.HexBytes `json:"key"`
	// base64-encode raw value
	Value []byte `json:"val"`
}

// CodeInfo is data for the uploaded contract WASM code
type CodeInfo struct {
	CodeHash []byte         `json:"code_hash"`
	Creator  sdk.AccAddress `json:"creator"`
	Source   string         `json:"source"`
	Builder  string         `json:"builder"`
}

// NewCodeInfo fills a new Contract struct
func NewCodeInfo(codeHash []byte, creator sdk.AccAddress, source string, builder string) CodeInfo {
	return CodeInfo{
		CodeHash: codeHash,
		Creator:  creator,
		Source:   source,
		Builder:  builder,
	}
}

// ContractInfo stores a WASM contract instance
type ContractInfo struct {
	CodeID  uint64          `json:"code_id"`
	Creator sdk.AccAddress  `json:"creator"`
	Label   string          `json:"label"`
	InitMsg json.RawMessage `json:"init_msg,omitempty"`
	// never show this in query results, just use for sorting
	// (Note: when using json tag "-" amino refused to serialize it...)
	Created *CreatedAt `json:"created,omitempty"`
}

// CreatedAt can be used to sort contracts
type CreatedAt struct {
	// BlockHeight is the block the contract was created at
	BlockHeight int64
	// TxIndex is a monotonic counter within the block (actual transaction index, or gas consumed)
	TxIndex uint64
}

// LessThan can be used to sort
func (a *CreatedAt) LessThan(b *CreatedAt) bool {
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}
	return a.BlockHeight < b.BlockHeight || (a.BlockHeight == b.BlockHeight && a.TxIndex < b.TxIndex)
}

// NewCreatedAt gets a timestamp from the context
func NewCreatedAt(ctx sdk.Context) *CreatedAt {
	// we must safely handle nil gas meters
	var index uint64
	meter := ctx.BlockGasMeter()
	if meter != nil {
		index = meter.GasConsumed()
	}
	return &CreatedAt{
		BlockHeight: ctx.BlockHeight(),
		TxIndex:     index,
	}
}

// NewContractInfo creates a new instance of a given WASM contract info
func NewContractInfo(codeID uint64, creator sdk.AccAddress, initMsg []byte, label string, createdAt *CreatedAt) ContractInfo {
	return ContractInfo{
		CodeID:  codeID,
		Creator: creator,
		InitMsg: initMsg,
		Label:   label,
		Created: createdAt,
	}
}

// NewEnv initializes the environment for a contract instance
func NewEnv(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAddr sdk.AccAddress) wasmTypes.Env {
	// safety checks before casting below
	if ctx.BlockHeight() < 0 {
		panic("Block height must never be negative")
	}
	if ctx.BlockTime().Unix() < 0 {
		panic("Block (unix) time must never be negative ")
	}
	env := wasmTypes.Env{
		Block: wasmTypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(ctx.BlockTime().Unix()),
			ChainID: ctx.ChainID(),
		},
		Message: wasmTypes.MessageInfo{
			Sender:    wasmTypes.CanonicalAddress(creator),
			SentFunds: NewWasmCoins(deposit),
		},
		Contract: wasmTypes.ContractInfo{
			Address: wasmTypes.CanonicalAddress(contractAddr),
		},
	}
	return env
}

// NewWasmCoins translates between Cosmos SDK coins and Wasm coins
func NewWasmCoins(cosmosCoins sdk.Coins) (wasmCoins []wasmTypes.Coin) {
	for _, coin := range cosmosCoins {
		wasmCoin := wasmTypes.Coin{
			Denom:  coin.Denom,
			Amount: ConvertSdkDecToString(coin.Amount),
		}
		wasmCoins = append(wasmCoins, wasmCoin)
	}
	return wasmCoins
}

const CustomEventType = "wasm"
const AttributeKeyContractAddr = "contract_address"

// CosmosResult converts from a Wasm Result type
func CosmosResult(wasmResult wasmTypes.Result, contractAddr sdk.AccAddress) sdk.Result {
	var events []sdk.Event
	if len(wasmResult.Log) > 0 {
		// we always tag with the contract address issuing this event
		attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddr, contractAddr.String())}
		for _, l := range wasmResult.Log {
			// and reserve the contract_address key for our use (not contract)
			if l.Key != AttributeKeyContractAddr {
				attr := sdk.NewAttribute(l.Key, l.Value)
				attrs = append(attrs, attr)
			}
		}
		events = []sdk.Event{sdk.NewEvent(CustomEventType, attrs...)}
	}
	return sdk.Result{
		Data:   []byte(wasmResult.Data),
		Events: events,
	}
}

// WasmConfig is the extra config required for wasm
type WasmConfig struct {
	SmartQueryGasLimit uint64 `mapstructure:"query_gas_limit"`
	CacheSize          uint64 `mapstructure:"lru_size"`
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWasmConfig() WasmConfig {
	return WasmConfig{
		SmartQueryGasLimit: defaultQueryGasLimit,
		CacheSize:          defaultLRUCacheSize,
	}
}

func ConvertSdkDecToString(amount sdk.Dec) string {
	value, err := strconv.ParseFloat(amount.String(), 64)
	if err != nil {
		value = 0
	}

	return strconv.FormatUint(uint64(value), 10)
}
