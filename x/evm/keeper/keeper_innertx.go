package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	cm_innertx "github.com/okex/exchain/libs/cosmos-sdk/types/innertx"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types/innertx"

	"github.com/ethereum/go-ethereum/common"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/spf13/viper"
)

func initInnerDB() error {
	innerTxPath := viper.GetString(cli.HomeFlag)
	dbBackend := viper.GetString("db_backend")

	return innertx.InitDB(innerTxPath, dbBackend)
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
		k.innerTxLock.RLock()
		existedTxs, ok := k.innerBlockData.TxMap[hash]
		k.innerTxLock.RUnlock()
		if !ok {
			//Initialization
			k.innerTxLock.Lock()
			k.innerBlockData.TxHashes = append(k.innerBlockData.TxHashes, hash)
			k.innerBlockData.TxMap[hash] = innerTxs
			k.innerTxLock.Unlock()
		} else {
			hasDeductTx := false
			//Determine if the pending innerTx contains deductTx
			//todo: binary search
			for _, tx := range innerTxs {
				//"To" address of deduct fee
				if tx.To == "0xf1829676DB577682E944fc3493d451B67Ff3E29F" {
					hasDeductTx = true
					break
				}
			}
			k.innerTxLock.Lock()
			if hasDeductTx {
				delete(k.innerBlockData.TxMap, hash)
				k.innerBlockData.TxMap[hash] = innerTxs
			} else {
				existedTxs = append(existedTxs, innerTxs...)
				k.innerBlockData.TxMap[hash] = existedTxs
			}
			k.innerTxLock.Unlock()
		}
	} else {
		panic("Invalid parameter types for evm")
	}
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
