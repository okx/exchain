package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/okexchain/x/evm/types"
)

// SetCodeDirectly commit code into db with no cache
func (k *Keeper) SetCodeDirectly(ctx sdk.Context, code []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCode)
	store.Set(ethcrypto.Keccak256Hash(code).Bytes(), code)
}

// SetStateDirectly commit one state into db with no cache
func (k *Keeper) SetStateDirectly(ctx sdk.Context, addr ethcmn.Address, key, value ethcmn.Hash) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AddressStoragePrefix(addr))
	store.Set(key.Bytes(), value.Bytes())
}

// SetTxLogsDirectly commit logs into db with no cache
func (k *Keeper) SetTxLogsDirectly(ctx sdk.Context, hash ethcmn.Hash, logs []*ethtypes.Log) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLogs)
	bz, err := types.MarshalLogs(logs)
	if err != nil {
		panic(err)
	}
	store.Set(hash.Bytes(), bz)
}
