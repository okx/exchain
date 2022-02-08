package state

// Executes block's transactions on proxyAppConn.
// Returns a list of transaction results and updates to the validator set
func execBlockOnProxyAppPartConcurrent(context *executionTask) (*ABCIResponses, error) {
	block := context.block

	
	abciResponses := NewABCIResponses(block)

	// Execute transactions and get hash.


	// Begin block


	// Run txs of block.


	// End block.


	return abciResponses, nil
}