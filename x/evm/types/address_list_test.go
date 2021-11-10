package types

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	addr           = "ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc"
	addr1          = "0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0"
	expectedOutput = `Address List:
ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc
ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc`
)

func TestAddressList_String(t *testing.T) {
	accAddr, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)

	addrList := AddressList{accAddr, accAddr}
	require.Equal(t, expectedOutput, addrList.String())
}

func TestBlockMethod(t *testing.T) {
	bcl := BlockedContractList{}
	accAddr1, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)
	accAddr2, err := sdk.AccAddressFromBech32(addr1)
	require.NoError(t, err)

	cmm := ContractMethods{}
	method1 := []byte("transfer")[:4]
	method2 := []byte("approve")[:4]
	cm1 := ContractMethod{Name: hexutil.Encode(method1), Extra: "test1"}
	cm2 := ContractMethod{Name: hexutil.Encode(method2), Extra: "test1"}
	cmm = append(cmm, cm1, cm2)
	bc1 := BlockedContract{Address: accAddr1, BlockMethods: cmm}
	bc2 := BlockedContract{Address: accAddr2, BlockMethods: cmm}
	bcl = append(bcl, bc1, bc2)

	//test decode and encode
	buff := ModuleCdc.MustMarshalJSON(bcl)
	t.Log(string(buff))
	nbcl := BlockedContractList{}
	ModuleCdc.MustUnmarshalJSON(buff, &nbcl)
}
