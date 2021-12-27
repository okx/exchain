package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func initInnerDB() error {
	return nil
}

type BlockInnerData = interface{}

func defaultBlockInnerData() BlockInnerData {
	return nil
}

// InitInnerBlock init inner block data
func (k *Keeper) InitInnerBlock(hash string) {}

func (k *Keeper) UpdateInnerBlockData(...interface{}) {}

// AddInnerTx add inner tx
func (k *Keeper) AddInnerTx(...interface{}) {}

// AddContract add erc20 contract
func (k *Keeper) AddContract(...interface{}) {}

func (k *Keeper) UpdateInnerTx(txBytes []byte, dept int64, from, to sdk.AccAddress, callType, name string, amt sdk.Coins, err error) {
}
