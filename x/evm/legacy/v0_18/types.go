package v018

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/legacy/v0_16"
)

type (
	// GenesisState defines the evm module genesis state
	GenesisState struct {
		Accounts                    []v0_16.GenesisAccount  `json:"accounts"`
		TxsLogs                     []v0_16.TransactionLogs `json:"txs_logs"`
		ContractDeploymentWhitelist AddressList             `json:"contract_deployment_whitelist"`
		ContractBlockedList         AddressList             `json:"contract_blocked_list"`
		ChainConfig                 v0_16.ChainConfig       `json:"chain_config"`
		Params                      Params                  `json:"params"`
	}

	AddressList []sdk.AccAddress

	Params struct {
		// EnableCreate toggles state transitions that use the vm.Create function
		EnableCreate bool `json:"enable_create" yaml:"enable_create"`
		// EnableCall toggles state transitions that use the vm.Call function
		EnableCall bool `json:"enable_call" yaml:"enable_call"`
		// ExtraEIPs defines the additional EIPs for the vm.Config
		ExtraEIPs []int `json:"extra_eips" yaml:"extra_eips"`
		// EnableContractDeploymentWhitelist controls the authorization of contract deployer
		EnableContractDeploymentWhitelist bool `json:"enable_contract_deployment_whitelist" yaml:"enable_contract_deployment_whitelist"`
		// EnableContractBlockedList controls the availability of contracts
		EnableContractBlockedList bool `json:"enable_contract_blocked_list" yaml:"enable_contract_blocked_list"`
		// MaxGasLimit defines the max gas limit in transaction
		MaxGasLimitPerTx uint64 `json:"max_gas_limit_per_tx" yaml:"max_gas_limit_per_tx"`
	}
)
