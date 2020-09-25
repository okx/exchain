package types

// farm module event types
const (
	EventTypeCreatePool = "create-pool"
	EventTypeProvide    = "provide"
	EventTypeLock       = "lock"
	EventTypeUnlock     = "unlock"
	EventTypeClaim      = "claim"

	// TODO: Create keys fo your events, the values will be derivided from the msg
	AttributeKeyAddress            = "address"
	AttributeKeyPool               = "pool"
	AttributeKeyStartHeightToYield = "start_height_to_yield"
	AttributeKeyYiledPerBlock      = "yield_per_block"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
)
