package types

type NodeMode string

const (
	// node mode values
	RpcNode       NodeMode = "rpc"
	ValidatorNode NodeMode = "validator"
	ArchiveNode   NodeMode = "archive"
	InnertxNode   NodeMode = "innertx"

	// node mode flag
	FlagNodeMode = "node-mode"
)
