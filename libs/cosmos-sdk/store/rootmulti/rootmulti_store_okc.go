package rootmulti

import "github.com/okex/exchain/libs/cosmos-sdk/types"

func (rs *Store) getCommitIDWithSupportVersion(infos map[string]storeInfo, name string) types.CommitID {
	info, ok := infos[name]
	if !ok {
		return types.CommitID{Version: 20}
	}
	return info.Core.CommitID
}
