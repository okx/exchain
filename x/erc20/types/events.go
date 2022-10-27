package types

const (
	EventTypDeployModuleERC20 = "deploy_erc20_contract"
	EventTypCallModuleERC20   = "call_erc20_contract"
	EventTypLock              = "erc20_lock"
	EventTypBurn              = "erc20_burn"

	AttributeKeyContractAddr   = "contract_address"
	AttributeKeyContractMethod = "contract_method"
	AttributeKeyFrom           = "from"
	AttributeKeyTo             = "to"

	InnerTxUnlock    = "erc20-unlock"
	InnerTxMint      = "erc20-mint"
	InnerTxSendToIbc = "erc20-send-to-ibc"
)
