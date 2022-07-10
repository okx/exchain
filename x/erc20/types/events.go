package types

const (
	EventTypDeployModuleERC20 = "deploy_erc20_contract"
	EventTypCallModuleERC20   = "call_erc20_contract"
	EventTypLock              = "lock"
	EventTypUnlock            = "unlock"
	EventTypMint              = "mint"
	EventTypBurn              = "burn"

	AttributeKeyContractAddr   = "contract_address"
	AttributeKeyContractMethod = "contract_method"
	AttributeKeyFrom           = "from"
	AttributeKeyTo             = "to"
)
