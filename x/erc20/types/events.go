package types

const (
	EventTypDeployModuleERC20 = "deploy_erc20_contract"
	EventTypCallModuleERC20   = "call_erc20_contract"
	EventTypLock              = "erc20_lock"
	EventTypUnlock            = "erc20_unlock"
	EventTypMint              = "erc20_mint"
	EventTypBurn              = "erc20_burn"

	AttributeKeyContractAddr   = "contract_address"
	AttributeKeyContractMethod = "contract_method"
	AttributeKeyFrom           = "from"
	AttributeKeyTo             = "to"
)
