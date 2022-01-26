package db

import (
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
)

func init() {
	dbCreator := func(dbName, dir string) (DB, error) {
		return NewBadgerDB(dbName, dir)
	}
	registerDBCreator(BadgerDBBackend, dbCreator, true)
}

type BadgerDB struct {
	db *badger.DB
}

var _ DB = (*BadgerDB)(nil)

// NewBadgerDB creates a Badger key-value store backed to the
// directory dir supplied. If dir does not exist, it will be created.
func NewBadgerDB(dbName, dir string) (*BadgerDB, error) {
	// Since Badger doesn't support database names, we join both to obtain
	// the final directory to use for the database.
	path := filepath.Join(dir, dbName+".db")

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	opts := badger.DefaultOptions(path)
	opts.SyncWrites = false // note that we have Sync methods
	opts.Logger = nil       // badger is too chatty by default
	return NewBadgerDBWithOptions(opts)
}

// NewBadgerDBWithOptions creates a BadgerDB key value store
// gives the flexibility of initializing a database with the
// respective options.
func NewBadgerDBWithOptions(opts badger.Options) (*BadgerDB, error) {
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerDB{db: db}, nil
}

func (b *BadgerDB) Get(key []byte) ([]byte, error) {
	key = nonNilBytes(key)
	var val []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		if err == nil && val == nil {
			val = []byte{}
		}
		return err
	})
	return val, err
}

func (b *BadgerDB) Has(key []byte) (bool, error) {
	key = nonNilBytes(key)
	var found bool
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		found = (err != badger.ErrKeyNotFound)
		return nil
	})
	return found, err
}

func (b *BadgerDB) Set(key, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (b *BadgerDB) SetSync(key, value []byte) error {
	err := b.Set(key, value)
	if err != nil {
		return err
	}
	return b.db.Sync()
}

func (b *BadgerDB) Delete(key []byte) error {
	key = nonNilBytes(key)
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BadgerDB) DeleteSync(key []byte) error {
	err := b.Delete(key)
	if err != nil {
		return err
	}
	return b.db.Sync()
}

func (db *BadgerDB) DB() *badger.DB {
	return db.db
}

func (b *BadgerDB) Close() error {
	return b.db.Close()
}

func (b *BadgerDB) Print() error {
	return nil
}

func (b *BadgerDB) Stats() map[string]string {
	return nil
}

func (b *BadgerDB) NewBatch() Batch {
	return newBadgerDBBatch(b.db)
}

func (b *BadgerDB) Iterator(start, end []byte) (Iterator, error) {
	opts := badger.DefaultIteratorOptions
	txn := b.db.NewTransaction(false)
	return newBadgerDBIterator(txn, start, end, opts), nil
}

func (b *BadgerDB) ReverseIterator(start, end []byte) (Iterator, error) {
	opts := badger.DefaultIteratorOptions
	opts.Reverse = true
	txn := b.db.NewTransaction(false)
	return newBadgerDBIterator(txn, end, start, opts), nil
}
