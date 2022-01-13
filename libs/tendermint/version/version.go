package version

import (
	"fmt"

	"github.com/tendermint/go-amino"
)

var (
	// GitCommit is the current HEAD set using ldflags.
	GitCommit string

	// Version is the built softwares version.
	Version = TMCoreSemVer
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}

const (
	// TMCoreSemVer is the current version of Tendermint Core.
	// It's the Semantic Version of the software.
	// Must be a string because scripts like dist.sh read this file.
	// XXX: Don't change the name of this variable or you will break
	// automation :)

	TMCoreSemVer = "0.33.9"

	// ABCISemVer is the semantic version of the ABCI library
	ABCISemVer  = "0.16.2"
	ABCIVersion = ABCISemVer
)

// Protocol is used for implementation agnostic versioning.
type Protocol uint64

// Uint64 returns the Protocol version as a uint64,
// eg. for compatibility with ABCI types.
func (p Protocol) Uint64() uint64 {
	return uint64(p)
}

var (
	// P2PProtocol versions all p2p behaviour and msgs.
	// This includes proposer selection.
	P2PProtocol Protocol = 7

	// BlockProtocol versions all block data structures and processing.
	// This includes validity of blocks and state updates.
	BlockProtocol Protocol = 10
)

//------------------------------------------------------------------------
// Version types

// App includes the protocol and software version for the application.
// This information is included in ResponseInfo. The App.Protocol can be
// updated in ResponseEndBlock.
type App struct {
	Protocol Protocol `json:"protocol"`
	Software string   `json:"software"`
}

// Consensus captures the consensus rules for processing a block in the blockchain,
// including all blockchain data structures and the rules of the application's
// state transition machine.
type Consensus struct {
	Block Protocol `json:"block"`
	App   Protocol `json:"app"`
}

func (c Consensus) AminoSize() int {
	var size int
	if c.Block != 0 {
		size += 1 + amino.UvarintSize(uint64(c.Block))
	}
	if c.App != 0 {
		size += 1 + amino.UvarintSize(uint64(c.App))
	}
	return size
}

func (c *Consensus) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	var dataLen uint64 = 0

	for {
		data = data[dataLen:]

		if len(data) == 0 {
			break
		}

		pos, aminoType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		data = data[1:]

		if aminoType != amino.Typ3_Varint {
			return fmt.Errorf("expected Varint, got %v", aminoType)
		}

		switch pos {
		case 1:
			var n int
			var uvint uint64
			uvint, n, err = amino.DecodeUvarint(data)
			c.Block = Protocol(uvint)
			dataLen = uint64(n)
		case 2:
			var n int
			var uvint uint64
			uvint, n, err = amino.DecodeUvarint(data)
			c.App = Protocol(uvint)
			dataLen = uint64(n)
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}
