package rest

import "github.com/okex/exchain/x/wasm/types"

type contractInfo struct {
	Address string `json:"address,omitempty"`
	CodeId  uint64 `json:"code_id,omitempty"`
	Creator string `json:"creator,omitempty"`
	Admin   string `json:"admin,omitempty"`
	Label   string `json:"label,omitempty"`
}

func fromGrpcContractInfo(resp *types.QueryContractInfoResponse) (info *contractInfo) {
	if resp == nil {
		return
	}
	info.Address = resp.Address
	info.CodeId = resp.CodeID
	info.Creator = resp.Creator
	info.Admin = resp.Admin
	info.Label = resp.Label

	return
}
