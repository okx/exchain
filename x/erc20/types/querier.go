package types

type ContractByDenomRequest struct {
	Denom string `json:"denom,omitempty"`
}

type DenomByContractRequest struct {
	Contract string `json:"contract,omitempty"`
}
