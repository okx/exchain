## StoreKey

### 1. ordersStoreKey 
| key                                    | value           |  number(keys)                 | value detail                                                                                                     | value size | clean up              | 备注                                            |
|----------------------------------------|-----------------|----------------------------|------------------------------------------------------------------------------------------------------------------|------------|-----------------------|-------------------------------------------------|
| orderNum:block({blockHeight})          | int64           |  区块高度                  | 每个区块的订单数量                                                                                                                 | <1k        | 每区块删除3天前的数据 | 某一区块的order数量                             |
| $orderid                               | order.Order     |  所有订单数量              | 每个订单                                                                                                                 | <1k        | 每区块删除3天前的数据 | 某一区块更新过的订单id列表                      |
| ${product}                              | []DepthBookItem |  币对数量           | 每个币对一个深度表，  假设某币对价格精度为当前价格的万分之一，<br>正常挂单都在当前价格+-5%以内，<br>则一个币对深度表中含有1000个表项   | >1k        |                       | 某一币对当前的深度表 <br>DepthBookItem数组      |
|** {product}-{price}-{side}             | []string        |  币对数量*价格可能取值数量 | 数组长度取决于某币对某价格的买/卖单数量<br> 平均值不好预估，峰值无上限                                           | >1k        |                       | 某一币对在某一价位的所有买单或卖单的订单id列表  |
| ${product}                             | sdk.Dec         |  币对数量                  |  当前价格                                                                                                                | <1k        |                       | 某一币对的最近成交价                            |
|expireBlockHeight:block(${blockHeight}) | []int64         |  区块高度                  | 在key高度，value里多少个区块的单是过期的                                                                                                                  |  < 1k      |                        |      某一区块应该处理的order过期的block        |
| productLockMap                         |types.ProductLockMap| 1                     | 所有被锁的pair
## Http api

| url              | method | 读key                         | 写key                                                                                                                                                       |
|------------------|--------|-------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| /order/new       | POST   | orderNum:block({blockHeight}) | orderNum:block({blockHeight})<br>ID{0-blockHeight}-${Num}<br>depthbook:{product}<br>{product}-{price}-{side}<br>lastprice:{product}<br> |
| /order/cancel    | POST   | ID{0-blockHeight}-${Num}      | ID{0-blockHeight}-${Num}<br>depthbook:{product}<br>{product}-{price}-{side}                                                                                     |
| /order/depthbook | GET    | depthbook:{product}           |                                                                                                                                                             |
| /order/{orderID} | GET    | ID{0-blockHeight}-${Num}      |                                                                                                                                                             |
