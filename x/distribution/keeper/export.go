package keeper

// GetFeeCollectorName returns the name of fee_collector
func (k Keeper) GetFeeCollectorName() string {
	return k.feeCollectorName
}
