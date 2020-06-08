package types

// Query endpoints supported by the margin querier
const (
	// TODO: Describe query parameters, update <action> with your query
	// Query<Action>    = "<action>"
	QueryAccount   = "account"
	QueryProducts  = "products"
	QuerySaving    = "saving"
	QueryBorrowing = "borrowing"

	QueryParameters = "params"
)

/*
Below you will be able how to set your own queries:


// QueryResList Queries Result Payload for a query
type QueryResList []string

// implement fmt.Stringer
func (n QueryResList) String() string {
	return strings.Join(n[:], "\n")
}

*/
