
1. v 节点内存中维护一个订单簿 B，一个订单队列 Q
2. v 节点通过 rpc 和 p2p 收取订单 o，将 o 按到达顺序加入到 Q 中，不进行撮合，B 不变
3. 当轮到此 v 节点出块时，从 Q 中按顺序读取订单 o，对每一个 o
   1. 读取 o 的链上最新状态(已撮合数量，是否已取消等)
   2. 将 o 与 B 进行撮合，根据撮合结果更新 B 并产生撮合交易 t，撮合结果涉及到的订单状态根据撮合数量进行冻结
4. 删除 Q 中被读取的订单，将所有的 t 打包到区块
4. 每执行完一个区块进行 mempool Update 时
   1. 如果 v 是当前块的出块者
      1. 解冻撮合过程中被冻结的订单
      2. 删除 Q 中已出块的订单
   2. 使用区块产生的 evm event 对 B 中的订单更新状态(已撮合数量，是否已取消等)



