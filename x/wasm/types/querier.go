package types

type QueryAddressWhitelistResponse struct {
	Whitelist []string `json:"whitelist,omitempty"`
}

func NewQueryAddressWhitelistResponse(whitelist []string) *QueryAddressWhitelistResponse {
	return &QueryAddressWhitelistResponse{
		Whitelist: whitelist,
	}
}
