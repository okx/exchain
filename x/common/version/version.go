package version

// ProtocolVersionType is the type of protocol version
type ProtocolVersionType int32

// const
const (
	ProtocolVersionV0      ProtocolVersionType = 0
	ProtocolVersionV1      ProtocolVersionType = 1
	CurrentProtocolVersion                     = ProtocolVersionV0
	Version                                    = "0"
)
