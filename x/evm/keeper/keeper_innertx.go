package keeper

import (
	"sync"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	cm_innertx "github.com/okex/exchain/libs/cosmos-sdk/types/innertx"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types/innertx"

	"github.com/ethereum/go-ethereum/common"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/spf13/viper"

	"github.com/okex/exchain/libs/tendermint/libs/cli"
)

var innerTxLock sync.Mutex

func initInnerDB() error {
	innerTxPath := viper.GetString(cli.HomeFlag)
	dbBackend := viper.GetString("db_backend")

	return innertx.InitDB(innerTxPath, dbBackend)
}

func closeInnerDB() []error {
	return innertx.CloseDB()
}

type BlockInnerData = ethvm.BlockInnerData

func defaultBlockInnerData() BlockInnerData {
	return BlockInnerData{
		BlockHash:    "",
		TxHashes:     make([]string, 0),
		TxMap:        make(map[string][]*ethvm.InnerTx),
		ContractList: make([]*ethvm.ERC20Contract, 0),
	}
}

// InitInnerBlock init inner block data
func (k *Keeper) InitInnerBlock(hash string) {
	k.innerBlockData = ethvm.BlockInnerData{
		BlockHash:    hash,
		TxHashes:     make([]string, 0),
		TxMap:        make(map[string][]*ethvm.InnerTx),
		ContractList: make([]*ethvm.ERC20Contract, 0),
	}
}

func (k *Keeper) UpdateInnerBlockData() {
	//Block write db
	if len(k.innerBlockData.TxHashes) > 0 {
		if err := ethvm.WriteBlockDB(k.innerBlockData.BlockHash, k.innerBlockData.TxHashes); err != nil {
			panic(err)
		}
	}
	//InnerTx write db
	if len(k.innerBlockData.TxMap) > 0 {
		for txHash, inTx := range k.innerBlockData.TxMap {
			if err := ethvm.WriteTx(txHash, inTx); err != nil {
				panic(err)
			}
		}
	}

	//Contract write db
	if len(k.innerBlockData.ContractList) > 0 {
		for _, contract := range k.innerBlockData.ContractList {
			if err := ethvm.WriteToken(contract.ContractAddr, contract.ContractCode); err != nil {
				panic(err)
			}
		}
	}
}

// AddInnerTx add inner tx
func (k *Keeper) AddInnerTx(hash string, txs interface{}) {
	if innerTxs, ok := txs.([]*ethvm.InnerTx); ok {
		innerTxLock.Lock()
		existedTxs, ok := k.innerBlockData.TxMap[hash]
		if !ok {
			//Empty the TxHashes and TxMap
			//Here, TxHashes and TxMap need to be updated synchronously
			k.innerBlockData.TxHashes = append(k.innerBlockData.TxHashes, hash)
			k.innerBlockData.TxMap[hash] = innerTxs
		} else {
			//Determine if the pending innerTx is deductTx.
			//If so, it means that re-commit has occurred in the mode with concurrent tx on,
			//then we need to delete the existing data and overwrite the innerTx data.
			if isDeductTx(innerTxs) {
				k.innerBlockData.TxMap[hash] = innerTxs
			} else {
				existedTxs = append(existedTxs, innerTxs...)
				k.innerBlockData.TxMap[hash] = existedTxs
			}
		}
		innerTxLock.Unlock()
	} else {
		panic("Invalid parameter types for evm")
	}
}

//deductTx is used to deduct fee, and it must be the first innerTx.
//Its "To" address is fixed, here it is hard-coded to determine if it is deductTx.
func isDeductTx(txs []*ethvm.InnerTx) bool {
	return len(txs) == 1 && txs[0].To == "0xf1829676DB577682E944fc3493d451B67Ff3E29F"
}

func (k *Keeper) UpdateInnerTx(txBytes []byte, blockHeight int64, dept int64, from, to sdk.AccAddress, callType, name string, amt sdk.Coins, err error) {
	txHash := tmtypes.Tx(txBytes).Hash(blockHeight)
	ethHash := common.BytesToHash(txHash)
	ethHashHex := ethHash.Hex()
	if txBytes == nil || len(txBytes) == 0 {
		ethHashHex = k.innerBlockData.BlockHash
	}

	innerTXValue := cm_innertx.BIG0
	if len(amt) != 0 {
		innerTXValue = amt[0].Amount.BigInt()
	}
	ethFrom := common.BytesToAddress(from.Bytes())
	ethTo := common.BytesToAddress(to.Bytes())
	innerTx := cm_innertx.CreateInnerTx(dept, ethFrom.String(), ethTo.String(), callType, name, innerTXValue, err)
	innerTxs := make([]*ethvm.InnerTx, 0)
	innerTxs = append(innerTxs, innerTx)
	k.AddInnerTx(ethHashHex, innerTxs)
}

// AddContract add erc20 contract
func (k *Keeper) AddContract(contracts interface{}) {
	if cs, ok := contracts.([]*ethvm.ERC20Contract); ok {
		k.innerBlockData.ContractList = append(k.innerBlockData.ContractList, cs...)
	} else {
		panic("Invalid parameter types")
	}
}
