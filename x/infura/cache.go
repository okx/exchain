package infura

import evm "github.com/okx/okbchain/x/evm/watcher"

const defaultCacheCap = 2000

type Cache struct {
	transactionReceipts []evm.TransactionReceipt
	block               *evm.Block
	transactions        []evm.Transaction
	contractCodes       map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		transactionReceipts: make([]evm.TransactionReceipt, 0, defaultCacheCap),
		block:               nil,
		transactions:        make([]evm.Transaction, 0, defaultCacheCap),
		contractCodes:       make(map[string][]byte, defaultCacheCap),
	}
}

func (c *Cache) Reset() {
	c.transactionReceipts = make([]evm.TransactionReceipt, 0, defaultCacheCap)
	c.block = nil
	c.transactions = make([]evm.Transaction, 0, defaultCacheCap)
	c.contractCodes = make(map[string][]byte, defaultCacheCap)
}

func (c *Cache) AddTransactionReceipt(tr evm.TransactionReceipt) {
	c.transactionReceipts = append(c.transactionReceipts, tr)
}

func (c *Cache) GetTransactionReceipts() []evm.TransactionReceipt {
	return c.transactionReceipts
}

func (c *Cache) AddBlock(b evm.Block) {
	c.block = &b
}

func (c *Cache) GetBlock() evm.Block {
	return *c.block
}

func (c *Cache) AddTransaction(t evm.Transaction) {
	c.transactions = append(c.transactions, t)
}

func (c *Cache) GetTransactions() []evm.Transaction {
	return c.transactions
}

func (c *Cache) AddContractCode(address string, code []byte) {
	c.contractCodes[address] = code
}

func (c *Cache) GetContractCodes() map[string][]byte {
	return c.contractCodes
}
