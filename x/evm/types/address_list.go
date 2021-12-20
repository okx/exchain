package types

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	ContractMethodBlockedCacheSize = 10000
)

var (
	// Map for quick access to contract method blocked.
	// txsMap: address.String() -> BlockedContract{}
	contractMethodBlockedCache = NewContractMethodBlockedCache() //Contract Method Blocked Cache
)

// AddressList is the type alias for []sdk.AccAddress
type AddressList []sdk.AccAddress

// String returns a human readable string representation of AddressList
func (al AddressList) String() string {
	var b strings.Builder
	b.WriteString("Address List:\n")
	for i := 0; i < len(al); i++ {
		b.WriteString(al[i].String())
		b.WriteByte('\n')
	}
	return strings.TrimSpace(b.String())
}

//BlockedContractList is the list of contract which method or all-method is blocked
type BlockedContractList []BlockedContract

// String returns string which is the list of blocked contract
func (bl BlockedContractList) String() string {
	var b strings.Builder
	b.WriteString("BlockedContractList List:\n")
	for i := 0; i < len(bl); i++ {
		b.WriteString(bl[i].String())
		b.WriteByte('\n')
	}

	return strings.TrimSpace(b.String())
}

// ValidateBasic validates the list of contract which method or all-method is blocked
func (bl BlockedContractList) ValidateBasic() sdk.Error {

	//check repeated contract address
	lenAddrs := len(bl)
	filter := make(map[string]struct{}, lenAddrs)
	for i := 0; i < lenAddrs; i++ {
		key := bl[i].Address.String()
		if _, ok := filter[key]; ok {
			return ErrDuplicatedAddr
		}
		if err := bl[i].ValidateBasic(); err != nil {
			return err
		}
		filter[key] = struct{}{}
	}
	return nil
}

//BlockedContract i the contract which method or all-method is blocked
type BlockedContract struct {
	//Contract Address
	Address sdk.AccAddress `json:"address" yaml:"address"`
	//the list of method which is blocked. If it's length equal to 0,it means all method is blocked.
	BlockMethods ContractMethods `json:"block_methods" yaml:"block_methods"`
}

// NewBlockContract return point of BlockedContract
func NewBlockContract(addr sdk.AccAddress, methods ContractMethods) *BlockedContract {
	bm := make([]ContractMethod,len(methods))
	copy(bm,methods)
	return &BlockedContract{Address: addr, BlockMethods: bm}
}

// ValidateBasic validates BlockedContract
func (bc BlockedContract) ValidateBasic() sdk.Error {
	if len(bc.Address) == 0 {
		return ErrEmptyAddressBlockedContract
	}
	return bc.BlockMethods.ValidateBasic()
}

// IsAllMethodBlocked return true if all method of contract is blocked.
func (bc BlockedContract) IsAllMethodBlocked() bool {
	return len(bc.BlockMethods) == 0
}

// IsMethodBlocked return true if the method of contract is blocked.
func (bc BlockedContract) IsMethodBlocked(method string) bool {
	return bc.BlockMethods.IsContain(method)
}

// String returns BlockedContract string
func (bc BlockedContract) String() string {
	var b strings.Builder
	b.WriteString("Address: ")
	b.WriteString(bc.Address.String())
	b.WriteByte('\n')
	b.WriteString(bc.BlockMethods.String())

	return strings.TrimSpace(b.String())
}

//ContractMethods is the list of blocked contract method
type ContractMethods []ContractMethod

func SortContractMethods(cms []ContractMethod)  {
	sort.Slice(cms, func(i, j int) bool {
		if cms[i].Sign == cms[j].Sign {
			return cms[i].Extra < cms[j].Extra
		}
		return cms[i].Sign < cms[j].Sign
	})
}
// String returns ContractMethods string
func (cms ContractMethods) String() string {
	var b strings.Builder
	b.WriteString("Method List:\n")
	for k, _ := range cms {
		b.WriteString(cms[k].String())
		b.WriteByte('\n')
	}

	return strings.TrimSpace(b.String())
}

// ValidateBasic validates the list of blocked contract method
func (cms ContractMethods) ValidateBasic() sdk.Error {
	methodMap := make(map[string]ContractMethod)
	for i, _ := range cms {
		if _, ok := methodMap[cms[i].Sign]; ok {
			return ErrDuplicatedMethod
		}
		if len(cms[i].Sign) == 0 {
			return ErrEmptyMethod
		}
		methodMap[cms[i].Sign] = cms[i]
	}
	return nil
}

// IsContain return true if the method of contract contains ContractMethods.
func (cms ContractMethods) IsContain(method string) bool {
	for i, _ := range cms {
		if strings.Compare(method, cms[i].Sign) == 0 {
			return true
		}
	}
	return false
}

// GetContractMethodsMap return map which key is method,value is ContractMethod.
func (cms ContractMethods) GetContractMethodsMap() map[string]ContractMethod {
	methodMap := make(map[string]ContractMethod)
	for i, _ := range cms {
		methodMap[cms[i].Sign] = cms[i]
	}
	return methodMap
}

// InsertContractMethods insert the list of ContractMethod into cms.
// if repeated,methods will cover cms
func (cms *ContractMethods) InsertContractMethods(methods ContractMethods) (ContractMethods,error) {
	methodMap := cms.GetContractMethodsMap()
	for i, _ := range methods {
		methodName := methods[i].Sign
		methodMap[methodName] = methods[i]
	}
	result := ContractMethods{}
	for k, _ := range methodMap {
		result = append(result, methodMap[k])
	}
	SortContractMethods(result)
	return result,nil
}

// DeleteContractMethodMap delete the list of ContractMethod from cms.
// if method is not exist,it can not be panic or error
func (cms *ContractMethods) DeleteContractMethodMap(methods ContractMethods) (ContractMethods,error) {
	methodMap := cms.GetContractMethodsMap()
	for i, _ := range methods {
		if _,ok := methodMap[methods[i].Sign]; !ok {
			return nil,errors.New(fmt.Sprintf("method(%s) is not exist",methods[i].Sign))
		}
		delete(methodMap, methods[i].Sign)
	}
	result := ContractMethods{}
	for k, _ := range methodMap {
		result = append(result, methodMap[k])
	}
	SortContractMethods(result)
	return result,nil
}

//ContractMethod is the blocked contract method
// Name is method  name
// Extra is a extend data is useless
type ContractMethod struct {
	Sign  string `json:"sign" yaml:"sign"`
	Extra string `json:"extra" yaml:"extra"`
}

func (cm ContractMethod) String() string {
	var b strings.Builder
	b.WriteString("Sign: ")
	b.WriteString(cm.Sign)
	b.WriteString("Extra: ")
	b.WriteString(cm.Extra)
	b.WriteString("\n")
	return strings.TrimSpace(b.String())
}

type ContractMethodBlockedCache struct {
	cache *lru.ARCCache
}

func NewContractMethodBlockedCache() *ContractMethodBlockedCache {
	cache, _ := lru.NewARC(ContractMethodBlockedCacheSize)
	return &ContractMethodBlockedCache{cache: cache}
}

func (cmbc *ContractMethodBlockedCache) GetContractMethod(keyData []byte) (ContractMethods, bool) {
	key := sha256.Sum256(keyData)
	value, success := cmbc.cache.Get(key)

	if success {
		cm, ok := value.(ContractMethods)
		return cm, ok
	}
	return ContractMethods{}, success
}

func (cmbc *ContractMethodBlockedCache) SetContractMethod(keyData []byte, bc ContractMethods) {
	key := sha256.Sum256(keyData)
	cmbc.cache.Add(key, bc)
}

func BlockedContractListIsEqual(t *testing.T,src, dst BlockedContractList) bool {
	expectedMap := make(map[string]ContractMethods, 0)
	actuallyMap := make(map[string]ContractMethods, 0)
	for i := range src {
		expectedMap[src[i].Address.String()] = src[i].BlockMethods
		actuallyMap[dst[i].Address.String()] = dst[i].BlockMethods
	}
	if len(expectedMap) != len(actuallyMap) {
		return false
	}

	for k, expected := range expectedMap {
		v, ok := actuallyMap[k]
		if !ok {
			return false
		}
		if !ContractMethodsIsEqual(expected, v) {
			return false
		}
		if expected != nil && v != nil {
			require.Equal(t, expected,v)
		}
	}
	return true
}

func ContractMethodsIsEqual(src, dst ContractMethods) bool {
	if len(src) != len(dst) {
		return false
	}
	srcMap := src.GetContractMethodsMap()
	for i, _ := range dst {
		if _, ok := srcMap[dst[i].Sign]; !ok {
			return false
		} else {
			delete(srcMap, dst[i].Sign)
		}
	}
	return true
}
