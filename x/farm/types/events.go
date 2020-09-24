package types

// farm module event types
const (
	// TODO: Create your event types
	EventTypeLock          = "lock"
	EventTypeUnlock            = "unlock"

	// TODO: Create keys fo your events, the values will be derivided from the msg
	AttributeKeyPool         = "pool"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
)
