package v0_16

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type (
	// GenesisState defines the evm module genesis state
	GenesisState struct {
		Accounts    []GenesisAccount  `json:"accounts"`
		TxsLogs     []TransactionLogs `json:"txs_logs"`
		ChainConfig ChainConfig       `json:"chain_config"`
		Params      Params            `json:"params"`
	}

	// GenesisAccount defines an account to be initialized in the genesis state.
	// Its main difference between with Geth's GenesisAccount is that it uses a custom
	// storage type and that it doesn't contain the private key field.
	// NOTE: balance is omitted as it is imported from the auth account balance.
	GenesisAccount struct {
		Address string        `json:"address"`
		Code    hexutil.Bytes `json:"code,omitempty"`
		Storage Storage       `json:"storage,omitempty"`
	}

	Storage []State

	State struct {
		Key   ethcmn.Hash `json:"key"`
		Value ethcmn.Hash `json:"value"`
	}

	TransactionLogs struct {
		Hash ethcmn.Hash     `json:"hash"`
		Logs []*ethtypes.Log `json:"logs"`
	}

	ChainConfig struct {
		HomesteadBlock sdk.Int `json:"homestead_block" yaml:"homestead_block"` // Homestead switch block (< 0 no fork, 0 = already homestead)

		DAOForkBlock   sdk.Int `json:"dao_fork_block" yaml:"dao_fork_block"`     // TheDAO hard-fork switch block (< 0 no fork)
		DAOForkSupport bool    `json:"dao_fork_support" yaml:"dao_fork_support"` // Whether the nodes supports or opposes the DAO hard-fork

		// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
		EIP150Block sdk.Int `json:"eip150_block" yaml:"eip150_block"` // EIP150 HF block (< 0 no fork)
		EIP150Hash  string  `json:"eip150_hash" yaml:"eip150_hash"`   // EIP150 HF hash (needed for header only clients as only gas pricing changed)

		EIP155Block sdk.Int `json:"eip155_block" yaml:"eip155_block"` // EIP155 HF block
		EIP158Block sdk.Int `json:"eip158_block" yaml:"eip158_block"` // EIP158 HF block

		ByzantiumBlock      sdk.Int `json:"byzantium_block" yaml:"byzantium_block"`           // Byzantium switch block (< 0 no fork, 0 = already on byzantium)
		ConstantinopleBlock sdk.Int `json:"constantinople_block" yaml:"constantinople_block"` // Constantinople switch block (< 0 no fork, 0 = already activated)
		PetersburgBlock     sdk.Int `json:"petersburg_block" yaml:"petersburg_block"`         // Petersburg switch block (< 0 same as Constantinople)
		IstanbulBlock       sdk.Int `json:"istanbul_block" yaml:"istanbul_block"`             // Istanbul switch block (< 0 no fork, 0 = already on istanbul)
		MuirGlacierBlock    sdk.Int `json:"muir_glacier_block" yaml:"muir_glacier_block"`     // Eip-2384 (bomb delay) switch block (< 0 no fork, 0 = already activated)

		//YoloV2Block sdk.Int `json:"yoloV2_block" yaml:"yoloV2_block"` // YOLO v1: https://github.com/ethereum/EIPs/pull/2657 (Ephemeral testnet)
		//EWASMBlock  sdk.Int `json:"ewasm_block" yaml:"ewasm_block"`   // EWASM switch block (< 0 no fork, 0 = already activated)
	}

	Params struct {
		// EVMDenom defines the token denomination used for state transitions on the
		// EVM module.
		EvmDenom string `json:"evm_denom" yaml:"evm_denom"`
		// EnableCreate toggles state transitions that use the vm.Create function
		EnableCreate bool `json:"enable_create" yaml:"enable_create"`
		// EnableCall toggles state transitions that use the vm.Call function
		EnableCall bool `json:"enable_call" yaml:"enable_call"`
		// ExtraEIPs defines the additional EIPs for the vm.Config
		ExtraEIPs []int `json:"extra_eips" yaml:"extra_eips"`
	}
)
