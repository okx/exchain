## storeKey

### staking 

| key                                 | value                              | number(key) | value details                   | value size | clean up                                                | 备注                    |
| ----------------------------------- | ---------------------------------- | ----------- | ------------------------------- | ---------- | ------------------------------------------------------- | ----------------------- |
| 0x11+OperatorAddr                   | Power                              | N/A         | 无数组                          | <1k        | 交易清理                                                | LastValidatorsPower     |
| 0x12                                | Total Power                        | 1           | 无数组                          | <1k        | 无需清零                                                | LastTotalPower          |
| 0x21+OperatorAddr                   | x/staking/types.Validator          | N/A         | 无数组                          | <1k        | 交易清理                                                | Validator               |
| 0x22+ConsensusAddr                  | OperatorAddr                       | N/A         | 无数组                          | <1k        | 交易清理                                                | Validator               |
| 0x23+Power+^OperatorAddr            | OperatorAddr                       | N/A         | 无数组                          | <1k        | 交易清理                                                | Validator               |
| 0x43+Time                           | x/staking/types.[]ValAddress       | N/A         | 数组长度最多为validator集合总数 | \>1k       | 每个区块清理到期                                        | ValidatorQueue          |
| 0x51+DelegatorAddr+ValidatorAddr    | x/staking/types.Shares             | N/A         | 无数组                          | <1k        | 取消投票时清理                                          | SharesKey                 |
| 0x52+DelegatorAddr                  | x/staking/types.Delegator          | N/A         | 无数组                          | <1k        | 当全部解委托tokens时清理                                | DelegatorKey            |
| 0x53+DelegatorAddr                  | x/staking/types.UndelegationInfo   | N/A         | 无数组                          | <1k        | 当解委托到期时清理                                      | UnDelegationInfoKey     |
| 0x54+Time                           | x/staking/[]types.UndelegationInfo | N/A         | 有数组                          | 可能会>1k  | 当[]UndelegationInfo中的UndelegationInfo都到期时        | UnDelegateQueueKey      |
| 0x55+ProxyAddr+DelegatorAddr        | []byte("")                         | N/A         | 无数组                          | <1k       | 当delegator发起解代理tx时                          | ProxyKey   |
| 0x60                                | x/staking/[]sdk.ValAddress         | 1           | 有数组                          | 可能会>1k  | 当存在要强制剔除出块集合的validator时，EndBlock时候清理 | ValidatorAbandonedKey   |




### params 

| key                       | value                  | number(key) | value details | value size | clean up | 备注   |
| ------------------------- | ---------------------- | ----------- | ------------- | ---------- | -------- | ------ |
| ParamsSubspace("staking") | x/staking/types.Params | 1           | 无数组        | <1k        | 无需清零 | Params |


​     