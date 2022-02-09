package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTreasure(t *testing.T) {
	treasure := NewTreasure(nil, sdk.NewDecWithPrec(1, 2))
	b, err := ModuleCdc.MarshalBinaryLengthPrefixed(treasure)
	require.NoError(t, err)
	treasure = &Treasure{}
	err = ModuleCdc.UnmarshalBinaryLengthPrefixed(b, treasure)
	require.NoError(t, err)
	b, err = ModuleCdc.MarshalBinaryBare(treasure)
}

func TestValidateBasic(t *testing.T) {
	treasure1 := NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(1, 2))
	treasure2 := NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(1, 2))
	treasure3 := NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(1, 2))

	//success treasures
	treasures := []Treasure{*treasure1, *treasure2, *treasure3}
	err := ValidateTreasures(treasures)
	require.NoError(t, err)

	// success treasure's proportion is equal one
	temp := Treasure{Address: sdk.AccAddress([]byte{0x00}), Proportion: sdk.NewDecWithPrec(1, 0)}
	treasures = []Treasure{temp}
	err = ValidateTreasures(treasures)
	require.NoError(t, err)

	// success the sum proportion of treasures is equal one
	temp = Treasure{Address: sdk.AccAddress([]byte{0x00}), Proportion: sdk.NewDecWithPrec(98, 2)}
	treasures = []Treasure{*treasure1, *treasure2, temp}
	err = ValidateTreasures(treasures)
	require.NoError(t, err)

	// error treasure's proportion is negative
	temp = Treasure{Address: sdk.AccAddress([]byte{0x00}), Proportion: sdk.NewDec(-1)}
	treasures = []Treasure{*treasure1, *treasure2, temp}
	err = ValidateTreasures(treasures)
	require.Error(t, err)
	require.Contains(t, err.Error(), "treasure proportion should non-negative")

	// error treasure's proportion is more than one
	temp = Treasure{Address: sdk.AccAddress([]byte{0x00}), Proportion: sdk.NewDecWithPrec(2, 0)}
	treasures = []Treasure{*treasure1, *treasure2, temp}
	err = ValidateTreasures(treasures)
	require.Error(t, err)
	require.Contains(t, err.Error(), "treasure proportion should non-negative and less than one")

	// error the sum proportion of treasures is more than one
	temp = Treasure{Address: sdk.AccAddress([]byte{0x00}), Proportion: sdk.NewDecWithPrec(99, 2)}
	treasures = []Treasure{*treasure1, *treasure2, temp}
	err = ValidateTreasures(treasures)
	require.Error(t, err)
	require.Contains(t, err.Error(), "the sum of treasure proportion should non-negative and less than one")
}

func TestSortTreasures(t *testing.T) {
	treasure1 := NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 0))
	treasure2 := NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 0))
	treasure3 := NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 0))
	treasure4 := NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 0))
	treasure5 := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 0))

	treasures := []Treasure{*treasure2, *treasure5, *treasure3, *treasure1, *treasure4}
	SortTreasures(treasures)
	for i, _ := range treasures {
		require.Equal(t, sdk.NewDec(int64(i)).Int64(), treasures[i].Proportion.Int64())
	}
}

func TestGetMapFromTreasures(t *testing.T) {
	treasure1 := NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 0))
	treasure2 := NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 0))
	treasure3 := NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 0))
	treasure4 := NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 0))
	treasure5 := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 0))

	excepted := make(map[string]Treasure, 0)
	excepted[treasure1.Address.String()] = *treasure1
	excepted[treasure2.Address.String()] = *treasure2
	excepted[treasure3.Address.String()] = *treasure3
	excepted[treasure4.Address.String()] = *treasure4
	excepted[treasure5.Address.String()] = *treasure5

	treasures := []Treasure{*treasure2, *treasure5, *treasure3, *treasure1, *treasure4}
	actual := GetMapFromTreasures(treasures)

	require.Equal(t, excepted, actual)
}

func TestGetTreasuresFromMap(t *testing.T) {
	treasure1 := NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 0))
	treasure2 := NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 0))
	treasure3 := NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 0))
	treasure4 := NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 0))
	treasure5 := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 0))

	temp := make(map[string]Treasure, 0)
	temp[treasure1.Address.String()] = *treasure1
	temp[treasure2.Address.String()] = *treasure2
	temp[treasure3.Address.String()] = *treasure3
	temp[treasure4.Address.String()] = *treasure4
	temp[treasure5.Address.String()] = *treasure5

	expected := []Treasure{*treasure5, *treasure4, *treasure3, *treasure2, *treasure1}
	actual := GetTreasuresFromMap(temp)
	require.Equal(t, expected, actual)
}

func TestIsTreasureDuplicated(t *testing.T) {
	treasure1 := NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 0))
	treasure2 := NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 0))
	treasure3 := NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 0))
	treasure4 := NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 0))
	treasure5 := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 0))
	treasures := []Treasure{*treasure5, *treasure4, *treasure3, *treasure2, *treasure1}

	// treasures is not Duplicated
	require.False(t, IsTreasureDuplicated(treasures))

	treasure5_duplicated := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 0))
	treasures = []Treasure{*treasure5, *treasure4, *treasure3, *treasure2, *treasure1, *treasure5_duplicated}
	require.True(t, IsTreasureDuplicated(treasures))
}

func TestInsertAndUpdateTreasures(t *testing.T) {
	treasure1 := NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 0))
	treasure2 := NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 0))
	treasure3 := NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 0))
	treasure4 := NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 0))
	treasure5 := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 0))

	// success insert treasures
	src := []Treasure{*treasure3, *treasure2, *treasure5}
	dst := []Treasure{*treasure1, *treasure4}
	expected := []Treasure{*treasure5, *treasure4, *treasure3, *treasure2, *treasure1}
	result := InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)

	// success update treasures
	src = []Treasure{*treasure3, *treasure2, *treasure5}
	treasure5_update := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(2, 0))
	dst = []Treasure{*treasure1, *treasure4, *treasure5_update}
	expected = []Treasure{*treasure5_update, *treasure4, *treasure3, *treasure2, *treasure1}
	result = InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)

	// success insert treasure
	src = []Treasure{*treasure3, *treasure2, *treasure5}
	dst = []Treasure{*treasure4}
	expected = []Treasure{*treasure5, *treasure4, *treasure3, *treasure2}
	result = InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)

	// success update treasure
	src = []Treasure{*treasure3, *treasure2, *treasure5}
	treasure5_update = NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(2, 0))
	dst = []Treasure{*treasure5_update}
	expected = []Treasure{*treasure5_update, *treasure3, *treasure2}
	result = InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)

	// success insert treasure from single treasure
	src = []Treasure{*treasure3}
	dst = []Treasure{*treasure4}
	expected = []Treasure{*treasure4, *treasure3}
	result = InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)

	// success update treasure from single treasure
	src = []Treasure{*treasure5}
	treasure5_update = NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(2, 0))
	dst = []Treasure{*treasure5_update}
	expected = []Treasure{*treasure5_update}
	result = InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)

	// success insert treasures from single treasure
	src = []Treasure{*treasure3}
	dst = []Treasure{*treasure4, *treasure2, *treasure1, *treasure5}
	expected = []Treasure{*treasure5, *treasure4, *treasure3, *treasure2, *treasure1}
	result = InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)

	// success update treasure from single treasure
	src = []Treasure{*treasure5}
	treasure5_update = NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(2, 0))
	dst = []Treasure{*treasure4, *treasure2, *treasure1, *treasure5_update}
	expected = []Treasure{*treasure5_update, *treasure4, *treasure2, *treasure1}
	result = InsertAndUpdateTreasures(src, dst)
	require.Equal(t, expected, result)
}

func TestDeleteTreasures(t *testing.T) {
	treasure1 := NewTreasure(sdk.AccAddress([]byte{0x01}), sdk.NewDecWithPrec(4, 0))
	treasure2 := NewTreasure(sdk.AccAddress([]byte{0x02}), sdk.NewDecWithPrec(3, 0))
	treasure3 := NewTreasure(sdk.AccAddress([]byte{0x03}), sdk.NewDecWithPrec(2, 0))
	treasure4 := NewTreasure(sdk.AccAddress([]byte{0x04}), sdk.NewDecWithPrec(1, 0))
	treasure5 := NewTreasure(sdk.AccAddress([]byte{0x05}), sdk.NewDecWithPrec(0, 0))

	src := make([]Treasure, 0)
	dst := make([]Treasure, 0)
	actual := make([]Treasure, 0)
	expected := make([]Treasure, 0)
	var err error
	testCases := []struct {
		msg     string
		prepare func()
		expPass bool
	}{
		{
			msg:     "delete one from one",
			expPass: true,
			prepare: func() {
				src = []Treasure{*treasure1}
				dst = []Treasure{*treasure1}
				expected = make([]Treasure, 0)
			},
		},
		{
			msg:     "delete one from one which is not exist",
			expPass: false,
			prepare: func() {
				src = []Treasure{*treasure1}
				dst = []Treasure{*treasure2}
			},
		},
		{
			msg:     "delete one from multi",
			expPass: true,
			prepare: func() {
				src = []Treasure{*treasure1, *treasure2}
				dst = []Treasure{*treasure1}
				expected = []Treasure{*treasure2}
			},
		},
		{
			msg:     "delete one from multi which is not exist",
			expPass: false,
			prepare: func() {
				src = []Treasure{*treasure3, *treasure2}
				dst = []Treasure{*treasure1}
			},
		},
		{
			msg:     "delete multi from multi result is empty",
			expPass: true,
			prepare: func() {
				src = []Treasure{*treasure1, *treasure2}
				dst = []Treasure{*treasure1, *treasure2}
				expected = []Treasure{}
			},
		},
		{
			msg:     "delete multi from multi result is  not empty",
			expPass: true,
			prepare: func() {
				src = []Treasure{*treasure1, *treasure2, *treasure4, *treasure5}
				dst = []Treasure{*treasure1, *treasure2}
				expected = []Treasure{*treasure5, *treasure4}
			},
		},
		{
			msg:     "delete multi from multi which has one is not exist",
			expPass: false,
			prepare: func() {
				src = []Treasure{*treasure1, *treasure2, *treasure4, *treasure5}
				dst = []Treasure{*treasure1, *treasure3}
			},
		},
		{
			msg:     "delete multi from multi which has multi is not exist",
			expPass: false,
			prepare: func() {
				src = []Treasure{*treasure1, *treasure2, *treasure5}
				dst = []Treasure{*treasure1, *treasure4, *treasure3}
			},
		},
		{
			msg:     "delete multi from multi which all is not exist",
			expPass: false,
			prepare: func() {
				src = []Treasure{*treasure1, *treasure2, *treasure5}
				dst = []Treasure{*treasure4, *treasure3}
			},
		},
		{
			msg:     "delete multi from one which all is not exist",
			expPass: false,
			prepare: func() {
				src = []Treasure{*treasure5}
				dst = []Treasure{*treasure4, *treasure3}
			},
		},
		{
			msg:     "delete multi from one ",
			expPass: false,
			prepare: func() {
				src = []Treasure{*treasure5}
				dst = []Treasure{*treasure5, *treasure3}
			},
		},
	}
	for _, tc := range testCases {
		tc.prepare()
		actual, err = DeleteTreasures(src, dst)
		if tc.expPass {
			require.NoError(t, err, tc.msg)
			require.Equal(t, expected, actual)
		} else {
			require.Error(t, err, tc.msg)
		}
	}
}
