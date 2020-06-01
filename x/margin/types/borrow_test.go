package types

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestBorrowInfoListSort(t *testing.T) {
	borrowInfo0 := &BorrowInfo{
		BlockHeight: 11,
		Rate:        sdk.NewDec(100),
	}

	borrowInfo1 := &BorrowInfo{
		BlockHeight: 10,
		Rate:        sdk.NewDec(100),
	}

	borrowInfo2 := &BorrowInfo{
		BlockHeight: 9,
		Rate:        sdk.NewDec(101),
	}

	wanted := BorrowInfoList{borrowInfo2, borrowInfo1, borrowInfo0}
	borrowInfoList := BorrowInfoList{borrowInfo0, borrowInfo1, borrowInfo2}
	require.NotEqual(t, wanted, borrowInfoList)
	sort.Sort(borrowInfoList)
	require.Equal(t, wanted, borrowInfoList)

}
