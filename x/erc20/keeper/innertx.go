package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/okex/exchain/libs/cosmos-sdk/types/innertx"
	"github.com/okex/exchain/x/erc20/types"
)

func (k Keeper) addEVMInnerTx(ethTxHash string, innertxs, contracts interface{}) {
	if innertxs != nil {
		k.evmKeeper.AddInnerTx(ethTxHash, innertxs)
	}
	if contracts != nil {
		k.evmKeeper.AddContract(contracts)
	}
}

func (k Keeper) addSendToIbcInnerTx(ethTxHash, from, sender, recipient, vouchers, ibcEvents string) {
	unlockTx := &vm.InnerTx{
		Dept:     *big.NewInt(0),
		From:     from,
		To:       sender,
		CallType: innertx.CosmosCallType,
		Name:     types.InnerTxUnlock,
		Input:    vouchers,
	}

	sendToIbcTx :=
		&vm.InnerTx{
			Dept:     *big.NewInt(0),
			From:     sender,
			To:       recipient,
			CallType: innertx.CosmosCallType,
			Name:     types.InnerTxSendToIbc,
			Input:    vouchers,
			Output:   ibcEvents,
		}

	k.evmKeeper.AddInnerTx(ethTxHash, []*vm.InnerTx{unlockTx, sendToIbcTx})
}

func (k Keeper) addSendNative20ToIbcInnerTx(ethTxHash, from, sender, recipient, native20s, ibcEvents string) {
	mintTx := &vm.InnerTx{
		Dept:     *big.NewInt(0),
		From:     from,
		To:       sender,
		CallType: innertx.CosmosCallType,
		Name:     types.InnerTxMint,
		Input:    native20s,
	}

	sendToIbcTx :=
		&vm.InnerTx{
			Dept:     *big.NewInt(0),
			From:     sender,
			To:       recipient,
			CallType: innertx.CosmosCallType,
			Name:     types.InnerTxSendToIbc,
			Input:    native20s,
			Output:   ibcEvents,
		}
	k.evmKeeper.AddInnerTx(ethTxHash, []*vm.InnerTx{mintTx, sendToIbcTx})
}
