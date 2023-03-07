package iavl

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/okx/okbchain/libs/iavl/mock"
	"github.com/stretchr/testify/require"
)

func TestIsFastStorageStrategy_True_GenesisVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock.NewMockDB(ctrl)

	rIter := mock.NewMockIterator(ctrl)
	dbMock.EXPECT().ReverseIterator(gomock.Any(), gomock.Any()).Return(rIter, nil).Times(1)
	rIter.EXPECT().Close()
	rIter.EXPECT().Valid().Return(false)

	isFss := IsFastStorageStrategy(dbMock)
	require.Equal(t, true, isFss)
}

func TestIsFastStorageStrategy_False_GetFssVersionFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock.NewMockDB(ctrl)

	const iavlVersion = 3
	rIter := mock.NewMockIterator(ctrl)
	dbMock.EXPECT().ReverseIterator(gomock.Any(), gomock.Any()).Return(rIter, nil).Times(1)
	rIter.EXPECT().Close()
	rIter.EXPECT().Valid().Return(true)
	rIter.EXPECT().Key().Return(rootKeyFormat.Key(iavlVersion)).Times(1)

	dbMock.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("get fss version error")).Times(1)

	isFss := IsFastStorageStrategy(dbMock)
	require.Equal(t, false, isFss)
}

func TestIsFastStorageStrategy_False_IAVLNotEqualFSS(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock.NewMockDB(ctrl)

	const expectedVersion = 3
	fssVersion := fastStorageVersionValue + fastStorageVersionDelimiter + strconv.Itoa(expectedVersion)

	dbMock.EXPECT().Get(gomock.Any()).Return([]byte(fssVersion), nil).Times(1)
	rIter := mock.NewMockIterator(ctrl)
	dbMock.EXPECT().ReverseIterator(gomock.Any(), gomock.Any()).Return(rIter, nil).Times(1)
	rIter.EXPECT().Close().Times(1)
	rIter.EXPECT().Valid().Return(true).Times(1)
	rIter.EXPECT().Key().Return(rootKeyFormat.Key(expectedVersion + 1)).Times(1)

	isFss := IsFastStorageStrategy(dbMock)
	require.Equal(t, false, isFss)
}

func TestIsFastStorageStrategy_True_IAVLEqualFSS(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbMock := mock.NewMockDB(ctrl)

	const expectedVersion = 3
	fssVersion := fastStorageVersionValue + fastStorageVersionDelimiter + strconv.Itoa(expectedVersion)

	dbMock.EXPECT().Get(gomock.Any()).Return([]byte(fssVersion), nil).Times(1)
	rIter := mock.NewMockIterator(ctrl)
	dbMock.EXPECT().ReverseIterator(gomock.Any(), gomock.Any()).Return(rIter, nil).Times(1)
	rIter.EXPECT().Close().Times(1)
	rIter.EXPECT().Valid().Return(true).Times(1)
	rIter.EXPECT().Key().Return(rootKeyFormat.Key(expectedVersion)).Times(1)

	isFss := IsFastStorageStrategy(dbMock)
	require.Equal(t, true, isFss)
}
