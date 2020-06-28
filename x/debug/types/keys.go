package types

const (
	// ModuleName is the name of the debug module
	ModuleName = "debug"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the debug module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the debug module
	QuerierRoute = ModuleName

	DumpStore   = "dump"
	SetLogLevel = "set-loglevel"
	SanityCheck = "sanity-check"
)
