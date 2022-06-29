package types

// staking module event types
const (
	EventTypeCompleteUnbonding           = "complete_unbonding"
	EventTypeCreateValidator             = "create_validator"
	EventTypeEditValidator               = "edit_validator"
	EventTypeEditValidatorCommissionRate = "edit_validator_commission_rate"
	EventTypeDelegate                    = "delegate"
	EventTypeUnbond                      = "unbond"

	AttributeKeyValidator         = "validator"
	AttributeKeyCommissionRate    = "commission_rate"
	AttributeKeyMinSelfDelegation = "min_self_delegation"
	AttributeKeyDelegator         = "delegator"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeValueCategory        = ModuleName

	EventTypeAddShares = "add_shares"

	AttributeKeyValidatorToAddShares = "validator_to_add_shares"
	AttributeKeyShares               = "shares"
)
