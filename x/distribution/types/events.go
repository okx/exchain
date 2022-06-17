package types

// distribution module event types
const (
	EventTypeSetWithdrawAddress = "set_withdraw_address"
	EventTypeCommission         = "commission"
	EventTypeWithdrawCommission = "withdraw_commission"
	EventTypeRewards            = "rewards"

	EventTypeWithdrawRewards = "withdraw_rewards"
	EventTypeProposerReward  = "proposer_reward"

	AttributeKeyWithdrawAddress = "withdraw_address"
	AttributeKeyValidator       = "validator"

	AttributeValueCategory = ModuleName
)
