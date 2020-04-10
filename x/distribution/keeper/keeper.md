## StoreKey

**声明**：
* sdk = cosmos/cosmos-sdk/types
* types = x/distribution/types
* 暂时保留：cosmos原有存在的，以后迭代开发会复用，暂时保留以供后面版本使用。

### 1. storeKey (distr)
|key|value|number(key)|value details|value size|clean up|备注|
|---|---|-------------|--------------|---------|--------|---|
|  FeePoolKey  | types.FeePool   | 1 |有数组，随币种种类增长| <1k | 不清理 | 基金池，暂时保留 |
|  ProposerKey | sdk.ConsAddress | 1 | 无数组           | <1k | 不清理 | 保存出块者地址   |
|  ValidatorOutstandingRewardsPrefix:${valAddr} | types.ValidatorOutstandingRewards | 验证者个数，默认21|有数组，随币种种类增长| 币种太多会超1k | 分红到账后清空 | 本周期超级节点所有奖励（包含委托者奖励） |     ||   
|  DelegatorWithdrawAddrPrefix:${delAddr} | sdk.AccAddress |     委托量    |无数组 |<1k| 只更新 | 用户取款地址 |    
|  DelegatorStartingInfoPrefix:${valAddr}:${delAddr} | types.DelegatorStartingInfo |     委托量     |有数组，随币种种类增长|币种太多会超1k| 分红到账后清空 |  委托开始时间，暂时保留      ||   
|  ValidatorHistoricalRewardsPrefix:${valAddr} | []byte sdk.ConsAddress |     验证者个数     | 有数组，随币种种类增长|币种太多会超1k| 分红到账后清空 | 出块者历史奖励，暂时保留    | | 
|  ValidatorCurrentRewardsPrefix:${valAddr} | types.ValidatorCurrentRewards |     验证者个数, 默认21     | 有数组，随币种种类增长|币种太多会超1k | 分红到账后清空|  委托者奖励池     |  | 
|  ValidatorAccumulatedCommissionPrefix:${valAddr} | types.ValidatorAccumulatedCommission |     验证者个数21     | 有数组，随币种种类增长 |币种太多会超1k | 分红到账后清空 | 委托费池    |  | 
|  ValidatorSlashEventPrefix:${valAddr} | types.ValidatorSlashEvent |     惩罚事件个数     |  无数组 |<1k | 执行后清理 | 暂时保留   | 
|  ParamStoreKeyCommunityTax | sdk.Dec |     1     | 无数组 | <1k | 不清理 |  基金池奖励比例， 暂时保留   | 
|  ParamStoreKeyBaseProposerReward | sdk.Dec |     1     | 无数组|<1k | 不清理|  出块者基本奖励，暂时保留    | 
|  ParamStoreKeyBonusProposerReward | sdk.Dec |     1     | 无数组|<1k |不清理|  出块者额外奖励，暂时保留   |
|  ParamStoreKeyWithdrawAddrEnabled | sdk.Dec |     1     | 无数组| <1k |不清理| 分红地址是否可修改配置项   |



## 备注
```sh
FeePoolKey                        = []byte{0x00} // key for global distribution state
ProposerKey                       = []byte{0x01} // key for the proposer operator address
ValidatorOutstandingRewardsPrefix = []byte{0x02} // key for outstanding rewards

DelegatorWithdrawAddrPrefix          = []byte{0x03} // key for delegator withdraw address
DelegatorStartingInfoPrefix          = []byte{0x04} // key for delegator starting info
ValidatorHistoricalRewardsPrefix     = []byte{0x05} // key for historical validators rewards / stake
ValidatorCurrentRewardsPrefix        = []byte{0x06} // key for current validator rewards
ValidatorAccumulatedCommissionPrefix = []byte{0x07} // key for accumulated validator commission
ValidatorSlashEventPrefix            = []byte{0x08} // key for validator slash fraction

ValidatorSnapshootPrefix  = []byte{0x80} // okdex, key for epoch validator snapshoot
DelegationSnapshootPrefix = []byte{0x81} //key for epoch delegation snapshoot

ParamStoreKeyCommunityTax        = []byte("communitytax")
ParamStoreKeyBaseProposerReward  = []byte("baseproposerreward")
ParamStoreKeyBonusProposerReward = []byte("bonusproposerreward")
ParamStoreKeyWithdrawAddrEnabled = []byte("withdrawaddrenabled")
```

