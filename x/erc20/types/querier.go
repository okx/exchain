package types

type ContractByDenomRequest struct {
	Denom string `json:"denom,omitempty"`
}

type DenomByContractRequest struct {
	Contract string `json:"contract,omitempty"`
}

type CurrentContractTemplate struct {
	Proxy     []byte `json:"proxy"`
	Implement []byte `json:"implement"`
}
