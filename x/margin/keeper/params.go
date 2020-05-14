package keeper

// TODO: Define if your module needs Parameters, if not this can be deleted

//// GetParams returns the total set of margin parameters.
//func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
//	k.paramspace.GetParamSet(ctx, &params)
//	return params
//}
//
//// SetParams sets the margin parameters to the param space.
//func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
//	k.paramspace.SetParamSet(ctx, &params)
//}
//

//var (
//	keyMarginDeposit = []byte("MarginDeposit")
//)
//
//// Params defines param object
//type Params struct {
//	MarginDeposit sdk.DecCoin `json:"list_fee"`
//}
//
//// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
//func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
//	return params.ParamSetPairs{
//		{Key: keyMarginDeposit, Value: &p.MarginDeposit},
//	}
//}
//
//// ParamKeyTable for auth module
//func ParamKeyTable() params.KeyTable {
//	return params.NewKeyTable().RegisterParamSet(&Params{})
//}
//
//// DefaultParams returns a default set of parameters.
//func DefaultParams() *Params {
//	defaultListFee := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultFeeList))
//	defaultTransferOwnershipFee := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultFeeTransferOwnership))
//	defaultDelistMinDeposit := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultDelistMinDeposit))
//	return &Params{
//		ListFee:                defaultListFee,
//		TransferOwnershipFee:   defaultTransferOwnershipFee,
//		DelistMaxDepositPeriod: time.Hour * 24,
//		DelistMinDeposit:       sdk.DecCoins{defaultDelistMinDeposit},
//		DelistVotingPeriod:     time.Hour * 72,
//		WithdrawPeriod:         DefaultWithdrawPeriod,
//	}
//}
//
//// String implements the stringer interface.
//func (p Params) String() string {
//	return fmt.Sprintf("Params: \nDexListFee:%s\nTransferOwnershipFee:%s\nDelistMaxDepositPeriod:%s\n"+
//		"DelistMinDeposit:%s\nDelistVotingPeriod:%s\nWithdrawPeriod:%d\n",
//		p.ListFee, p.TransferOwnershipFee, p.DelistMaxDepositPeriod, p.DelistMinDeposit, p.DelistVotingPeriod, p.WithdrawPeriod)
//}
