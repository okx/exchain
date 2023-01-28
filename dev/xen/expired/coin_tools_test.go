package expired

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"testing"
)

func Test_getIndex(t *testing.T) {
	getIndex(41550, "0xd7a68dcd1f39b42f1b617fa14df83c2342a38dfd")
}

func TestDyncPack(t *testing.T) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "byte32", nil)
	addressTy, _ := abi.NewType("address", "address", nil)

	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: uint256Ty,
		},
	}

	bytes, _ := arguments.Pack(
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		[32]byte{'I', 'D', '1'},
		big.NewInt(42),
	)

	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	log.Println(hexutil.Encode(buf))

}
