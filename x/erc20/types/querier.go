package types

// Supported endpoints
const (
	QueryParameters = "params"
	QueryAllMapping = "all-mapping"
)

type QueryResAllMapping struct {
	Mapping map[string]string `json:"mapping"`
}
