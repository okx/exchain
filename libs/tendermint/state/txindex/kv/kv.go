package kv

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/orderedcode"
	"github.com/pkg/errors"

	dbm "github.com/okex/exchain/libs/tm-db"

	"github.com/okex/exchain/libs/tendermint/libs/pubsub/query"
	"github.com/okex/exchain/libs/tendermint/libs/pubsub/query/syntax"
	tmstring "github.com/okex/exchain/libs/tendermint/libs/strings"
	"github.com/okex/exchain/libs/tendermint/state/txindex"
	"github.com/okex/exchain/libs/tendermint/types"
)

const (
	tagKeySeparator = "/"
)

var _ txindex.TxIndexer = (*TxIndex)(nil)

// TxIndex is the simplest possible indexer, backed by key-value storage (levelDB).
type TxIndex struct {
	store                dbm.DB
	compositeKeysToIndex []string
	indexAllEvents       bool
}

// NewTxIndex creates new KV indexer.
func NewTxIndex(store dbm.DB, options ...func(*TxIndex)) *TxIndex {
	txi := &TxIndex{store: store, compositeKeysToIndex: make([]string, 0), indexAllEvents: false}
	for _, o := range options {
		o(txi)
	}
	return txi
}

// IndexEvents is an option for setting which composite keys to index.
func IndexEvents(compositeKeys []string) func(*TxIndex) {
	return func(txi *TxIndex) {
		txi.compositeKeysToIndex = compositeKeys
	}
}

// IndexAllEvents is an option for indexing all events.
func IndexAllEvents() func(*TxIndex) {
	return func(txi *TxIndex) {
		txi.indexAllEvents = true
	}
}

func getTxResultFromBytes(rawBytes []byte) (*types.TxResult, error) {
	if rawBytes == nil {
		return nil, nil
	}
	txResult := new(types.TxResult)
	err := txResult.UnmarshalFromAmino(cdc, rawBytes)
	if err != nil {
		txResult = new(types.TxResult)
		err = cdc.UnmarshalBinaryBare(rawBytes, &txResult)
		if err != nil {
			return nil, fmt.Errorf("error reading TxResult: %v", err)
		}
	}
	return txResult, nil
}

// Get gets transaction from the TxIndex storage and returns it or nil if the
// transaction is not found.
func (txi *TxIndex) Get(hash []byte) (*types.TxResult, error) {
	if len(hash) == 0 {
		return nil, txindex.ErrorEmptyHash
	}

	v, err := txi.store.GetUnsafeValue(hash, func(rawBytes []byte) (interface{}, error) {
		return getTxResultFromBytes(rawBytes)
	})
	if err != nil {
		txResult, ok := v.(*types.TxResult)
		if !ok {
			panic(err)
		} else {
			return txResult, err
		}
	}

	txResult := v.(*types.TxResult)
	return txResult, err
}

// AddBatch indexes a batch of transactions using the given list of events. Each
// key that indexed from the tx's events is a composite of the event type and
// the respective attribute's key delimited by a "." (eg. "account.number").
// Any event with an empty type is not indexed.
func (txi *TxIndex) AddBatch(b *txindex.Batch) error {
	storeBatch := txi.store.NewBatch()
	defer storeBatch.Close()

	for _, result := range b.Ops {
		hash := result.Tx.Hash(result.Height)

		// index tx by events
		txi.indexEvents(result, hash, storeBatch)

		// index tx by height
		if txi.indexAllEvents || tmstring.StringInSlice(types.TxHeightKey, txi.compositeKeysToIndex) {
			storeBatch.Set(keyForHeight(result), hash)
		}

		// index tx by hash
		rawBytes, err := cdc.MarshalBinaryBare(result)
		if err != nil {
			return err
		}
		storeBatch.Set(hash, rawBytes)
	}

	storeBatch.WriteSync()
	return nil
}

// Index indexes a single transaction using the given list of events. Each key
// that indexed from the tx's events is a composite of the event type and the
// respective attribute's key delimited by a "." (eg. "account.number").
// Any event with an empty type is not indexed.
func (txi *TxIndex) Index(result *types.TxResult) error {
	b := txi.store.NewBatch()
	defer b.Close()

	hash := result.Tx.Hash(result.Height)

	// index tx by events
	txi.indexEvents(result, hash, b)

	// index tx by height
	if txi.indexAllEvents || tmstring.StringInSlice(types.TxHeightKey, txi.compositeKeysToIndex) {
		b.Set(keyForHeight(result), hash)
	}

	// index tx by hash
	rawBytes, err := cdc.MarshalBinaryBare(result)
	if err != nil {
		return err
	}

	b.Set(hash, rawBytes)
	b.WriteSync()

	return nil
}

func (txi *TxIndex) indexEvents(result *types.TxResult, hash []byte, store dbm.SetDeleter) {
	for _, event := range result.Result.Events {
		// only index events with a non-empty type
		if len(event.Type) == 0 {
			continue
		}

		for _, attr := range event.Attributes {
			if len(attr.Key) == 0 {
				continue
			}

			compositeTag := fmt.Sprintf("%s.%s", event.Type, string(attr.Key))
			if txi.indexAllEvents || tmstring.StringInSlice(compositeTag, txi.compositeKeysToIndex) {
				store.Set(keyForEvent(compositeTag, attr.Value, result), hash)
			}
		}
	}
}

// Search performs a search using the given query.
//
// It breaks the query into conditions (like "tx.height > 5"). For each
// condition, it queries the DB index. One special use cases here: (1) if
// "tx.hash" is found, it returns tx result for it (2) for range queries it is
// better for the client to provide both lower and upper bounds, so we are not
// performing a full scan. Results from querying indexes are then intersected
// and returned to the caller, in no particular order.
//
// Search will exit early and return any result fetched so far,
// when a message is received on the context chan.
func (txi *TxIndex) Search(ctx context.Context, q *query.Query) ([]*types.TxResult, error) {
	// Potentially exit early.
	select {
	case <-ctx.Done():
		results := make([]*types.TxResult, 0)
		return results, nil
	default:
	}

	var hashesInitialized bool
	filteredHashes := make(map[string][]byte)

	// get a list of conditions (like "tx.height > 5")
	conditions := q.Syntax()

	// if there is a hash condition, return the result immediately
	hash, ok, err := lookForHash(conditions)
	if err != nil {
		return nil, errors.Wrap(err, "error during searching for a hash in the query")
	} else if ok {
		res, err := txi.Get(hash)
		switch {
		case err != nil:
			return []*types.TxResult{}, errors.Wrap(err, "error while retrieving the result")
		case res == nil:
			return []*types.TxResult{}, nil
		default:
			return []*types.TxResult{res}, nil
		}
	}

	// conditions to skip because they're handled before "everything else"
	skipIndexes := make([]int, 0)

	// extract ranges
	// if both upper and lower bounds exist, it's better to get them in order not
	// no iterate over kvs that are not within range.
	ranges, rangeIndexes := lookForRanges(conditions)
	if len(ranges) > 0 {
		skipIndexes = append(skipIndexes, rangeIndexes...)

		for _, r := range ranges {
			if !hashesInitialized {
				filteredHashes = txi.matchRange(ctx, r, prefixFromCompositeKey(r.key), filteredHashes, true)
				hashesInitialized = true

				// Ignore any remaining conditions if the first condition resulted
				// in no matches (assuming implicit AND operand).
				if len(filteredHashes) == 0 {
					break
				}
			} else {
				filteredHashes = txi.matchRange(ctx, r, prefixFromCompositeKey(r.key), filteredHashes, false)
			}
		}
	}

	// if there is a height condition ("tx.height=3"), extract it
	height := lookForHeight(conditions)

	// for all other conditions
	for i, c := range conditions {
		if intInSlice(i, skipIndexes) {
			continue
		}

		if !hashesInitialized {
			filteredHashes = txi.match(ctx, c, prefixForCondition(c, height), filteredHashes, true)
			hashesInitialized = true

			// Ignore any remaining conditions if the first condition resulted
			// in no matches (assuming implicit AND operand).
			if len(filteredHashes) == 0 {
				break
			}
		} else {
			filteredHashes = txi.match(ctx, c, prefixForCondition(c, height), filteredHashes, false)
		}
	}

	results := make([]*types.TxResult, 0, len(filteredHashes))
hashes:
	for _, h := range filteredHashes {
		res, err := txi.Get(h)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get Tx{%X}", h)
		}
		results = append(results, res)

		// Potentially exit early.
		select {
		case <-ctx.Done():
			break hashes
		default:
		}
	}

	return results, nil
}

func lookForHash(conditions []syntax.Condition) (hash []byte, ok bool, err error) {
	for _, c := range conditions {
		if c.Tag == types.TxHashKey {
			decoded, err := hex.DecodeString(c.Arg.Value())
			return decoded, true, err
		}
	}
	return
}

// lookForHeight returns a height if there is an "height=X" condition.
func lookForHeight(conditions []syntax.Condition) (height int64) {
	for _, c := range conditions {
		if c.Tag == types.TxHeightKey && c.Op == syntax.TEq {
			return int64(c.Arg.Number())
		}
	}
	return 0
}

// special map to hold range conditions
// Example: account.number => queryRange{lowerBound: 1, upperBound: 5}
type queryRanges map[string]queryRange

type queryRange struct {
	lowerBound        interface{} // int || time.Time
	upperBound        interface{} // int || time.Time
	key               string
	includeLowerBound bool
	includeUpperBound bool
}

func (r queryRange) lowerBoundValue() interface{} {
	if r.lowerBound == nil {
		return nil
	}

	if r.includeLowerBound {
		return r.lowerBound
	}

	switch t := r.lowerBound.(type) {
	case int64:
		return t + 1
	case time.Time:
		return t.Unix() + 1
	default:
		panic("not implemented")
	}
}

func (r queryRange) AnyBound() interface{} {
	if r.lowerBound != nil {
		return r.lowerBound
	}

	return r.upperBound
}

func (r queryRange) upperBoundValue() interface{} {
	if r.upperBound == nil {
		return nil
	}

	if r.includeUpperBound {
		return r.upperBound
	}

	switch t := r.upperBound.(type) {
	case int64:
		return t - 1
	case time.Time:
		return t.Unix() - 1
	default:
		panic("not implemented")
	}
}

func lookForRanges(conditions []syntax.Condition) (ranges queryRanges, indexes []int) {
	ranges = make(queryRanges)
	for i, c := range conditions {
		if isRangeOperation(c.Op) {
			r, ok := ranges[c.Tag]
			if !ok {
				r = queryRange{key: c.Tag}
			}
			switch c.Op {
			case syntax.TGt:
				r.lowerBound = conditionArg(c)
			case syntax.TGeq:
				r.includeLowerBound = true
				r.lowerBound = conditionArg(c)
			case syntax.TLt:
				r.upperBound = conditionArg(c)
			case syntax.TLeq:
				r.includeUpperBound = true
				r.upperBound = conditionArg(c)
			}
			ranges[c.Tag] = r
			indexes = append(indexes, i)
		}
	}
	return ranges, indexes
}

func isRangeOperation(op syntax.Token) bool {
	switch op {
	case syntax.TGt, syntax.TGeq, syntax.TLt, syntax.TLeq:
		return true
	default:
		return false
	}
}

func conditionArg(c syntax.Condition) interface{} {
	if c.Arg == nil {
		return nil
	}
	switch c.Arg.Type {
	case syntax.TNumber:
		return int64(c.Arg.Number())
	case syntax.TTime, syntax.TDate:
		return c.Arg.Time()
	default:
		return c.Arg.Value() // string
	}
}

// match returns all matching txs by hash that meet a given condition and start
// key. An already filtered result (filteredHashes) is provided such that any
// non-intersecting matches are removed.
//
// NOTE: filteredHashes may be empty if no previous condition has matched.
func (txi *TxIndex) match(
	ctx context.Context,
	c syntax.Condition,
	startKeyBz []byte,
	filteredHashes map[string][]byte,
	firstRun bool,
) map[string][]byte {
	// A previous match was attempted but resulted in no matches, so we return
	// no matches (assuming AND operand).
	if !firstRun && len(filteredHashes) == 0 {
		return filteredHashes
	}

	tmpHashes := make(map[string][]byte)

	switch {
	case c.Op == syntax.TEq:
		it, err := dbm.IteratePrefix(txi.store, startKeyBz)
		if err != nil {
			panic(err)
		}
		defer it.Close()

	iterEqual:
		for ; it.Valid(); it.Next() {
			tmpHashes[string(it.Value())] = it.Value()

			// Potentially exit early.
			select {
			case <-ctx.Done():
				break iterEqual
			default:
			}
		}
		if err := it.Error(); err != nil {
			panic(err)
		}

	case c.Op == syntax.TExists:
		// XXX: can't use startKeyBz here because c.Operand is nil
		// (e.g. "account.owner/<nil>/" won't match w/ a single row)
		it, err := dbm.IteratePrefix(txi.store, prefixFromCompositeKey(c.Tag))
		if err != nil {
			panic(err)
		}
		defer it.Close()

	iterExists:
		for ; it.Valid(); it.Next() {
			tmpHashes[string(it.Value())] = it.Value()

			// Potentially exit early.
			select {
			case <-ctx.Done():
				break iterExists
			default:
			}
		}
		if err := it.Error(); err != nil {
			panic(err)
		}

	case c.Op == syntax.TContains:
		// XXX: startKey does not apply here.
		// For example, if startKey = "account.owner/an/" and search query = "account.owner CONTAINS an"
		// we can't iterate with prefix "account.owner/an/" because we might miss keys like "account.owner/Ulan/"
		it, err := dbm.IteratePrefix(txi.store, prefixFromCompositeKey(c.Tag))
		if err != nil {
			panic(err)
		}
		defer it.Close()

	iterContains:
		for ; it.Valid(); it.Next() {
			value, err := parseValueFromKey(it.Key())
			if err != nil {
				continue
			}
			if strings.Contains(value, c.Arg.Value()) {
				tmpHashes[string(it.Value())] = it.Value()
			}

			// Potentially exit early.
			select {
			case <-ctx.Done():
				break iterContains
			default:
			}
		}
		if err := it.Error(); err != nil {
			panic(err)
		}
	default:
		panic("other operators should be handled already")
	}

	if len(tmpHashes) == 0 || firstRun {
		// Either:
		//
		// 1. Regardless if a previous match was attempted, which may have had
		// results, but no match was found for the current condition, then we
		// return no matches (assuming AND operand).
		//
		// 2. A previous match was not attempted, so we return all results.
		return tmpHashes
	}

	// Remove/reduce matches in filteredHashes that were not found in this
	// match (tmpHashes).
	for k := range filteredHashes {
		if tmpHashes[k] == nil {
			delete(filteredHashes, k)

			// Potentially exit early.
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}

	return filteredHashes
}

// matchRange returns all matching txs by hash that meet a given queryRange and
// start key. An already filtered result (filteredHashes) is provided such that
// any non-intersecting matches are removed.
//
// NOTE: filteredHashes may be empty if no previous condition has matched.
func (txi *TxIndex) matchRange(
	ctx context.Context,
	r queryRange,
	startKey []byte,
	filteredHashes map[string][]byte,
	firstRun bool,
) map[string][]byte {
	// A previous match was attempted but resulted in no matches, so we return
	// no matches (assuming AND operand).
	if !firstRun && len(filteredHashes) == 0 {
		return filteredHashes
	}

	tmpHashes := make(map[string][]byte)
	lowerBound := r.lowerBoundValue()
	upperBound := r.upperBoundValue()

	it, err := dbm.IteratePrefix(txi.store, startKey)
	if err != nil {
		panic(err)
	}
	defer it.Close()

LOOP:
	for ; it.Valid(); it.Next() {
		if !isTagKey(it.Key()) {
			continue
		}

		if _, ok := r.AnyBound().(int64); ok {
			v, err := strconv.ParseInt(extractValueFromKey(it.Key()), 10, 64)
			if err != nil {
				continue LOOP
			}

			include := true
			if lowerBound != nil && v < lowerBound.(int64) {
				include = false
			}

			if upperBound != nil && v > upperBound.(int64) {
				include = false
			}

			if include {
				tmpHashes[string(it.Value())] = it.Value()
			}

			// XXX: passing time in a ABCI Events is not yet implemented
			// case time.Time:
			// 	v := strconv.ParseInt(extractValueFromKey(it.Key()), 10, 64)
			// 	if v == r.upperBound {
			// 		break
			// 	}
		}

		// Potentially exit early.
		select {
		case <-ctx.Done():
			break
		default:
		}
	}

	if len(tmpHashes) == 0 || firstRun {
		// Either:
		//
		// 1. Regardless if a previous match was attempted, which may have had
		// results, but no match was found for the current condition, then we
		// return no matches (assuming AND operand).
		//
		// 2. A previous match was not attempted, so we return all results.
		return tmpHashes
	}

	// Remove/reduce matches in filteredHashes that were not found in this
	// match (tmpHashes).
	for k := range filteredHashes {
		if tmpHashes[k] == nil {
			delete(filteredHashes, k)

			// Potentially exit early.
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}

	return filteredHashes
}

///////////////////////////////////////////////////////////////////////////////
// Keys

func isTagKey(key []byte) bool {
	return strings.Count(string(key), tagKeySeparator) == 3
}

func extractValueFromKey(key []byte) string {
	parts := strings.SplitN(string(key), tagKeySeparator, 3)
	return parts[1]
}

func keyForEvent(key string, value []byte, result *types.TxResult) []byte {
	return []byte(fmt.Sprintf("%s/%s/%d/%d",
		key,
		value,
		result.Height,
		result.Index,
	))
}

func keyForHeight(result *types.TxResult) []byte {
	return []byte(fmt.Sprintf("%s/%d/%d/%d",
		types.TxHeightKey,
		result.Height,
		result.Height,
		result.Index,
	))
}

func startKeyForCondition(c query.Condition, height int64) []byte {
	if height > 0 {
		return startKey(c.CompositeKey, c.Operand, height)
	}
	return startKey(c.CompositeKey, c.Operand)
}

func startKey(fields ...interface{}) []byte {
	var b bytes.Buffer
	for _, f := range fields {
		b.Write([]byte(fmt.Sprintf("%v", f) + tagKeySeparator))
	}
	return b.Bytes()
}

// Prefixes: these represent an initial part of the key and are used by iterators to iterate over a small
// section of the kv store during searches.

func prefixFromCompositeKey(compositeKey string) []byte {
	key, err := orderedcode.Append(nil, compositeKey)
	if err != nil {
		panic(err)
	}
	return key
}

func prefixFromCompositeKeyAndValue(compositeKey, value string) []byte {
	key, err := orderedcode.Append(nil, compositeKey, value)
	if err != nil {
		panic(err)
	}
	return key
}

// a small utility function for getting a keys prefix based on a condition and a height
func prefixForCondition(c syntax.Condition, height int64) []byte {
	key := prefixFromCompositeKeyAndValue(c.Tag, c.Arg.Value())
	if height > 0 {
		var err error
		key, err = orderedcode.Append(key, height)
		if err != nil {
			panic(err)
		}
	}
	return key
}

// parseValueFromKey parses an event key and extracts out the value, returning an error if one arises.
// This will also involve ensuring that the key has the correct format.
// CONTRACT: function doesn't check that the prefix is correct. This should have already been done by the iterator
func parseValueFromKey(key []byte) (string, error) {
	var (
		compositeKey, value string
		height, index       int64
	)
	remaining, err := orderedcode.Parse(string(key), &compositeKey, &value, &height, &index)
	if err != nil {
		return "", err
	}
	if len(remaining) != 0 {
		return "", fmt.Errorf("unexpected remainder in key: %s", remaining)
	}
	return value, nil
}
