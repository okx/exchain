package types

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"path/filepath"
	"sync"
	"sync/atomic"
)

var (
	indexer           *Indexer
	enableBloomFilter bool
	once              sync.Once
)

type Keeper interface {
	GetBlockBloom(ctx sdk.Context, height int64) ethtypes.Bloom
	GetHeightHash(ctx sdk.Context, height uint64) common.Hash
}

func init() {
	server.TrapSignal(func() {
		if indexer != nil && indexer.backend.db != nil {
			indexer.backend.db.Close()
		}
	})
}

func GetEnableBloomFilter() bool {
	once.Do(func() {
		enableBloomFilter = viper.GetBool(FlagEnableBloomFilter)
	})
	return enableBloomFilter
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
	backend bloomIndexer // Background processor generating the index data content

	update chan sdk.Context // Notification channel that headers should be processed
	quit   chan struct{}    // Quit channel to tear down running goroutines

	storedSections uint64 // Number of sections successfully indexed into the database
	processing     uint32 // Atomic flag whether indexer is processing or not
}

func InitIndexer(db dbm.DB) {
	if !enableBloomFilter {
		return
	}

	indexer = &Indexer{
		backend: initBloomIndexer(db),
		update:  make(chan sdk.Context),
		quit:    make(chan struct{}),
	}
	indexer.setValidSections(indexer.GetValidSections())
}

func BloomDb() dbm.DB {
	dataDir := filepath.Join(viper.GetString("home"), "data")
	var err error
	db, err := sdk.NewLevelDB(bloomDir, dataDir)
	if err != nil {
		panic(err)
	}
	return db
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

func (i *Indexer) IsProcessing() bool {
	return atomic.LoadUint32(&i.processing) == 1
}

func (i *Indexer) ProcessSection(ctx sdk.Context, k Keeper, interval uint64) {
	if atomic.SwapUint32(&i.processing, 1) == 1 {
		ctx.Logger().Error("matcher is already running")
		return
	}
	defer atomic.StoreUint32(&i.processing, 0)
	knownSection := interval / BloomBitsBlocks
	for i.storedSections < knownSection {
		section := i.storedSections
		var lastHead common.Hash
		if section > 0 {
			lastHead = i.sectionHead(section - 1)
		}
		ctx.Logger().Debug("Processing new chain section", "section", section)

		// Reset and partial processing
		if err := i.backend.Reset(section); err != nil {
			i.setValidSections(0)
			ctx.Logger().Error(err.Error())
			return
		}

		begin := section*BloomBitsBlocks + uint64(tmtypes.GetStartBlockHeight())
		end := (section+1)*BloomBitsBlocks + uint64(tmtypes.GetStartBlockHeight())

		for number := begin; number < end; number++ {
			var (
				bloom ethtypes.Bloom
				hash  common.Hash
			)
			ctx = i.updateCtx(ctx)
			// the initial height is 1 but it on ethereum is 0. so set the bloom and hash of the block 0 to empty.
			if number == uint64(tmtypes.GetStartBlockHeight()) {
				bloom = ethtypes.Bloom{}
				hash = common.Hash{}
			} else {
				hash = k.GetHeightHash(ctx, number)
				if hash == (common.Hash{}) {
					ctx.Logger().Error("canonical block #%d unknown", number)
					return
				}
				bloom = k.GetBlockBloom(ctx, int64(number))
			}
			if err := i.backend.Process(hash, number, bloom); err != nil {
				ctx.Logger().Error(err.Error())
				return
			}
			lastHead = hash
		}
		if err := i.backend.Commit(); err != nil {
			ctx.Logger().Error(err.Error())
			return
		}
		i.setSectionHead(section, lastHead)
		i.setValidSections(section + 1)
	}
}

// GetDB get db of bloomIndexer
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
		i.removeSectionHead(i.storedSections)
	}
	i.storedSections = sections // needed if new > old
}

// GetValidSections reads the number of valid sections from the index database
// and caches is into the local state.
func (i *Indexer) GetValidSections() uint64 {
	data, _ := i.backend.db.Get([]byte("count"))
	if len(data) == 8 {
		return binary.BigEndian.Uint64(data)
	}
	return 0
}

// sectionHead retrieves the last block hash of a processed section from the
// index database.
func (i *Indexer) sectionHead(section uint64) common.Hash {
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

func (i *Indexer) NotifyNewHeight(ctx sdk.Context) {
	i.update <- ctx
}

func (i *Indexer) updateCtx(oldCtx sdk.Context) sdk.Context {
	newCtx := oldCtx

	select {
	case newCtx = <-i.update:
	default:
	}

	return newCtx
}