package types

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	headerPrefix     = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	headerHashSuffix = []byte("n") // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
)

var indexer *Indexer

func init() {
	server.TrapSignal(func() {
		if indexer.backend.db != nil {
			indexer.backend.db.Close()
		}
	})
}

type Keeper interface {
	GetBlockBloom(ctx sdk.Context, height int64) ethtypes.Bloom
	GetHeightHash(ctx sdk.Context, height uint64) common.Hash
}

// Indexer does a post-processing job for equally sized sections of the
// canonical chain (like BlooomBits and CHT structures). A Indexer is
// connected to the blockchain through the event system by starting a
// ChainHeadEventLoop in a goroutine.
//
// Further child ChainIndexers can be added which use the output of the parent
// section indexer. These child indexers receive new head notifications only
// after an entire section has been finished or in case of rollbacks that might
// affect already finished sections.
type Indexer struct {
	backend BloomIndexer // Background processor generating the index data content

	update chan struct{} // Notification channel that headers should be processed
	quit   chan struct{} // Quit channel to tear down running goroutines

	sectionSize    uint64 // Number of blocks in a single chain segment to process
	storedSections uint64 // Number of sections successfully indexed into the database
}

func InitIndexer() {
	enableBloomFilter = viper.GetBool(FlagEnableBloomFilter)
	if !enableBloomFilter {
		return
	}

	indexer = &Indexer{
		backend:     initBloomIndexer(),
		update:      make(chan struct{}, 1),
		quit:        make(chan struct{}),
		sectionSize: BloomBitsBlocks,
	}
	indexer.setValidSections(indexer.getValidSections())
}

func GetIndexer() *Indexer {
	return indexer
}

func (i *Indexer) StoredSection() uint64 {
	if i != nil {
		return i.storedSections
	}
	return 0
}

func (i *Indexer) ProcessSection(ctx sdk.Context, k Keeper, height int64) error {
	interval := uint64(height - tmtypes.GetStartBlockHeight())
	// the hash of current block is stored when executing BeginBlock of next block.
	// so update section in the next block.
	if interval%i.sectionSize != 0 {
		return nil
	}
	section := i.storedSections
	lastHead := i.SectionHead(section)

	ctx.Logger().Debug("Processing new chain section", "section", section)

	// Reset and partial processing
	if err := i.backend.Reset(section); err != nil {
		i.setValidSections(0)
		return err
	}

	begin := section*i.sectionSize + uint64(tmtypes.GetStartBlockHeight())
	end := (section+1)*i.sectionSize + uint64(tmtypes.GetStartBlockHeight())

	for number := begin; number < end; number++ {
		var (
			bloom ethtypes.Bloom
			hash  common.Hash
		)
		// the initial height is 1 but it on ethereum is 0. so set the bloom and hash of the block 0 to empty.
		if number == uint64(tmtypes.GetStartBlockHeight()) {
			bloom = ethtypes.Bloom{}
			hash = common.Hash{}
		} else {
			hash = k.GetHeightHash(ctx, number)
			if hash == (common.Hash{}) {
				return fmt.Errorf("canonical block #%d unknown", number)
			}
			bloom = k.GetBlockBloom(ctx, int64(number))
		}
		if err := i.backend.Process(hash, number, bloom); err != nil {
			return err
		}
		lastHead = hash
	}
	if err := i.backend.Commit(); err != nil {
		return err
	}
	i.setSectionHead(section, lastHead)
	i.setValidSections(section + 1)
	return nil
}

// GetDB get db of BloomIndexer
func (b *Indexer) GetDB() dbm.DB {
	if b != nil {
		return b.backend.db
	}
	return nil
}

// setValidSections writes the number of valid sections to the index database
func (i *Indexer) setValidSections(sections uint64) {
	// Set the current number of valid sections in the database
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], sections)
	i.backend.db.Set([]byte("count"), data[:])

	// Remove any reorged sections, caching the valids in the mean time
	for i.storedSections > sections {
		i.storedSections--
		i.RemoveSectionHead(i.storedSections)
	}
	i.storedSections = sections // needed if new > old
}

// loadValidSections reads the number of valid sections from the index database
// and caches is into the local state.
func (i *Indexer) getValidSections() uint64 {
	data, _ := i.backend.db.Get([]byte("count"))
	if len(data) == 8 {
		return binary.BigEndian.Uint64(data)
	}
	return 0
}

// RemoveSectionHead removes the reference to a processed section from the index
// database.
func (i *Indexer) RemoveSectionHead(section uint64) {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], section)

	i.backend.db.Delete(append([]byte("shead"), data[:]...))
}

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// SectionHead retrieves the last block hash of a processed section from the
// index database.
func (i *Indexer) SectionHead(section uint64) common.Hash {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], section)

	hash, _ := i.backend.db.Get(append([]byte("shead"), data[:]...))
	if len(hash) == len(common.Hash{}) {
		return common.BytesToHash(hash)
	}
	return common.Hash{}
}

// setSectionHead writes the last block hash of a processed section to the index
// database.
func (i *Indexer) setSectionHead(section uint64, hash common.Hash) {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], section)

	i.backend.db.Set(append([]byte("shead"), data[:]...), hash.Bytes())
}

// removeSectionHead removes the reference to a processed section from the index
// database.
func (i *Indexer) removeSectionHead(section uint64) {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], section)

	i.backend.db.Delete(append([]byte("shead"), data[:]...))
}
