package types

import (
	"bytes"
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/pkg/errors"
	"sort"
)

//Treasure is the struct which has address and proportion of mint reward.
type Treasure struct {
	//Treasure Address
	Address sdk.AccAddress `json:"address" yaml:"address"`
	// proportion of minted for treasure
	Proportion sdk.Dec `json:"proportion" yaml:"proportion"`
}

func NewTreasure(address sdk.AccAddress, proportion sdk.Dec) *Treasure {
	return &Treasure{
		Address:    address,
		Proportion: proportion,
	}
}

func (t Treasure) ValidateBasic() error {
	// proportion must (0,1]. if proportion <= 0 or  > 1
	if t.Proportion.LTE(sdk.ZeroDec()) || t.Proportion.GT(sdk.OneDec()) {
		return errors.New(fmt.Sprintf("treasure proportion should non-negative and less than one: %s", t.Proportion))
	}
	return nil
}

func ValidateTreasures(treasures []Treasure) error {
	sumProportion := sdk.ZeroDec()
	for i, _ := range treasures {
		if err := treasures[i].ValidateBasic(); err != nil {
			return err
		}
		sumProportion = sumProportion.Add(treasures[i].Proportion)
	}
	if sumProportion.IsNegative() || sumProportion.GT(sdk.OneDec()) {
		return errors.New(fmt.Sprintf("the sum of treasure proportion should non-negative and less than one: %s", sumProportion))
	}
	return nil
}

func SortTreasures(treasures []Treasure) {
	sort.Slice(treasures, func(i, j int) bool {
		return bytes.Compare(treasures[i].Address.Bytes(), treasures[j].Address.Bytes()) > 0
	})
}

func GetMapFromTreasures(treasures []Treasure) map[string]Treasure {
	temp := make(map[string]Treasure, 0)
	for i, _ := range treasures {
		temp[treasures[i].Address.String()] = treasures[i]
	}
	return temp
}

func GetTreasuresFromMap(src map[string]Treasure) []Treasure {
	result := make([]Treasure, 0)
	for k, _ := range src {
		result = append(result, src[k])
	}
	SortTreasures(result)
	return result
}

func IsTreasureDuplicated(treasures []Treasure) bool {
	temp := make(map[string]Treasure, 0)
	for i, _ := range treasures {
		key := treasures[i].Address.String()
		if _, ok := temp[key]; ok {
			return true
		}
		temp[key] = treasures[i]
	}
	return false
}

func InsertAndUpdateTreasures(src, dst []Treasure) []Treasure {
	temp := GetMapFromTreasures(src)
	for i, _ := range dst {
		key := dst[i].Address.String()
		temp[key] = dst[i]
	}
	return GetTreasuresFromMap(temp)
}

func DeleteTreasures(src, dst []Treasure) ([]Treasure, error) {
	temp := GetMapFromTreasures(src)
	for i, _ := range dst {
		key := dst[i].Address.String()
		if _, ok := temp[key]; !ok {
			return nil, errors.New(fmt.Sprintf("can not delete %s,because it's not exist from treasures", key))
		}
		delete(temp, key)
	}
	return GetTreasuresFromMap(temp), nil
}
