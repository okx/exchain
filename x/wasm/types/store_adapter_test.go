package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/flatkv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/gaskv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	db2 "github.com/okex/exchain/libs/tm-db"
	"testing"
)

func TestNewStoreAdapter(t *testing.T) {
	db := db2.NewMemDB()
	prefixDb := flatkv.NewStore(db)
	h := gaskv.NewStore()
	dbAdapter := NewStoreAdapter(prefix.NewStore(h, []byte("")))
}
