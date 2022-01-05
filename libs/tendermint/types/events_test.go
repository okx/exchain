package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryTxFor(t *testing.T) {
	tx := Tx("foo")
	height := int64(0)
	assert.Equal(t,
		fmt.Sprintf("tm.event='Tx' AND tx.hash='%X'", tx.Hash(height)),
		EventQueryTxFor(tx, height).String(),
	)
}

func TestQueryForEvent(t *testing.T) {
	assert.Equal(t,
		"tm.event='NewBlock'",
		QueryForEvent(EventNewBlock).String(),
	)
}
